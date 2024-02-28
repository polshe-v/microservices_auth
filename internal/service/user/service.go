package user

import (
	"github.com/polshe-v/microservices_auth/internal/client/db"

	"github.com/polshe-v/microservices_auth/internal/repository"
	"github.com/polshe-v/microservices_auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
	logRepository  repository.LogRepository
	txManager      db.TxManager
}

// NewService creates new object of service layer.
func NewService(userRepository repository.UserRepository, logRepository repository.LogRepository, txManager db.TxManager) service.UserService {
	return &serv{
		userRepository: userRepository,
		logRepository:  logRepository,
		txManager:      txManager,
	}
}
