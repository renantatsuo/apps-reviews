package apps

import (
	"fmt"

	"github.com/renantatsuo/app-review/server/internal/models"
)

// ErrAppNotFound is an error type for when an app is not found.
type ErrAppNotFound struct {
	AppID string
}

func (e ErrAppNotFound) Error() string {
	return fmt.Sprintf("app not found: %s", e.AppID)
}

// GetAppData gets the app data from the Apple API.
func (a *AppsClient) GetAppData(appID string) (models.App, error) {
	appsResponse, err := a.appleClient.GetAppData(appID)
	if err != nil {
		return models.App{}, err
	}

	if appsResponse.ResultCount == 0 {
		return models.App{}, ErrAppNotFound{AppID: appID}
	}

	return models.AppFromAppleApp(appsResponse.Results[0]), nil
}
