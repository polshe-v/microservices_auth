package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	config "github.com/polshe-v/microservices_auth/internal/config"
	env "github.com/polshe-v/microservices_auth/internal/config/env"
	desc "github.com/polshe-v/microservices_auth/pkg/user_v1"
)

var configPath string
var errQueryBuild = errors.New("failed to build query")

func init() {
	flag.StringVar(&configPath, "config", ".env", "Path to config file")
}

const (
	bcryptCost = 12
	delim      = "---"
)

type server struct {
	desc.UnimplementedUserV1Server
	pool *pgxpool.Pool
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	user := req.GetUser()
	log.Printf("\n%s\nName: %s\nEmail: %s\nPassword: %s\nPassword confirm: %s\nRole: %v\n%s", delim, user.GetName(), user.GetEmail(), user.GetPassword(), user.GetPasswordConfirm(), user.GetRole(), delim)

	// Hashing the password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.GetPassword()), bcryptCost)
	if err != nil {
		log.Printf("%v", err)
		return nil, errors.New("failed to process password")
	}

	builderInsert := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "role", "email", "password").
		Values(user.GetName(), user.GetRole(), user.GetEmail(), hashedPassword).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("%v", err)
		return nil, errQueryBuild
	}

	var id int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		log.Printf("%v", err)
		return nil, errors.New("failed to create user")
	}

	return &desc.CreateResponse{
		Id: id,
	}, nil
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("\n%s\nID: %d\n%s", delim, req.GetId(), delim)

	builderSelect := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()}).
		Limit(1)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		log.Printf("%v", err)
		return nil, errQueryBuild
	}

	var id int64
	var name, email, role string
	var createdAt time.Time
	var updatedAt sql.NullTime

	err = s.pool.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("%v", err)
			return nil, errors.New("no user with given id")
		}
		log.Printf("%v", err)
		return nil, errors.New("failed to read user info")
	}

	var updatedAtTime *timestamppb.Timestamp
	if updatedAt.Valid {
		updatedAtTime = timestamppb.New(updatedAt.Time)
	}

	return &desc.GetResponse{
		User: &desc.User{
			Id:        id,
			Name:      name,
			Email:     email,
			Role:      desc.Role(desc.Role_value[role]),
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: updatedAtTime,
		},
	}, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*empty.Empty, error) {
	user := req.GetUser()
	log.Printf("\n%s\nID: %d\nName: %s\nEmail: %s\nRole: %v\n%s", delim, user.GetId(), user.GetName().GetValue(), user.GetEmail().GetValue(), user.GetRole(), delim)

	builderUpdate := sq.Update("users").
		SetMap(map[string]interface{}{
			"name":       user.GetName().GetValue(),
			"email":      user.GetEmail().GetValue(),
			"role":       user.GetRole(),
			"updated_at": sq.Expr("NOW()"),
		}).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": user.GetId()})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Printf("%v", err)
		return nil, errQueryBuild
	}

	res, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("%v", err)
		return nil, errors.New("failed to update user info")
	}
	log.Printf("result: %v", res)

	return &empty.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	log.Printf("\n%s\nID: %d\n%s", delim, req.GetId(), delim)

	builderDelete := sq.Delete("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		log.Printf("%v", err)
		return nil, errQueryBuild
	}

	res, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("%v", err)
		return nil, errors.New("failed to delete user")
	}
	log.Printf("result: %v", res)

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

	// Create gRPC *Server which has no service registered and has not started to accept requests yet.
	s := grpc.NewServer()

	// Upon the client's request, the server will automatically provide information on the supported methods.
	reflection.Register(s)

	// Register service with corresponded interface.
	desc.RegisterUserV1Server(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
