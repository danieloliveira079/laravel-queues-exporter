package redis

import "errors"

type RedisCommandDispatcher struct {
	connector Connector
}

func NewRedisCommandDispatcher(connector Connector) (*RedisCommandDispatcher, error) {
	if connector == nil {
		return nil, errors.New("connector can't be nil")
	}

	return &RedisCommandDispatcher{connector: connector}, nil
}

func (d *RedisCommandDispatcher) Do(command string, args ...interface{}) (results interface{}, err error) {
	return d.connector.Do(command, args...)
}
