package key

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/polshe-v/microservices_auth/internal/repository"
	"github.com/polshe-v/microservices_common/pkg/db"
)

const (
	tableName = "hmac_keys"

	keyNameColumn  = "key"
	keyValueColumn = "value"
)

type repo struct {
	db db.Client
}

// NewRepository creates new object of repository layer.
func NewRepository(db db.Client) repository.KeyRepository {
	return &repo{db: db}
}

func (r *repo) GetKey(ctx context.Context, keyName string) (string, error) {
	builderSelect := sq.Select(keyValueColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{keyNameColumn: keyName}).
		Limit(1)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return "", err
	}

	q := db.Query{
		Name:     "key_repository.GetKey",
		QueryRaw: query,
	}

	var keyValue string
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&keyValue)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", err
		}
		return "", err
	}

	return keyValue, nil
}
