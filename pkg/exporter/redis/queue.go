package redis

import (
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/queue"
	"strings"
)

const (
	QUEUE_ROOT_NODE = "queues"
)

type RedisQueue struct {
	queueItem *queue.QueueItem
}

func (q *RedisQueue) Name() string {
	return q.queueItem.Name
}

func (q *RedisQueue) GetQueueType() string {
	return q.queueItem.Type
}

func (q *RedisQueue) SetQueueType(queueType string) {
	q.queueItem.Type = queueType
}

func (q *RedisQueue) GetCurrentJobsCount() int64 {
	return q.queueItem.Jobs
}

func (q *RedisQueue) SetCurrentJobsCount(count int64) {
	q.queueItem.Jobs = count
}

func (q *RedisQueue) LaravelQueueName() string {
	var laravelName string

	if q.queueItem == nil {
		return laravelName
	}

	if q.queueItem.HasQueueName() == false {
		return q.queueItem.Name
	}

	parts, partsCount := q.laravelQueueNameSplit()

	switch {
	case partsCount == 0:
		return laravelName
	case partsCount >= 1:
		tmpParts := []string{
			QUEUE_ROOT_NODE,
		}

		for _, p := range parts {
			if len(p) > 0 && p != QUEUE_ROOT_NODE {
				tmpParts = append(tmpParts, p)
			}
		}

		laravelName = strings.Join(tmpParts, ":")
	}

	return laravelName
}

func (q *RedisQueue) laravelQueueNameSplit() ([]string, int) {
	parts := strings.Split(q.queueItem.Name, ":")
	return parts, len(parts)
}
