package reviews

import (
	"time"

	"github.com/renantatsuo/app-review/server/internal/models"
)

// FindReviewsByAppID returns all the reviews for a given app ID.
func (r *ReviewsClient) FindReviewsByAppID(appID string, since time.Time) ([]models.Review, error) {
	reviews := []models.Review{}

	rows, err := r.db.Query("SELECT id, app_id, author, title, content, rating, sent_at, created_at, updated_at FROM reviews WHERE app_id = ? AND sent_at > ? ORDER BY sent_at DESC", appID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var review models.Review
		err := rows.Scan(&review.ID, &review.AppID, &review.Author, &review.Title,
			&review.Content, &review.Rating, &review.SentAt,
			&review.CreatedAt, &review.UpdatedAt)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}

// FindLatestReviewByAppID finds the latest review for a given app ID.
func (r *ReviewsClient) FindLatestReviewByAppID(appID string) (models.Review, error) {
	review := models.Review{}

	row := r.db.QueryRow("SELECT id, app_id, author, title, content, rating, sent_at, created_at, updated_at FROM reviews WHERE app_id = ? ORDER BY created_at DESC LIMIT 1", appID)
	if err := row.Scan(&review.ID, &review.AppID, &review.Author, &review.Title,
		&review.Content, &review.Rating, &review.SentAt,
		&review.CreatedAt, &review.UpdatedAt); err != nil {
		return models.Review{}, err
	}

	return review, nil
}

// AddReview adds a new review to the database.
func (r *ReviewsClient) AddReview(review models.Review) error {
	_, err := r.db.Exec(
		"INSERT INTO reviews (id, app_id, author, title, content, rating, sent_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		review.ID, review.AppID, review.Author, review.Title, review.Content, review.Rating, review.SentAt)
	if err != nil {
		return err
	}
	return nil
}
