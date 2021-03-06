package redis

import "github.com/danieloliveira079/laravel-queues-exporter/pkg/config"

type RedisExporterBuilder struct {
}

//TODO Create tests for Builder
func (b *RedisExporterBuilder) Build(appConfig *config.AppConfig) (*Exporter, error) {
	connectionConfig := createRedisConnectionConfig(appConfig)
	connector, err := createRedisConnector(connectionConfig)
	if err != nil {
		return nil, err
	}

	dispatcher, err := createRedisCommandDispatcher(connector)
	if err != nil {
		return nil, err
	}

	extractor, err := createRedisExtractor(dispatcher)
	if err != nil {
		return nil, err
	}

	exporterConfig := createRedisExporterConfig(appConfig, connectionConfig)

	return NewRedisExporter(exporterConfig, connector, extractor)
}

func createRedisConnectionConfig(appConfig *config.AppConfig) *ConnectionConfig {
	return &ConnectionConfig{
		Host: appConfig.RedisHost,
		Port: appConfig.RedisPort,
		DB:   appConfig.RedisDB,
	}
}

func createRedisConnector(connectionConfig *ConnectionConfig) (Connector, error) {
	return NewRedisConnector(connectionConfig)
}

func createRedisCommandDispatcher(connector Connector) (CommandDispatcher, error) {
	return NewRedisCommandDispatcher(connector)
}

func createRedisExtractor(dispatcher CommandDispatcher) (Extractor, error) {
	return NewRedisExtractor(dispatcher)
}

func createRedisExporterConfig(appConfig *config.AppConfig, connectionConfig *ConnectionConfig) *ExporterConfig {
	return &ExporterConfig{
		QueueNames:       appConfig.QueuesNames,
		CollectInterval:  appConfig.CollectInterval,
		ConnectionConfig: connectionConfig,
	}
}
