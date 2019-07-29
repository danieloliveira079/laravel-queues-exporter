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
	connector     Connector
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
	ListQueues() ([]QueueItem, error)
	CountJobsForQueue(queue *QueueItem) error
}

func NewRedisExporter(targetHost string,
	targetPort string,
	targetDB int,
	checkInterval int,
	queueNames string,
	extractor Extractor,
	connector Connector) *RedisExporter {

	if connector == nil {
		connector = NewRedisConnector(targetHost, targetPort, targetDB)
	}

	if extractor == nil {
		extractor = NewRedisExtractor(targetHost, targetPort, targetDB, connector)
	}

	return &RedisExporter{targetHost: targetHost,
		targetPort:    targetPort,
		targetDB:      targetDB,
		checkInterval: checkInterval,
		queuesNames:   queueNames,
		extractor:     extractor,
		connector:     connector,
	}
}

func (r *RedisExporter) Stop(done chan os.Signal) {
	log.Println("Stopping exporter")
	r.interrupt = true
	err := r.connector.Close()
	if err != nil {
		log.Println("error closing connector:", err)
	}
	log.Println("Exporter stopped")
	close(done)
}

func (r *RedisExporter) Scan() {
	ticker := time.NewTicker(time.Duration(r.checkInterval) * time.Second)
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
				err = r.extractor.CountJobsForQueue(&queue)
				if err != nil {
					log.Println(fmt.Sprintf("error getting metrics for %s: %v", queue.Name, err))
					continue
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
