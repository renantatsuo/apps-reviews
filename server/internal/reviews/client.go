package reviews

import (
	"log/slog"
	"time"

	"github.com/renantatsuo/app-review/server/internal/store"
	"github.com/renantatsuo/app-review/server/pkg/apple"
)

type ReviewsClient struct {
	logger           *slog.Logger
	store            *store.Store
	apple            *apple.AppleClient
	pollingInterval  time.Duration
	reviewsTimeLimit time.Duration
	storeDir         string
	appleAppIDs      []string
}

func New(logger *slog.Logger, store *store.Store, apple *apple.AppleClient, pollingInterval time.Duration, reviewsTimeLimit time.Duration, storeDir string, appleAppIDs []string) *ReviewsClient {
	return &ReviewsClient{
		logger:           logger,
		store:            store,
		apple:            apple,
		pollingInterval:  pollingInterval,
		reviewsTimeLimit: reviewsTimeLimit,
		storeDir:         storeDir,
		appleAppIDs:      appleAppIDs,
	}
}
