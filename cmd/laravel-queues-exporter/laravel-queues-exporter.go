package main

import (
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/exporter/redis"
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
	exporter := redis.NewRedisExporter(
		config.redisHost,
		config.redisPort,
		config.redisDB,
		config.checkInterval,
		config.queuesNames,
		nil,
	)

	exporter.Scan()

	for {
		select {
		case signalReceived := <-signals:
			switch signalReceived {
			case syscall.SIGTERM:
			case syscall.SIGINT:
				exporter.Stop(done)
				<-done
				runtime.GC()
				os.Exit(0)
			}
		}
	}
}
