package apps

import (
	"github.com/renantatsuo/app-review/server/internal/models"
)

func (a *AppsClient) GetAppData(appID string) (models.App, error) {
	appsResponse, err := a.appleClient.GetAppData(appID)
	if err != nil {
		return models.App{}, err
	}

	return models.AppFromAppleApp(appsResponse.Results[0]), nil
}
