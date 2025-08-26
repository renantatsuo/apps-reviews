package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/renantatsuo/app-review/server/internal/apps"
	"github.com/renantatsuo/app-review/server/internal/config"
	"github.com/renantatsuo/app-review/server/internal/db"
	"github.com/renantatsuo/app-review/server/internal/queue"
	"github.com/renantatsuo/app-review/server/internal/reviews"
	"github.com/renantatsuo/app-review/server/internal/server"
	"github.com/renantatsuo/app-review/server/pkg/apple"
)

const (
	serverShutdownTimeout = 15 * time.Second
)

func main() {
	config, err := config.LoadConfigFromEnv()
	if err != nil {
		slog.Error("error loading config", "error", err)
		os.Exit(1)
	}

	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: config.LogLevel,
	})).With(slog.String("service", "server"))

	l.Info("initializing server", "port", config.Port, "logLevel", config.LogLevel)

	appleClient := apple.New()
	db := db.New(config.DatabaseConnStr).Connect()
	reviewsClient := reviews.New(l, appleClient, db, config)
	appsClient := apps.New(db, appleClient)
	queue := queue.New(config.QueueConnStr)

	s := server.New(config.Port, l, reviewsClient, appsClient, queue, config)

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

	db.Close()

	killCtx, killCancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer killCancel()

	if err := s.Stop(killCtx); err != nil {
		l.Error("error stopping server", "error", err)
		os.Exit(1)
	}
}
