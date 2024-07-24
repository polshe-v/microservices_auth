package env

import (
	"net"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/polshe-v/microservices_auth/internal/config"
)

const (
	redisHostEnvName              = "REDIS_HOST"
	redisPortEnvName              = "REDIS_PORT"
	redisConnectionTimeoutEnvName = "REDIS_CONNECTION_TIMEOUT"
	redisIdleTimeoutEnvName       = "REDIS_IDLE_TIMEOUT"
	redisMaxIdleEnvName           = "REDIS_MAX_IDLE"

	defaultConnectionTimeout = 5   // seconds
	defaultIdleTimeout       = 300 // seconds
	defaultMaxIdle           = 10  // connections
)

type redisConfig struct {
	host              string
	port              string
	connectionTimeout time.Duration
	idleTimeout       time.Duration
	maxIdle           int
}

var _ config.RedisConfig = (*redisConfig)(nil)

// NewRedisConfig creates new object of RedisConfig interface.
func NewRedisConfig() (config.RedisConfig, error) {
	host := os.Getenv(redisHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("redis host not found")
	}

	port := os.Getenv(redisPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("redis port not found")
	}

	var connectionTimeout int
	connectionTimeoutStr := os.Getenv(redisConnectionTimeoutEnvName)
	if len(connectionTimeoutStr) == 0 {
		connectionTimeout = defaultConnectionTimeout
	} else {
		res, err := strconv.ParseUint(connectionTimeoutStr, 10, 32)
		if err != nil {
			return nil, errors.Errorf("failed to process %s setting", redisConnectionTimeoutEnvName)
		}
		connectionTimeout = int(res)
	}

	var idleTimeout int
	idleTimeoutStr := os.Getenv(redisIdleTimeoutEnvName)
	if len(idleTimeoutStr) == 0 {
		idleTimeout = defaultIdleTimeout
	} else {
		res, err := strconv.ParseUint(idleTimeoutStr, 10, 32)
		if err != nil {
			return nil, errors.Errorf("failed to process %s setting", redisIdleTimeoutEnvName)
		}
		idleTimeout = int(res)
	}

	var maxIdle int
	maxIdleStr := os.Getenv(redisMaxIdleEnvName)
	if len(maxIdleStr) == 0 {
		maxIdle = defaultMaxIdle
	} else {
		res, err := strconv.ParseUint(maxIdleStr, 10, 32)
		if err != nil {
			return nil, errors.Errorf("failed to process %s setting", redisMaxIdleEnvName)
		}
		maxIdle = int(res)
	}

	return &redisConfig{
		host:              host,
		port:              port,
		connectionTimeout: time.Duration(connectionTimeout) * time.Second,
		idleTimeout:       time.Duration(idleTimeout) * time.Second,
		maxIdle:           maxIdle,
	}, nil
}

func (cfg *redisConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func (cfg *redisConfig) ConnectionTimeout() time.Duration {
	return cfg.connectionTimeout
}

func (cfg *redisConfig) IdleTimeout() time.Duration {
	return cfg.idleTimeout
}

func (cfg *redisConfig) MaxIdle() int {
	return cfg.maxIdle
}
