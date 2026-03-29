package internal

import (
	"testing"
)

func TestDecodeInternalResponseBody_topLevelArray(t *testing.T) {
	body := []byte(`[{"id":"bf71","name":"A","owner_email":"x@y.com","restricted_access":false}]`)
	r, err := decodeInternalResponseBody(body, 200)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Items) != 1 {
		t.Fatalf("Items len = %d", len(r.Items))
	}
	m, ok := r.Items[0].(map[string]interface{})
	if !ok || m["id"] != "bf71" {
		t.Fatalf("item = %#v", r.Items[0])
	}
}

func TestDecodeInternalResponseBody_envelope(t *testing.T) {
	body := []byte(`{"data":{"id":"1"},"has_next_page":true}`)
	r, err := decodeInternalResponseBody(body, 200)
	if err != nil {
		t.Fatal(err)
	}
	data, ok := r.Data.(map[string]interface{})
	if !ok || data["id"] != "1" {
		t.Fatalf("Data = %#v", r.Data)
	}
	if !r.HasNext {
		t.Fatal("expected HasNext")
	}
}

func TestDecodeInternalResponseBody_bareObject(t *testing.T) {
	body := []byte(`{"id":"p1","name":"Proj","owner_email":"a@b"}`)
	r, err := decodeInternalResponseBody(body, 200)
	if err != nil {
		t.Fatal(err)
	}
	data, ok := r.Data.(map[string]interface{})
	if !ok || data["name"] != "Proj" {
		t.Fatalf("Data = %#v", r.Data)
	}
}

func TestDecodeInternalResponseBody_errorShape(t *testing.T) {
	body := []byte(`{"code":7117,"message":"nope"}`)
	r, err := decodeInternalResponseBody(body, 404)
	if err != nil {
		t.Fatal(err)
	}
	if string(r.Code) != "7117" || r.Message != "nope" {
		t.Fatalf("code=%q msg=%q", r.Code, r.Message)
	}
}
