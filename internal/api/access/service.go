package access

import (
	"github.com/polshe-v/microservices_auth/internal/service"
	desc "github.com/polshe-v/microservices_auth/pkg/access_v1"
)

// Implementation structure describes API layer.
type Implementation struct {
	desc.UnimplementedAccessV1Server
	accessService service.AccessService
}

// NewImplementation creates new object of API layer.
func NewImplementation(accessService service.AccessService) *Implementation {
	return &Implementation{
		accessService: accessService,
	}
}
