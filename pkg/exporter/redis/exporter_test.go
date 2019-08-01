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

func (f *FakeRedisExtractor) ListAllQueues() (queueItems []QueueItem, err error) {
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

type FakeRedisConnector struct {
}

func (c *FakeRedisConnector) Connect() (err error) {
	panic("implement me")
}

func (c *FakeRedisConnector) Close() (err error) {
	panic("implement me")
}

func (c *FakeRedisConnector) Do(command string, args ...interface{}) (results interface{}, err error) {
	panic("implement me")
}

func TestShouldSelectAllQueuesToScan(t *testing.T) {
	extractor := &FakeRedisExtractor{
		Queues: []string{"queue1", "queue2", "queue3"},
	}

	connector := &FakeRedisConnector{}

	config := RedisExporterConfig{
		ConnectionConfig: ConnectionConfig{
			Host: "none",
			Port: "none",
			DB:   0,
		},
		ScanInterval: 5,
		Extractor:    extractor,
		Connector:    connector,
	}

	exporter, _ := NewRedisExporter(config)

	selected, _ := exporter.SelectQueuesToScan()

	//TODO Implement test cases
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

	connector := &FakeRedisConnector{}

	filtered := []string{"queue4", "queue5", "queue6"}

	config := RedisExporterConfig{
		ConnectionConfig: ConnectionConfig{
			Host: "none",
			Port: "none",
			DB:   0,
		},
		ScanInterval: 5,
		QueueNames:   strings.Join(filtered, ","),
		Extractor:    extractor,
		Connector:    connector,
	}

	exporter, _ := NewRedisExporter(config)

	selected, _ := exporter.SelectQueuesToScan()

	//TODO Implement test cases
	for i := range extractor.Queues {
		if selected[i].Name != filtered[i] {
			t.Errorf("expected %s, actual: %s", selected[i].Name, filtered[i])
		}
	}
}
