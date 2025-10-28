package rest

import (
	"net/http"
	"net/url"
)

type restClient struct {
	baseURL url.URL
	client  *http.Client
}

type RESTClientConfig struct {
	BaseURL    *url.URL
	HTTPClient *http.Client
}

func NewRESTClient(cfg *RESTClientConfig) *restClient {
	return &restClient{
		baseURL: *cfg.BaseURL,
		client:  cfg.HTTPClient,
	}
}

func (c *restClient) Get(path string) (*http.Response, error) {
	reqURL := c.baseURL.ResolveReference(&url.URL{Path: path})
	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, err
	}
	response, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return response, nil
}
