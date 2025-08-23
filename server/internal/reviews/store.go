package reviews

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/renantatsuo/app-review/server/internal/store"
	"github.com/renantatsuo/app-review/server/pkg/apple"
)

const (
	appIdsFile     = "app-ids.json"
	reviewsFileFmt = "reviews-%s.json"
	ReviewsPrefix  = "app-id:"
	AppIDsKey      = "app-ids"
)

// Load loads the app ids and reviews from the file system.
func (c *ReviewsClient) Load() {
	if _, err := os.Stat(c.storeDir); os.IsNotExist(err) {
		os.MkdirAll(c.storeDir, 0755)
		return
	}

	store.Set(AppIDsKey, c.appleAppIDs, c.store)

	for _, appID := range c.appleAppIDs {
		reviewsBytes, err := os.ReadFile(path.Join(c.storeDir, fmt.Sprintf(reviewsFileFmt, appID)))
		if err != nil {
			c.logger.Info("no reviews file found, creating one")
			continue
		}

		var reviews []apple.Review
		if err := json.Unmarshal(reviewsBytes, &reviews); err != nil {
			c.logger.Error("error unmarshalling reviews", "error", err)
			continue
		}

		store.Set(ReviewsPrefix+appID, reviews, c.store)
	}
}

// SaveReviews sets the reviews for a given app ID.
// It also writes the reviews to the file system.
func (c *ReviewsClient) SaveReviews(appID string, reviews []apple.Review) error {
	store.Set(ReviewsPrefix+appID, reviews, c.store)

	reviewsBytes, err := json.Marshal(reviews)
	if err != nil {
		return err
	}

	path := path.Join(c.storeDir, fmt.Sprintf(reviewsFileFmt, appID))
	os.WriteFile(path, reviewsBytes, 0644)

	return nil
}

// GetReviews gets the reviews for a given app ID.
func (c *ReviewsClient) GetReviews(appID string) ([]apple.Review, error) {
	reviews, err := store.Get[[]apple.Review](ReviewsPrefix+appID, c.store)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}
