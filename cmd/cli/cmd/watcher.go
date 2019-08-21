package cmd

import (
	"context"
	queuesmetrics "github.com/danieloliveira079/laravel-queues-exporter/pkg/proto"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)

type Watcher struct {
	client        queuesmetrics.QueuesMetricsClient
	metrics       *queuesmetrics.Metrics
	mutex         *sync.Mutex
	watchInterval time.Duration
}

func NewWatcher(serverAddress string, watchInterval time.Duration) (*Watcher, error) {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := queuesmetrics.NewQueuesMetricsClient(conn)
	return &Watcher{
		client:        client,
		mutex:         &sync.Mutex{},
		watchInterval: watchInterval,
	}, nil
}

func (w *Watcher) Run() {
	ticker := time.NewTicker(w.watchInterval)
	w.getData()

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				w.getData()
			}
		}
	}()
}

func (w *Watcher) getData() {
	metrics, err := w.client.Get(context.Background(), &queuesmetrics.Query{Queues: "dummy"})
	if err != nil {
		log.Fatal(err)
	}
	w.mutex.Lock()
	w.metrics = metrics
	w.mutex.Unlock()
}

func (w *Watcher) GetMetrics() map[string]*queuesmetrics.Metric {
	metrics := map[string]*queuesmetrics.Metric{}
	if w.metrics != nil {
		for _, m := range w.metrics.Metric {
			if len(m.Queue) > 0 {
				metrics[m.Queue] = m
			}
		}
	}

	return metrics
}
