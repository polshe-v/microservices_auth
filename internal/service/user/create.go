package user

import (
	"context"
	"fmt"

	"github.com/polshe-v/microservices_auth/internal/model"
)

func (s *serv) Create(ctx context.Context, user *model.UserCreate) (int64, error) {
	var id int64

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.userRepository.Create(ctx, user)
		if errTx != nil {
			return errTx
		}

		errTx = s.logRepository.Log(ctx, &model.Log{
			Log: fmt.Sprintf("Created user with id: %d", id),
		})
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
