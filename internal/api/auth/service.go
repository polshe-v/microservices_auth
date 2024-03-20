package auth

import (
	"github.com/polshe-v/microservices_auth/internal/service"
	desc "github.com/polshe-v/microservices_auth/pkg/auth_v1"
)

// Implementation structure describes API layer.
type Implementation struct {
	desc.UnimplementedAuthV1Server
	authService service.AuthService
}

// NewImplementation creates new object of API layer.
func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{
		authService: authService,
	}
}
