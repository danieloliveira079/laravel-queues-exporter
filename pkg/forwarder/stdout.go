package forwarder

import (
	"fmt"
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/metric"
)

type Stdout struct {
}

func (s *Stdout) Forward(metrics []metric.Metric) {
	for _, m := range metrics {
		fmt.Println(m.Name, m.Value)
	}
}
