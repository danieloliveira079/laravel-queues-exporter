package redis

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_RedisConnectionConfig_ShouldCheckRequiredConnectionInfo(t *testing.T) {

	testCases := []struct {
		desc     string
		config   ConnectionConfig
		expected bool
	}{
		{
			desc: "All connection info is fulfilled",
			config: ConnectionConfig{
				Host: "localhost",
				Port: "6379",
			},
			expected: true,
		}, {
			desc: "Host info is missing",
			config: ConnectionConfig{
				Host: "",
				Port: "6379",
			},
			expected: false,
		}, {
			desc: "Port info is missing",
			config: ConnectionConfig{
				Host: "localhost",
				Port: "",
			},
			expected: false,
		}, {
			desc: "Host and Port info are missing",
			config: ConnectionConfig{
				Host: "",
				Port: "",
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			hasRequiredInfo, err := tc.config.HasRequiredConnectionInfo()
			require.Equal(t, tc.expected, hasRequiredInfo, err)
		})
	}

}

func Test_RedisConnectionConfig_ShouldCheckHasValidDB(t *testing.T) {

	testCases := []struct {
		desc     string
		config   ConnectionConfig
		expected bool
	}{
		{
			desc: "DB index equals 0",
			config: ConnectionConfig{
				DB: 0,
			},
			expected: true,
		},
		{
			desc: "DB index lower than 0",
			config: ConnectionConfig{
				DB: -1,
			},
			expected: false,
		},
		{
			desc: "DB index greater than 0",
			config: ConnectionConfig{
				DB: 1,
			},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			hasValidDB, err := tc.config.HasValidDB()
			require.Equal(t, tc.expected, hasValidDB, err)
		})
	}

}

func Test_RedisConnectionConfig_ShouldCheckHasValidPort(t *testing.T) {

	testCases := []struct {
		desc     string
		config   ConnectionConfig
		expected bool
	}{
		{
			desc: "Port is higher than zero",
			config: ConnectionConfig{
				Port: "6379",
			},
			expected: true,
		},
		{
			desc: "Port is equals zero",
			config: ConnectionConfig{
				Port: "0",
			},
			expected: true,
		},
		{
			desc: "Port is lower than zero",
			config: ConnectionConfig{
				Port: "-6379",
			},
			expected: false,
		},
		{
			desc: "Port is blank",
			config: ConnectionConfig{
				Port: "",
			},
			expected: false,
		},
		{
			desc:     "Port is null",
			config:   ConnectionConfig{},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			hasValidPort, err := tc.config.HasValidPort()
			require.Equal(t, tc.expected, hasValidPort, err)
		})
	}

}

func Test_RedisConnectionConfig_ShouldCheckHasValidConnectionInfo(t *testing.T) {

	testCases := []struct {
		desc     string
		config   ConnectionConfig
		expected bool
	}{
		{
			desc: "All configs set properly",
			config: ConnectionConfig{
				Host: "localhost",
				Port: "6379",
				DB:   0,
			},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			hasValidInfo, err := tc.config.HasValidConnectionInfo()
			require.Equal(t, tc.expected, hasValidInfo, err)
		})
	}

}
