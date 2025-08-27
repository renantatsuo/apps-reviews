package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/renantatsuo/app-review/server/internal/apps"
	"github.com/renantatsuo/app-review/server/internal/config"
	"github.com/renantatsuo/app-review/server/internal/queue"
	"github.com/renantatsuo/app-review/server/internal/reviews"
)

type server struct {
	port          int
	logger        *slog.Logger
	server        *http.Server
	reviewsClient *reviews.ReviewsClient
	appsClient    *apps.AppsClient
	queue         queue.Queue
	config        config.Config
}

type ResponseData[T any] struct {
	Data T `json:"data"`
}

func New(port int, logger *slog.Logger, reviewsClient *reviews.ReviewsClient, appsClient *apps.AppsClient, queue queue.Queue, config config.Config) *server {
	return &server{
		port:          port,
		logger:        logger,
		reviewsClient: reviewsClient,
		appsClient:    appsClient,
		queue:         queue,
		config:        config,
	}
}

func (s *server) Start() error {
	router := http.NewServeMux()
	router.Handle("GET /reviews/{appID}", corsMiddleware(s.getReviewsHandler))
	router.Handle("GET /apps", corsMiddleware(s.getAppsHandler))
	router.Handle("POST /apps/{appID}", corsMiddleware(s.postAppsHandler))
	router.Handle("GET /apps/{appID}", corsMiddleware(s.getAppHandler))

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: router,
	}

	s.logger.Info("starting server", "port", s.port)
	return s.server.ListenAndServe()
}

func (s *server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
