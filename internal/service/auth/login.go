package auth

import (
	"context"
	"log"

	"github.com/pkg/errors"

	"github.com/polshe-v/microservices_auth/internal/model"
	"github.com/polshe-v/microservices_auth/internal/utils"
)

func (s *serv) Login(ctx context.Context, creds *model.UserCreds) (string, error) {
	// Get role and hashed password by username from storage
	authInfo, err := s.userRepository.GetAuthInfo(ctx, creds.Username)
	if err != nil {
		log.Print(err)
		return "", errors.New("no user found")
	}

	if !utils.VerifyPassword(authInfo.Password, creds.Password) {
		log.Print(err)
		return "", errors.New("wrong password")
	}

	// Get secret key from storage for refresh token HMAC
	refreshTokenSecretKey, err := s.keyRepository.GetKey(ctx, refreshTokenSecretKeyName)
	if err != nil {
		log.Print(err)
		return "", errors.New("failed to generate token")
	}

	refreshToken, err := utils.GenerateToken(model.User{
		Name: authInfo.Username,
		Role: authInfo.Role,
	},
		[]byte(refreshTokenSecretKey),
		refreshTokenExpiration,
	)
	if err != nil {
		log.Print(err)
		return "", errors.New("failed to generate token")
	}

	return refreshToken, nil
}
