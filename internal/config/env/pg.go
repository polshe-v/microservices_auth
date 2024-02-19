package env

import (
	"errors"
	"os"

	"github.com/polshe-v/microservices_auth/internal/config"
)

const (
	dsnEnvName = "PG_DSN"
)

type pgConfig struct {
	dsn string
}

var _ config.PgConfig = (*pgConfig)(nil)

// NewPgConfig creates new object of PgConfig interface.
func NewPgConfig() (config.PgConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn not found")
	}

	return &pgConfig{
		dsn: dsn,
	}, nil
}

func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
