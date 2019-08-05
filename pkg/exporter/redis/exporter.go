package redis

import (
	"errors"
	"fmt"
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/queue"
	"log"
	"os"
	"strings"
	"time"
)

type Exporter struct {
	Config        ExporterConfig
	interruptScan bool
}

type ExporterConfig struct {
	ConnectionConfig ConnectionConfig
	ScanInterval     int
	QueueNames       string
	Extractor        Extractor
	Connector        Connector
}

type Extractor interface {
	ListAllQueuesFromDB() ([]*RedisQueue, error)
	CountJobsForQueue(queue *RedisQueue) error
	SetQueueTypeForQueues(queues []*RedisQueue)
}

type Connector interface {
	Connect() (err error)
	Close() (err error)
	Do(command string, args ...interface{}) (results interface{}, err error)
}

func NewRedisExporter(config ExporterConfig) (*Exporter, error) {
	if config.Connector == nil {
		return nil, errors.New("connector can't be nil")
	}

	if config.Extractor == nil {
		return nil, errors.New("extractor can't be nil")
	}

	return &Exporter{
		Config: config,
	}, nil
}

func (xp *Exporter) Stop(done chan os.Signal) {
	log.Println("Stopping exporter")
	xp.interruptScan = true
	err := xp.CloseConnector()
	if err != nil {
		log.Println("error closing connector:", err)
	}
	log.Println("Exporter stopped")
	close(done)
}

func (xp *Exporter) CloseConnector() error {
	return xp.Connector().Close()
}

func (xp *Exporter) Connector() Connector {
	return xp.Config.Connector
}

func (xp *Exporter) Scan() {
	ticker := time.NewTicker(time.Duration(xp.Config.ScanInterval) * time.Second)
	go func() {
		defer ticker.Stop()
		log.Println("Starting scanner")

		for _ = range ticker.C {
			if xp.interruptScan == true {
				log.Println("Stopping scanner")
				ticker.Stop()
				break
			}

			queues, err := xp.SelectQueuesToScan()
			xp.SetQueuesType(queues)

			if err != nil {
				log.Fatal(err)
			}

			for _, queue := range queues {
				err = xp.CountJobsForQueue(queue)
				if err != nil {
					log.Println(fmt.Sprintf("error getting metrics for %s: %v", queue.queueItem.Name, err))
					continue
				}

				//TODO Implement RedisQueueMetricsFormatter to output metrics
				log.Println(strings.Replace(queue.Name(), fmt.Sprintf("%s:", QUEUE_ROOT_NODE), "", 1), queue.GetCurrentJobsCount())
			}
		}
	}()
}

func (xp *Exporter) CountJobsForQueue(queue *RedisQueue) error {
	return xp.Extractor().CountJobsForQueue(queue)
}

func (xp *Exporter) SelectQueuesToScan() ([]*RedisQueue, error) {
	var err error
	queueItems := []*RedisQueue{}

	if len(xp.Config.QueueNames) > 0 {
		queueItems = parsedQueuesNames(xp.Config.QueueNames)
	} else {
		queueItems, err = xp.Extractor().ListAllQueuesFromDB()
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
	xp.Extractor().SetQueueTypeForQueues(queues)
}

func (xp *Exporter) Extractor() Extractor {
	return xp.Config.Extractor
}
