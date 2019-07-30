package redis

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type ConnectionConfig struct {
	Host string
	Port string
	DB   int
}

func (c *ConnectionConfig) HasRequiredConnectionInfo() (bool, error) {
	var err error
	var missingInfo []string
	hasRequiredInfo := true

	if len(c.Host) == 0 {
		missingInfo = append(missingInfo, "host")
		hasRequiredInfo = false
	}

	if len(c.Port) == 0 {
		missingInfo = append(missingInfo, "port")
		hasRequiredInfo = false
	}

	errMessage := fmt.Sprintf("connection info has missing information for fields: %s", strings.Join(missingInfo, ", "))
	err = errors.New(errMessage)

	return hasRequiredInfo, err
}

func (c *ConnectionConfig) HasValidConnectionInfo() (bool, error) {
	var err error
	var invalidFieldMessages []string
	hasValidInfo := true

	if v, _ := c.HasValidDB(); v == false {
		invalidFieldMessages = append(invalidFieldMessages, fmt.Sprintf("DB: %d", c.DB))
		hasValidInfo = false
	}

	_, convErr := c.PortToInt()
	if convErr != nil {
		invalidFieldMessages = append(invalidFieldMessages, fmt.Sprintf("Port: %d", c.DB))
		hasValidInfo = false
	}

	errMessage := fmt.Sprintf("connection info has invalid fields: %s", strings.Join(invalidFieldMessages, ", "))
	err = errors.New(errMessage)

	return hasValidInfo, err
}

func (c *ConnectionConfig) HasValidDB() (bool, error) {
	if c.DB < 0 {
		return false, errors.New(fmt.Sprintf("db can't be lower than zero: %d", c.DB))
	}
	return true, nil
}

func (c *ConnectionConfig) HasValidPort() (bool, error) {
	if len(c.Port) < 0 {
		return false, errors.New(fmt.Sprintf("port can't be blank or null: %s", c.Port))
	}

	port, err := c.PortToInt()

	if err != nil {
		return false, err
	}

	if port < 0 {
		return false, errors.New(fmt.Sprintf("port can't be lower than zero: %s", c.Port))
	}

	return true, nil
}

func (c *ConnectionConfig) PortToInt() (int, error) {
	var err error
	var port int

	port, err = strconv.Atoi(c.Port)
	if err != nil {
		return -1, errors.New(fmt.Sprintf("error converting port: %s", c.Port))
	}

	return port, err
}
