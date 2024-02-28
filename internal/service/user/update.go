package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/polshe-v/microservices_auth/internal/model"
)

func (s *serv) Update(ctx context.Context, user *model.UserUpdate) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		rowsNumber, errTx := s.userRepository.Update(ctx, user)
		if errTx != nil {
			return errTx
		}
		if rowsNumber == 0 {
			return errors.New("no user found to update")
		}

		errTx = s.logRepository.Log(ctx, &model.Log{
			Log: fmt.Sprintf("Updated user with id: %d", user.ID),
		})
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
