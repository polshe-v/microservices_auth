package auth

import (
	"context"

	desc "github.com/polshe-v/microservices_auth/pkg/auth_v1"
)

// GetRefreshToken updates refresh token.
func (i *Implementation) GetRefreshToken(ctx context.Context, req *desc.GetRefreshTokenRequest) (*desc.GetRefreshTokenResponse, error) {
	refreshToken, err := i.authService.GetRefreshToken(ctx, req.GetOldRefreshToken())
	if err != nil {
		return nil, err
	}

	return &desc.GetRefreshTokenResponse{
		RefreshToken: refreshToken,
	}, nil
}
