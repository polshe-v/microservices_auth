package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/polshe-v/microservices_auth/internal/model"
	"github.com/polshe-v/microservices_auth/internal/repository"
	repositoryMocks "github.com/polshe-v/microservices_auth/internal/repository/mocks"
	authService "github.com/polshe-v/microservices_auth/internal/service/auth"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type keyRepositoryMockFunc func(mc *minimock.Controller) repository.KeyRepository

	type args struct {
		ctx context.Context
		req *model.UserCreds
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		username     = "username"
		password     = "password"
		role         = "USER"
		refreshToken = "refresh_token"
		keyName      = "key_name"
		key          = "key"

		repositoryErr = fmt.Errorf("repository error")

		req = &model.UserCreds{
			Username: username,
			Password: password,
		}

		res = refreshToken

		authInfo = &model.AuthInfo{
			Username: username,
			Password: password,
			Role:     role,
		}
	)

	tests := []struct {
		name               string
		args               args
		want               string
		err                error
		userRepositoryMock userRepositoryMockFunc
		keyRepositoryMock  keyRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.GetAuthInfoMock.Expect(minimock.AnyContext, username).Return(authInfo, nil)
				return mock
			},
			keyRepositoryMock: func(mc *minimock.Controller) repository.KeyRepository {
				mock := repositoryMocks.NewKeyRepositoryMock(mc)
				mock.GetKeyMock.Expect(minimock.AnyContext, keyName).Return(key, nil)
				return mock
			},
		},
		{
			name: "user repository error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: "",
			err:  repositoryErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.GetAuthInfoMock.Expect(minimock.AnyContext, username).Return(nil, repositoryErr)
				return mock
			},
			keyRepositoryMock: func(mc *minimock.Controller) repository.KeyRepository {
				mock := repositoryMocks.NewKeyRepositoryMock(mc)
				return mock
			},
		},
		{
			name: "key repository error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: "",
			err:  repositoryErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.GetAuthInfoMock.Expect(minimock.AnyContext, username).Return(authInfo, nil)
				return mock
			},
			keyRepositoryMock: func(mc *minimock.Controller) repository.KeyRepository {
				mock := repositoryMocks.NewKeyRepositoryMock(mc)
				mock.GetKeyMock.Expect(minimock.AnyContext, keyName).Return("", repositoryErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepositoryMock := tt.userRepositoryMock(mc)
			keyRepositoryMock := tt.keyRepositoryMock(mc)
			srv := authService.NewService(userRepositoryMock, keyRepositoryMock)

			res, err := srv.Login(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
