package models

import (
	"strconv"
	"time"

	"github.com/renantatsuo/app-review/server/pkg/apple"
)

type Review struct {
	ID      string       `json:"id"`
	Author  ReviewAuthor `json:"author"`
	Title   string       `json:"title"`
	Content string       `json:"content"`
	Rating  int          `json:"rating"`
	Updated time.Time    `json:"updated"`
}

type ReviewAuthor struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
}

// ReviewFromAppleReview transforms the apple review to the models.Review.
func ReviewFromAppleReview(review apple.Review) (Review, error) {
	rating, err := strconv.Atoi(review.Rating.Label)
	if err != nil {
		return Review{}, err
	}

	updated, err := time.Parse("2006-01-02T15:04:05-07:00", review.Updated.Label)
	if err != nil {
		return Review{}, err
	}

	res := Review{
		ID: review.ID.Label,
		Author: ReviewAuthor{
			Name: review.Author.Name.Label,
			URI:  review.Author.Uri.Label,
		},
		Title:   review.Title.Label,
		Content: review.Content.Label,
		Rating:  rating,
		Updated: updated,
	}

	return res, nil
}
