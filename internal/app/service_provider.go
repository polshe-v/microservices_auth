package app

import (
	"context"

	redigo "github.com/gomodule/redigo/redis"
	"go.uber.org/zap"

	"github.com/polshe-v/microservices_auth/internal/api/access"
	"github.com/polshe-v/microservices_auth/internal/api/auth"
	"github.com/polshe-v/microservices_auth/internal/api/user"
	"github.com/polshe-v/microservices_auth/internal/config"
	"github.com/polshe-v/microservices_auth/internal/config/env"
	"github.com/polshe-v/microservices_auth/internal/repository"
	accessRepository "github.com/polshe-v/microservices_auth/internal/repository/access"
	cacheRepository "github.com/polshe-v/microservices_auth/internal/repository/cache"
	keyRepository "github.com/polshe-v/microservices_auth/internal/repository/key"
	logRepository "github.com/polshe-v/microservices_auth/internal/repository/log"
	userRepository "github.com/polshe-v/microservices_auth/internal/repository/user"
	"github.com/polshe-v/microservices_auth/internal/service"
	accessService "github.com/polshe-v/microservices_auth/internal/service/access"
	authService "github.com/polshe-v/microservices_auth/internal/service/auth"
	userService "github.com/polshe-v/microservices_auth/internal/service/user"
	"github.com/polshe-v/microservices_auth/internal/tokens"
	"github.com/polshe-v/microservices_auth/internal/tokens/jwt"
	"github.com/polshe-v/microservices_common/pkg/cache"
	"github.com/polshe-v/microservices_common/pkg/cache/redis"
	"github.com/polshe-v/microservices_common/pkg/closer"
	"github.com/polshe-v/microservices_common/pkg/db"
	"github.com/polshe-v/microservices_common/pkg/db/pg"
	"github.com/polshe-v/microservices_common/pkg/db/transaction"
	"github.com/polshe-v/microservices_common/pkg/logger"
)

type serviceProvider struct {
	pgConfig         config.PgConfig
	grpcConfig       config.GrpcConfig
	httpConfig       config.HTTPConfig
	swaggerConfig    config.SwaggerConfig
	prometheusConfig config.PrometheusConfig
	logConfig        config.LogConfig
	tracingConfig    config.TracingConfig
	redisConfig      config.RedisConfig

	dbClient  db.Client
	txManager db.TxManager

	redisPool   *redigo.Pool
	redisClient cache.Client

	userRepository   repository.UserRepository
	keyRepository    repository.KeyRepository
	accessRepository repository.AccessRepository
	logRepository    repository.LogRepository
	cacheRepository  repository.CacheRepository

	userService   service.UserService
	authService   service.AuthService
	accessService service.AccessService

	userImpl   *user.Implementation
	authImpl   *auth.Implementation
	accessImpl *access.Implementation

	tokenOperations tokens.TokenOperations
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PgConfig() config.PgConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPgConfig()
		if err != nil {
			logger.Fatal("failed to get pg config: ", zap.Error(err))
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GrpcConfig() config.GrpcConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGrpcConfig()
		if err != nil {
			logger.Fatal("failed to get grpc config: ", zap.Error(err))
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := env.NewHTTPConfig()
		if err != nil {
			logger.Fatal("failed to get http config: ", zap.Error(err))
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := env.NewSwaggerConfig()
		if err != nil {
			logger.Fatal("failed to get swagger config: ", zap.Error(err))
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

func (s *serviceProvider) PrometheusConfig() config.PrometheusConfig {
	if s.prometheusConfig == nil {
		cfg, err := env.NewPrometheusConfig()
		if err != nil {
			logger.Fatal("failed to get prometheus config: ", zap.Error(err))
		}

		s.prometheusConfig = cfg
	}

	return s.prometheusConfig
}

func (s *serviceProvider) LogConfig() config.LogConfig {
	if s.logConfig == nil {
		cfg, err := env.NewLogConfig()
		if err != nil {
			logger.Fatal("failed to get log config: ", zap.Error(err))
		}

		s.logConfig = cfg
	}

	return s.logConfig
}

func (s *serviceProvider) TracingConfig() config.TracingConfig {
	if s.tracingConfig == nil {
		cfg, err := env.NewTracingConfig()
		if err != nil {
			logger.Fatal("failed to get tracing config: ", zap.Error(err))
		}

		s.tracingConfig = cfg
	}

	return s.tracingConfig
}

func (s *serviceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := env.NewRedisConfig()
		if err != nil {
			logger.Fatal("failed to get redis config: ", zap.Error(err))
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		c, err := pg.New(ctx, s.PgConfig().DSN())
		if err != nil {
			logger.Fatal("failed to create db client: ", zap.Error(err))
		}

		err = c.DB().Ping(ctx)
		if err != nil {
			logger.Fatal("failed to ping database: ", zap.Error(err))
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

func (s *serviceProvider) RedisClient() cache.Client {
	if s.redisClient == nil {
		s.redisClient = redis.NewClient(s.RedisPool(), s.RedisConfig().ConnectionTimeout())
	}

	return s.redisClient
}

func (s *serviceProvider) RedisPool() *redigo.Pool {
	if s.redisPool == nil {
		s.redisPool = &redigo.Pool{
			MaxIdle:     s.RedisConfig().MaxIdle(),
			IdleTimeout: s.RedisConfig().IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", s.RedisConfig().Address())
			},
		}
	}

	return s.redisPool
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

func (s *serviceProvider) CacheRepository(ctx context.Context) repository.CacheRepository {
	if s.cacheRepository == nil {
		s.cacheRepository = cacheRepository.NewRepository(s.RedisClient())
	}
	return s.cacheRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(s.UserRepository(ctx), s.CacheRepository(ctx), s.LogRepository(ctx), s.TxManager(ctx))
	}
	return s.userService
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(s.UserRepository(ctx), s.KeyRepository(ctx), s.TokenOperations(ctx))
	}
	return s.authService
}

func (s *serviceProvider) AccessService(ctx context.Context) service.AccessService {
	if s.accessService == nil {
		s.accessService = accessService.NewService(s.AccessRepository(ctx), s.KeyRepository(ctx), s.TokenOperations(ctx))
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

func (s *serviceProvider) TokenOperations(_ context.Context) tokens.TokenOperations {
	if s.tokenOperations == nil {
		s.tokenOperations = jwt.NewTokenOperations()
	}
	return s.tokenOperations
}
