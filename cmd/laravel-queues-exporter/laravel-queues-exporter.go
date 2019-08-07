package main

import (
	"flag"
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/config"
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

type Exporter interface {
	Run(collected chan []metric.Metric)
}
type ExporterBuilder interface {
	Build(appConfig config.AppConfig) (Exporter, error)
}

func main() {
	appConfig := getConfig()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGINT)
	done := make(chan os.Signal, 1)

	redisExporterBuilder := new(redis.RedisExporterBuilder)
	exporter, err := redisExporterBuilder.Build(appConfig)
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
			//TODO abstract forwarder and use concurrency
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

func getConfig() *config.AppConfig {
	var appConfig config.AppConfig

	flag.StringVar(&appConfig.RedisHost, "redis-host", config.GetEnv("REDIS_HOST", "127.0.0.1"), "Redis host where queues are stored")
	flag.StringVar(&appConfig.RedisPort, "redis-port", config.GetEnv("REDIS_PORT", "6379"), "Redis target port open for connections")
	flag.IntVar(&appConfig.RedisDB, "redis-db", config.GetEnvInt("REDIS_DB", 0), "Redis DB used by Laravel")
	flag.StringVar(&appConfig.StatsDHost, "statsd-host", config.GetEnv("STATSD_HOST", "127.0.0.1"), "StatsD target to where metrics must be sent")
	flag.StringVar(&appConfig.StatsDPort, "statsd-port", config.GetEnv("STATSD_PORT", "8125"), "StatsD target port open for connections")
	flag.StringVar(&appConfig.MetricsPrefix, "metrics-prefix", config.GetEnv("METRICS_PREFIX", "exporter"), "Prefix to be added to every metric")
	flag.IntVar(&appConfig.CollectInterval, "collect-interval", config.GetEnvInt("SCAN_INTERVAL", 60), "Interval in seconds between each metrics collect")
	flag.StringVar(&appConfig.QueuesNames, "queues-names", config.GetEnv("QUEUES_NAMES", ""), "Names of the queues to be scanned separated by comma. I.e: queue1,queue2")

	flag.Parse()
	return &appConfig
}
