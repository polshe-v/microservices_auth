package app

import (
	"context"
	"flag"
	"io"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/natefinch/lumberjack"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/polshe-v/microservices_auth/internal/config"
	"github.com/polshe-v/microservices_auth/internal/interceptor"
	"github.com/polshe-v/microservices_auth/internal/metrics"
	"github.com/polshe-v/microservices_auth/internal/tracing"
	descAccess "github.com/polshe-v/microservices_auth/pkg/access_v1"
	descAuth "github.com/polshe-v/microservices_auth/pkg/auth_v1"
	descUser "github.com/polshe-v/microservices_auth/pkg/user_v1"
	_ "github.com/polshe-v/microservices_auth/statik" // Not used for importing, only init() needed.
	"github.com/polshe-v/microservices_common/pkg/closer"
	"github.com/polshe-v/microservices_common/pkg/logger"
)

// App structure contains main application structures.
type App struct {
	serviceProvider  *serviceProvider
	grpcServer       *grpc.Server
	httpServer       *http.Server
	swaggerServer    *http.Server
	prometheusServer *http.Server
}

var configPath string

func init() {
	flag.StringVar(&configPath, "config", ".env", "Path to config file")
}

// NewApp creates new App object.
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}
	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Run executes the application.
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}
	wg.Add(3) // gRPC, HTTP and Swagger servers

	go func() {
		defer wg.Done()

		err := a.runGrpcServer()
		if err != nil {
			logger.Fatal("failed to run gRPC server: ", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			logger.Fatal("failed to run HTTP server: ", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runSwaggerServer()
		if err != nil {
			logger.Fatal("failed to run Swagger server: ", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runPrometheusServer()
		if err != nil {
			logger.Fatal("failed to run Prometheus server: ", zap.Error(err))
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initLogger,
		a.initTracing,
		a.initGrpcServer,
		a.initHTTPServer,
		a.initPrometheusServer,
		a.initSwaggerServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initConfig(_ context.Context) error {
	// Parse the command-line flags from os.Args[1:].
	flag.Parse()

	err := config.Load(configPath)
	if err != nil {
		logger.Fatal("failed to load config: ", zap.Error(err))
	}
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGrpcServer(ctx context.Context) error {
	cfg := a.serviceProvider.GrpcConfig()
	creds, err := credentials.NewServerTLSFromFile(cfg.CertPath(), cfg.KeyPath())
	if err != nil {
		return err
	}

	a.grpcServer = grpc.NewServer(
		grpc.Creds(creds),
		grpc.ChainUnaryInterceptor(
			interceptor.LogInterceptor,
			interceptor.MetricsInterceptor,
			interceptor.ValidateInterceptor,
			interceptor.TracingInterceptor,
		),
	)

	// Upon the client's request, the server will automatically provide information on the supported methods.
	reflection.Register(a.grpcServer)

	// Register service with corresponded interface.
	descUser.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserImpl(ctx))
	descAuth.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.AuthImpl(ctx))
	descAccess.RegisterAccessV1Server(a.grpcServer, a.serviceProvider.AccessImpl(ctx))

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	cfg := a.serviceProvider.GrpcConfig()
	creds, err := credentials.NewClientTLSFromFile(cfg.CertPath(), "")
	if err != nil {
		return err
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	mux := runtime.NewServeMux()
	err = descUser.RegisterUserV1HandlerFromEndpoint(ctx, mux, cfg.Address(), opts)
	if err != nil {
		return err
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		Addr:              a.serviceProvider.HTTPConfig().Address(),
		Handler:           corsMiddleware.Handler(mux),
		ReadHeaderTimeout: 15 * time.Second,
	}

	return nil
}

func (a *App) initSwaggerServer(_ context.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc("/api.swagger.json", serveSwaggerFile("/api.swagger.json"))

	a.swaggerServer = &http.Server{
		Addr:              a.serviceProvider.SwaggerConfig().Address(),
		Handler:           mux,
		ReadHeaderTimeout: 15 * time.Second,
	}

	return nil
}

func (a *App) initPrometheusServer(ctx context.Context) error {
	err := metrics.Init(ctx)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	a.prometheusServer = &http.Server{
		Addr:              a.serviceProvider.PrometheusConfig().Address(),
		Handler:           mux,
		ReadHeaderTimeout: 15 * time.Second,
	}

	return nil
}

func (a *App) initLogger(_ context.Context) error {
	cfg := a.serviceProvider.LogConfig()

	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.LogFilePath(),
		MaxSize:    cfg.LogMaxSize(),
		MaxBackups: cfg.LogMaxFiles(),
		MaxAge:     cfg.LogMaxAge(),
	})

	level, err := zapcore.ParseLevel(cfg.LogLevel())
	if err != nil {
		return err
	}

	cfgConsoleLog := zap.NewProductionEncoderConfig()
	cfgConsoleLog.TimeKey = "timestamp"
	cfgConsoleLog.EncodeTime = zapcore.ISO8601TimeEncoder

	cfgFileLog := zap.NewDevelopmentEncoderConfig()
	cfgFileLog.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(cfgConsoleLog)
	fileEncoder := zapcore.NewJSONEncoder(cfgFileLog)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)

	logger.Init(core)
	return nil
}

func (a *App) initTracing(ctx context.Context) error {
	cfg := a.serviceProvider.TracingConfig()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceName(cfg.ServiceName()),
		),
	)
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(cfg.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return err
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// Set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	closer.Add(func() error {
		return tracerProvider.Shutdown(ctx)
	})

	err = tracing.InitGlobalTracer(cfg.ServiceName())
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runGrpcServer() error {
	// Open IP and port for server.
	lis, err := net.Listen(a.serviceProvider.GrpcConfig().Transport(), a.serviceProvider.GrpcConfig().Address())
	if err != nil {
		return err
	}

	logger.Info("gRPC server running on ", zap.String("address", a.serviceProvider.GrpcConfig().Address()))

	err = a.grpcServer.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runHTTPServer() error {
	logger.Info("HTTP server running on ", zap.String("address", a.serviceProvider.HTTPConfig().Address()))

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runPrometheusServer() error {
	logger.Info("Prometheus server running on ", zap.String("address", a.serviceProvider.PrometheusConfig().Address()))

	err := a.prometheusServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runSwaggerServer() error {
	logger.Info("Swagger server running on ", zap.String("address", a.serviceProvider.SwaggerConfig().Address()))

	err := a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		logger.Info("Serving swagger file: ", zap.String("path", path))

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info("Open swagger file: ", zap.String("path", path))

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		closer.Add(file.Close)

		logger.Info("Read swagger file: ", zap.String("path", path))

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info("Write swagger file: ", zap.String("path", path))

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info("Served swagger file: ", zap.String("path", path))
	}
}
