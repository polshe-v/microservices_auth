package user

import (
	"context"

	"github.com/polshe-v/microservices_auth/internal/model"
)

func (s *serv) Update(ctx context.Context, user *model.UserUpdate) error {
	err := s.userRepository.Update(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
