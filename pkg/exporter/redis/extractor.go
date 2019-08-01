package redis

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
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

func (xt *RedisExtractor) ListAllQueuesFromDB() ([]QueueItem, error) {
	var err error
	queueItems := []QueueItem{}

	//TODO Warning when key is not available
	list, err := xt.Dispatcher().Do("keys", fmt.Sprintf("%s:*", QUEUE_ROOT_NODE))

	if err != nil {
		return nil, err
	}

	parsedList, err := redis.Strings(list, nil)

	if err != nil {
		return nil, err
	}

	for _, q := range parsedList {
		queueItems = append(queueItems, QueueItem{
			Name: q,
		})
	}
	return queueItems, err
}

func (xt *RedisExtractor) Dispatcher() CommandDispatcher {
	return xt.Config.Dispatcher
}

func (xt *RedisExtractor) CountJobsForQueue(queue *QueueItem) error {
	queueName := queue.LaravelQueueName()

	//TODO Extract queue type check to its own method
	queueType, err := redis.String(xt.Dispatcher().Do("type", queueName))

	if err != nil {
		return err
	}

	var jobsCount int64
	redisCmd := xt.CountJobCmdNameByQueueType(queueType)

	jobsCount, err = redis.Int64(xt.Dispatcher().Do(redisCmd, queueName))
	if err != nil {
		return err
	}

	queue.Jobs = jobsCount
	return err
}

func (xt *RedisExtractor) CountJobCmdNameByQueueType(queueType string) string {
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
