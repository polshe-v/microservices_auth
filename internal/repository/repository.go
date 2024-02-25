package repository

import (
	"context"

	desc "github.com/polshe-v/microservices_auth/pkg/user_v1"
)

// UserRepository is the interface for repository communication.
type UserRepository interface {
	Create(ctx context.Context, user *desc.UserCreate) (int64, error)
	Get(ctx context.Context, id int64) (*desc.User, error)
	Update(ctx context.Context, user *desc.UserUpdate) error
	Delete(ctx context.Context, id int64) error
}
