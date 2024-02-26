package user

import (
	"context"
	"log"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/polshe-v/microservices_auth/internal/converter"
	desc "github.com/polshe-v/microservices_auth/pkg/user_v1"
)

// Update is used for updating user info.
func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*empty.Empty, error) {
	user := req.GetUser()
	log.Printf("\n%s\nID: %d\nName: %s\nEmail: %s\nRole: %v\n%s", delim, user.GetId(), user.GetName().GetValue(), user.GetEmail().GetValue(), user.GetRole(), delim)

	err := i.userService.Update(ctx, converter.ToUserUpdateFromDesc(user))
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
