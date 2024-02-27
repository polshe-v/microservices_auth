package pg

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/polshe-v/microservices_auth/internal/client/db"
)

type pgClient struct {
	masterDBC db.DB
}

// New creates DB client object.
func New(ctx context.Context, dsn string) (db.Client, error) {
	dbc, err := pgxpool.New(ctx, dsn)

	if err != nil {
		log.Printf("%v", err)
		return nil, errors.New("failed to connect to db")
	}

	return &pgClient{
		masterDBC: &pg{dbc: dbc},
	}, nil
}

func (c *pgClient) DB() db.DB {
	return c.masterDBC
}

func (c *pgClient) Close() error {
	if c.masterDBC != nil {
		c.masterDBC.Close()
	}
	return nil
}
