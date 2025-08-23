package reviews

import (
	"log/slog"
	"time"

	"github.com/renantatsuo/app-review/server/internal/store"
	"github.com/renantatsuo/app-review/server/pkg/apple"
)

// ReviewsClient is the client for the reviews.
type ReviewsClient struct {
	logger           *slog.Logger
	store            *store.Store
	apple            *apple.AppleClient
	pollingInterval  time.Duration // The interval for polling the reviews.
	reviewsTimeLimit time.Duration // The time limit to go back in time to fetch the reviews.
	storePath        string        // The directory for the store persistence.
	appleAppIDs      []string      // The list of apple apps.
}

func New(logger *slog.Logger, store *store.Store, apple *apple.AppleClient, pollingInterval time.Duration, reviewsTimeLimit time.Duration, storeDir string, appleAppIDs []string) *ReviewsClient {
	return &ReviewsClient{
		logger:           logger,
		store:            store,
		apple:            apple,
		pollingInterval:  pollingInterval,
		reviewsTimeLimit: reviewsTimeLimit,
		storePath:        storeDir,
		appleAppIDs:      appleAppIDs,
	}
}
