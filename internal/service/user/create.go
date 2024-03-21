package user

import (
	"context"
	"errors"
	"fmt"

	errorsExt "github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/polshe-v/microservices_auth/internal/model"
)

// ErrUserExists - custom error for user name duplicate.
var ErrUserExists = errors.New("user with provided name or email already exists")

func (s *serv) Create(ctx context.Context, user *model.UserCreate) (int64, error) {
	if user.Password != user.PasswordConfirm {
		return 0, errors.New("passwords don't match")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, errors.New("failed to process password")
	}
	user.Password = string(hashedPassword)

	var id int64
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.userRepository.Create(ctx, user)
		if errTx != nil {
			return errTx
		}

		errTx = s.logRepository.Log(ctx, &model.Log{
			Text: fmt.Sprintf("Created user with id: %d", id),
		})
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		if errorsExt.Cause(err) == ErrUserExists {
			return 0, ErrUserExists
		}
		return 0, errors.New("failed to create user")
	}

	return id, nil
}
