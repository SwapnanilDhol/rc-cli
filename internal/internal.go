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

// FlexStr unmarshals JSON fields that may be string or number (RevenueCat error payloads use numeric "code").
type FlexStr string

func (f *FlexStr) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		*f = ""
		return nil
	}
	switch b[0] {
	case '"':
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		*f = FlexStr(s)
		return nil
	case 'n': // null
		*f = ""
		return nil
	default:
		var n json.Number
		if err := json.Unmarshal(b, &n); err != nil {
			return err
		}
		*f = FlexStr(n)
		return nil
	}
}

const (
	InternalBaseURL   = "https://app.revenuecat.com/internal/v1"
	LoginURL         = "https://app.revenuecat.com/v1/developers/login"
	RefreshTokenURL  = "https://app.revenuecat.com/v1/developers/login/refresh-token"
)

type Client struct {
	httpClient *http.Client
	cfg        *config.Config
	// onTokenChange persists a refreshed token. Optional.
	onTokenChange func(*config.Config) error
}

func NewClient(cfg *config.Config) (*Client, error) {
	if cfg.AuthToken == "" {
		return nil, fmt.Errorf("not logged in. Run: rc login")
	}

	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		cfg:        cfg,
	}, nil
}

// EnsureAuthenticated returns a client for the dashboard API. It does not probe the
// network: an expired token is detected from the first real request's 401 and
// refreshed transparently by doRequest. The bool reports whether cfg.AuthToken was
// changed here (by an eager email/password login) and so needs persisting.
func EnsureAuthenticated(cfg *config.Config) (*Client, bool, error) {
	changed := false
	if cfg.AuthToken == "" {
		if cfg.Email == "" || cfg.Password == "" {
			return nil, false, fmt.Errorf("not logged in. Run: rc login")
		}
		loginResp, err := Login(cfg.Email, cfg.Password)
		if err != nil {
			return nil, false, fmt.Errorf("login failed: %w", err)
		}
		cfg.AuthToken = loginResp.AuthenticationToken
		changed = true
	}

	return &Client{
		httpClient:    &http.Client{Timeout: 30 * time.Second},
		cfg:           cfg,
		onTokenChange: config.SaveConfig,
	}, changed, nil
}

// reauthenticate obtains a fresh token after a 401, preferring the refresh endpoint
// and falling back to a full email/password login. It persists the new token.
func (c *Client) reauthenticate() error {
	if refreshResp, err := RefreshToken(c.cfg.AuthToken); err == nil && refreshResp.AuthenticationToken != "" {
		c.cfg.AuthToken = refreshResp.AuthenticationToken
	} else {
		if c.cfg.Email == "" || c.cfg.Password == "" {
			return fmt.Errorf("session expired. Run: rc login")
		}
		loginResp, loginErr := Login(c.cfg.Email, c.cfg.Password)
		if loginErr != nil {
			return fmt.Errorf("session expired and re-login failed: %w", loginErr)
		}
		c.cfg.AuthToken = loginResp.AuthenticationToken
	}
	if c.onTokenChange != nil {
		if err := c.onTokenChange(c.cfg); err != nil {
			return fmt.Errorf("failed to save refreshed token: %w", err)
		}
	}
	return nil
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

// PatchWithPackages patches an offering at the offerings_with_packages endpoint,
// which allows setting both metadata AND attaching/detaching packages in one call.
func (c *Client) PatchWithPackages(path string, data interface{}) (*Response, error) {
	return c.doRequest("PATCH", path, data)
}

func (c *Client) doRequest(method, path string, data interface{}) (*Response, error) {
	var payload []byte
	if data != nil {
		var err error
		payload, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request: %w", err)
		}
	}

	response, err := c.attempt(method, path, payload)
	if err != nil {
		return nil, err
	}

	// A 401 is the only reliable signal that the session cookie has expired — the
	// body still decodes cleanly, so it must be detected from the status code.
	if response.StatusCode == http.StatusUnauthorized {
		if err := c.reauthenticate(); err != nil {
			return nil, err
		}
		return c.attempt(method, path, payload)
	}

	return response, nil
}

func (c *Client) attempt(method, path string, payload []byte) (*Response, error) {
	var body io.Reader
	if payload != nil {
		body = bytes.NewReader(payload)
	}

	url := InternalBaseURL + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Cookie", "rc_auth_token="+c.cfg.AuthToken)

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

	return decodeInternalResponseBody(respBody, resp.StatusCode)
}

// DoRaw performs an arbitrary request against the internal dashboard API and returns
// the status and body exactly as the server sent them.
func (c *Client) DoRaw(method, path string, body []byte) (int, []byte, error) {
	if len(path) == 0 || path[0] != '/' {
		return 0, nil, fmt.Errorf("path must start with /, got %q", path)
	}

	do := func() (int, []byte, error) {
		var rdr io.Reader
		if len(body) > 0 {
			rdr = bytes.NewReader(body)
		}
		req, err := http.NewRequest(method, InternalBaseURL+path, rdr)
		if err != nil {
			return 0, nil, err
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Requested-With", "XMLHttpRequest")
		req.Header.Set("Cookie", "rc_auth_token="+c.cfg.AuthToken)
		if method != "GET" && method != "DELETE" {
			req.Header.Set("Origin", "https://app.revenuecat.com")
			req.Header.Set("Referer", "https://app.revenuecat.com/")
		}
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return 0, nil, err
		}
		defer resp.Body.Close()
		raw, err := io.ReadAll(resp.Body)
		return resp.StatusCode, raw, err
	}

	status, raw, err := do()
	if err != nil {
		return status, raw, err
	}
	if status == http.StatusUnauthorized {
		if err := c.reauthenticate(); err != nil {
			return status, raw, err
		}
		return do()
	}
	return status, raw, nil
}

// decodeInternalResponseBody parses the dashboard internal API JSON body. Some endpoints
// return a top-level array (e.g. list projects) or a bare object instead of { "data": ... }.
func decodeInternalResponseBody(respBody []byte, statusCode int) (*Response, error) {
	trim := bytes.TrimSpace(respBody)
	if len(trim) == 0 {
		return &Response{StatusCode: statusCode}, nil
	}
	switch trim[0] {
	case '[':
		var items []interface{}
		if err := json.Unmarshal(trim, &items); err != nil {
			return nil, fmt.Errorf("error unmarshaling response: %w (body: %s)", err, string(trim[:min(200, len(trim))]))
		}
		return &Response{StatusCode: statusCode, Items: items}, nil
	case '{':
		var response Response
		if err := json.Unmarshal(trim, &response); err != nil {
			return nil, fmt.Errorf("error unmarshaling response: %w (body: %s)", err, string(trim[:min(200, len(trim))]))
		}
		response.StatusCode = statusCode
		if response.Data == nil && len(response.Items) == 0 && response.Code == "" && response.Message == "" && !response.HasNext && response.NextPage == "" {
			var raw map[string]interface{}
			if err := json.Unmarshal(trim, &raw); err != nil {
				return nil, fmt.Errorf("error unmarshaling response: %w (body: %s)", err, string(trim[:min(200, len(trim))]))
			}
			response.Data = raw
		}
		return &response, nil
	default:
		return nil, fmt.Errorf("error unmarshaling response: unexpected JSON (body: %s)", string(trim[:min(200, len(trim))]))
	}
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
	Code       FlexStr         `json:"code,omitempty"`
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
	Code                     FlexStr `json:"code,omitempty"`
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
