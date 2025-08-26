package models

import (
	"strconv"
	"time"

	"github.com/renantatsuo/app-review/server/pkg/apple"
)

type Review struct {
	ID        string    `json:"id"`
	AppID     string    `json:"app_id"`
	Author    string    `json:"author"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Rating    int       `json:"rating"`
	SentAt    time.Time `json:"sent_at"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// ReviewFromAppleReview transforms the apple review to the models.Review.
func ReviewFromAppleReview(review apple.Review, appID string) (Review, error) {
	rating, err := strconv.Atoi(review.Rating.Label)
	if err != nil {
		return Review{}, err
	}

	updated, err := time.Parse("2006-01-02T15:04:05-07:00", review.Updated.Label)
	if err != nil {
		return Review{}, err
	}

	res := Review{
		ID:      review.ID.Label,
		AppID:   appID,
		Author:  review.Author.Name.Label,
		Title:   review.Title.Label,
		Content: review.Content.Label,
		Rating:  rating,
		SentAt:  updated,
	}

	return res, nil
}
