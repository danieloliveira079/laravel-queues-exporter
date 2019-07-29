package redis

type RedisCommandDispatcher struct {
	connector Connector
}

type Connector interface {
	Connect() (err error)
	Close() (err error)
	Do(command string, args ...interface{}) (results interface{}, err error)
}

func NewRedisCommandDispatcher(connector Connector) *RedisCommandDispatcher {
	if connector != nil {
		return &RedisCommandDispatcher{connector: connector}
	}
	return nil
}

func (d *RedisCommandDispatcher) Do(command string, args ...interface{}) (results interface{}, err error) {
	return d.connector.Do(command, args)
}
