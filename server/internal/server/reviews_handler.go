package server

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"sort"

	"github.com/renantatsuo/app-review/server/internal/models"
	"github.com/renantatsuo/app-review/server/internal/store"
	"github.com/renantatsuo/app-review/server/pkg/apple"
)

// reviewsHandler is the handler for the /reviews/{appID} endpoint.
// It returns the reviews for the given appID.
// TODO: support pagination.
// TODO: support live fetch if the appID is not found in the store.
func (s *server) reviewsHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Data []models.Review `json:"data"`
	}

	appID := r.PathValue("appID")

	reviews, err := s.reviewsClient.GetReviews(appID)
	if err != nil {
		if errors.Is(err, store.ErrKeyNotFound) {
			s.logger.Error("key not found", "appID", appID)
			http.Error(w, "key not found", http.StatusNotFound)
			return
		}
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

	// ensure the reviews are sorted by latest to oldest
	sort.Slice(reviewsModel, func(i, j int) bool {
		return reviewsModel[i].Updated.After(reviewsModel[j].Updated)
	})

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
