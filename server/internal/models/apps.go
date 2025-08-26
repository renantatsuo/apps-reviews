package models

import (
	"strconv"

	"github.com/renantatsuo/app-review/server/pkg/apple"
)

type App struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ThumbnailURL string `json:"thumbnail_url"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func AppFromAppleApp(app apple.App) App {
	return App{
		ID:           strconv.Itoa(app.TrackID),
		Name:         app.TrackName,
		ThumbnailURL: app.ArtworkURL512,
	}
}
