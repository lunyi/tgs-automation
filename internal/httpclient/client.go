package httpclient

import "net/http"

// HTTPClient defines the interface for an HTTP client.
type HttpClient interface {
	Get(url string) (*http.Response, error)
}

type StandardHTTPClient struct {
	Client *http.Client
}

func (c *StandardHTTPClient) Get(url string) (*http.Response, error) {
	return c.Client.Get(url)
}
