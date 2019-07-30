package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

type RedisConnector struct {
	Config ConnectorConfig
	conn   redis.Conn
}

type ConnectorConfig struct {
	ConnConfig ConnectionConfig
}

func NewRedisConnector(config ConnectorConfig) (*RedisConnector, error) {
	return &RedisConnector{Config: config}, nil
}

func (c *RedisConnector) Connect() (err error) {
	if c.conn != nil && c.conn.Err() == nil {
		return nil
	}

	conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", c.Config.ConnConfig.Host, c.Config.ConnConfig.Port), redis.DialDatabase(c.Config.ConnConfig.DB), redis.DialConnectTimeout(15*time.Second))
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

func (c *RedisConnector) Close() (err error) {
	log.Println("Closing connector")
	if c.conn != nil {
		err = c.conn.Close()
	}

	return err
}

func (c *RedisConnector) Do(command string, args ...interface{}) (reply interface{}, err error) {
	err = c.Connect()
	if err != nil {
		return nil, err
	}

	reply, err = c.conn.Do(command, args...)

	return reply, err
}
