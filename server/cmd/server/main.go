package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/renantatsuo/app-review/server/internal/reviews"
	"github.com/renantatsuo/app-review/server/internal/server"
	"github.com/renantatsuo/app-review/server/internal/store"
	"github.com/renantatsuo/app-review/server/pkg/apple"
	"github.com/renantatsuo/envv"
)

const (
	serverShutdownTimeout = 15 * time.Second
)

func main() {
	port := envv.Get("PORT").Int().Default(8080).Parse()
	logLevelStr := envv.Get("LOG_LEVEL").String().Default("debug").Parse()
	// defaults to hevy and instagram
	appleAppIDsStr := envv.Get("APP_IDS").String().Default("1458862350,389801252").Parse()
	appleAppIDs := strings.Split(appleAppIDsStr, ",")
	storeDir := envv.Get("STORE_DIR").String().Default("data").Parse()
	reviewsTimeLimit := envv.Get("REVIEWS_TIME_LIMIT").Duration().Default(48 * time.Hour).Parse()
	pollingInterval := envv.Get("POLLING_INTERVAL").Duration().Default(30 * time.Second).Parse()

	logLevel, err := parseLogLevel(logLevelStr)
	if err != nil {
		panic(err)
	}

	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))

	l.Info("initializing server", "port", port, "logLevel", logLevel)

	store := store.New()
	appleClient := apple.New()
	reviewsClient := reviews.New(l, store, appleClient, pollingInterval, reviewsTimeLimit, storeDir, appleAppIDs)
	reviewsClient.Load()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		reviewsClient.StartPolling(ctx)
	}()

	s := server.New(port, l, reviewsClient)

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

func parseLogLevel(logLevel string) (l slog.Level, err error) {
	err = l.UnmarshalText([]byte(logLevel))
	return
}
