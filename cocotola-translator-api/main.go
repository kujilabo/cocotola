package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	libG "github.com/kujilabo/cocotola/lib/gateway"
)

const readHeaderTimeout = time.Duration(30) * time.Second
const serverPort = 8180
const metricsPort = 8081
const gracefulShutdownTime = 10

// @securityDefinitions.basic BasicAuth
func main() {
	result := run(context.Background())

	gracefulShutdownTime2 := time.Duration(gracefulShutdownTime) * time.Second
	time.Sleep(gracefulShutdownTime2)
	logrus.Info("exited")
	os.Exit(result)
}

func run(ctx context.Context) int {
	var eg *errgroup.Group
	eg, ctx = errgroup.WithContext(ctx)

	eg.Go(func() error {
		return httpServer(ctx)
	})
	eg.Go(func() error {
		return libG.MetricsServerProcess(ctx, metricsPort, gracefulShutdownTime)
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

func httpServer(ctx context.Context) error {
	router := gin.Default()

	httpServer := http.Server{
		Addr:              ":" + strconv.Itoa(serverPort),
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
		gracefulShutdownTime1 := time.Duration(gracefulShutdownTime) * time.Second
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
