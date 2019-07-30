package redis

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strings"
)

type RedisExtractor struct {
	Config ExtractorConfig
}

type ExtractorConfig struct {
	ConnConfig ConnectionConfig
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

func (re *RedisExtractor) ListAllQueues() ([]QueueItem, error) {
	var err error
	queueItems := []QueueItem{}

	list, err := re.Dispatcher().Do("KEYS", fmt.Sprintf("%s:*", QUEUE_ROOT_NODE))

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

func (re *RedisExtractor) Dispatcher() CommandDispatcher {
	return re.Config.Dispatcher
}

func checkQueueName(name string) string {
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

func (re *RedisExtractor) CountJobsForQueue(queue *QueueItem) error {
	name := checkQueueName(queue.Name)
	queueType, err := redis.String(re.Dispatcher().Do("type", name))

	if err != nil {
		return err
	}

	var jobsCount int64
	var redisCmd string

	switch queueType {
	case "zset":
		redisCmd = "zcard"
	case "list":
		redisCmd = "llen"
	case "none":
		redisCmd = "llen"
	}

	jobsCount, err = redis.Int64(re.Dispatcher().Do(redisCmd, name))
	if err != nil {
		return err
	}

	queue.Jobs = jobsCount
	return err
}
