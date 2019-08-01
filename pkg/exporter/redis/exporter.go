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
	ListAllQueuesFromDB() ([]QueueItem, error)
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

			if err != nil {
				log.Fatal(err)
			}

			for _, queue := range queues {
				err = xp.CountJobsForQueue(&queue)
				if err != nil {
					log.Println(fmt.Sprintf("error getting metrics for %s: %v", queue.Name, err))
					continue
				}

				log.Println(strings.Replace(queue.Name, fmt.Sprintf("%s:", QUEUE_ROOT_NODE), "", 1), queue.Jobs)
			}
		}
	}()
}

func (xp *Exporter) CountJobsForQueue(queue *QueueItem) error {
	return xp.Extractor().CountJobsForQueue(queue)
}

func (xp *Exporter) SelectQueuesToScan() ([]QueueItem, error) {

	var err error
	queueItems := []QueueItem{}

	if len(xp.Config.QueueNames) > 0 {
		queueItems = parseQueueNames(xp.Config.QueueNames)
	} else {
		queueItems, err = xp.Extractor().ListAllQueuesFromDB()
	}

	return queueItems, err
}

func (xp *Exporter) Extractor() Extractor {
	return xp.Config.Extractor
}

func parseQueueNames(queueNames string) []QueueItem {
	queueItems := []QueueItem{}
	names := strings.Split(queueNames, ",")
	for _, n := range names {
		queueItems = append(queueItems, QueueItem{Name: n})
	}

	return queueItems
}

func (q *QueueItem) LaravelQueueName() string {
	if len(q.Name) == 0 {
		return q.Name
	}

	name := q.Name
	parts := strings.Split(name, ":")

	if len(parts) == 1 {
		name = fmt.Sprintf("%s:%s", QUEUE_ROOT_NODE, name)
	} else if len(parts) > 1 {
		if parts[0] != QUEUE_ROOT_NODE {
			name = fmt.Sprintf("%s:%s", QUEUE_ROOT_NODE, name)
		}
	}

	return name
}
