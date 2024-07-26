package converter

import (
	"database/sql"
	"time"

	model "github.com/polshe-v/microservices_auth/internal/model"
	modelRepo "github.com/polshe-v/microservices_auth/internal/repository/cache/model"
)

// ToUserFromRepo converts repository layer model to structure of service layer.
func ToUserFromRepo(user *modelRepo.User) *model.User {
	var updatedAt sql.NullTime
	if user.UpdatedAtNs != nil {
		updatedAt = sql.NullTime{
			Time:  time.Unix(0, *user.UpdatedAtNs),
			Valid: true,
		}
	}

	return &model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: time.Unix(0, user.CreatedAtNs),
		UpdatedAt: updatedAt,
	}
}
