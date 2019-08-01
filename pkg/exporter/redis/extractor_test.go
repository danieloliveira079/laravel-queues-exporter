package redis

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

type FakeDispatcher struct {
	mock.Mock
}

func (d *FakeDispatcher) Do(command string, args ...interface{}) (reply interface{}, err error) {
	argsDo := d.Called(command, args)
	return argsDo.Get(0), nil
}

func Test_Extractor_ShouldNotCreateNewRedisExtractorWithNilDispatcher(t *testing.T) {
	config := ExtractorConfig{
		Dispatcher: nil,
	}

	extractor, err := NewRedisExtractor(config)
	require.Nil(t, extractor)
	require.Error(t, err)
}

func Test_Extractor_ShouldCreateNewRedisExtractorWithDispatcher(t *testing.T) {
	dispatcher := &FakeDispatcher{}

	config := ExtractorConfig{
		Dispatcher: dispatcher,
	}

	extractor, err := NewRedisExtractor(config)
	require.NotNil(t, extractor)
	require.Nil(t, err)
}

func Test_RedisExtractor_ShouldListAllQueuesFromDB(t *testing.T) {
	dispatcher := new(FakeDispatcher)

	cmd := "keys"
	args := []interface{}{
		fmt.Sprintf("%s:*", QUEUE_ROOT_NODE),
	}

	queuesFromDB := []interface{}{
		"queue1",
		"queue2",
		"queue3",
	}

	queuesMatch := func(onDB []interface{}, fromDB []QueueItem) bool {
		allMatch := true
		for _, queue := range onDB {
			found := false
			for _, item := range fromDB {
				if queue == item.Name {
					found = true
					break
				}
			}

			if found == false {
				allMatch = false
				break
			}
		}

		return allMatch
	}

	dispatcher.On("Do", cmd, args).Return(queuesFromDB)

	config := ExtractorConfig{
		Dispatcher: dispatcher,
	}

	extractor, err := NewRedisExtractor(config)
	require.Nil(t, err)

	queueItems, err := extractor.ListAllQueuesFromDB()
	require.Nil(t, err)
	assert.Equal(t, queuesMatch(queuesFromDB, queueItems), true)
}

func Test_RedisExtractor_ShouldReturnCommandByQueueType(t *testing.T) {
	dispatcher := new(FakeDispatcher)

	config := ExtractorConfig{
		Dispatcher: dispatcher,
	}

	extractor, err := NewRedisExtractor(config)
	require.Nil(t, err)

	testCases := []struct {
		desc      string
		queueType string
		expected  string
	}{
		{
			desc:      "Return zcard command for zset queue type",
			queueType: "zset",
			expected:  "zcard",
		},
		{
			desc:      "Return llen command for list queue type",
			queueType: "list",
			expected:  "llen",
		},
		{
			desc:      "Return llen command for none queue type",
			queueType: "none",
			expected:  "llen",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			cmd := extractor.CountJobCmdNameByQueueType(tc.queueType)
			assert.Equal(t, tc.expected, cmd)
		})
	}

}
