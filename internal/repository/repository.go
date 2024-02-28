package repository

import (
	"context"

	"github.com/polshe-v/microservices_auth/internal/model"
)

// UserRepository is the interface for user info repository communication.
type UserRepository interface {
	Create(ctx context.Context, user *model.UserCreate) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.UserUpdate) (int64, error)
	Delete(ctx context.Context, id int64) (int64, error)
}

// LogRepository is the interface for transaction log repository communication.
type LogRepository interface {
	Log(ctx context.Context, log *model.Log) error
}
