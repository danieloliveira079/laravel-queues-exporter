package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strings"
)

type RedisExtractor struct {
	targetHost string
	targetPort string
	targetDB   int
	dispatcher CommandDispatcher
}

type CommandDispatcher interface {
	Do(command string, args ...interface{}) (reply interface{}, err error)
}

func NewRedisExtractor(targetHost string, targetPort string, targetDB int, dispatcher CommandDispatcher) *RedisExtractor {
	return &RedisExtractor{targetHost: targetHost, targetPort: targetPort, targetDB: targetDB, dispatcher: dispatcher}
}

func (re *RedisExtractor) ListQueues() (queueItems []QueueItem, err error) {
	list, err := re.dispatcher.Do("KEYS", fmt.Sprintf("%s:*", QUEUE_ROOT_NODE))

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

func (re *RedisExtractor) CountJobsForQueue(queue *QueueItem) (err error) {
	name := checkQueueName(queue.Name)
	queueType, err := redis.String(re.dispatcher.Do("type", name))

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

	jobsCount, err = redis.Int64(re.dispatcher.Do(redisCmd, name))
	if err != nil {
		return err
	}

	queue.Jobs = jobsCount
	return err
}
