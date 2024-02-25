package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/polshe-v/microservices_auth/internal/repository/user/model"
	desc "github.com/polshe-v/microservices_auth/pkg/user_v1"
)

// ToUserFromRepo convers model from repository layer to structure of service layer.
func ToUserFromRepo(user *model.User) *desc.User {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      desc.Role(desc.Role_value[user.Role]),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}
