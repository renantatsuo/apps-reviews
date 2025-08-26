package apps

import "github.com/renantatsuo/app-review/server/internal/models"

// AddApp adds a new app to the database.
func (a *AppsClient) AddApp(app models.App) error {
	_, err := a.db.Exec("INSERT INTO apps (id, name, thumbnail_url) VALUES (?, ?, ?)", app.ID, app.Name, app.ThumbnailURL)
	if err != nil {
		return err
	}
	return nil
}

// GetAllApps returns all the app IDs from the database.
func (a *AppsClient) GetAllApps() ([]models.App, error) {
	apps := []models.App{}

	rows, err := a.db.Query("SELECT * FROM apps")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var app models.App
		err := rows.Scan(&app.ID, &app.Name, &app.ThumbnailURL, &app.CreatedAt, &app.UpdatedAt)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}

	return apps, nil
}
