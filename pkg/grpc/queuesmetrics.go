package grpc

import (
	"context"
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/metric"
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/proto"
	"log"
)

type queuesMetricsServer struct {
	metrics []metric.Metric
}

func (q *queuesMetricsServer) Get(ctx context.Context, query *queuesmetrics.Query) (*queuesmetrics.Metrics, error) {
	log.Println("grpc request GET:", query.Queues)
	response := &queuesmetrics.Metrics{}

	for _, m := range q.metrics {
		response.Metric = append(response.Metric, &queuesmetrics.Metric{
			Queue: m.Name,
			Jobs:  m.Value,
		})
	}
	return response, nil
}

func (q *queuesMetricsServer) Process(metrics []metric.Metric) {
	q.metrics = metrics
}
