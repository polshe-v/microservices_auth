package jwt

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"

	"github.com/polshe-v/microservices_auth/internal/model"
	"github.com/polshe-v/microservices_auth/internal/tokens"
)

type tokenOperations struct{}

var _ tokens.TokenOperations = (*tokenOperations)(nil)

// NewTokenOperations creates new object for using token functions.
func NewTokenOperations() tokens.TokenOperations {
	return &tokenOperations{}
}

// Generate creates JWT for user.
func (t *tokenOperations) Generate(user model.User, secretKey []byte, duration time.Duration) (string, error) {
	claims := model.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
		Username: user.Name,
		Role:     user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

// Verify checks validity of provided JWT.
func (t *tokenOperations) Verify(tokenStr string, secretKey []byte) (*model.UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&model.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.Errorf("unexpected token signing method")
			}

			return secretKey, nil
		},
	)
	if err != nil {
		return nil, errors.Errorf("invalid token: %s", err.Error())
	}

	claims, ok := token.Claims.(*model.UserClaims)
	if !ok {
		return nil, errors.Errorf("invalid token claims")
	}

	return claims, nil
}
