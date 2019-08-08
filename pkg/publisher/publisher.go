package publisher

import (
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/consumer"
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/metric"
)

type MetricsPublisher struct {
	consumers []consumer.Consumer
}

func (p *MetricsPublisher) SubscribeConsumers(consumer ...consumer.Consumer) {
	if consumer == nil {
		return
	}

	p.consumers = append(p.consumers, consumer...)
}

func (p *MetricsPublisher) Publish(metrics []metric.Metric) {
	if len(metrics) == 0 {
		return
	}

	for _, c := range p.consumers {
		c.Process(metrics)
	}
}
