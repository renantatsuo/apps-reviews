package apple

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	AppleRSSURLFmt  = "https://itunes.apple.com/us/rss/customerreviews/id=%s/sortBy=mostRecent/json"
	AppleTimeFormat = "2006-01-02T15:04:05-07:00"
)

var ErrNoNextPage = fmt.Errorf("no next page")

type ReviewsResponse[T any] struct {
	Feed struct {
		Entry []T          `json:"entry"`
		Link  []ReviewLink `json:"link"`
	} `json:"feed"`
	httpClient *http.Client `json:"-"`
}

type Review struct {
	ID struct {
		Label string `json:"label"`
	} `json:"id"`
	Author ReviewAuthor `json:"author"`
	Title  struct {
		Label string `json:"label"`
	} `json:"title"`
	Content struct {
		Label string `json:"label"`
	} `json:"content"`
	Rating struct {
		Label string `json:"label"`
	} `json:"im:rating"`
	Updated struct {
		Label string `json:"label"`
	} `json:"updated"`
}

type ReviewAuthor struct {
	Name struct {
		Label string `json:"label"`
	} `json:"name"`
	Uri struct {
		Label string `json:"label"`
	} `json:"uri"`
}

type ReviewLink struct {
	Attributes struct {
		Rel  string `json:"rel"`
		HREF string `json:"href"`
	} `json:"attributes"`
}

func getAppleRSSURL(appID string) string {
	return fmt.Sprintf(AppleRSSURLFmt, appID)
}

// GetLatestReviews returns the latest reviews for a given app ID
func (c *AppleClient) GetLatestReviews(appID string) (ReviewsResponse[Review], error) {
	url := getAppleRSSURL(appID)

	response, err := c.httpClient.Get(url)
	if err != nil {
		return ReviewsResponse[Review]{}, err
	}
	defer response.Body.Close()

	var reviewsResponse ReviewsResponse[Review]
	if err := json.NewDecoder(response.Body).Decode(&reviewsResponse); err != nil {
		return ReviewsResponse[Review]{}, err
	}

	reviewsResponse.httpClient = c.httpClient

	return reviewsResponse, nil
}

// HasNext returns true if there is a next page of reviews.
func (r *ReviewsResponse[Review]) HasNext() bool {
	if len(r.Feed.Link) == 0 || len(r.Feed.Entry) == 0 {
		return false
	}

	for _, link := range r.Feed.Link {
		if link.Attributes.Rel == "next" {
			return true
		}
	}

	return false
}

// Next returns the next page of reviews.
func (r *ReviewsResponse[Review]) Next() (ReviewsResponse[Review], error) {
	if !r.HasNext() {
		return ReviewsResponse[Review]{}, ErrNoNextPage
	}

	var nextPageURL string

	for _, link := range r.Feed.Link {
		if link.Attributes.Rel == "next" {
			nextPageURL = fixNextPageURL(link.Attributes.HREF)
			break
		}
	}

	if nextPageURL == "" {
		return ReviewsResponse[Review]{}, ErrNoNextPage
	}

	response, err := r.httpClient.Get(nextPageURL)
	if err != nil {
		return ReviewsResponse[Review]{}, err
	}
	defer response.Body.Close()

	var nextPageResponse ReviewsResponse[Review]
	if err := json.NewDecoder(response.Body).Decode(&nextPageResponse); err != nil {
		return ReviewsResponse[Review]{}, err
	}

	nextPageResponse.httpClient = r.httpClient

	return nextPageResponse, nil
}

// fixNextPageURL fixes the next page URL
// sometimes apple returns a xml url, but we want a json url
func fixNextPageURL(url string) string {
	return strings.ReplaceAll(url, "/xml", "/json")
}
