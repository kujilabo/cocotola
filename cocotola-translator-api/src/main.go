package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"flag"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	"github.com/kujilabo/cocotola/cocotola-translator-api/docs"
	"github.com/kujilabo/cocotola/cocotola-translator-api/src/config"
	"github.com/kujilabo/cocotola/cocotola-translator-api/src/controller"
	"github.com/kujilabo/cocotola/cocotola-translator-api/src/gateway"
	"github.com/kujilabo/cocotola/cocotola-translator-api/src/sqls"
	"github.com/kujilabo/cocotola/cocotola-translator-api/src/usecase"
	libconfig "github.com/kujilabo/cocotola/lib/config"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	libG "github.com/kujilabo/cocotola/lib/gateway"
	pb "github.com/kujilabo/cocotola/proto/src/proto"
)

const readHeaderTimeout = time.Duration(30) * time.Second

// @securityDefinitions.basic BasicAuth
func main() {
	logrus.Infof("Starting cocotola-translator-api")

	ctx := context.Background()
	env := flag.String("env", "", "environment")
	flag.Parse()
	if len(*env) == 0 {
		appEnv := os.Getenv("APP_ENV")
		if len(appEnv) == 0 {
			*env = "local"
		} else {
			*env = appEnv
		}
	}

	logrus.Infof("env: %s", *env)

	liberrors.UseXerrorsErrorf()

	cfg, db, sqlDB, tp, err := initialize(ctx, *env)
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()
	defer tp.ForceFlush(ctx) // flushes any pending spans

	azureTranslationClient := gateway.NewAzureTranslationClient(cfg.Azure.SubscriptionKey)
	rf, err := gateway.NewRepositoryFactory(ctx, db, cfg.DB.DriverName)
	if err != nil {
		panic(err)
	}

	adminUsecase := usecase.NewAdminUsecase(rf)
	userUsecase := usecase.NewUserUsecase(rf, azureTranslationClient)

	result := run(context.Background(), cfg, db, adminUsecase, userUsecase)

	gracefulShutdownTime2 := time.Duration(cfg.Shutdown.TimeSec2) * time.Second
	time.Sleep(gracefulShutdownTime2)
	logrus.Info("exited")
	os.Exit(result)
}

func run(ctx context.Context, cfg *config.Config, db *gorm.DB, adminUsecase usecase.AdminUsecase, userUsecase usecase.UserUsecase) int {
	var eg *errgroup.Group
	eg, ctx = errgroup.WithContext(ctx)

	eg.Go(func() error {
		return httpServer(ctx, cfg, db, adminUsecase, userUsecase)
	})
	eg.Go(func() error {
		return grpcServer(ctx, cfg, db, adminUsecase, userUsecase)
	})
	eg.Go(func() error {
		return libG.MetricsServerProcess(ctx, cfg.App.MetricsPort, cfg.Shutdown.TimeSec1)
	})
	eg.Go(func() error {
		return libG.SignalWatchProcess(ctx)
	})
	eg.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})

	if err := eg.Wait(); err != nil {
		logrus.Error(err)
		return 1
	}
	return 0
}

func httpServer(ctx context.Context, cfg *config.Config, db *gorm.DB, adminUsecase usecase.AdminUsecase, userUsecase usecase.UserUsecase) error {
	// cors
	corsConfig := libconfig.InitCORS(cfg.CORS)
	logrus.Infof("cors: %+v", corsConfig)

	if err := corsConfig.Validate(); err != nil {
		return err
	}

	router := controller.NewRouter(adminUsecase, userUsecase, corsConfig, cfg.App, cfg.Auth, cfg.Debug)

	if cfg.Swagger.Enabled {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		docs.SwaggerInfo.Title = cfg.App.Name
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Host = cfg.Swagger.Host
		docs.SwaggerInfo.Schemes = []string{cfg.Swagger.Schema}
	}

	httpServer := http.Server{
		Addr:              ":" + strconv.Itoa(cfg.App.HTTPPort),
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	logrus.Printf("http server listening at %v", httpServer.Addr)

	errCh := make(chan error)
	go func() {
		defer close(errCh)
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logrus.Infof("failed to ListenAndServe. err: %v", err)
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		gracefulShutdownTime1 := time.Duration(cfg.Shutdown.TimeSec1) * time.Second
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), gracefulShutdownTime1)
		defer shutdownCancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logrus.Infof("Server forced to shutdown. err: %v", err)
			return err
		}
		return nil
	case err := <-errCh:
		return err
	}
}

func grpcServer(ctx context.Context, cfg *config.Config, db *gorm.DB, adminUsecase usecase.AdminUsecase, userUsecase usecase.UserUsecase) error {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(cfg.App.GRPCPort))
	if err != nil {
		logrus.Fatalf("failed to Listen: %v", err)
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			otelgrpc.UnaryServerInterceptor(),
			grpc_auth.UnaryServerInterceptor(NewAuthFunc(cfg.Auth.Username, cfg.Auth.Password)),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_recovery.UnaryServerInterceptor(),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			otelgrpc.StreamServerInterceptor(),
			grpc_auth.StreamServerInterceptor(NewAuthFunc(cfg.Auth.Username, cfg.Auth.Password)),
			grpc_prometheus.StreamServerInterceptor,
			grpc_recovery.StreamServerInterceptor(),
		)),
	)
	reflection.Register(grpcServer)

	userServer := controller.NewTranslatorUserServer(userUsecase)
	pb.RegisterTranslatorUserServer(grpcServer, userServer)

	logrus.Printf("grpc server listening at %v", lis.Addr())

	errCh := make(chan error)
	go func() {
		defer close(errCh)
		if err := grpcServer.Serve(lis); err != nil {
			logrus.Fatalf("failed to Serve: %v", err)
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		grpcServer.GracefulStop()
		return nil
	case err := <-errCh:
		return err
	}
}

func initialize(ctx context.Context, env string) (*config.Config, *gorm.DB, *sql.DB, *sdktrace.TracerProvider, error) {
	cfg, err := config.LoadConfig(env)
	if err != nil {
		return nil, nil, nil, nil, liberrors.Errorf("failed to config.LoadConfig in main.initialize. err: %w", err)
	}

	// init log
	if err := libconfig.InitLog(env, cfg.Log); err != nil {
		return nil, nil, nil, nil, err
	}

	// init tracer
	tp, err := libconfig.InitTracerProvider(cfg.App.Name, cfg.Trace)
	if err != nil {
		return nil, nil, nil, nil, liberrors.Errorf("failed to config.InitTracerProvider in main.initialize. err: %w", err)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// init db
	db, sqlDB, err := libconfig.InitDB(cfg.DB, sqls.SQL)
	if err != nil {
		return nil, nil, nil, nil, liberrors.Errorf("failed to configInitDB in main.initialize. err: %w", err)
	}

	return cfg, db, sqlDB, tp, nil
}

func NewAuthFunc(username, password string) func(ctx context.Context) (context.Context, error) {
	data := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	return func(ctx context.Context) (context.Context, error) {
		basic, err := grpc_auth.AuthFromMD(ctx, "basic")
		if err != nil {
			return nil, err
		}

		if basic != data {
			return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
		}

		return ctx, nil
	}
}
