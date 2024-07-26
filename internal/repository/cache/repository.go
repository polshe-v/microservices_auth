package cache

import (
	"context"
	"strconv"

	redigo "github.com/gomodule/redigo/redis"

	"github.com/polshe-v/microservices_auth/internal/model"
	"github.com/polshe-v/microservices_auth/internal/repository"
	"github.com/polshe-v/microservices_auth/internal/repository/cache/converter"
	modelRepo "github.com/polshe-v/microservices_auth/internal/repository/cache/model"
	"github.com/polshe-v/microservices_common/pkg/cache"
)

type repo struct {
	client cache.Client
}

// NewRepository creates new object of repository layer.
func NewRepository(client cache.Client) repository.CacheRepository {
	return &repo{client: client}
}

func (r *repo) CreateRecord(ctx context.Context, user *model.User) error {
	var updatedAtNs *int64

	if user.UpdatedAt.Valid {
		updatedAt := user.UpdatedAt.Time.UnixNano()
		updatedAtNs = &updatedAt
	}

	userRecord := modelRepo.User{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Role:        user.Role,
		CreatedAtNs: user.CreatedAt.UnixNano(),
		UpdatedAtNs: updatedAtNs,
	}

	idStr := strconv.FormatInt(user.ID, 10)
	err := r.client.HSet(ctx, idStr, userRecord)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) GetRecord(ctx context.Context, id int64) (*model.User, error) {
	idStr := strconv.FormatInt(id, 10)
	values, err := r.client.HGetAll(ctx, idStr)
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, model.ErrorUserNotFound
	}

	var user modelRepo.User
	err = redigo.ScanStruct(values, &user)
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}
