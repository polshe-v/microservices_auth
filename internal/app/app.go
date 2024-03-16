package app

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"github.com/polshe-v/microservices_auth/internal/config"
	"github.com/polshe-v/microservices_auth/internal/interceptor"
	descAccess "github.com/polshe-v/microservices_auth/pkg/access_v1"
	descAuth "github.com/polshe-v/microservices_auth/pkg/auth_v1"
	descUser "github.com/polshe-v/microservices_auth/pkg/user_v1"
	_ "github.com/polshe-v/microservices_auth/statik" // Not used for importing, only init() needed.
	"github.com/polshe-v/microservices_common/pkg/closer"
)

// App structure contains main application structures.
type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
	swaggerServer   *http.Server
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
			log.Fatalf("failed to run gRPC server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			log.Fatalf("failed to run HTTP server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runSwaggerServer()
		if err != nil {
			log.Fatalf("failed to run Swagger server: %v", err)
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGrpcServer,
		a.initHTTPServer,
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
		log.Fatalf("failed to load config: %v", err)
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
		grpc.UnaryInterceptor(interceptor.ValidateInterceptor),
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

func (a *App) runGrpcServer() error {
	// Open IP and port for server.
	lis, err := net.Listen(a.serviceProvider.GrpcConfig().Transport(), a.serviceProvider.GrpcConfig().Address())
	if err != nil {
		return err
	}

	log.Printf("gRPC server running on %v", a.serviceProvider.GrpcConfig().Address())

	err = a.grpcServer.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runHTTPServer() error {
	log.Printf("HTTP server running on %v", a.serviceProvider.HTTPConfig().Address())

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runSwaggerServer() error {
	log.Printf("Swagger server running on %v", a.serviceProvider.SwaggerConfig().Address())

	err := a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		log.Printf("Serving swagger file: %s", path)

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Open swagger file: %s", path)

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		closer.Add(file.Close)

		log.Printf("Read swagger file: %s", path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Write swagger file: %s", path)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Served swagger file: %s", path)
	}
}
