package internal

import (
	"net/url"
	"testing"
)

// GetWithParams used to concatenate the query by hand: values went unescaped and
// key order changed run to run because Go randomises map iteration.
func TestQueryEncodingIsEscapedAndStable(t *testing.T) {
	params := map[string]string{"platform": "IOS", "q": "a b&c=d", "z": "1"}
	q := make(url.Values, len(params))
	for k, v := range params {
		q.Set(k, v)
	}
	first := q.Encode()
	for i := 0; i < 50; i++ {
		if q.Encode() != first {
			t.Fatal("query encoding is not stable")
		}
	}
	if first != "platform=IOS&q=a+b%26c%3Dd&z=1" {
		t.Errorf("got %q; values must be escaped and keys sorted", first)
	}
}
