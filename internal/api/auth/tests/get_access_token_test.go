package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	authAPI "github.com/polshe-v/microservices_auth/internal/api/auth"
	"github.com/polshe-v/microservices_auth/internal/service"
	serviceMocks "github.com/polshe-v/microservices_auth/internal/service/mocks"
	desc "github.com/polshe-v/microservices_auth/pkg/auth_v1"
)

func TestGetAccessToken(t *testing.T) {
	t.Parallel()

	type authServiceMockFunc func(mc *minimock.Controller) service.AuthService

	type args struct {
		ctx context.Context
		req *desc.GetAccessTokenRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		refreshToken = "refresh_token"
		accessToken  = "access_token"

		serviceErr = fmt.Errorf("service error")

		req = &desc.GetAccessTokenRequest{
			RefreshToken: refreshToken,
		}

		res = &desc.GetAccessTokenResponse{
			AccessToken: accessToken,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.GetAccessTokenResponse
		err             error
		authServiceMock authServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.GetAccessTokenMock.Expect(minimock.AnyContext, refreshToken).Return(accessToken, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.GetAccessTokenMock.Expect(minimock.AnyContext, refreshToken).Return("", serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authServiceMock := tt.authServiceMock(mc)
			api := authAPI.NewImplementation(authServiceMock)

			res, err := api.GetAccessToken(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
