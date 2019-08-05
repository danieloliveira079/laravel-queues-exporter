package queue

type QueueItem struct {
	Name string
	Jobs int64
}

func (q *QueueItem) HasQueueName() bool {
	return len(q.Name) > 0
}
