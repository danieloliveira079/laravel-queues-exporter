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
	client        *statsd.Client
	//TODO add init message
}

func New(config *config.AppConfig) (*StatsD, error) {
	conn := fmt.Sprintf("%s:%s", config.StatsDHost, config.StatsDPort)
	client, err := statsd.New(conn)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &StatsD{
		host:          config.StatsDHost,
		port:          config.StatsDPort,
		metricsPrefix: config.MetricsPrefix,
		client:        client,
	}, nil
}

func (s *StatsD) Process(metrics []metric.Metric) {
	for _, metric := range metrics {
		err := s.client.Gauge(metric.WithPrefix(s.metricsPrefix), metric.ValueToFloat64(), []string{}, 1)
		if err != nil {
			log.Println(err)
		}
	}
}
