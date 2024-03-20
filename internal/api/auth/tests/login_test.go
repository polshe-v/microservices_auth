package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	authAPI "github.com/polshe-v/microservices_auth/internal/api/auth"
	"github.com/polshe-v/microservices_auth/internal/model"
	"github.com/polshe-v/microservices_auth/internal/service"
	serviceMocks "github.com/polshe-v/microservices_auth/internal/service/mocks"
	desc "github.com/polshe-v/microservices_auth/pkg/auth_v1"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	type authServiceMockFunc func(mc *minimock.Controller) service.AuthService

	type args struct {
		ctx context.Context
		req *desc.LoginRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		username     = "username"
		password     = "password"
		refreshToken = "refresh_token"

		serviceErr = fmt.Errorf("service error")

		req = &desc.LoginRequest{
			Creds: &desc.Creds{
				Username: username,
				Password: password,
			},
		}

		creds = &model.UserCreds{
			Username: username,
			Password: password,
		}

		res = &desc.LoginResponse{
			RefreshToken: refreshToken,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.LoginResponse
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
				mock.LoginMock.Expect(minimock.AnyContext, creds).Return(refreshToken, nil)
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
				mock.LoginMock.Expect(minimock.AnyContext, creds).Return("", serviceErr)
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

			res, err := api.Login(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
