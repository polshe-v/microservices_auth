package auth

import (
	"context"
	"errors"

	"github.com/polshe-v/microservices_auth/internal/model"
)

func (s *serv) GetRefreshToken(ctx context.Context, oldRefreshToken string) (string, error) {
	// Get secret key from storage for refresh token HMAC
	refreshTokenSecretKey, err := s.keyRepository.GetKey(ctx, refreshTokenSecretKeyName)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	claims, err := s.tokenOperations.Verify(oldRefreshToken, []byte(refreshTokenSecretKey))
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	refreshToken, err := s.tokenOperations.Generate(model.User{
		Name: claims.Username,
		Role: claims.Role,
	},
		[]byte(refreshTokenSecretKey),
		refreshTokenExpiration,
	)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return refreshToken, nil
}
