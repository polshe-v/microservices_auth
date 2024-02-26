package user

import (
	"github.com/polshe-v/microservices_auth/internal/repository"
	"github.com/polshe-v/microservices_auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
}

// NewService creates new object of service layer.
func NewService(userRepository repository.UserRepository) service.UserService {
	return &serv{
		userRepository: userRepository,
	}
}
