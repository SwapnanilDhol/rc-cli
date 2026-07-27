package internalapi

import (
	"strings"
	"testing"
)

func TestLooksLikeID(t *testing.T) {
	cases := []struct {
		ref  string
		want bool
	}{
		// Projects use bare hex — assuming a "proj" prefix here was a real bug.
		{"1439b090", true},
		{"bf71a0a6", true},
		{"e2d63c06", true},
		// Offerings and other entities use prefixes.
		{"ofrng56fe7fb1eb", true},
		{"prodabc123", true},
		{"app99f0aa", true},
		// Human names must not be mistaken for IDs.
		{"Recur", false},
		{"Money Tracker", false},
		{"Palette Buddy", false},
		{"pro-access-weekly-yearly-lifetime", false},
		{"", false},
		// Too short to be a hex ID.
		{"abc", false},
		// Bare prefix with nothing after it is not an ID.
		{"ofrng", false},
	}
	for _, c := range cases {
		if got := looksLikeID(c.ref); got != c.want {
			t.Errorf("looksLikeID(%q) = %v, want %v", c.ref, got, c.want)
		}
	}
}

func items(vals ...map[string]interface{}) []interface{} {
	out := make([]interface{}, 0, len(vals))
	for _, v := range vals {
		out = append(out, v)
	}
	return out
}

var projects = items(
	map[string]interface{}{"id": "1439b090", "name": "Recur"},
	map[string]interface{}{"id": "bf71a0a6", "name": "Aeronautical"},
	map[string]interface{}{"id": "8f1fd662", "name": "Money Tracker"},
)

func TestResolveByRefMatchesLiteralID(t *testing.T) {
	// The regression: a literal ID must resolve even though it looks like
	// nothing in the name fields.
	got, err := resolveByRef(projects, "1439b090", "project")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "1439b090" {
		t.Errorf("got %q, want %q", got, "1439b090")
	}
}

func TestResolveByRefMatchesName(t *testing.T) {
	for _, ref := range []string{"Recur", "recur", "RECUR"} {
		got, err := resolveByRef(projects, ref, "project")
		if err != nil {
			t.Fatalf("resolveByRef(%q): %v", ref, err)
		}
		if got != "1439b090" {
			t.Errorf("resolveByRef(%q) = %q, want 1439b090", ref, got)
		}
	}
}

func TestResolveByRefMatchesIdentifier(t *testing.T) {
	offerings := items(
		map[string]interface{}{"id": "ofrng1", "identifier": "pro-weekly", "display_name": "Weekly"},
		map[string]interface{}{"id": "ofrng2", "identifier": "pro-yearly", "display_name": "Yearly"},
	)
	got, err := resolveByRef(offerings, "pro-yearly", "offering")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "ofrng2" {
		t.Errorf("got %q, want ofrng2", got)
	}
}

func TestResolveByRefUnknownListsCandidates(t *testing.T) {
	_, err := resolveByRef(projects, "Nope", "project")
	if err == nil {
		t.Fatal("expected an error for an unknown reference")
	}
	// The error must be actionable — it is the user's only way to discover IDs.
	for _, want := range []string{"Recur", "1439b090"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q should list %q", err, want)
		}
	}
}

func TestResolveByRefAmbiguousIsAnError(t *testing.T) {
	dupes := items(
		map[string]interface{}{"id": "ofrng1", "display_name": "Default"},
		map[string]interface{}{"id": "ofrng2", "display_name": "Default"},
	)
	_, err := resolveByRef(dupes, "Default", "offering")
	if err == nil {
		t.Fatal("ambiguous reference must error rather than pick the first match")
	}
	if !strings.Contains(err.Error(), "ambiguous") {
		t.Errorf("error should say it is ambiguous, got: %v", err)
	}
}

func TestResolveByRefPrefersExactIDOverAmbiguousName(t *testing.T) {
	// An entity whose *name* collides with another entity's ID must still be
	// reachable by its own ID.
	tricky := items(
		map[string]interface{}{"id": "ofrng1", "display_name": "ofrng2"},
		map[string]interface{}{"id": "ofrng2", "display_name": "Real"},
	)
	got, err := resolveByRef(tricky, "ofrng2", "offering")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "ofrng2" {
		t.Errorf("got %q, want ofrng2 (exact id match must win)", got)
	}
}

func TestResolveByRefEmpty(t *testing.T) {
	if _, err := resolveByRef(projects, "", "project"); err == nil {
		t.Fatal("empty reference must error")
	}
}
