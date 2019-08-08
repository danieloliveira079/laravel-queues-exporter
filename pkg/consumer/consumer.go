package consumer

import "github.com/danieloliveira079/laravel-queues-exporter/pkg/metric"

type Consumer interface {
	Process(metrics []metric.Metric)
}
