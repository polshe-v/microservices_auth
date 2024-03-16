package auth

import (
	"context"

	desc "github.com/polshe-v/microservices_auth/pkg/auth_v1"
)

// GetAccessToken returns access token for later operations.
func (i *Implementation) GetAccessToken(ctx context.Context, req *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {
	accessToken, err := i.authService.GetAccessToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	return &desc.GetAccessTokenResponse{
		AccessToken: accessToken,
	}, nil
}
