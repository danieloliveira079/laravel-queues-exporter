package redis

import (
	"errors"
	"fmt"
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/queue"
	"github.com/gomodule/redigo/redis"
	"log"
)

type RedisExtractor struct {
	Config ExtractorConfig
}

type ExtractorConfig struct {
	Dispatcher CommandDispatcher
}

type CommandDispatcher interface {
	Do(command string, args ...interface{}) (reply interface{}, err error)
}

func NewRedisExtractor(config ExtractorConfig) (*RedisExtractor, error) {
	if config.Dispatcher == nil {
		return nil, errors.New("dispatcher can't be nil")
	}
	return &RedisExtractor{Config: config}, nil
}

func (xt *RedisExtractor) ListAllQueuesFromDB() ([]*RedisQueue, error) {
	var err error
	queueItems := []*RedisQueue{}

	list, err := xt.Dispatcher().Do("keys", fmt.Sprintf("%s:*", QUEUE_ROOT_NODE))

	if err != nil {
		return nil, err
	}

	parsedList, err := redis.Strings(list, nil)

	if err != nil {
		return nil, err
	}

	for _, q := range parsedList {
		queueItems = append(queueItems, &RedisQueue{
			queueItem: &queue.QueueItem{
				Name: q,
			},
		})
	}
	return queueItems, err
}

func (xt *RedisExtractor) Dispatcher() CommandDispatcher {
	return xt.Config.Dispatcher
}

func (xt *RedisExtractor) CountJobsForQueues(queues []*RedisQueue) error {
	var err error

	for _, q := range queues {
		queueName := q.ToLaravel()

		var jobsCount int64
		cmdName := xt.CountJobsCmdNameByQueueType(q.GetQueueType())

		//TODO Implement a parser instead of using package directly
		jobsCount, err := redis.Int64(xt.Dispatcher().Do(cmdName, queueName))
		if err != nil {
			return err
		}

		q.SetCurrentJobsCount(jobsCount)
	}

	return err
}

func (xt *RedisExtractor) SetQueueTypeForQueues(queues []*RedisQueue) {
	for i, q := range queues {
		queueType, err := redis.String(xt.Dispatcher().Do("type", q.Name()))

		if err != nil {
			log.Printf("error: type could not be defined for queue %s", q.Name())
		}

		queues[i].SetQueueType(queueType)
	}
}

//TODO ImplementCountJobs extrctor per queue type
func (xt *RedisExtractor) CountJobsCmdNameByQueueType(queueType string) string {
	var cmd string

	switch queueType {
	case "zset":
		cmd = "zcard"
	case "list":
		cmd = "llen"
	case "none":
		cmd = "llen"
	}

	return cmd
}
