package internalapi

import (
	"strings"
	"testing"

	rcinternal "revenuecat-cli/internal"
)

// A 401/403 body decodes cleanly into an empty Response, so without a status
// check the caller printed "no results found" and exited 0.
func TestCheckResponseRejectsAuthFailures(t *testing.T) {
	for _, status := range []int{401, 403} {
		err := CheckResponse(&rcinternal.Response{StatusCode: status, Message: "unauthorized"})
		if err == nil {
			t.Fatalf("HTTP %d must produce an error", status)
		}
		if !strings.Contains(err.Error(), "rc login") {
			t.Errorf("HTTP %d error should tell the user to log in, got: %v", status, err)
		}
	}
}

func TestCheckResponseRejectsOtherHTTPErrors(t *testing.T) {
	err := CheckResponse(&rcinternal.Response{StatusCode: 500, Message: "boom"})
	if err == nil {
		t.Fatal("HTTP 500 must produce an error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should mention the status, got: %v", err)
	}
}

// Some endpoints return 200 with an error code in the body.
func TestCheckResponseRejectsAPIErrorCode(t *testing.T) {
	err := CheckResponse(&rcinternal.Response{StatusCode: 200, Code: "7001", Message: "nope"})
	if err == nil {
		t.Fatal("an API error code must produce an error even on HTTP 200")
	}
}

func TestCheckResponseAcceptsSuccess(t *testing.T) {
	if err := CheckResponse(&rcinternal.Response{StatusCode: 200}); err != nil {
		t.Errorf("200 with no error code should pass, got: %v", err)
	}
	if err := CheckResponse(&rcinternal.Response{StatusCode: 204}); err != nil {
		t.Errorf("204 should pass, got: %v", err)
	}
}

func TestCheckResponseNil(t *testing.T) {
	if err := CheckResponse(nil); err == nil {
		t.Fatal("a nil response must error rather than panic downstream")
	}
}
