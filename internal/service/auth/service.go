package auth

import (
	"time"

	"github.com/polshe-v/microservices_auth/internal/repository"
	"github.com/polshe-v/microservices_auth/internal/service"
	"github.com/polshe-v/microservices_common/pkg/db"
)

const (
	refreshTokenExpiration    = 60 * time.Minute
	accessTokenExpiration     = 5 * time.Minute
	refreshTokenSecretKeyName = "refresh"
	accessTokenSecretKeyName  = "access"
)

type serv struct {
	userRepository repository.UserRepository
	keyRepository  repository.KeyRepository
	logRepository  repository.LogRepository
	txManager      db.TxManager
}

// NewService creates new object of service layer.
func NewService(userRepository repository.UserRepository, keyRepository repository.KeyRepository, logRepository repository.LogRepository, txManager db.TxManager) service.AuthService {
	return &serv{
		userRepository: userRepository,
		keyRepository:  keyRepository,
		logRepository:  logRepository,
		txManager:      txManager,
	}
}
