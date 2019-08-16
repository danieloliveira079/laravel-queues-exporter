package log

import (
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/config"
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/metric"
	"log"
)

type Log struct {
}

func New(config *config.AppConfig) (*Log, error) {
	return &Log{}, nil
}

func (l *Log) Process(metrics []metric.Metric) {
	for _, m := range metrics {
		log.Printf("%s %v", m.Name, m.Value)
	}
}
