package config

import (
	"log"
	"os"
	"strconv"
)

type AppConfig struct {
	RedisHost       string
	RedisPort       string
	RedisDB         int
	StatsDHost      string
	StatsDPort      string
	MetricsPrefix   string
	CollectInterval int
	QueuesNames     string
	ExportTo        string
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetEnvBool(key string, fallback bool) bool {
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

func GetEnvInt(key string, fallback int) int {
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
