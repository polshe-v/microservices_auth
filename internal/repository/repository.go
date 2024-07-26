package repository

import (
	"context"

	"github.com/polshe-v/microservices_auth/internal/model"
)

// UserRepository is the interface for user info repository communication.
type UserRepository interface {
	Create(ctx context.Context, user *model.UserCreate) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.UserUpdate) error
	Delete(ctx context.Context, id int64) error
	GetAuthInfo(ctx context.Context, username string) (*model.AuthInfo, error)
}

// KeyRepository is the interface for key info repository communication.
type KeyRepository interface {
	GetKey(ctx context.Context, keyName string) (string, error)
}

// AccessRepository is the interface for access policies repository communication.
type AccessRepository interface {
	GetRoleEndpoints(ctx context.Context) ([]*model.EndpointPermissions, error)
}

// LogRepository is the interface for transaction log repository communication.
type LogRepository interface {
	Log(ctx context.Context, log *model.Log) error
}

// CacheRepository is the interface for cache communication.
type CacheRepository interface {
	CreateRecord(ctx context.Context, user *model.User) error
	GetRecord(ctx context.Context, id int64) (*model.User, error)
}
