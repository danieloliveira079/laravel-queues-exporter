package redis

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_RedisConnectionConfig_ShouldCheckRequiredConnectionInfo(t *testing.T) {

	testCases := []struct {
		desc     string
		config   RedisConnectionConfig
		expected bool
	}{
		{
			desc: "All connection info is fulfilled",
			config: RedisConnectionConfig{
				Host: "localhost",
				Port: "6379",
			},
			expected: true,
		}, {
			desc: "Host info is missing",
			config: RedisConnectionConfig{
				Host: "",
				Port: "6379",
			},
			expected: false,
		}, {
			desc: "Port info is missing",
			config: RedisConnectionConfig{
				Host: "localhost",
				Port: "",
			},
			expected: false,
		}, {
			desc: "Host and Port info are missing",
			config: RedisConnectionConfig{
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
		config   RedisConnectionConfig
		expected bool
	}{
		{
			desc: "DB index equals 0",
			config: RedisConnectionConfig{
				DB: 0,
			},
			expected: true,
		},
		{
			desc: "DB index lower than 0",
			config: RedisConnectionConfig{
				DB: -1,
			},
			expected: false,
		},
		{
			desc: "DB index greater than 0",
			config: RedisConnectionConfig{
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
		config   RedisConnectionConfig
		expected bool
	}{
		{
			desc: "Port is higher than zero",
			config: RedisConnectionConfig{
				Port: "6379",
			},
			expected: true,
		},
		{
			desc: "Port is equals zero",
			config: RedisConnectionConfig{
				Port: "0",
			},
			expected: true,
		},
		{
			desc: "Port is lower than zero",
			config: RedisConnectionConfig{
				Port: "-6379",
			},
			expected: false,
		},
		{
			desc: "Port is blank",
			config: RedisConnectionConfig{
				Port: "",
			},
			expected: false,
		},
		{
			desc:     "Port is null",
			config:   RedisConnectionConfig{},
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
		config   RedisConnectionConfig
		expected bool
	}{
		{
			desc: "All configs set properly",
			config: RedisConnectionConfig{
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
