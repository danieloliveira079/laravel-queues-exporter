package consumer

import (
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/config"
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/consumer/log"
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/consumer/statsd"
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/consumer/stdout"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func Test_Consumer_ShouldReturnStdoutConsumerGivenConfiguration(t *testing.T) {
	appConfig := &config.AppConfig{
		RedisHost:  "0.0.0.0",
		RedisPort:  "6379",
		RedisDB:    0,
		StatsDHost: "0.0.0.0",
		StatsDPort: "8125",
	}

	testCase := []struct {
		desc         string
		exportTo     string
		consumerType interface{}
		expected     bool
	}{
		{
			desc:         "Return stdout consumer",
			exportTo:     "stdout",
			consumerType: reflect.TypeOf(&stdout.Stdout{}),
			expected:     true,
		},
		{
			desc:         "Return statsd consumer",
			exportTo:     "statsd",
			consumerType: reflect.TypeOf(&statsd.StatsD{}),
			expected:     true,
		},
		{
			desc:         "Return log consumer",
			exportTo:     "log",
			consumerType: reflect.TypeOf(&log.Log{}),
			expected:     true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.desc, func(t *testing.T) {
			appConfig.ExportTo = tc.exportTo
			consumer, err := BuildConsumersListFromConfig(appConfig)
			require.Nil(t, err)
			assert.Equal(t, tc.consumerType, reflect.TypeOf(consumer[0]))
		})
	}
}
