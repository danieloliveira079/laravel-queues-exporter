package consumer

import (
	"fmt"
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/metric"
)

type Stdout struct {
}

func (s *Stdout) Process(metrics []metric.Metric) {
	for _, m := range metrics {
		fmt.Println(m.Name, m.Value)
	}
}
