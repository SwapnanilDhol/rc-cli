package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"revenuecat-cli/config"
)

const (
	InternalBaseURL   = "https://app.revenuecat.com/internal/v1"
	LoginURL         = "https://app.revenuecat.com/v1/developers/login"
	RefreshTokenURL  = "https://app.revenuecat.com/v1/developers/login/refresh-token"
)

type Client struct {
	httpClient *http.Client
	authToken  string
}

func NewClient(cfg *config.Config) (*Client, error) {
	if cfg.AuthToken == "" {
		return nil, fmt.Errorf("not logged in. Run: rc login")
	}

	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		authToken: cfg.AuthToken,
	}, nil
}

// EnsureAuthenticated returns a client with a valid auth token, refreshing if needed
func EnsureAuthenticated(cfg *config.Config) (*Client, error) {
	if cfg.AuthToken == "" {
		if cfg.Email == "" || cfg.Password == "" {
			return nil, fmt.Errorf("not logged in. Run: rc login")
		}
		// Need to login with email/password
		loginResp, err := Login(cfg.Email, cfg.Password)
		if err != nil {
			return nil, fmt.Errorf("login failed: %w", err)
		}
		cfg.AuthToken = loginResp.AuthenticationToken
		// Note: caller should save config if they want to persist
	}

	client := &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		authToken: cfg.AuthToken,
	}

	// Try a simple request to verify token is valid
	_, err := client.Get("/projects")
	if err != nil {
		// Token might be expired, try to refresh
		refreshResp, refreshErr := RefreshToken(cfg.AuthToken)
		if refreshErr != nil {
			// Refresh failed, try re-logging in with email/password
			if cfg.Email == "" || cfg.Password == "" {
				return nil, fmt.Errorf("session expired. Run: rc login")
			}
			loginResp, loginErr := Login(cfg.Email, cfg.Password)
			if loginErr != nil {
				return nil, fmt.Errorf("re-login failed: %w", loginErr)
			}
			cfg.AuthToken = loginResp.AuthenticationToken
			client.authToken = cfg.AuthToken
		} else {
			// Refresh succeeded, update token
			cfg.AuthToken = refreshResp.AuthenticationToken
			client.authToken = cfg.AuthToken
		}
	}

	return client, nil
}

func (c *Client) Get(path string) (*Response, error) {
	return c.doRequest("GET", path, nil)
}

func (c *Client) GetWithParams(path string, params map[string]string) (*Response, error) {
	if len(params) > 0 {
		query := ""
		for key, value := range params {
			if query != "" {
				query += "&"
			}
			query += key + "=" + value
		}
		if path != "" {
			path += "?"
		}
		path += query
	}
	return c.doRequest("GET", path, nil)
}

func (c *Client) Post(path string, data interface{}) (*Response, error) {
	return c.doRequest("POST", path, data)
}

func (c *Client) Put(path string, data interface{}) (*Response, error) {
	return c.doRequest("PUT", path, data)
}

func (c *Client) Delete(path string) (*Response, error) {
	return c.doRequest("DELETE", path, nil)
}

func (c *Client) Patch(path string, data interface{}) (*Response, error) {
	return c.doRequest("PATCH", path, data)
}

func (c *Client) doRequest(method, path string, data interface{}) (*Response, error) {
	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	url := InternalBaseURL + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Cookie", "rc_auth_token="+c.authToken)

	if method != "GET" && method != "DELETE" {
		req.Header.Set("Origin", "https://app.revenuecat.com")
		req.Header.Set("Referer", "https://app.revenuecat.com/")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var response Response
	if err := json.Unmarshal(respBody, &response); err != nil {
		// Check if it's an HTML error page
		return nil, fmt.Errorf("error unmarshaling response: %w (body: %s)", err, string(respBody[:min(200, len(respBody))]))
	}

	response.StatusCode = resp.StatusCode
	return &response, nil
}

// Login authenticates with email/password and returns the auth token
func Login(email, password string) (*LoginResponse, error) {
	data := map[string]string{
		"email":    email,
		"password": password,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", LoginURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Origin", "https://app.revenuecat.com")
	req.Header.Set("Referer", "https://app.revenuecat.com/login")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(respBody, &loginResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	if loginResp.Message != "" {
		return nil, fmt.Errorf("login failed: %s", loginResp.Message)
	}

	return &loginResp, nil
}

// RefreshToken uses the existing auth token to get a new one
func RefreshToken(authToken string) (*LoginResponse, error) {
	req, err := http.NewRequest("POST", RefreshTokenURL, bytes.NewBuffer([]byte("{}")))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Origin", "https://app.revenuecat.com")
	req.Header.Set("Referer", "https://app.revenuecat.com/login")
	req.Header.Set("Cookie", "rc_auth_token="+authToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(respBody, &loginResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	if loginResp.Message != "" {
		return nil, fmt.Errorf("refresh failed: %s", loginResp.Message)
	}

	return &loginResp, nil
}

type Response struct {
	StatusCode int             `json:"-"`
	Code       string          `json:"code,omitempty"`
	Message    string          `json:"message,omitempty"`
	Data       interface{}     `json:"data,omitempty"`
	Items      []interface{}   `json:"items,omitempty"`
	HasNext    bool            `json:"has_next_page,omitempty"`
	NextPage   string          `json:"next_page,omitempty"`
}

type LoginResponse struct {
	AuthenticationToken       string `json:"authentication_token"`
	AuthenticationTokenExpiration string `json:"authentication_token_expiration"`
	DistinctID               string `json:"distinct_id"`
	Email                    string `json:"email"`
	Message                  string `json:"message"`
	Code                     string `json:"code,omitempty"`
}

// Project represents a RevenueCat project
type Project struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	OwnerEmail       string `json:"owner_email"`
	RestrictedAccess bool   `json:"restricted_access"`
	IconURL          string `json:"icon_url"`
}

// Entitlement represents an entitlement
type Entitlement struct {
	ID          string    `json:"id"`
	Identifier  string    `json:"identifier"`
	DisplayName string    `json:"display_name"`
	IsArchived  bool      `json:"is_archived"`
	CreatedAt  time.Time `json:"created_at"`
	Products    []Product `json:"products,omitempty"`
}

// Offering represents an offering
type Offering struct {
	ID          string     `json:"id"`
	Identifier  string     `json:"identifier"`
	DisplayName string     `json:"display_name"`
	IsArchived  bool       `json:"is_archived"`
	IsCurrent   bool       `json:"is_current"`
	CreatedAt   time.Time  `json:"created_at"`
	Metadata    interface{} `json:"metadata"`
	Packages    []Package  `json:"packages,omitempty"`
}

// Package represents a package within an offering
type Package struct {
	ID          string     `json:"id"`
	Identifier  string     `json:"identifier"`
	DisplayName string    `json:"display_name"`
	Position    int        `json:"position"`
	CreatedAt   time.Time  `json:"created_at"`
	Products    []Product  `json:"products,omitempty"`
}

// Product represents a product
type Product struct {
	ID           string    `json:"id"`
	Identifier   string    `json:"identifier"`
	DisplayName  string    `json:"display_name"`
	IsArchived   bool      `json:"is_archived"`
	IsSubscription bool    `json:"is_subscription"`
	ProductType  string    `json:"product_type"`
	CreatedAt    time.Time `json:"created_at"`
	App          *App      `json:"app,omitempty"`
	Entitlements []Entitlement `json:"entitlements,omitempty"`
	ProductGroup *ProductGroup `json:"product_group,omitempty"`
}

// ProductGroup represents a product group
type ProductGroup struct {
	ID         string    `json:"id"`
	Identifier string    `json:"identifier"`
	CreatedAt  time.Time `json:"created_at"`
}

// App represents an app
type App struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID              string                 `json:"id"`
	ActionType      string                 `json:"action_type"`
	ActorIdentifier string                 `json:"actor_identifier"`
	ActorType       string                 `json:"actor_type"`
	TargetIdentifier string                `json:"target_identifier"`
	TargetType      string                 `json:"target_type"`
	OccurredAt      time.Time              `json:"occurred_at"`
	AdditionalData  map[string]interface{} `json:"additional_data"`
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Helper to parse next page URL
func ParseNextPageURL(nextPage string) string {
	if nextPage == "" {
		return ""
	}
	u, err := url.Parse(nextPage)
	if err != nil {
		return ""
	}
	return u.Query().Get("starting_after")
}
