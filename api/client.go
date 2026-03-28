package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"revenuecat-cli/config"
)

const BaseURL = "https://api.revenuecat.com/v2"

type Client struct {
	httpClient *http.Client
	apiKey    string
}

func NewClient(cfg *config.Config) (*Client, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("API key not configured. Run: rc config\nOr provide with: rc <command> --api-key <your-api-key>")
	}

	return &Client{
		httpClient: &http.Client{},
		apiKey:     cfg.APIKey,
	}, nil
}

func (c *Client) Get(path string) (*Response, error) {
	url := BaseURL + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	response.StatusCode = resp.StatusCode
	return &response, nil
}

func (c *Client) GetWithParams(path string, params map[string]string) (*Response, error) {
	url := BaseURL + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	response.StatusCode = resp.StatusCode
	return &response, nil
}

type Response struct {
	StatusCode int             `json:"-"`
	Items      []interface{}  `json:"items,omitempty"`
	NextPage   string          `json:"next_page,omitempty"`
	Data       interface{}     `json:"data,omitempty"`
	Error      string          `json:"error,omitempty"`
	Message    string          `json:"message,omitempty"`
}
