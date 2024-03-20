package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"

	"github.com/polshe-v/microservices_auth/internal/model"
	"github.com/polshe-v/microservices_auth/internal/repository"
	repositoryMocks "github.com/polshe-v/microservices_auth/internal/repository/mocks"
	accessService "github.com/polshe-v/microservices_auth/internal/service/access"
	"github.com/polshe-v/microservices_auth/internal/tokens"
	tokenMocks "github.com/polshe-v/microservices_auth/internal/tokens/mocks"
)

func TestCheck(t *testing.T) {
	t.Parallel()

	type keyRepositoryMockFunc func(mc *minimock.Controller) repository.KeyRepository
	type accessRepositoryMockFunc func(mc *minimock.Controller) repository.AccessRepository
	type tokenOperationsMockFunc func(mc *minimock.Controller) tokens.TokenOperations

	type args struct {
		ctx context.Context
		req string
	}

	var (
		mdNoAuthHeader = metadata.New(map[string]string{"header": "access_token"})
		mdNoAuthPrefix = metadata.New(map[string]string{"Authorization": "access_token"})
		md             = metadata.New(map[string]string{"Authorization": "Bearer access_token"})

		ctxNoMd         = context.Background()
		ctx             = metadata.NewIncomingContext(ctxNoMd, md)
		ctxNoAuthHeader = metadata.NewIncomingContext(ctxNoMd, mdNoAuthHeader)
		ctxNoAuthPrefix = metadata.NewIncomingContext(ctxNoMd, mdNoAuthPrefix)

		mc = minimock.NewController(t)

		endpointCreate      = "/chat_v1.ChatV1/Create"
		endpointDelete      = "/chat_v1.ChatV1/Delete"
		endpointSendMessage = "/chat_v1.ChatV1/SendMessage"
		endpointNotExists   = "/chat_v1.ChatV1/NotExists"

		username  = "username"
		roleUser  = "USER"
		roleAdmin = "ADMIN"

		accessKeyName  = "access"
		accessKey      = "access_key"
		accessKeyBytes = []byte("access_key")

		accessToken = "access_token"

		endpointPermissions = []*model.EndpointPermissions{
			{
				Endpoint: endpointCreate,
				Roles:    []string{roleAdmin},
			},
			{
				Endpoint: endpointDelete,
				Roles:    []string{roleAdmin},
			},
			{
				Endpoint: endpointSendMessage,
				Roles:    []string{roleAdmin, roleUser},
			},
		}

		claimsAdmin = &model.UserClaims{
			Username: username,
			Role:     roleAdmin,
		}

		claimsUser = &model.UserClaims{
			Username: username,
			Role:     roleUser,
		}

		noMdErr          = fmt.Errorf("metadata is not provided")
		noAuthHeaderErr  = fmt.Errorf("authorization header is not provided")
		noAuthPrefixErr  = fmt.Errorf("invalid authorization header format")
		keyRepositoryErr = fmt.Errorf("failed to generate token")
		noEndpointErr    = fmt.Errorf("failed to find endpoint")
		tokenInvalidErr  = fmt.Errorf("access token is invalid")
		accessDeniedErr  = fmt.Errorf("access denied")

		req = endpointCreate
	)

	tests := []struct {
		name                 string
		args                 args
		err                  error
		keyRepositoryMock    keyRepositoryMockFunc
		accessRepositoryMock accessRepositoryMockFunc
		tokenOperationsMock  tokenOperationsMockFunc
	}{
		{
			name: "metadata not provided error case",
			args: args{
				ctx: ctxNoMd,
				req: req,
			},
			err: noMdErr,
			keyRepositoryMock: func(mc *minimock.Controller) repository.KeyRepository {
				mock := repositoryMocks.NewKeyRepositoryMock(mc)
				return mock
			},
			accessRepositoryMock: func(mc *minimock.Controller) repository.AccessRepository {
				mock := repositoryMocks.NewAccessRepositoryMock(mc)
				return mock
			},
			tokenOperationsMock: func(mc *minimock.Controller) tokens.TokenOperations {
				mock := tokenMocks.NewTokenOperationsMock(mc)
				return mock
			},
		},
		{
			name: "authorization header not provided error case",
			args: args{
				ctx: ctxNoAuthHeader,
				req: req,
			},
			err: noAuthHeaderErr,
			keyRepositoryMock: func(mc *minimock.Controller) repository.KeyRepository {
				mock := repositoryMocks.NewKeyRepositoryMock(mc)
				return mock
			},
			accessRepositoryMock: func(mc *minimock.Controller) repository.AccessRepository {
				mock := repositoryMocks.NewAccessRepositoryMock(mc)
				return mock
			},
			tokenOperationsMock: func(mc *minimock.Controller) tokens.TokenOperations {
				mock := tokenMocks.NewTokenOperationsMock(mc)
				return mock
			},
		},
		{
			name: "authorization header format error case",
			args: args{
				ctx: ctxNoAuthPrefix,
				req: req,
			},
			err: noAuthPrefixErr,
			keyRepositoryMock: func(mc *minimock.Controller) repository.KeyRepository {
				mock := repositoryMocks.NewKeyRepositoryMock(mc)
				return mock
			},
			accessRepositoryMock: func(mc *minimock.Controller) repository.AccessRepository {
				mock := repositoryMocks.NewAccessRepositoryMock(mc)
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
			err: keyRepositoryErr,
			keyRepositoryMock: func(mc *minimock.Controller) repository.KeyRepository {
				mock := repositoryMocks.NewKeyRepositoryMock(mc)
				mock.GetKeyMock.Expect(minimock.AnyContext, accessKeyName).Return("", keyRepositoryErr)
				return mock
			},
			accessRepositoryMock: func(mc *minimock.Controller) repository.AccessRepository {
				mock := repositoryMocks.NewAccessRepositoryMock(mc)
				return mock
			},
			tokenOperationsMock: func(mc *minimock.Controller) tokens.TokenOperations {
				mock := tokenMocks.NewTokenOperationsMock(mc)
				return mock
			},
		},
		{
			name: "endpoint not found error case",
			args: args{
				ctx: ctx,
				req: endpointNotExists,
			},
			err: noEndpointErr,
			keyRepositoryMock: func(mc *minimock.Controller) repository.KeyRepository {
				mock := repositoryMocks.NewKeyRepositoryMock(mc)
				mock.GetKeyMock.Expect(minimock.AnyContext, accessKeyName).Return(accessKey, nil)
				return mock
			},
			accessRepositoryMock: func(mc *minimock.Controller) repository.AccessRepository {
				mock := repositoryMocks.NewAccessRepositoryMock(mc)
				mock.GetRoleEndpointsMock.Expect(minimock.AnyContext).Return(endpointPermissions, nil)
				return mock
			},
			tokenOperationsMock: func(mc *minimock.Controller) tokens.TokenOperations {
				mock := tokenMocks.NewTokenOperationsMock(mc)
				return mock
			},
		},
		{
			name: "token verify error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: tokenInvalidErr,
			keyRepositoryMock: func(mc *minimock.Controller) repository.KeyRepository {
				mock := repositoryMocks.NewKeyRepositoryMock(mc)
				mock.GetKeyMock.Expect(minimock.AnyContext, accessKeyName).Return(accessKey, nil)
				return mock
			},
			accessRepositoryMock: func(mc *minimock.Controller) repository.AccessRepository {
				mock := repositoryMocks.NewAccessRepositoryMock(mc)
				return mock
			},
			tokenOperationsMock: func(mc *minimock.Controller) tokens.TokenOperations {
				mock := tokenMocks.NewTokenOperationsMock(mc)
				mock.VerifyMock.Expect(accessToken, accessKeyBytes).Return(nil, tokenInvalidErr)
				return mock
			},
		},
		{
			name: "access denied error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: accessDeniedErr,
			keyRepositoryMock: func(mc *minimock.Controller) repository.KeyRepository {
				mock := repositoryMocks.NewKeyRepositoryMock(mc)
				mock.GetKeyMock.Expect(minimock.AnyContext, accessKeyName).Return(accessKey, nil)
				return mock
			},
			accessRepositoryMock: func(mc *minimock.Controller) repository.AccessRepository {
				mock := repositoryMocks.NewAccessRepositoryMock(mc)
				return mock
			},
			tokenOperationsMock: func(mc *minimock.Controller) tokens.TokenOperations {
				mock := tokenMocks.NewTokenOperationsMock(mc)
				mock.VerifyMock.Expect(accessToken, accessKeyBytes).Return(claimsUser, nil)
				return mock
			},
		},
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: nil,
			keyRepositoryMock: func(mc *minimock.Controller) repository.KeyRepository {
				mock := repositoryMocks.NewKeyRepositoryMock(mc)
				mock.GetKeyMock.Expect(minimock.AnyContext, accessKeyName).Return(accessKey, nil)
				return mock
			},
			accessRepositoryMock: func(mc *minimock.Controller) repository.AccessRepository {
				mock := repositoryMocks.NewAccessRepositoryMock(mc)
				return mock
			},
			tokenOperationsMock: func(mc *minimock.Controller) tokens.TokenOperations {
				mock := tokenMocks.NewTokenOperationsMock(mc)
				mock.VerifyMock.Expect(accessToken, accessKeyBytes).Return(claimsAdmin, nil)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			keyRepositoryMock := tt.keyRepositoryMock(mc)
			accessRepositoryMock := tt.accessRepositoryMock(mc)
			tokenOperationsMock := tt.tokenOperationsMock(mc)
			srv := accessService.NewService(accessRepositoryMock, keyRepositoryMock, tokenOperationsMock)

			err := srv.Check(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
		})
	}
}
