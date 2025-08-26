package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/renantatsuo/app-review/server/internal/config"
	"github.com/renantatsuo/app-review/server/internal/reviews"
	"github.com/renantatsuo/app-review/server/internal/server"
	"github.com/renantatsuo/app-review/server/internal/store"
	"github.com/renantatsuo/app-review/server/pkg/apple"
)

const (
	serverShutdownTimeout = 15 * time.Second
)

func main() {
	config, err := config.LoadConfigFromEnv()
	if err != nil {
		log.Fatalf("error loading config: %w", err)
	}

	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: config.LogLevel,
	}))

	l.Info("initializing server", "port", config.Port, "logLevel", config.LogLevel)

	store := store.New()
	appleClient := apple.New()
	reviewsClient := reviews.New(l, store, appleClient, config.PollingInterval, config.ReviewsTimeLimit, config.StoreDir, config.AppIDs)
	reviewsClient.Load()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		reviewsClient.StartPolling(ctx)
	}()

	s := server.New(config.Port, l, reviewsClient)

	go func() {
		if err := s.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Error("error starting server", "error", err)
			os.Exit(1)
		}
	}()

	kill := make(chan os.Signal, 1)
	signal.Notify(kill, syscall.SIGINT, syscall.SIGTERM)

	<-kill

	l.Info("received shutdown signal, gracefully shutting down server")

	cancel()

	killCtx, killCancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer killCancel()

	if err := s.Stop(killCtx); err != nil {
		l.Error("error stopping server", "error", err)
		os.Exit(1)
	}
}
