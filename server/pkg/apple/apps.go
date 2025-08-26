package apple

import (
	"encoding/json"
	"fmt"
)

const AppleAppsURLFmt = "https://itunes.apple.com/lookup?id=%s"

type App struct {
	TrackID       int    `json:"trackId"`
	TrackName     string `json:"trackName"`
	ArtworkURL512 string `json:"artworkUrl512"`
}

type AppsResponse struct {
	ResultCount int    `json:"resultCount"`
	Results     [1]App `json:"results"`
}

// GetAppData returns the app data for a given app ID.
func (c *AppleClient) GetAppData(appID string) (AppsResponse, error) {
	url := fmt.Sprintf(AppleAppsURLFmt, appID)

	response, err := c.httpClient.Get(url)
	if err != nil {
		return AppsResponse{}, err
	}
	defer response.Body.Close()

	var appsResponse AppsResponse
	if err := json.NewDecoder(response.Body).Decode(&appsResponse); err != nil {
		return AppsResponse{}, err
	}

	return appsResponse, nil
}
