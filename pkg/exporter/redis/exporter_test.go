package redis

import (
	"strings"
	"testing"
)

type FakeRedisExtractor struct {
	Queues []string
}

func (f *FakeRedisExtractor) Connect() error {
	return nil
}

func (f *FakeRedisExtractor) Close() error {
	return nil
}

func (f *FakeRedisExtractor) ListQueues() (queueItems []QueueItem, err error) {
	for _, q := range f.Queues {
		queueItems = append(queueItems, QueueItem{
			Name: q,
		})
	}

	return queueItems, err
}

func (f *FakeRedisExtractor) CountJobsForQueue(queue *QueueItem) error {
	panic("implement me")
}

func TestShouldSelectAllQueuesToScan(t *testing.T) {
	extractor := &FakeRedisExtractor{
		Queues: []string{"queue1", "queue2", "queue3"},
	}

	exporter := NewRedisExporter("none", "none", 0, 5, "", extractor)

	selected, _ := exporter.SelectQueuesToScan()

	for i := range extractor.Queues {
		if selected[i].Name != extractor.Queues[i] {
			t.Errorf("expected %s, actual: %s", selected[i].Name, extractor.Queues[i])
		}
	}
}

func TestShouldSelectFilteredQueuesToScan(t *testing.T) {
	extractor := &FakeRedisExtractor{
		Queues: []string{"queue1", "queue2", "queue3"},
	}

	filtered := []string{"queue4", "queue5", "queue6"}

	exporter := NewRedisExporter("none", "none", 0, 5, strings.Join(filtered, ","), extractor)

	selected, _ := exporter.SelectQueuesToScan()

	for i := range extractor.Queues {
		if selected[i].Name != filtered[i] {
			t.Errorf("expected %s, actual: %s", selected[i].Name, filtered[i])
		}
	}
}
