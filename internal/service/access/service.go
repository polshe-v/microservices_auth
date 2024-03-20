package access

import (
	"time"

	"github.com/polshe-v/microservices_auth/internal/repository"
	"github.com/polshe-v/microservices_auth/internal/service"
	"github.com/polshe-v/microservices_auth/internal/tokens"
)

const (
	accessTokenExpiration    = 5 * time.Minute
	accessTokenSecretKeyName = "access"
)

type serv struct {
	accessRepository repository.AccessRepository
	keyRepository    repository.KeyRepository
	tokenOperations  tokens.TokenOperations
}

// NewService creates new object of service layer.
func NewService(accessRepository repository.AccessRepository, keyRepository repository.KeyRepository, tokenOperations tokens.TokenOperations) service.AccessService {
	return &serv{
		accessRepository: accessRepository,
		keyRepository:    keyRepository,
		tokenOperations:  tokenOperations,
	}
}
