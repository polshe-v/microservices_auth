package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math"
	"math/big"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/polshe-v/microservices_auth/pkg/user_v1"
)

const (
	grpcTransport = "tcp"
	grpcIP        = "0.0.0.0"
	grpcPort      = 50000
	delim         = "---"
)

type server struct {
	desc.UnimplementedUserV1Server
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("\n%s\nName: %s\nEmail: %s\nPassword: %s\nPassword confirm: %s\nRole: %v\n%s", delim, req.GetName(), req.GetEmail(), req.GetPassword(), req.GetPasswordConfirm(), req.GetRole(), delim)

	// Generate random ID.
	id, err := rand.Int(rand.Reader, big.NewInt(int64(math.MaxInt64)))
	if err != nil {
		log.Fatalf("failed to generate ID: %v", err)
	}

	return &desc.CreateResponse{
		Id: id.Int64(),
	}, nil
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("\n%s\nID: %d\n%s", delim, req.GetId(), delim)

	return &desc.GetResponse{
		Id:        req.GetId(),
		Name:      "<unimplemented>",
		Email:     "<unimplemented>",
		Role:      desc.Role_UNKNOWN,
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*empty.Empty, error) {
	log.Printf("\n%s\nID: %d\nName: %s\nEmail: %s\nRole: %v\n%s", delim, req.GetId(), req.GetName(), req.GetEmail(), req.GetRole(), delim)

	return &empty.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	log.Printf("\n%s\nID: %d\n%s", delim, req.GetId(), delim)

	return &empty.Empty{}, nil
}

func main() {
	// Open IP and port for server.
	lis, err := net.Listen(grpcTransport, fmt.Sprintf("%s:%d", grpcIP, grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create gRPC *Server which has no service registered and has not started to accept requests yet.
	s := grpc.NewServer()

	// Upon the client's request, the server will automatically provide information on the supported methods.
	reflection.Register(s)

	// Register service with corresponded interface.
	desc.RegisterUserV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
