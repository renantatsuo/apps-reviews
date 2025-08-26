package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/renantatsuo/app-review/server/internal/config"
	"github.com/renantatsuo/app-review/server/internal/consumer"
	"github.com/renantatsuo/app-review/server/internal/db"
	"github.com/renantatsuo/app-review/server/internal/queue"
	"github.com/renantatsuo/app-review/server/internal/reviews"
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
	})).With(slog.String("service", "consumer"))

	ctx, cancel := context.WithCancel(context.Background())

	queue := queue.New()
	appleClient := apple.New()
	db := db.Connect()
	reviewsClient := reviews.New(l, appleClient, db, config)
	consumer := consumer.New(l, queue, config, reviewsClient)
	consumer.Start(ctx)

	kill := make(chan os.Signal, 1)
	signal.Notify(kill, syscall.SIGINT, syscall.SIGTERM)

	<-kill

	cancel()
	queue.Close()

	l.Info("consumer stopped")
}
