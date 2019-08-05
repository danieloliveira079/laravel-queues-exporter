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
	queueType string
}

func (q *RedisQueue) Name() string {
	return q.queueItem.Name
}

func (q *RedisQueue) GetQueueType() string {
	return q.queueType
}

func (q *RedisQueue) SetQueueType(queueType string) {
	q.queueType = queueType
}

func (q *RedisQueue) GetCurrentJobsCount() int64 {
	return q.queueItem.Jobs
}

func (q *RedisQueue) SetCurrentJobsCount(count int64) {
	q.queueItem.Jobs = count
}

func (q *RedisQueue) ToLaravel() string {
	var laravelName string

	if q.queueItem == nil || q.queueItem.HasQueueName() == false {
		return laravelName
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