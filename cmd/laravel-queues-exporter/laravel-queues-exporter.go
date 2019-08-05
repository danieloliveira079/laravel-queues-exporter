package main

import (
	"fmt"
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
		Dispatcher: dispatcher,
	}
	extractor, err := redis.NewRedisExtractor(extractorConfig)
	if err != nil {
		log.Fatal(err)
	}

	exporterConfig := redis.ExporterConfig{
		ConnectionConfig: connectionConfig,
		QueueNames:       config.queuesNames,
		ScanInterval:     config.scanInterval,
	}
	exporter, err := redis.NewRedisExporter(exporterConfig, connector, extractor)

	if err != nil {
		log.Fatal(err)
	}

	collected := make(chan string)
	exporter.Run(collected)

	for {
		select {
		case forwardChan := <-collected:
			fmt.Println(forwardChan)
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
