package reviews

import (
	"database/sql"
	"log/slog"

	"github.com/renantatsuo/app-review/server/internal/config"
	"github.com/renantatsuo/app-review/server/pkg/apple"
)

// ReviewsClient is the client for the reviews.
type ReviewsClient struct {
	logger *slog.Logger
	apple  *apple.AppleClient
	db     *sql.DB
	config config.Config
}

func New(logger *slog.Logger, apple *apple.AppleClient, db *sql.DB, config config.Config) *ReviewsClient {
	return &ReviewsClient{
		logger: logger,
		apple:  apple,
		db:     db,
		config: config,
	}
}
