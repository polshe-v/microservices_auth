package access

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"

	"github.com/polshe-v/microservices_auth/internal/model"
	"github.com/polshe-v/microservices_auth/internal/repository"
	"github.com/polshe-v/microservices_auth/internal/repository/access/converter"
	modelRepo "github.com/polshe-v/microservices_auth/internal/repository/access/model"
	"github.com/polshe-v/microservices_common/pkg/db"
)

const (
	tableName = "policies"

	endpointColumn     = "endpoint"
	allowedRolesColumn = "allowed_roles"
)

type repo struct {
	db db.Client
}

// NewRepository creates new object of repository layer.
func NewRepository(db db.Client) repository.AccessRepository {
	return &repo{db: db}
}

func (r *repo) GetRoleEndpoints(ctx context.Context) ([]*model.EndpointPermissions, error) {
	builderSelect := sq.Select(endpointColumn, allowedRolesColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "access_repository.GetRoleEndpoints",
		QueryRaw: query,
	}

	var endpointPermissions []*modelRepo.EndpointPermissions
	err = r.db.DB().ScanAllContext(ctx, &endpointPermissions, q, args...)
	if err != nil {
		return nil, err
	}

	return converter.ToEndpointPermissionsFromRepo(endpointPermissions), nil
}
