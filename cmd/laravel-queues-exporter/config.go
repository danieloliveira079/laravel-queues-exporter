package main

import (
	"flag"
	"log"
	"os"
	"strconv"
)

// helper to get string from ENV with default
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// help to get bool from ENV with default
func getEnvBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		var boolValue bool
		var err error

		if boolValue, err = strconv.ParseBool(value); err != nil {
			log.Printf("WARN: environment variable \"%s\"=\"%v\" can not be parsed to bool: %v", key, value, err)
			return fallback
		}
		return boolValue
	}
	return fallback
}

// help to get int from ENV with default
func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		var err error
		var intValue int64
		if intValue, err = strconv.ParseInt(value, 10, 8); err != nil {
			log.Printf("WARN: environment variable \"%s\"=\"%v\" can not be parsed to integer: %v", key, value, err)
			return fallback
		}
		return int(intValue)
	}
	return fallback
}

type globalConfig struct {
	redisHost       string
	redisPort       string
	redisDB         int
	statsdHost      string
	statsdPort      string
	metricsPrefix   string
	collectInterval int
	queuesNames     string
}

var config globalConfig

func getConfig() {
	flag.StringVar(&config.redisHost, "redis-host", getEnv("REDIS_HOST", "127.0.0.1"), "Redis host where queues are stored")
	flag.StringVar(&config.redisPort, "redis-port", getEnv("REDIS_PORT", "6379"), "Redis target port open for connections")
	flag.IntVar(&config.redisDB, "redis-db", getEnvInt("REDIS_DB", 0), "Redis DB used by Laravel")
	flag.StringVar(&config.statsdHost, "statsd-host", getEnv("STATSD_HOST", "127.0.0.1"), "StatsD target to where metrics must be sent")
	flag.StringVar(&config.statsdPort, "statsd-port", getEnv("STATSD_PORT", "8125"), "StatsD target port open for connections")
	flag.StringVar(&config.metricsPrefix, "metrics-prefix", getEnv("METRICS_PREFIX", "exporter"), "Prefix to be added to every metric")
	flag.IntVar(&config.collectInterval, "collect-interval", getEnvInt("SCAN_INTERVAL", 60), "Interval in seconds between each metrics collect")
	flag.StringVar(&config.queuesNames, "queues-names", getEnv("QUEUES_NAMES", ""), "Names of the queues to be scanned separated by comma. I.e: queue1,queue2")

	flag.Parse()
}
