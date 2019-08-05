package forwarder

import (
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/metric"
	"log"
)

type Log struct {
}

func (s *Log) Forward(metrics []metric.Metric) {
	for _, m := range metrics {
		log.Printf("%s %v", m.Name, m.Value)
	}
}
