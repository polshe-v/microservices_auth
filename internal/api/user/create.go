package user

import (
	"context"
	"log"

	"github.com/polshe-v/microservices_auth/internal/converter"
	desc "github.com/polshe-v/microservices_auth/pkg/user_v1"
)

// Create is used for creating new user.
func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	user := req.GetUser()
	log.Printf("\n%s\nName: %s\nEmail: %s\nPassword: %s\nPassword confirm: %s\nRole: %v\n%s", delim, user.GetName(), user.GetEmail(), user.GetPassword(), user.GetPasswordConfirm(), user.GetRole(), delim)

	id, err := i.userService.Create(ctx, converter.ToUserCreateFromDesc(user))
	if err != nil {
		return nil, err
	}

	log.Printf("Created user with id: %d", id)

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
