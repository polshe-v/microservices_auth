package tests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"

	userAPI "github.com/polshe-v/microservices_auth/internal/api/user"
	"github.com/polshe-v/microservices_auth/internal/model"
	"github.com/polshe-v/microservices_auth/internal/service"
	serviceMocks "github.com/polshe-v/microservices_auth/internal/service/mocks"
	desc "github.com/polshe-v/microservices_auth/pkg/user_v1"
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.UpdateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = int64(1)
		name  = "name"
		email = "email"
		role  = desc.Role_USER

		serviceErr = fmt.Errorf("service error")

		req = &desc.UpdateRequest{
			User: &desc.UserUpdate{
				Id:    id,
				Name:  wrapperspb.String(name),
				Email: wrapperspb.String(email),
				Role:  role,
			},
		}

		userUpdate = &model.UserUpdate{
			ID: id,
			Name: sql.NullString{
				String: name,
				Valid:  true,
			},
			Email: sql.NullString{
				String: email,
				Valid:  true,
			},
			Role: desc.Role_name[int32(role)],
		}

		res = &empty.Empty{}
	)

	tests := []struct {
		name            string
		args            args
		want            *empty.Empty
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(minimock.AnyContext, userUpdate).Return(nil)
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
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(minimock.AnyContext, userUpdate).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.userServiceMock(mc)
			api := userAPI.NewImplementation(userServiceMock)

			res, err := api.Update(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
