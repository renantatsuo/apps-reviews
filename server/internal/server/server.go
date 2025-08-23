package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/renantatsuo/app-review/server/internal/reviews"
)

type server struct {
	port          int
	logger        *slog.Logger
	server        *http.Server
	reviewsClient *reviews.ReviewsClient
}

func New(port int, logger *slog.Logger, reviewsClient *reviews.ReviewsClient) *server {
	return &server{
		port:          port,
		logger:        logger,
		reviewsClient: reviewsClient,
	}
}

func (s *server) Start() error {
	router := http.NewServeMux()
	router.Handle("GET /reviews/{appID}", corsMiddleware(s.reviewsHandler))

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
