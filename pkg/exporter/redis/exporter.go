package redis

import (
	"errors"
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/metric"
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/queue"
	"log"
	"os"
	"strings"
	"time"
)

type Exporter struct {
	Config           *ExporterConfig
	Extractor        Extractor
	Connector        Connector
	Queues           []*RedisQueue
	interruptCollect bool
}

type ExporterConfig struct {
	ConnectionConfig *ConnectionConfig
	CollectInterval  int
	QueueNames       string
}

type Extractor interface {
	ListAllQueuesFromDB() ([]*RedisQueue, error)
	CountJobsForQueues(queue []*RedisQueue) error
	SetQueueTypeForQueues(queues []*RedisQueue)
}

type Connector interface {
	Connect() (err error)
	Close() (err error)
	Do(command string, args ...interface{}) (results interface{}, err error)
}

func NewRedisExporter(config *ExporterConfig, connector Connector, extractor Extractor) (*Exporter, error) {
	if connector == nil {
		return nil, errors.New("connector can't be nil")
	}

	if extractor == nil {
		return nil, errors.New("extractor can't be nil")
	}

	return &Exporter{
		Config:    config,
		Connector: connector,
		Extractor: extractor,
	}, nil
}

func (xp *Exporter) Stop(done chan os.Signal) {
	log.Println("Stopping exporter")
	xp.interruptCollect = true
	err := xp.CloseConnector()
	if err != nil {
		log.Println("error closing connector:", err)
	}
	log.Println("Exporter stopped")
	close(done)
}

func (xp *Exporter) CloseConnector() error {
	return xp.Connector.Close()
}

func (xp *Exporter) Run(collected chan []metric.Metric) {
	log.Println("Starting exporter")
	ticker := time.NewTicker(time.Duration(xp.Config.CollectInterval) * time.Second)

	go func() {
		defer ticker.Stop()
		defer close(collected)

		for _ = range ticker.C {
			if xp.interruptCollect == true {
				log.Println("Stopping scanner")
				break
			}

			collected <- xp.Collect()
		}
	}()
}

func (xp *Exporter) Collect() []metric.Metric {
	var err error

	xp.Queues, err = xp.SelectQueuesToScan()
	if err != nil {
		log.Fatal(err)
	}

	xp.SetQueuesType(xp.Queues)

	err = xp.CountJobsForQueues(xp.Queues)
	if err != nil {
		log.Fatal(err)
	}

	metrics := []metric.Metric{}
	if len(xp.Queues) == 0 {
		emptyMetrics := make([]metric.Metric, 1)
		return emptyMetrics
	}

	for _, q := range xp.Queues {
		metrics = append(metrics, metric.Metric{
			Name:  q.Name(),
			Value: q.GetCurrentJobsCount(),
		})
	}

	return metrics
}

func (xp *Exporter) CountJobsForQueues(queues []*RedisQueue) error {
	return xp.Extractor.CountJobsForQueues(queues)
}

func (xp *Exporter) SelectQueuesToScan() ([]*RedisQueue, error) {
	var err error
	queueItems := []*RedisQueue{}

	if len(xp.Config.QueueNames) > 0 {
		queueItems = parsedQueuesNames(xp.Config.QueueNames)
	} else {
		queueItems, err = xp.Extractor.ListAllQueuesFromDB()
	}

	return queueItems, err
}

func parsedQueuesNames(queueNames string) []*RedisQueue {
	queueItems := []*RedisQueue{}
	names := strings.Split(queueNames, ",")
	for _, n := range names {
		queueItems = append(queueItems, &RedisQueue{queueItem: &queue.QueueItem{Name: n}})
	}

	return queueItems
}

func (xp *Exporter) SetQueuesType(queues []*RedisQueue) {
	xp.Extractor.SetQueueTypeForQueues(queues)
}
