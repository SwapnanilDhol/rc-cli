package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func testClient(t *testing.T, h http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return &Client{httpClient: srv.Client(), apiKey: "sk_test", baseURL: srv.URL}
}

func TestGetSendsBearerAuth(t *testing.T) {
	var gotAuth string
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.Write([]byte(`{"items":[]}`))
	})
	if _, err := c.Get("/projects"); err != nil {
		t.Fatal(err)
	}
	if gotAuth != "Bearer sk_test" {
		t.Errorf("Authorization = %q, want Bearer sk_test", gotAuth)
	}
}

func TestGetWithParamsEscapesAndSorts(t *testing.T) {
	var gotQuery string
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		w.Write([]byte(`{}`))
	})
	_, err := c.GetWithParams("/projects", map[string]string{"limit": "20", "q": "a b&c"})
	if err != nil {
		t.Fatal(err)
	}
	if gotQuery != "limit=20&q=a+b%26c" {
		t.Errorf("query = %q; values must be escaped and keys sorted", gotQuery)
	}
}

func TestGetReportsStatusCode(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
		w.Write([]byte(`{"message":"nope"}`))
	})
	resp, err := c.Get("/projects")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 403 || resp.Message != "nope" {
		t.Errorf("got %d %q, want 403 nope", resp.StatusCode, resp.Message)
	}
}

// An empty body must not be treated as a decode failure.
func TestGetToleratesEmptyBody(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	})
	resp, err := c.Get("/projects")
	if err != nil {
		t.Fatalf("empty body should not error: %v", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("got %d, want 204", resp.StatusCode)
	}
}

func TestDoRawRejectsBadPath(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {})
	if _, _, err := c.DoRaw("GET", "projects", nil, nil); err == nil {
		t.Fatal("a path without a leading slash must error")
	}
}
