package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	"github.com/kujilabo/cocotola/cocotola-tatoeba-api/docs"
	"github.com/kujilabo/cocotola/cocotola-tatoeba-api/src/config"
	"github.com/kujilabo/cocotola/cocotola-tatoeba-api/src/controller"
	"github.com/kujilabo/cocotola/cocotola-tatoeba-api/src/gateway"
	"github.com/kujilabo/cocotola/cocotola-tatoeba-api/src/service"
	"github.com/kujilabo/cocotola/cocotola-tatoeba-api/src/sqls"
	"github.com/kujilabo/cocotola/cocotola-tatoeba-api/src/usecase"
	libconfig "github.com/kujilabo/cocotola/lib/config"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
	libG "github.com/kujilabo/cocotola/lib/gateway"
)

const readHeaderTimeout = time.Duration(30) * time.Second

// @securityDefinitions.basic BasicAuth
func main() {
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

	rfFunc := func(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error) {
		return gateway.NewRepositoryFactory(ctx, db, cfg.DB.DriverName)
	}

	result := run(context.Background(), cfg, db, rfFunc)

	gracefulShutdownTime2 := time.Duration(cfg.Shutdown.TimeSec2) * time.Second
	time.Sleep(gracefulShutdownTime2)
	logrus.Info("exited")
	os.Exit(result)
}

func run(ctx context.Context, cfg *config.Config, db *gorm.DB, rfFunc service.RepositoryFactoryFunc) int {
	var eg *errgroup.Group
	eg, ctx = errgroup.WithContext(ctx)

	eg.Go(func() error {
		return httpServer(ctx, cfg, db, rfFunc)
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

func httpServer(ctx context.Context, cfg *config.Config, db *gorm.DB, rfFunc service.RepositoryFactoryFunc) error {
	// cors
	corsConfig := libconfig.InitCORS(cfg.CORS)
	logrus.Infof("cors: %+v", corsConfig)

	if err := corsConfig.Validate(); err != nil {
		return liberrors.Errorf("corsConfig.Validate. err: %w", err)
	}

	if !cfg.Debug.GinMode {
		gin.SetMode(gin.ReleaseMode)
	}

	adminUsecase := usecase.NewAdminUsecase(db, rfFunc)
	userUsecase := usecase.NewUserUsecase(db, rfFunc)

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

func initialize(ctx context.Context, env string) (*config.Config, *gorm.DB, *sql.DB, *sdktrace.TracerProvider, error) {
	cfg, err := config.LoadConfig(env)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// init log
	if err := libconfig.InitLog(env, cfg.Log); err != nil {
		return nil, nil, nil, nil, err
	}

	// init tracer
	tp, err := libconfig.InitTracerProvider(cfg.App.Name, cfg.Trace)
	if err != nil {
		return nil, nil, nil, nil, liberrors.Errorf("failed to InitTracerProvider. err: %w", err)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// init db
	db, sqlDB, err := libconfig.InitDB(cfg.DB, sqls.SQL)
	if err != nil {
		return nil, nil, nil, nil, liberrors.Errorf("failed to InitDB. err: %w", err)
	}

	return cfg, db, sqlDB, tp, nil
}
