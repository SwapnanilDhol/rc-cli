package internal

import (
	"encoding/json"
	"testing"
)

func emit(t *testing.T, r *Response) string {
	t.Helper()
	b, err := json.Marshal(r.Payload())
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// The bug this prevents: an empty list used to fall through to a status object,
// so `rc … list --json | jq '.[]'` failed only when a project had no results.
func TestPayloadEmptyListStaysAList(t *testing.T) {
	got := emit(t, &Response{StatusCode: 200, Items: []interface{}{}, IsList: true})
	if got != "[]" {
		t.Errorf("empty list rendered as %s, want []", got)
	}
}

func TestPayloadNilItemsOnListStaysAList(t *testing.T) {
	got := emit(t, &Response{StatusCode: 200, IsList: true})
	if got != "[]" {
		t.Errorf("nil items on a list response rendered as %s, want []", got)
	}
}

func TestPayloadListShapeIsStableAcrossSizes(t *testing.T) {
	empty := emit(t, &Response{StatusCode: 200, Items: []interface{}{}, IsList: true})
	full := emit(t, &Response{StatusCode: 200, IsList: true,
		Items: []interface{}{map[string]interface{}{"id": "x"}}})
	if empty[0] != '[' || full[0] != '[' {
		t.Errorf("same command must emit the same JSON type: empty=%s full=%s", empty, full)
	}
}

func TestPayloadObjectResponsePassesThrough(t *testing.T) {
	got := emit(t, &Response{StatusCode: 200,
		Data: map[string]interface{}{"has_more": false}})
	if got != `{"has_more":false}` {
		t.Errorf("got %s, want the object verbatim", got)
	}
}

// A decoded top-level array must be flagged, otherwise Payload cannot tell an
// empty list from a bodyless response.
func TestDecodeMarksTopLevelArrays(t *testing.T) {
	r, err := decodeInternalResponseBody([]byte(`[]`), 200)
	if err != nil {
		t.Fatal(err)
	}
	if !r.IsList {
		t.Error("a top-level [] must set IsList")
	}
	if got := emit(t, r); got != "[]" {
		t.Errorf("decoded empty array rendered as %s, want []", got)
	}
}

func TestDecodeDoesNotMarkObjects(t *testing.T) {
	r, err := decodeInternalResponseBody([]byte(`{"products":[]}`), 200)
	if err != nil {
		t.Fatal(err)
	}
	if r.IsList {
		t.Error("a top-level object must not set IsList")
	}
}
