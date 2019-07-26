package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

type RedisConnector struct {
	targetHost string
	targetPort string
	targetDB   int
	conn       redis.Conn
}

func NewRedisConnector(targetHost string, targetPort string, targetDB int) *RedisConnector {
	return &RedisConnector{targetHost: targetHost, targetPort: targetPort, targetDB: targetDB}
}

func (c *RedisConnector) Connect() (err error) {
	if c.conn != nil && c.conn.Err() == nil {
		return nil
	}

	conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", c.targetHost, c.targetPort), redis.DialDatabase(c.targetDB), redis.DialConnectTimeout(15*time.Second))
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
