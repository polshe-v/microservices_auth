package env

import (
	"errors"
	"net"
	"os"

	"github.com/polshe-v/microservices_auth/internal/config"
)

const (
	grpcHostEnvName      = "GRPC_HOST"
	grpcPortEnvName      = "GRPC_PORT"
	grpcTransportEnvName = "GRPC_TRANSPORT"
)

type grpcConfig struct {
	host      string
	port      string
	transport string
}

var _ config.GrpcConfig = (*grpcConfig)(nil)

// NewGrpcConfig creates new object of GrpcConfig interface.
func NewGrpcConfig() (config.GrpcConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("grpc host not found")
	}

	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("grpc port not found")
	}

	transport := os.Getenv(grpcTransportEnvName)
	if len(transport) == 0 {
		return nil, errors.New("grpc transport not found")
	}

	return &grpcConfig{
		host:      host,
		port:      port,
		transport: transport,
	}, nil
}

func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func (cfg *grpcConfig) Transport() string {
	return cfg.transport
}
