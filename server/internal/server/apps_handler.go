package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/renantatsuo/app-review/server/internal/models"
)

// getAppsHandler is the handler for the /apps endpoint.
func (s *server) getAppsHandler(w http.ResponseWriter, r *http.Request) {
	apps, err := s.appsClient.GetAllApps()
	if err != nil {
		s.logger.Error("error getting apps", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseData[[]models.App]{
		Data: apps,
	})
}

// getAppHandler is the handler for the /apps/:appID endpoint.
func (s *server) getAppHandler(w http.ResponseWriter, r *http.Request) {
	appID := r.PathValue("appID")
	app, err := s.appsClient.GetAppData(appID)
	if err != nil {
		s.logger.Error("error getting app data", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseData[models.App]{
		Data: app,
	})
}

// postAppsHandler is the handler for the /apps endpoint.
// It creates a new app.
func (s *server) postAppsHandler(w http.ResponseWriter, r *http.Request) {
	appID := r.PathValue("appID")
	appID, err := validateAppID(appID)
	if err != nil {
		s.logger.Error("error validating appID", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app, err := s.appsClient.GetAppData(appID)
	if err != nil {
		s.logger.Error("error getting app data", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = s.appsClient.AddApp(app)
	if err != nil {
		s.logger.Error("error creating app", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = s.queue.Enqueue([]byte(appID))
	if err != nil {
		s.logger.Error("error enqueuing appID", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(appID)
}

// validateAppID is a helper function to validate the appID.
func validateAppID(appID string) (string, error) {
	if appID == "" {
		return "", errors.New("appID is required")
	}

	// only numbers are allowed
	if _, err := strconv.ParseInt(appID, 10, 64); err != nil {
		return "", errors.New("appID must be a number")
	}

	return appID, nil
}
