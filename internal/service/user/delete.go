package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/polshe-v/microservices_auth/internal/model"
)

func (s *serv) Delete(ctx context.Context, id int64) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		rowsNumber, errTx := s.userRepository.Delete(ctx, id)
		if errTx != nil {
			return errTx
		}
		if rowsNumber == 0 {
			return errors.New("no user found to delete")
		}

		errTx = s.logRepository.Log(ctx, &model.Log{
			Log: fmt.Sprintf("Deleted user with id: %d", id),
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
