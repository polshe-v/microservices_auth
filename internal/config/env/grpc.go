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
	grpcCertPathEnvName  = "GRPC_CERT_PATH"
	grpcKeyPathEnvName   = "GRPC_KEY_PATH"
	grpcCaPathEnvName    = "GRPC_CA_PATH"
)

type grpcConfig struct {
	host      string
	port      string
	transport string
	certPath  string
	keyPath   string
	caPath    string
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

	certPath := os.Getenv(grpcCertPathEnvName)
	if len(certPath) == 0 {
		return nil, errors.New("grpc certificate not found")
	}

	keyPath := os.Getenv(grpcKeyPathEnvName)
	if len(keyPath) == 0 {
		return nil, errors.New("grpc private key not found")
	}

	caPath := os.Getenv(grpcCaPathEnvName)
	if len(caPath) == 0 {
		return nil, errors.New("grpc CA certificate not found")
	}

	return &grpcConfig{
		host:      host,
		port:      port,
		transport: transport,
		certPath:  certPath,
		keyPath:   keyPath,
		caPath:    caPath,
	}, nil
}

func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func (cfg *grpcConfig) Transport() string {
	return cfg.transport
}

func (cfg *grpcConfig) CertPath() string {
	return cfg.certPath
}

func (cfg *grpcConfig) KeyPath() string {
	return cfg.keyPath
}

func (cfg *grpcConfig) CaPath() string {
	return cfg.caPath
}
