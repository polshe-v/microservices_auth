package auth

import (
	"time"

	"github.com/polshe-v/microservices_auth/internal/repository"
	"github.com/polshe-v/microservices_auth/internal/service"
	"github.com/polshe-v/microservices_auth/internal/tokens"
)

const (
	refreshTokenExpiration    = 60 * time.Minute
	accessTokenExpiration     = 5 * time.Minute
	refreshTokenSecretKeyName = "refresh"
	accessTokenSecretKeyName  = "access"
)

type serv struct {
	userRepository  repository.UserRepository
	keyRepository   repository.KeyRepository
	tokenOperations tokens.TokenOperations
}

// NewService creates new object of service layer.
func NewService(userRepository repository.UserRepository, keyRepository repository.KeyRepository, tokenOperations tokens.TokenOperations) service.AuthService {
	return &serv{
		userRepository:  userRepository,
		keyRepository:   keyRepository,
		tokenOperations: tokenOperations,
	}
}
