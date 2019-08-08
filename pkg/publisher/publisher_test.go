package publisher

import (
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/consumer"
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/metric"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type fakeConsumer struct {
	mock.Mock
}

func (c *fakeConsumer) Process(metrics []metric.Metric) {
	c.Called(metrics)
}

func Test_Publisher_Should_Not_Notify_Consumers_When_Metrics_Are_NilOrEmpty(t *testing.T) {
	testCases := []struct {
		desc            string
		metrics         []metric.Metric
		metricsConsumer *fakeConsumer
		expected        bool
	}{
		{
			desc:            "Metrics are nil",
			metrics:         nil,
			metricsConsumer: new(fakeConsumer),
			expected:        true,
		},
		{
			desc:            "Metrics are empty",
			metrics:         []metric.Metric{},
			metricsConsumer: new(fakeConsumer),
			expected:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			publisher := new(MetricsPublisher)

			publisher.SubscribeConsumers(tc.metricsConsumer)
			tc.metricsConsumer.On("Process", tc.metrics)

			publisher.Publish(tc.metrics)

			notCalled := tc.metricsConsumer.AssertNotCalled(t, "Process")
			assert.Equal(t, tc.expected, notCalled)
		})
	}
}

func Test_Publisher_Should_Notify_Consumers_When_Metrics_Are_Present(t *testing.T) {
	testCases := []struct {
		desc     string
		metrics  []metric.Metric
		expected bool
	}{
		{
			desc: "Single metric",
			metrics: []metric.Metric{
				{
					Name:  "queue1",
					Value: 1,
				},
			},
			expected: true,
		},
		{
			desc: "Multiple metrics",
			metrics: []metric.Metric{
				{
					Name:  "queue1",
					Value: 1,
				},
				{
					Name:  "queue2",
					Value: 2,
				},
			},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			publisher := new(MetricsPublisher)
			metricsConsumer := new(fakeConsumer)

			metricsConsumer.On("Process", tc.metrics)

			publisher.SubscribeConsumers(metricsConsumer)
			publisher.Publish(tc.metrics)

			notCalled := metricsConsumer.AssertCalled(t, "Process", tc.metrics)
			assert.Equal(t, tc.expected, notCalled)
		})
	}
}

func Test_Publisher_Should_Subscribe_Consumers(t *testing.T) {
	testCases := []struct {
		desc             string
		metricsConsumers []consumer.Consumer
		expected         int
	}{
		{
			desc: "Single metrics consumer",
			metricsConsumers: []consumer.Consumer{
				new(fakeConsumer),
			},
			expected: 1,
		},
		{
			desc: "Multiple metrics consumers",
			metricsConsumers: []consumer.Consumer{
				new(fakeConsumer),
				new(fakeConsumer),
			},
			expected: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			publisher := new(MetricsPublisher)
			publisher.SubscribeConsumers(tc.metricsConsumers...)

			assert.Equal(t, tc.expected, len(publisher.consumers))
		})
	}
}

func Test_Publisher_Should_Not_Subscribe_Consumers(t *testing.T) {
	testCases := []struct {
		desc             string
		metricsConsumers []consumer.Consumer
		expected         int
	}{
		{
			desc:             "Empty metrics consumers",
			metricsConsumers: []consumer.Consumer{},
			expected:         0,
		},
		{
			desc:             "Nil metrics consumers",
			metricsConsumers: nil,
			expected:         0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			publisher := new(MetricsPublisher)
			publisher.SubscribeConsumers(tc.metricsConsumers...)

			assert.Equal(t, tc.expected, len(publisher.consumers))
		})
	}
}
