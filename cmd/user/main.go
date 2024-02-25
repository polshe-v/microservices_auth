package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	config "github.com/polshe-v/microservices_auth/internal/config"
	env "github.com/polshe-v/microservices_auth/internal/config/env"
	"github.com/polshe-v/microservices_auth/internal/repository"
	"github.com/polshe-v/microservices_auth/internal/repository/user"
	desc "github.com/polshe-v/microservices_auth/pkg/user_v1"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", ".env", "Path to config file")
}

const (
	bcryptCost = 12
	delim      = "---"
)

type server struct {
	desc.UnimplementedUserV1Server
	userRepository repository.UserRepository
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	user := req.GetUser()
	log.Printf("\n%s\nName: %s\nEmail: %s\nPassword: %s\nPassword confirm: %s\nRole: %v\n%s", delim, user.GetName(), user.GetEmail(), user.GetPassword(), user.GetPasswordConfirm(), user.GetRole(), delim)

	id, err := s.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	log.Printf("Created user with id: %d", id)

	return &desc.CreateResponse{
		Id: id,
	}, nil
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("\n%s\nID: %d\n%s", delim, req.GetId(), delim)

	user, err := s.userRepository.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &desc.GetResponse{
		User: user,
	}, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*empty.Empty, error) {
	user := req.GetUser()
	log.Printf("\n%s\nID: %d\nName: %s\nEmail: %s\nRole: %v\n%s", delim, user.GetId(), user.GetName().GetValue(), user.GetEmail().GetValue(), user.GetRole(), delim)

	err := s.userRepository.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	log.Printf("\n%s\nID: %d\n%s", delim, req.GetId(), delim)

	err := s.userRepository.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func main() {
	// Parse the command-line flags from os.Args[1:].
	flag.Parse()
	ctx := context.Background()

	// Read config file.
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGrpcConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	// Open IP and port for server.
	lis, err := net.Listen(grpcConfig.Transport(), grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pgConfig, err := env.NewPgConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	// Create database connections pool.
	pool, err := pgxpool.New(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Create repository layer.
	userRepo := user.NewRepository(pool)

	// Create gRPC *Server which has no service registered and has not started to accept requests yet.
	s := grpc.NewServer()

	// Upon the client's request, the server will automatically provide information on the supported methods.
	reflection.Register(s)

	// Register service with corresponded interface.
	desc.RegisterUserV1Server(s, &server{userRepository: userRepo})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
