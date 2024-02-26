package user

import (
	"context"
	"log"

	"github.com/golang/protobuf/ptypes/empty"

	desc "github.com/polshe-v/microservices_auth/pkg/user_v1"
)

// Delete is used for deleting user.
func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	log.Printf("\n%s\nID: %d\n%s", delim, req.GetId(), delim)

	err := i.userService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
