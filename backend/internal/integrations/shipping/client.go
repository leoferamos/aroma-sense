package shipping

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"path"
	"strings"
	"time"
)

// Client handles HTTP calls to the external shipping provider API.
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient builds a Client from explicit Config.
func NewClient(cfg Config) (*Client, error) {
	if cfg.BaseURL == "" {
		return nil, errors.New("shipping provider base URL not set")
	}
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	return &Client{
		httpClient: &http.Client{Timeout: timeout},
		baseURL:    cfg.BaseURL,
	}, nil
}

// NewJSONRequest creates an HTTP request with JSON body.
func (c *Client) NewJSONRequest(method, relPath string, body any) (*http.Request, error) {
	var rdr *bytes.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		rdr = bytes.NewReader(b)
	} else {
		rdr = bytes.NewReader(nil)
	}
	url := strings.TrimRight(c.baseURL, "/") + path.Clean(relPath)
	req, err := http.NewRequest(method, url, rdr)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// Do executes the given HTTP request.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}
