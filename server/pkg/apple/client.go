package apple

import (
	"net/http"
	"time"
)

type AppleClient struct {
	httpClient *http.Client
}

type Option func(*AppleClient)

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *AppleClient) {
		c.httpClient = httpClient
	}
}

func New(opts ...Option) *AppleClient {
	client := &AppleClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}
