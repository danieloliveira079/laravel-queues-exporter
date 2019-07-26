package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"strings"
	"time"
)

type RedisExtractor struct {
	targetHost string
	targetPort string
	targetDB   int
	conn       redis.Conn
}

func NewRedisExtractor(targetHost string, targetPort string, targetDB int) *RedisExtractor {
	return &RedisExtractor{targetHost: targetHost, targetPort: targetPort, targetDB: targetDB}
}

func (r *RedisExtractor) Close() (err error) {
	log.Println("Closing extractor connection")
	if r.conn != nil {
		err = r.conn.Close()
	}

	return err
}

func (r *RedisExtractor) Connect() (err error) {
	if r.conn != nil && r.conn.Err() == nil {
		return nil
	}

	conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", r.targetHost, r.targetPort), redis.DialDatabase(r.targetDB), redis.DialConnectTimeout(15*time.Second))
	if err != nil {
		return err
	}

	r.conn = conn
	return nil
}

func (re *RedisExtractor) ListQueues() (queueItems []QueueItem, err error) {
	err = re.Connect()
	if err != nil {
		return nil, err
	}

	list, err := re.conn.Do("KEYS", fmt.Sprintf("%s:*", QUEUE_ROOT_NODE))

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
	err = re.Connect()
	if err != nil {
		return err
	}

	name := checkQueueName(queue.Name)
	queueType, err := redis.String(re.conn.Do("type", name))

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

	jobsCount, err = redis.Int64(re.conn.Do(redisCmd, name))
	if err != nil {
		return err
	}

	queue.Jobs = jobsCount
	return err
}
