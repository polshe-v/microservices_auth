package app

import (
	"context"
	"log"

	"github.com/polshe-v/microservices_auth/internal/api/access"
	"github.com/polshe-v/microservices_auth/internal/api/auth"
	"github.com/polshe-v/microservices_auth/internal/api/user"
	"github.com/polshe-v/microservices_auth/internal/config"
	"github.com/polshe-v/microservices_auth/internal/config/env"
	"github.com/polshe-v/microservices_auth/internal/repository"
	accessRepository "github.com/polshe-v/microservices_auth/internal/repository/access"
	keyRepository "github.com/polshe-v/microservices_auth/internal/repository/key"
	logRepository "github.com/polshe-v/microservices_auth/internal/repository/log"
	userRepository "github.com/polshe-v/microservices_auth/internal/repository/user"
	"github.com/polshe-v/microservices_auth/internal/service"
	accessService "github.com/polshe-v/microservices_auth/internal/service/access"
	authService "github.com/polshe-v/microservices_auth/internal/service/auth"
	userService "github.com/polshe-v/microservices_auth/internal/service/user"
	"github.com/polshe-v/microservices_common/pkg/closer"
	"github.com/polshe-v/microservices_common/pkg/db"
	"github.com/polshe-v/microservices_common/pkg/db/pg"
	"github.com/polshe-v/microservices_common/pkg/db/transaction"
)

type serviceProvider struct {
	pgConfig      config.PgConfig
	grpcConfig    config.GrpcConfig
	httpConfig    config.HTTPConfig
	swaggerConfig config.SwaggerConfig

	dbClient  db.Client
	txManager db.TxManager

	userRepository   repository.UserRepository
	keyRepository    repository.KeyRepository
	accessRepository repository.AccessRepository
	logRepository    repository.LogRepository
	userService      service.UserService
	authService      service.AuthService
	accessService    service.AccessService
	userImpl         *user.Implementation
	authImpl         *auth.Implementation
	accessImpl       *access.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PgConfig() config.PgConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPgConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %v", err)
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GrpcConfig() config.GrpcConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGrpcConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := env.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %v", err)
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := env.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %v", err)
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		c, err := pg.New(ctx, s.PgConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = c.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("failed to ping database: %v", err)
		}

		closer.Add(c.Close)

		s.dbClient = c
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}
	return s.txManager
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}
	return s.userRepository
}

func (s *serviceProvider) KeyRepository(ctx context.Context) repository.KeyRepository {
	if s.keyRepository == nil {
		s.keyRepository = keyRepository.NewRepository(s.DBClient(ctx))
	}
	return s.keyRepository
}

func (s *serviceProvider) AccessRepository(ctx context.Context) repository.AccessRepository {
	if s.accessRepository == nil {
		s.accessRepository = accessRepository.NewRepository(s.DBClient(ctx))
	}
	return s.accessRepository
}

func (s *serviceProvider) LogRepository(ctx context.Context) repository.LogRepository {
	if s.logRepository == nil {
		s.logRepository = logRepository.NewRepository(s.DBClient(ctx))
	}
	return s.logRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(s.UserRepository(ctx), s.LogRepository(ctx), s.TxManager(ctx))
	}
	return s.userService
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(s.UserRepository(ctx), s.KeyRepository(ctx))
	}
	return s.authService
}

func (s *serviceProvider) AccessService(ctx context.Context) service.AccessService {
	if s.accessService == nil {
		s.accessService = accessService.NewService(s.AccessRepository(ctx), s.KeyRepository(ctx))
	}
	return s.accessService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *user.Implementation {
	if s.userImpl == nil {
		s.userImpl = user.NewImplementation(s.UserService(ctx))
	}
	return s.userImpl
}

func (s *serviceProvider) AuthImpl(ctx context.Context) *auth.Implementation {
	if s.authImpl == nil {
		s.authImpl = auth.NewImplementation(s.AuthService(ctx))
	}
	return s.authImpl
}

func (s *serviceProvider) AccessImpl(ctx context.Context) *access.Implementation {
	if s.accessImpl == nil {
		s.accessImpl = access.NewImplementation(s.AccessService(ctx))
	}
	return s.accessImpl
}
