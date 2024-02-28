package user

import (
	"context"
	"fmt"

	"github.com/polshe-v/microservices_auth/internal/model"
)

func (s *serv) Get(ctx context.Context, id int64) (*model.User, error) {
	var user *model.User
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		user, errTx = s.userRepository.Get(ctx, id)
		if errTx != nil {
			return errTx
		}

		errTx = s.logRepository.Log(ctx, &model.Log{
			Log: fmt.Sprintf("Read info about user with id: %d", id),
		})
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return user, nil
}
