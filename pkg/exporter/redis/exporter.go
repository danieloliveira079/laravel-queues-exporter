package redis

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	QUEUE_ROOT_NODE = "queues"
)

type Exporter struct {
	Config        RedisExporterConfig
	interruptScan bool
	queueItems    []QueueItem
}

type RedisExporterConfig struct {
	ConnectionConfig ConnectionConfig
	ScanInterval     int
	QueueNames       string
	Extractor        Extractor
	Connector        Connector
}

type QueueItem struct {
	Name string
	Jobs int64
}

type Extractor interface {
	ListAllQueues() ([]QueueItem, error)
	CountJobsForQueue(queue *QueueItem) error
}

type Connector interface {
	Connect() (err error)
	Close() (err error)
	Do(command string, args ...interface{}) (results interface{}, err error)
}

func NewRedisExporter(config RedisExporterConfig) (*Exporter, error) {
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

func (r *Exporter) Stop(done chan os.Signal) {
	log.Println("Stopping exporter")
	r.interrupt = true
	err := r.CloseConnector()
	if err != nil {
		log.Println("error closing connector:", err)
	}
	log.Println("Exporter stopped")
	close(done)
}

func (r *Exporter) CloseConnector() error {
	return r.Connector().Close()
}

func (r *Exporter) Connector() Connector {
	return r.Config.Connector
}

func (r *Exporter) Scan() {
	ticker := time.NewTicker(time.Duration(r.Config.CheckInterval) * time.Second)
	go func() {
		defer ticker.Stop()
		log.Println("Starting scanner")

		for _ = range ticker.C {
			if r.interrupt == true {
				log.Println("Stopping scanner")
				ticker.Stop()
				break
			}

			queues, err := r.SelectQueuesToScan()

			if err != nil {
				log.Fatal(err)
			}

			for _, queue := range queues {
				err = r.CountJobsForQueue(&queue)
				if err != nil {
					log.Println(fmt.Sprintf("error getting metrics for %s: %v", queue.Name, err))
					continue
				}

				log.Println(strings.Replace(queue.Name, fmt.Sprintf("%s:", QUEUE_ROOT_NODE), "", 1), queue.Jobs)
			}
		}
	}()
}

func (r *Exporter) CountJobsForQueue(queue *QueueItem) error {
	return r.Extractor().CountJobsForQueue(queue)
}

func (r *Exporter) SelectQueuesToScan() ([]QueueItem, error) {

	var err error
	queueItems := []QueueItem{}

	if len(r.Config.QueueNames) > 0 {
		queueItems = parseQueueNames(r.Config.QueueNames)
	} else {
		queueItems, err = r.Extractor().ListAllQueues()
	}

	return queueItems, err
}

func (r *Exporter) Extractor() Extractor {
	return r.Config.Extractor
}

func parseQueueNames(queueNames string) []QueueItem {
	queueItems := []QueueItem{}
	names := strings.Split(queueNames, ",")
	for _, n := range names {
		queueItems = append(queueItems, QueueItem{Name: n})
	}

	return queueItems
}
