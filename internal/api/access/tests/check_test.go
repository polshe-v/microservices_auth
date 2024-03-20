package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/require"

	accessAPI "github.com/polshe-v/microservices_auth/internal/api/access"
	"github.com/polshe-v/microservices_auth/internal/service"
	serviceMocks "github.com/polshe-v/microservices_auth/internal/service/mocks"
	desc "github.com/polshe-v/microservices_auth/pkg/access_v1"
)

func TestCheck(t *testing.T) {
	t.Parallel()

	type accessServiceMockFunc func(mc *minimock.Controller) service.AccessService

	type args struct {
		ctx context.Context
		req *desc.CheckRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		endpoint = "/chat_v1.ChatV1/Create"

		serviceErr = fmt.Errorf("service error")

		req = &desc.CheckRequest{
			Endpoint: endpoint,
		}

		res = &empty.Empty{}
	)

	tests := []struct {
		name              string
		args              args
		want              *empty.Empty
		err               error
		accessServiceMock accessServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			accessServiceMock: func(mc *minimock.Controller) service.AccessService {
				mock := serviceMocks.NewAccessServiceMock(mc)
				mock.CheckMock.Expect(minimock.AnyContext, endpoint).Return(nil)
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
			accessServiceMock: func(mc *minimock.Controller) service.AccessService {
				mock := serviceMocks.NewAccessServiceMock(mc)
				mock.CheckMock.Expect(minimock.AnyContext, endpoint).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			accessServiceMock := tt.accessServiceMock(mc)
			api := accessAPI.NewImplementation(accessServiceMock)

			res, err := api.Check(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
