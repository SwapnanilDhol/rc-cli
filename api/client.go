package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"revenuecat-cli/config"
)

const BaseURL = "https://api.revenuecat.com/v2"

type Client struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

func NewClient(cfg *config.Config) (*Client, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("API key not configured. Run: rc config\nOr provide with: rc <command> --api-key <your-api-key>")
	}

	return &Client{
		httpClient: &http.Client{},
		apiKey:     cfg.APIKey,
		baseURL:    BaseURL,
	}, nil
}

// DoRaw performs an HTTP request against the v2 API and returns status and body as
// returned by the server. path must start with "/" (e.g. "/projects/foo/customers").
// query may be nil. Every other method here is built on this one.
func (c *Client) DoRaw(method, path string, query url.Values, body []byte) (status int, respBody []byte, err error) {
	if len(path) == 0 || path[0] != '/' {
		return 0, nil, fmt.Errorf("path must start with /, got %q", path)
	}
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return 0, nil, err
	}
	if query != nil {
		u.RawQuery = query.Encode()
	}

	var rdr io.Reader
	if len(body) > 0 {
		rdr = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, u.String(), rdr)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}
	return resp.StatusCode, raw, nil
}

// get issues a GET and decodes the standard v2 envelope.
func (c *Client) get(path string, query url.Values) (*Response, error) {
	status, body, err := c.DoRaw(http.MethodGet, path, query, nil)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	var response Response
	if len(body) > 0 {
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, fmt.Errorf("error unmarshaling response: %w", err)
		}
	}
	response.StatusCode = status
	return &response, nil
}

func (c *Client) Get(path string) (*Response, error) {
	return c.get(path, nil)
}

func (c *Client) GetWithParams(path string, params map[string]string) (*Response, error) {
	q := make(url.Values, len(params))
	for k, v := range params {
		q.Set(k, v)
	}
	return c.get(path, q)
}

type Response struct {
	StatusCode int           `json:"-"`
	Items      []interface{} `json:"items,omitempty"`
	NextPage   string        `json:"next_page,omitempty"`
	Data       interface{}   `json:"data,omitempty"`
	Error      string        `json:"error,omitempty"`
	Message    string        `json:"message,omitempty"`
}
