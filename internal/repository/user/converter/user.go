package converter

import (
	model "github.com/polshe-v/microservices_auth/internal/model"
	modelRepo "github.com/polshe-v/microservices_auth/internal/repository/user/model"
)

// ToUserFromRepo converts repository layer model to structure of API layer.
func ToUserFromRepo(user *modelRepo.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
