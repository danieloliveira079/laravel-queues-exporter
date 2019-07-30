package main

import (
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/exporter/redis"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

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
		ConnConfig: connectionConfig,
		Dispatcher: dispatcher,
	}
	extractor, err := redis.NewRedisExtractor(extractorConfig)
	if err != nil {
		log.Fatal(err)
	}

	exporterConfig := redis.RedisExporterConfig{
		ConnectionConfig: connectionConfig,
		QueueNames:       config.queuesNames,
		CheckInterval:    config.checkInterval,
		Connector:        connector,
		Extractor:        extractor,
	}
	exporter, err := redis.NewRedisExporter(exporterConfig)

	if err != nil {
		log.Fatal(err)
	}

	exporter.Scan()

	for {
		select {
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
