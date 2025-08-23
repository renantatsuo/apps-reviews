package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/renantatsuo/app-review/server/internal/models"
	"github.com/renantatsuo/app-review/server/pkg/apple"
)

func (s *server) reviewsHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Data []models.Review `json:"data"`
	}

	appID := r.PathValue("appID")

	reviews, err := s.reviewsClient.GetReviews(appID)
	if err != nil {
		s.logger.Error("error getting reviews", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(reviews) == 0 {
		s.logger.Error("no reviews found", "appID", appID)
		http.Error(w, "no reviews found", http.StatusNotFound)
		return
	}

	reviewsModel := transformReviewsModel(reviews, s.logger)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response{
		Data: reviewsModel,
	})
}

// transformReviewsModel transforms the reviews from the apple api to the models.Review.
func transformReviewsModel(res []apple.Review, l *slog.Logger) []models.Review {
	reviews := make([]models.Review, len(res))
	for i, r := range res {
		review, err := models.ReviewFromAppleReview(r)
		if err != nil {
			l.Error("error transforming review", "error", err)
			continue
		}
		reviews[i] = review
	}
	return reviews
}
