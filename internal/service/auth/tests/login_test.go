package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/polshe-v/microservices_auth/internal/model"
	"github.com/polshe-v/microservices_auth/internal/repository"
	repositoryMocks "github.com/polshe-v/microservices_auth/internal/repository/mocks"
	authService "github.com/polshe-v/microservices_auth/internal/service/auth"
	"github.com/polshe-v/microservices_auth/internal/tokens"
	tokenMocks "github.com/polshe-v/microservices_auth/internal/tokens/mocks"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type keyRepositoryMockFunc func(mc *minimock.Controller) repository.KeyRepository
	type tokenOperationsMockFunc func(mc *minimock.Controller) tokens.TokenOperations

	type args struct {
		ctx context.Context
		req *model.UserCreds
	}

	var password = "password"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Print("failed to process password")
		return
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		username               = "username"
		passwordWrong          = "passwordWrong"
		role                   = "USER"
		refreshKeyName         = "refresh"
		refreshKey             = "refresh_key"
		refreshKeyBytes        = []byte("refresh_key")
		refreshTokenExpiration = 60 * time.Minute
		refreshToken           = "refresh_token"

		keyRepositoryErr  = fmt.Errorf("failed to generate token")
		userRepositoryErr = fmt.Errorf("user not found")
		wrongPasswordErr  = fmt.Errorf("wrong password")

		authInfo = &model.AuthInfo{
			Username: username,
			Password: string(hashedPassword),
			Role:     role,
		}

		user = model.User{
			Name: username,
			Role: role,
		}

		req = &model.UserCreds{
			Username: username,
			Password: password,
		}

		reqWrongPass = &model.UserCreds{
			Username: username,
			Password: passwordWrong,
		}

		res = refreshToken
	)

	tests := []struct {
		name                string
		args                args
		want                string
		err                 error
		userRepositoryMock  userRepositoryMockFunc
		keyRepositoryMock   keyRepositoryMockFunc
		tokenOperationsMock tokenOperationsMockFunc
	}{
		{
			name: "user repository error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: "",
			err:  userRepositoryErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.GetAuthInfoMock.Expect(minimock.AnyContext, username).Return(nil, userRepositoryErr)
				return mock
			},
			keyRepositoryMock: func(mc *minimock.Controller) repository.KeyRepository {
				mock := repositoryMocks.NewKeyRepositoryMock(mc)
				return mock
			},
			tokenOperationsMock: func(mc *minimock.Controller) tokens.TokenOperations {
				mock := tokenMocks.NewTokenOperationsMock(mc)
				return mock
			},
		},
		{
			name: "wrong password error case",
			args: args{
				ctx: ctx,
				req: reqWrongPass,
			},
			want: "",
			err:  wrongPasswordErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.GetAuthInfoMock.Expect(minimock.AnyContext, username).Return(authInfo, nil)
				return mock
			},
			keyRepositoryMock: func(mc *minimock.Controller) repository.KeyRepository {
				mock := repositoryMocks.NewKeyRepositoryMock(mc)
				return mock
			},
			tokenOperationsMock: func(mc *minimock.Controller) tokens.TokenOperations {
				mock := tokenMocks.NewTokenOperationsMock(mc)
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
			err:  keyRepositoryErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.GetAuthInfoMock.Expect(minimock.AnyContext, username).Return(authInfo, nil)
				return mock
			},
			keyRepositoryMock: func(mc *minimock.Controller) repository.KeyRepository {
				mock := repositoryMocks.NewKeyRepositoryMock(mc)
				mock.GetKeyMock.Expect(minimock.AnyContext, refreshKeyName).Return("", keyRepositoryErr)
				return mock
			},
			tokenOperationsMock: func(mc *minimock.Controller) tokens.TokenOperations {
				mock := tokenMocks.NewTokenOperationsMock(mc)
				return mock
			},
		},
		{
			name: "token generate error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: "",
			err:  keyRepositoryErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.GetAuthInfoMock.Expect(minimock.AnyContext, username).Return(authInfo, nil)
				return mock
			},
			keyRepositoryMock: func(mc *minimock.Controller) repository.KeyRepository {
				mock := repositoryMocks.NewKeyRepositoryMock(mc)
				mock.GetKeyMock.Expect(minimock.AnyContext, refreshKeyName).Return(refreshKey, nil)
				return mock
			},
			tokenOperationsMock: func(mc *minimock.Controller) tokens.TokenOperations {
				mock := tokenMocks.NewTokenOperationsMock(mc)
				mock.GenerateMock.Expect(user, refreshKeyBytes, refreshTokenExpiration).Return("", keyRepositoryErr)
				return mock
			},
		},
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
				mock.GetKeyMock.Expect(minimock.AnyContext, refreshKeyName).Return(refreshKey, nil)
				return mock
			},
			tokenOperationsMock: func(mc *minimock.Controller) tokens.TokenOperations {
				mock := tokenMocks.NewTokenOperationsMock(mc)
				mock.GenerateMock.Expect(user, refreshKeyBytes, refreshTokenExpiration).Return(refreshToken, nil)
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
			tokenOperationsMock := tt.tokenOperationsMock(mc)
			srv := authService.NewService(userRepositoryMock, keyRepositoryMock, tokenOperationsMock)

			res, err := srv.Login(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
