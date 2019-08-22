package stdout

import (
	"fmt"
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/config"
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/metric"
)

type Stdout struct {
}

func New(config *config.AppConfig) (*Stdout, error) {
	return new(Stdout), nil
}

func (s *Stdout) Process(metrics []metric.Metric) {
	for _, m := range metrics {
		if len(m.Name) > 0 {
			fmt.Println(m.Name, m.Value)
		}
	}
}
