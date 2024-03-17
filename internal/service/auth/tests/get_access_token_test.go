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
		accessToken    = "access_token"
		/*		refreshToken   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTA2MjI2ODUsInVzZXJuYW1lIjoidGVzdDQiLCJyb2xlIjoiVVNFUiJ9.Blv3fYn_ZWzHRstw7tlTxXJQJGPjC_pPKLwda3Sgoso"
				accessToken    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTA2MTk0MTQsInVzZXJuYW1lIjoidGVzdDQiLCJyb2xlIjoiVVNFUiJ9.Ml33vpFMcvc7jdhjNsfMke_9EPlyLxAw7NBSutPx5z8"
				refreshKeyName = "refresh"
				refreshKey     = "u4CxvwDQMQDJnwAwu5xuPJVuv5azKLsv3m9FV5PSqqekZdAfjRqV2pcfjfXpCgPW"
				accessKeyName  = "access"
				accessKey      = "5LbCPS9xmAFWdhjPUqrfq2FQtpKPtUaNCjTX5brWPhyfZNRWedZowD2doC3QRd2E"
		*/
		repositoryErr = fmt.Errorf("failed to generate token")

		req = refreshToken

		res = accessToken
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
				return mock
			},
			keyRepositoryMock: func(mc *minimock.Controller) repository.KeyRepository {
				mock := repositoryMocks.NewKeyRepositoryMock(mc)
				mock.GetKeyMock.Expect(minimock.AnyContext, refreshKeyName).Return(refreshKey, nil)
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
				return mock
			},
			keyRepositoryMock: func(mc *minimock.Controller) repository.KeyRepository {
				mock := repositoryMocks.NewKeyRepositoryMock(mc)
				mock.GetKeyMock.Expect(minimock.AnyContext, refreshKeyName).Return("", repositoryErr)
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
