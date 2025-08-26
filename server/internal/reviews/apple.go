package reviews

import (
	"fmt"
	"time"

	"github.com/renantatsuo/app-review/server/pkg/apple"
)

// GetLatestReviewsFromApple fetches the latest reviews for a given app ID.
// It returns a slice of reviews that were updated after the since time.
func (c *ReviewsClient) GetLatestReviewsFromApple(appID string, since time.Time) (res []apple.Review, err error) {
	reviews, err := c.apple.GetLatestReviews(appID)
	if err != nil {
		return nil, err
	}

	if reviews.Feed.Entry == nil {
		return nil, fmt.Errorf("no reviews found")
	}

	latest := time.Now()

	for latest.After(since) {
		for _, review := range reviews.Feed.Entry {
			updated, err := time.Parse(apple.AppleTimeFormat, review.Updated.Label)
			if err != nil {
				return nil, err
			}

			latest = updated

			if updated.After(since) {
				res = append(res, review)
			} else {
				return res, nil
			}
		}

		if reviews.HasNext() {
			reviews, err = reviews.Next()
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return res, nil
}
