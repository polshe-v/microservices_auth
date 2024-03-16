package auth

import (
	"context"
	"log"

	"github.com/pkg/errors"

	"github.com/polshe-v/microservices_auth/internal/model"
	"github.com/polshe-v/microservices_auth/internal/utils"
)

func (s *serv) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	// Get secret key from storage for refresh token HMAC
	refreshTokenSecretKey, err := s.keyRepository.GetKey(ctx, refreshTokenSecretKeyName)
	if err != nil {
		log.Print(err)
		return "", errors.New("failed to generate token")
	}

	claims, err := utils.VerifyToken(refreshToken, []byte(refreshTokenSecretKey))
	if err != nil {
		log.Print(err)
		return "", errors.New("invalid refresh token")
	}

	// Get secret key from storage for access token HMAC
	accessTokenSecretKey, err := s.keyRepository.GetKey(ctx, accessTokenSecretKeyName)
	if err != nil {
		log.Print(err)
		return "", errors.New("failed to generate token")
	}

	accessToken, err := utils.GenerateToken(model.User{
		Name: claims.Username,
		Role: claims.Role,
	},
		[]byte(accessTokenSecretKey),
		accessTokenExpiration,
	)
	if err != nil {
		log.Print(err)
		return "", errors.New("failed to generate token")
	}

	return accessToken, nil
}
