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

func (xt *RedisExtractor) CountJobsForQueue(queue *RedisQueue) error {
	queueName := queue.LaravelQueueName()

	var jobsCount int64
	cmdName := xt.CountJobsCmdNameByQueueType(queue.GetQueueType())

	//TODO Implement a parser instead of using package directly
	jobsCount, err := redis.Int64(xt.Dispatcher().Do(cmdName, queueName))
	if err != nil {
		return err
	}

	queue.SetCurrentJobsCount(jobsCount)
	return err
}

func (xt *RedisExtractor) SetQueueTypeForQueues(queues []*RedisQueue) {
	for i, queue := range queues {
		queueType, err := redis.String(xt.Dispatcher().Do("type", queue.Name()))

		if err != nil {
			log.Printf("error: type could not defined for queue %s", queue.Name())
		}

		queues[i].SetQueueType(queueType)
	}
}

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
