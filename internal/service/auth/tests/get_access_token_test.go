package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/polshe-v/microservices_auth/internal/repository"
	repositoryMocks "github.com/polshe-v/microservices_auth/internal/repository/mocks"
	authService "github.com/polshe-v/microservices_auth/internal/service/auth"
)

func TestGetAccessToken(t *testing.T) {
	t.Parallel()

	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type keyRepositoryMockFunc func(mc *minimock.Controller) repository.KeyRepository

	type args struct {
		ctx context.Context
		req string
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		refreshKeyName = "refresh"
		refreshKey     = "refresh_key"
		refreshToken   = "refresh_token"
		accessKeyName  = "access"

		repositoryErr = fmt.Errorf("failed to generate token")

		req = refreshToken
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
			name: "refresh key repository error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: "",
			err:  repositoryErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				return mock
			},
			keyRepositoryMock: func(mc *minimock.Controller) repository.KeyRepository {
				mock := repositoryMocks.NewKeyRepositoryMock(mc)
				mock.GetKeyMock.Expect(minimock.AnyContext, refreshKeyName).Return("", repositoryErr)
				return mock
			},
		},
		{
			name: "access key repository error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: "",
			err:  repositoryErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				return mock
			},
			keyRepositoryMock: func(mc *minimock.Controller) repository.KeyRepository {
				mock := repositoryMocks.NewKeyRepositoryMock(mc)
				mock.GetKeyMock.When(minimock.AnyContext, refreshKeyName).Then(refreshKey, nil)
				mock.GetKeyMock.When(minimock.AnyContext, accessKeyName).Then("", repositoryErr)
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

			res, err := srv.GetAccessToken(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
