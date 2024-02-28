package user

import (
	"context"
	"errors"
	"log"

	sq "github.com/Masterminds/squirrel"

	"github.com/polshe-v/microservices_auth/internal/client/db"
	"github.com/polshe-v/microservices_auth/internal/model"
	"github.com/polshe-v/microservices_auth/internal/repository"
)

const (
	tableName = "transaction_log"

	idColumn        = "id"
	timestampColumn = "timestamp"
	logColumn       = "log"
)

var errQueryBuild = errors.New("failed to build query")

type repo struct {
	db db.Client
}

// NewRepository creates new object of repository layer.
func NewRepository(db db.Client) repository.LogRepository {
	return &repo{db: db}
}

func (r *repo) Log(ctx context.Context, text *model.Log) error {
	builderInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(logColumn).
		Values(text.Log).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("%v", err)
		return errQueryBuild
	}

	q := db.Query{
		Name:     "log_repository.Log",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		log.Printf("%v", err)
		return errors.New("failed to create transaction log record")
	}

	return nil
}
