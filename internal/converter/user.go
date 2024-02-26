package converter

import (
	"database/sql"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/polshe-v/microservices_auth/internal/model"
	desc "github.com/polshe-v/microservices_auth/pkg/user_v1"
)

// ToUserFromService converts service layer model to structure of API layer.
func ToUserFromService(user *model.User) *desc.User {
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

// ToUserFromDesc converts structure of API layer to service layer model.
func ToUserFromDesc(user *desc.User) *model.User {
	updatedAt := sql.NullTime{
		Time:  user.UpdatedAt.AsTime(),
		Valid: true,
	}

	return &model.User{
		ID:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Role:      desc.Role_name[int32(user.Role)],
		CreatedAt: user.CreatedAt.AsTime(),
		UpdatedAt: updatedAt,
	}
}

// ToUserCreateFromDesc converts structure of API layer to service layer model.
func ToUserCreateFromDesc(user *desc.UserCreate) *model.UserCreate {
	return &model.UserCreate{
		Name:            user.Name,
		Email:           user.Email,
		Password:        user.Password,
		PasswordConfirm: user.PasswordConfirm,
		Role:            desc.Role_name[int32(user.Role)],
	}
}

// ToUserUpdateFromDesc converts structure of API layer to service layer model.
func ToUserUpdateFromDesc(user *desc.UserUpdate) *model.UserUpdate {
	return &model.UserUpdate{
		ID:    user.Id,
		Name:  user.Name.GetValue(),
		Email: user.Email.GetValue(),
		Role:  desc.Role_name[int32(user.Role)],
	}
}
