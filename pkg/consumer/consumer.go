package consumer

import (
	"errors"
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/config"
	log_consumer "github.com/danieloliveira079/laravel-queues-exporter/pkg/consumer/log"
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/consumer/statsd"
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/consumer/stdout"
	"github.com/danieloliveira079/laravel-queues-exporter/pkg/metric"
	"strings"
)

type Consumer interface {
	Process(metrics []metric.Metric)
}

type ConsumerFactory struct {
}

var errorConsumersConfigIsEmpty = errors.New("The provided list of consumers from config is empty")
var errorConsumerTypeNotRegistered = errors.New("Consumer type not registered")

func BuildConsumersListFromConfig(appConfig *config.AppConfig) ([]Consumer, error) {
	var consumers []Consumer
	var err error

	if len(appConfig.ExportTo) == 0 {
		return nil, errorConsumersConfigIsEmpty
	}

	splitConsumers := splitFromConfig(appConfig.ExportTo)
	factory := new(ConsumerFactory)

	for _, consumerType := range splitConsumers {
		consumer, err := factory.NewConsumer(appConfig, consumerType)
		if err != nil {
			return nil, err
		}

		consumers = append(consumers, consumer)
	}

	return consumers, err

}

func splitFromConfig(config string) []string {
	return strings.Split(config, ",")
}

func (f *ConsumerFactory) NewConsumer(config *config.AppConfig, consumerType string) (Consumer, error) {
	var consumer Consumer
	var err error

	switch consumerType {
	case "stdout":
		consumer, err = stdout.New(config)
	case "statsd":
		consumer, err = statsd.New(config)
	case "log":
		consumer, err = log_consumer.New(config)
	default:
		err = errorConsumerTypeNotRegistered
	}

	return consumer, err
}
