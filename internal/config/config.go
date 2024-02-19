package config

import (
	"github.com/joho/godotenv"
)

// GrpcConfig provides gRPC settings from config file.
type GrpcConfig interface {
	Address() string
	Transport() string
}

// PgConfig provides PostgreSQL settings from config file.
type PgConfig interface {
	DSN() string
}

// Load reads .env config file.
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}
