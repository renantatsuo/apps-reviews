package reviews

import (
	"context"
	"time"
)

// StartPolling starts the polling process.
func (c *ReviewsClient) StartPolling(ctx context.Context) {
	c.logger.Info("starting polling", "pollingInterval", c.pollingInterval)

	c.initialize()

	ticker := time.Tick(0)
	go func() {
		for {
			select {
			case <-ctx.Done():
				c.logger.Info("stopping polling")
				return
			case <-ticker:
				c.fetchReviews()
			}
		}
	}()
}

// initialize initializes the reviews client.
func (c *ReviewsClient) initialize() {
	c.fetchReviews()
}

// fetchReviews polls the reviews for all the apple app ids.
func (c *ReviewsClient) fetchReviews() {
	for _, appID := range c.appleAppIDs {
		go func(appID string) {
			c.logger.Info("fetching reviews for app", "appID", appID)

			since := time.Now().Add(-c.reviewsTimeLimit)
			reviews, err := c.getLatestReviewsFromApple(appID, since)
			if err != nil {
				c.logger.Error("error getting latest reviews", "error", err)
				return
			}

			if err := c.SaveReviews(appID, reviews); err != nil {
				c.logger.Error("error saving reviews", "error", err)
				return
			}
		}(appID)
	}
}
