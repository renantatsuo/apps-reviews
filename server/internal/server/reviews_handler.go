package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/renantatsuo/app-review/server/internal/models"
)

// getReviewsHandler is the handler for the /reviews/{appID} endpoint.
// It returns the reviews for the given appID.
// TODO: support pagination.
func (s *server) getReviewsHandler(w http.ResponseWriter, r *http.Request) {
	appID := r.PathValue("appID")

	since := time.Now().Add(-s.config.ReviewsTimeLimit)
	reviews, err := s.reviewsClient.FindReviewsByAppID(appID, since)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.logger.Error("error getting reviews", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if errors.Is(err, sql.ErrNoRows) {
		s.logger.Error("no reviews found", "appID", appID)
		http.Error(w, "no reviews found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseData[[]models.Review]{
		Data: reviews,
	})
}
