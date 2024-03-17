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
	accessService "github.com/polshe-v/microservices_auth/internal/service/access"
)

func TestCheck(t *testing.T) {
	t.Parallel()

	type keyRepositoryMockFunc func(mc *minimock.Controller) repository.KeyRepository
	type accessRepositoryMockFunc func(mc *minimock.Controller) repository.AccessRepository

	type args struct {
		ctx context.Context
		req string
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		endpointCreate      = "/chat_v1.ChatV1/Create"
		endpointDelete      = "/chat_v1.ChatV1/Delete"
		endpointSendMessage = "/chat_v1.ChatV1/SendMessage"
		roleUser            = "USER"
		roleAdmin           = "ADMIN"
		keyName             = "key_name"
		key                 = "key"

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

		repositoryErr = fmt.Errorf("repository error")

		req = endpointCreate
	)

	tests := []struct {
		name                 string
		args                 args
		err                  error
		keyRepositoryMock    keyRepositoryMockFunc
		accessRepositoryMock accessRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: nil,
			keyRepositoryMock: func(mc *minimock.Controller) repository.KeyRepository {
				mock := repositoryMocks.NewKeyRepositoryMock(mc)
				mock.GetKeyMock.Expect(minimock.AnyContext, keyName).Return(key, nil)
				return mock
			},
			accessRepositoryMock: func(mc *minimock.Controller) repository.AccessRepository {
				mock := repositoryMocks.NewAccessRepositoryMock(mc)
				mock.GetRoleEndpointsMock.Expect(minimock.AnyContext).Return(endpointPermissions, nil)
				return mock
			},
		},
		{
			name: "key repository error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: repositoryErr,
			keyRepositoryMock: func(mc *minimock.Controller) repository.KeyRepository {
				mock := repositoryMocks.NewKeyRepositoryMock(mc)
				mock.GetKeyMock.Expect(minimock.AnyContext, keyName).Return("", repositoryErr)
				return mock
			},
			accessRepositoryMock: func(mc *minimock.Controller) repository.AccessRepository {
				mock := repositoryMocks.NewAccessRepositoryMock(mc)
				return mock
			},
		},
		{
			name: "access repository error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: repositoryErr,
			keyRepositoryMock: func(mc *minimock.Controller) repository.KeyRepository {
				mock := repositoryMocks.NewKeyRepositoryMock(mc)
				mock.GetKeyMock.Expect(minimock.AnyContext, keyName).Return(key, nil)
				return mock
			},
			accessRepositoryMock: func(mc *minimock.Controller) repository.AccessRepository {
				mock := repositoryMocks.NewAccessRepositoryMock(mc)
				mock.GetRoleEndpointsMock.Expect(minimock.AnyContext).Return(nil, repositoryErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			keyRepositoryMock := tt.keyRepositoryMock(mc)
			accessRepositoryMock := tt.accessRepositoryMock(mc)
			srv := accessService.NewService(accessRepositoryMock, keyRepositoryMock)

			err := srv.Check(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
		})
	}
}
