package service

import (
	"context"

	"github.com/polshe-v/microservices_auth/internal/model"
)

// UserService is the interface for service communication.
type UserService interface {
	Create(ctx context.Context, user *model.UserCreate) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.UserUpdate) error
	Delete(ctx context.Context, id int64) error
}
