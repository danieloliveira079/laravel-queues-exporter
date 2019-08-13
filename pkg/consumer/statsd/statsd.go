package statsd

import (
	"fmt"
	"github.com/DataDog/datadog-go/statsd"
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/config"
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/metric"
	"log"
)

type StatsD struct {
	host          string
	port          string
	metricsPrefix string
}

func New(config *config.AppConfig) *StatsD {
	return &StatsD{
		host:          config.StatsDHost,
		port:          config.StatsDPort,
		metricsPrefix: config.MetricsPrefix,
	}
}

func (s *StatsD) Process(metrics []metric.Metric) {
	client, err := statsd.New(fmt.Sprintf("%s:%s", s.host, s.port))

	if err != nil {
		log.Println(err)
		return
	}

	for _, metric := range metrics {
		err = client.Gauge(metric.WithPrefix(s.metricsPrefix), metric.ValueToFloat64(), []string{}, 1)
	}
}
