package auth

import (
	"context"
	"log"

	"github.com/pkg/errors"

	"github.com/polshe-v/microservices_auth/internal/model"
	"github.com/polshe-v/microservices_auth/internal/utils"
)

func (s *serv) GetRefreshToken(ctx context.Context, oldRefreshToken string) (string, error) {
	// Get secret key from storage for refresh token HMAC
	refreshTokenSecretKey, err := s.keyRepository.GetKey(ctx, refreshTokenSecretKeyName)
	if err != nil {
		log.Print(err)
		return "", errors.New("failed to generate token")
	}

	claims, err := utils.VerifyToken(oldRefreshToken, []byte(refreshTokenSecretKey))
	if err != nil {
		log.Print(err)
		return "", errors.New("invalid refresh token")
	}

	refreshToken, err := utils.GenerateToken(model.User{
		Name: claims.Username,
		Role: claims.Role,
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
