package user

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/polshe-v/microservices_auth/internal/model"
	"github.com/polshe-v/microservices_common/pkg/logger"
)

func (s *serv) Get(ctx context.Context, id int64) (*model.User, error) {
	var user *model.User

	// Check cache first
	user, err := s.cacheRepository.GetRecord(ctx, id)
	if err == nil {
		return user, nil
	}
	if err != model.ErrorUserNotFound {
		logger.Error("failed to read from cache: ", zap.Error(err))
	}

	// If record not found in cache, then get it from DB and put into cache
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		user, errTx = s.userRepository.Get(ctx, id)
		if errTx != nil {
			return errTx
		}

		errTx = s.logRepository.Log(ctx, &model.Log{
			Text: fmt.Sprintf("Read info about user with id: %d", id),
		})
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return nil, errors.New("failed to read user info")
	}

	err = s.cacheRepository.CreateRecord(ctx, user)
	if err != nil {
		logger.Error("failed to write to cache: ", zap.Error(err))
	}

	return user, nil
}
