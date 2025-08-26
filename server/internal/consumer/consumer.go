package consumer

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/mattdeak/gopq"
	"github.com/renantatsuo/app-review/server/internal/config"
	"github.com/renantatsuo/app-review/server/internal/models"
	"github.com/renantatsuo/app-review/server/internal/queue"
	"github.com/renantatsuo/app-review/server/internal/reviews"
)

type Consumer struct {
	l             *slog.Logger
	queue         queue.Queue
	config        config.Config
	reviewsClient *reviews.ReviewsClient
}

func New(l *slog.Logger, queue queue.Queue, config config.Config, reviewsClient *reviews.ReviewsClient) *Consumer {
	return &Consumer{l: l, queue: queue, config: config, reviewsClient: reviewsClient}
}

func (c *Consumer) Start(ctx context.Context) {
	c.l.Info("starting consumer")
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				item, err := c.queue.Dequeue()
				if err != nil {
					if errors.Is(err, &gopq.ErrNoItemsWaiting{}) {
						c.l.Info("queue is empty, waiting for next tick")
						continue
					}
					c.l.Error("error dequeuing item", "error", err)
					continue
				}

				appID := string(item)
				var latestTime time.Time
				latestReview, err := c.reviewsClient.FindLatestReviewByAppID(appID)
				if err != nil && !errors.Is(err, sql.ErrNoRows) {
					c.l.Error("error finding latest review", "error", err)
					continue
				}

				if err != nil && errors.Is(err, sql.ErrNoRows) {
					c.l.Info("no latest review found, fetching all reviews")
					latestTime = time.Now().Add(-c.config.ReviewsTimeLimit)
				} else {
					latestTime = latestReview.SentAt
				}

				reviews, err := c.reviewsClient.GetLatestReviewsFromApple(appID, latestTime)
				if err != nil {
					c.l.Error("error getting latest reviews", "error", err)
					continue
				}

				if len(reviews) == 0 {
					c.l.Info("no new reviews found, skipping")
					continue
				}

				for _, review := range reviews {
					r, err := models.ReviewFromAppleReview(review, appID)
					if err != nil {
						c.l.Error("error converting apple review to model", "error", err)
						continue
					}
					if err := c.reviewsClient.AddReview(r); err != nil {
						c.l.Error("error adding review", "error", err)
						continue
					}
				}
			}
		}
	}()
}
