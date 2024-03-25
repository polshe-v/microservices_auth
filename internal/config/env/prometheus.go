package env

import (
	"errors"
	"net"
	"os"

	"github.com/polshe-v/microservices_auth/internal/config"
)

const (
	prometheusHostEnvName = "PROMETHEUS_HTTP_HOST"
	prometheusPortEnvName = "PROMETHEUS_HTTP_PORT"
)

type prometheusConfig struct {
	host string
	port string
}

var _ config.PrometheusConfig = (*prometheusConfig)(nil)

// NewPrometheusConfig creates new object of PrometheusConfig interface.
func NewPrometheusConfig() (config.PrometheusConfig, error) {
	host := os.Getenv(prometheusHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("prometheus host not found")
	}

	port := os.Getenv(prometheusPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("prometheus port not found")
	}

	return &prometheusConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *prometheusConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
