package user

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	desc "github.com/polshe-v/microservices_auth/pkg/user_v1"
)

// Delete is used for deleting user.
func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	err := i.userService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
