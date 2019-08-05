package main

import (
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/exporter/redis"
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/forwarder"
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/metric"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

type Forwarder interface {
	Forward(metrics []metric.Metric)
}

func main() {
	getConfig()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGINT)
	done := make(chan os.Signal, 1)

	connectionConfig := redis.ConnectionConfig{
		Host: config.redisHost,
		Port: config.redisPort,
		DB:   config.redisDB,
	}

	connectorConfig := redis.ConnectorConfig{
		ConnConfig: connectionConfig,
	}
	connector, err := redis.NewRedisConnector(connectorConfig)
	if err != nil {
		log.Fatal(err)
	}

	dispatcher, err := redis.NewRedisCommandDispatcher(connector)
	if err != nil {
		log.Fatal(err)
	}

	extractorConfig := redis.ExtractorConfig{
		Dispatcher: dispatcher,
	}
	extractor, err := redis.NewRedisExtractor(extractorConfig)
	if err != nil {
		log.Fatal(err)
	}

	exporterConfig := redis.ExporterConfig{
		ConnectionConfig: connectionConfig,
		QueueNames:       config.queuesNames,
		CollectInterval:  config.collectInterval,
	}
	exporter, err := redis.NewRedisExporter(exporterConfig, connector, extractor)

	if err != nil {
		log.Fatal(err)
	}

	collected := make(chan []metric.Metric)
	exporter.Run(collected)

	forwarders := []Forwarder{
		&forwarder.Stdout{},
		&forwarder.Log{},
	}

	for {
		select {
		case metrics := <-collected:
			for _, f := range forwarders {
				f.Forward(metrics)
			}
		case signalReceived := <-signals:
			switch signalReceived {
			case syscall.SIGTERM:
			case syscall.SIGINT:
				exporter.Stop(done) //TODO Implement ticker to avoid waiting forever
				<-done
				runtime.GC()
				os.Exit(0)
			}
		}
	}
}
