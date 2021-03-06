package redis

import (
	"errors"
	"fmt"
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/queue"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

type FakeRedisExtractor struct {
	QueuesOnDB []string
}

func (f *FakeRedisExtractor) Connect() error {
	return nil
}

func (f *FakeRedisExtractor) Close() error {
	return nil
}

func (f *FakeRedisExtractor) ListAllQueuesFromDB() (queueItems []*RedisQueue, err error) {
	for _, q := range f.QueuesOnDB {
		queueItems = append(queueItems, &RedisQueue{
			queueItem: &queue.QueueItem{
				Name: q,
			},
		})
	}

	return queueItems, err
}

func (f *FakeRedisExtractor) CountJobsForQueues(queue []*RedisQueue) error {
	panic("implement me")
}

func (f *FakeRedisExtractor) SetQueueTypeForQueues(queues []*RedisQueue) {
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

func Test_Exporter_ShouldReturnAllQueuesFromDB(t *testing.T) {
	extractor := &FakeRedisExtractor{
		QueuesOnDB: []string{"queue1", "queue2", "queue3"},
	}
	connector := &FakeRedisConnector{}
	config := &ExporterConfig{}

	exporter, _ := NewRedisExporter(config, connector, extractor)

	actual := func(queueItems []*RedisQueue) string {
		var names []string
		for _, q := range queueItems {
			names = append(names, q.Name())
		}
		return strings.Join(names, ",")
	}

	dbMatchesSelected := func(fromDB []string, selected []*RedisQueue) (bool, error) {
		allMatch := true
		var err error
		for _, queue := range fromDB {
			found := false
			for _, item := range selected {
				if queue == item.Name() {
					found = true
					break
				}
			}

			if found == false {
				allMatch = false
				err = errors.New(fmt.Sprintf("expected: %s \nactual: %s", strings.Join(fromDB, ","), actual(selected)))
				break
			}
		}

		return allMatch, err
	}

	testCases := []struct {
		desc                string
		remoteQueues        []string
		validateQueuesMatch func(fromDB []string, selected []*RedisQueue) (bool, error)
		expected            bool
	}{
		{
			desc:                "Return all queues from DB",
			remoteQueues:        []string{"queue1", "queue2", "queue3"},
			validateQueuesMatch: dbMatchesSelected,
			expected:            true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			selected, err := exporter.SelectQueuesToScan()
			require.Nil(t, err)
			match, err := tc.validateQueuesMatch(tc.remoteQueues, selected)
			require.Equal(t, tc.expected, match, err)
		})
	}
}

func Test_Exporter_ShouldSelectFilteredQueuesFromDB(t *testing.T) {
	extractor := &FakeRedisExtractor{}
	connector := &FakeRedisConnector{}
	config := &ExporterConfig{}

	exporter, _ := NewRedisExporter(config, connector, extractor)

	actual := func(queueItems []*RedisQueue) string {
		var names []string
		for _, q := range queueItems {
			names = append(names, q.Name())
		}
		return strings.Join(names, ",")
	}

	configMatchesSelected := func(fromConfig string, selected []*RedisQueue) (bool, error) {
		allMatch := true
		var err error
		for _, queue := range strings.Split(fromConfig, ",") {
			found := false
			for _, item := range selected {
				if queue == item.Name() {
					found = true
					break
				}
			}

			if found == false {
				allMatch = false
				err = errors.New(fmt.Sprintf("expected: %s \nactual: %s", fromConfig, actual(selected)))
				break
			}
		}

		return allMatch, err
	}

	testCases := []struct {
		desc                string
		remoteQueues        []string
		queuesFromConfig    string
		validateQueuesMatch func(fromConfig string, actual []*RedisQueue) (bool, error)
		expected            bool
	}{
		{
			desc:                "Remote and config queues match",
			remoteQueues:        []string{"queue1", "queue2", "queue3"},
			queuesFromConfig:    "queue1,queue2,queue3",
			validateQueuesMatch: configMatchesSelected,
			expected:            true,
		},
		{
			desc:                "Remote holds extra queues than config",
			remoteQueues:        []string{"queue1", "queue2", "queue3", "queue4"},
			queuesFromConfig:    "queue1,queue2,queue3",
			validateQueuesMatch: configMatchesSelected,
			expected:            true,
		},
		{
			desc:                "RedisQueue from config is not available on remote",
			remoteQueues:        []string{"queue1", "queue2", "queue3", "queue5"},
			queuesFromConfig:    "queue1,queue2,queue3,queue4",
			validateQueuesMatch: configMatchesSelected,
			expected:            true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			extractor.QueuesOnDB = tc.remoteQueues
			config.QueueNames = tc.queuesFromConfig

			selected, err := exporter.SelectQueuesToScan()
			require.Nil(t, err)
			match, err := tc.validateQueuesMatch(tc.queuesFromConfig, selected)
			require.Equal(t, tc.expected, match, err)
		})
	}
}

func Test_Exporter_ShouldReturnLaravelQueueNameForGivenQueueName(t *testing.T) {
	nameWithRootNode := func(name string) string {
		return fmt.Sprintf("%s:%s", LARAVEL_QUEUE_ROOT_NODE, name)
	}

	testCases := []struct {
		desc      string
		queueItem RedisQueue
		expected  string
	}{
		{
			desc: "RedisQueue name already contains queue root node",
			queueItem: RedisQueue{
				queueItem: &queue.QueueItem{
					Name: nameWithRootNode("queueTest"),
				},
			},
			expected: nameWithRootNode("queueTest"),
		},
		{
			desc: "RedisQueue name does not contain queue root node",
			queueItem: RedisQueue{
				queueItem: &queue.QueueItem{
					Name: "queueTest",
				},
			},
			expected: nameWithRootNode("queueTest"),
		},
		{
			desc: "Reserved queue's name does not contain queue root node",
			queueItem: RedisQueue{
				queueItem: &queue.QueueItem{
					Name: "queueTest:reserved",
				},
			},
			expected: nameWithRootNode("queueTest:reserved"),
		},
		{
			desc: "Delayed queue's name does not contain queue root node",
			queueItem: RedisQueue{
				queueItem: &queue.QueueItem{
					Name: "queueTest:delayed",
				},
			},
			expected: nameWithRootNode("queueTest:delayed"),
		},
		{
			desc: "RedisQueue name is empty",
			queueItem: RedisQueue{
				queueItem: &queue.QueueItem{
					Name: "",
				},
			},
			expected: "",
		},
		{
			desc: "RedisQueue name contains root but name is empty",
			queueItem: RedisQueue{
				queueItem: &queue.QueueItem{
					Name: nameWithRootNode(":"),
				},
			},
			expected: LARAVEL_QUEUE_ROOT_NODE,
		},
		{
			desc: "RedisQueue name has not parent node",
			queueItem: RedisQueue{
				queueItem: &queue.QueueItem{
					Name: ":queueTest:",
				},
			},
			expected: nameWithRootNode("queueTest"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			name := tc.queueItem.FullName()
			require.Equal(t, tc.expected, name)
		})
	}
}
