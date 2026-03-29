package internal

import (
	"encoding/json"
	"testing"
)

func TestFlexStr_numericCode(t *testing.T) {
	const payload = `{"code":7117,"message":"something went wrong"}`
	var r Response
	if err := json.Unmarshal([]byte(payload), &r); err != nil {
		t.Fatal(err)
	}
	if string(r.Code) != "7117" {
		t.Fatalf("Code = %q, want 7117", r.Code)
	}
	if r.Message != "something went wrong" {
		t.Fatalf("Message = %q", r.Message)
	}
}

func TestFlexStr_stringCode(t *testing.T) {
	const payload = `{"code":"ERR_1","message":"nope"}`
	var r Response
	if err := json.Unmarshal([]byte(payload), &r); err != nil {
		t.Fatal(err)
	}
	if string(r.Code) != "ERR_1" {
		t.Fatalf("Code = %q", r.Code)
	}
}

func TestFlexStr_loginResponse(t *testing.T) {
	const payload = `{"code":401,"message":"bad"}`
	var lr LoginResponse
	if err := json.Unmarshal([]byte(payload), &lr); err != nil {
		t.Fatal(err)
	}
	if string(lr.Code) != "401" {
		t.Fatalf("Code = %q", lr.Code)
	}
}
