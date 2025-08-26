package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/renantatsuo/app-review/server/internal/apps"
	"github.com/renantatsuo/app-review/server/internal/config"
	"github.com/renantatsuo/app-review/server/internal/db"
	"github.com/renantatsuo/app-review/server/internal/queue"
	"github.com/renantatsuo/app-review/server/internal/scheduler"
	"github.com/renantatsuo/app-review/server/pkg/apple"
)

func main() {
	config, err := config.LoadConfigFromEnv()
	if err != nil {
		slog.Error("error loading config", "error", err)
		os.Exit(1)
	}

	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: config.LogLevel,
	})).With(slog.String("service", "scheduler"))

	db := db.Connect()
	appleClient := apple.New()
	appsClient := apps.New(db, appleClient)
	queue := queue.New()

	s := scheduler.New(l, appsClient, queue, config)

	ctx, cancel := context.WithCancel(context.Background())

	s.Start(ctx)

	kill := make(chan os.Signal, 1)
	signal.Notify(kill, syscall.SIGINT, syscall.SIGTERM)

	<-kill

	cancel()
	db.Close()

	l.Info("scheduler stopped")
}
