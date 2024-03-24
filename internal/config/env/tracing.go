package env

import (
	"errors"
	"net"
	"os"

	"github.com/polshe-v/microservices_auth/internal/config"
)

const (
	tracingHostEnvName        = "JAEGER_GRPC_EXPORTER_HOST"
	tracingPortEnvName        = "JAEGER_GRPC_EXPORTER_PORT"
	tracingServiceNameEnvName = "JAEGER_SERVICE_NAME"
)

type tracingConfig struct {
	host        string
	port        string
	serviceName string
}

var _ config.TracingConfig = (*tracingConfig)(nil)

// NewTracingConfig creates new object of TracingConfig interface.
func NewTracingConfig() (config.TracingConfig, error) {
	host := os.Getenv(tracingHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("tracing exporter host not found")
	}

	port := os.Getenv(tracingPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("tracing exporter port not found")
	}

	serviceName := os.Getenv(tracingServiceNameEnvName)
	if len(serviceName) == 0 {
		return nil, errors.New("tracing service name not found")
	}

	return &tracingConfig{
		host:        host,
		port:        port,
		serviceName: serviceName,
	}, nil
}

func (cfg *tracingConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func (cfg *tracingConfig) ServiceName() string {
	return cfg.serviceName
}
