package user

import (
	"context"
	"fmt"
	"log"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/polshe-v/microservices_auth/internal/model"
)

// ErrNameExists - custom error for user name duplicate.
var ErrNameExists = errors.New("user with provided name already exists")

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
		log.Print(err)
		if errors.Cause(err) == ErrNameExists {
			return 0, ErrNameExists
		}
		return 0, errors.New("failed to create user")
	}

	return id, nil
}
