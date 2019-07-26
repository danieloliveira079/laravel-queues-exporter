package redis

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	QUEUE_ROOT_NODE = "queues"
)

type RedisExporter struct {
	targetHost    string
	targetPort    string
	targetDB      int
	extractor     Extractor
	interrupt     bool
	checkInterval int
	queuesNames   string
	queueItems    []QueueItem
}

type QueueItem struct {
	Name string
	Jobs int64
}

type Extractor interface {
	Connect() error
	Close() error
	ListQueues() ([]QueueItem, error)
	CountJobsForQueue(queue *QueueItem) error
}

func NewRedisExporter(targetHost string,
	targetPort string,
	targetDB int,
	checkInterval int,
	queueNames string,
	extractor Extractor) *RedisExporter {

	if extractor == nil {
		extractor = NewRedisExtractor(targetHost, targetPort, targetDB)
	}

	return &RedisExporter{targetHost: targetHost,
		targetPort:    targetPort,
		targetDB:      targetDB,
		checkInterval: checkInterval,
		queuesNames:   queueNames,
		extractor:     extractor,
	}
}

func (r *RedisExporter) Stop(done chan os.Signal) {
	log.Println("Stopping exporter")
	r.interrupt = true
	_ = r.extractor.Close()
	log.Println("Exporter stopped")
	close(done)
}

func (r *RedisExporter) Scan() {
	err := r.extractor.Connect()
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(time.Duration(r.checkInterval) * time.Second)
	go func() {
		defer ticker.Stop()
		log.Println("Starting scanner")

		for _ = range ticker.C {
			if r.interrupt == true {
				ticker.Stop()
				break
			}

			queues, err := r.SelectQueuesToScan()

			if err != nil {
				log.Fatal(err)
			}

			for _, queue := range queues {
				err = r.extractor.CountJobsForQueue(&queue)
				if err != nil {
					log.Println(fmt.Sprintf("error getting metrics for %s: %v", queue.Name, err))
				}

				log.Println(strings.Replace(queue.Name, fmt.Sprintf("%s:", QUEUE_ROOT_NODE), "", 1), queue.Jobs)
			}
		}
	}()
}

func (r *RedisExporter) SelectQueuesToScan() (queueItems []QueueItem, err error) {

	if len(r.queuesNames) > 0 {
		queueItems = parseQueueNames(r.queuesNames)
	} else {
		queueItems, err = r.extractor.ListQueues()
	}

	return queueItems, err
}

func parseQueueNames(queueNames string) (queueItems []QueueItem) {

	names := strings.Split(queueNames, ",")
	for _, n := range names {
		queueItems = append(queueItems, QueueItem{Name: n})
	}

	return queueItems
}
