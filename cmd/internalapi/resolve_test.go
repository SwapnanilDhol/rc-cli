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
		{"a1b2c3d4", true},
		{"e5f6a7b8", true},
		{"0f1e2d3c", true},
		// Offerings and other entities use prefixes.
		{"ofrng0123456789", true},
		{"prodabc123", true},
		{"app99f0aa", true},
		// Human names must not be mistaken for IDs.
		{"Habits", false},
		{"Budget", false},
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
	map[string]interface{}{"id": "a1b2c3d4", "name": "Habits"},
	map[string]interface{}{"id": "e5f6a7b8", "name": "Weather"},
	map[string]interface{}{"id": "c9d0e1f2", "name": "Budget"},
)

func TestResolveByRefMatchesLiteralID(t *testing.T) {
	// The regression: a literal ID must resolve even though it looks like
	// nothing in the name fields.
	got, err := resolveByRef(projects, "a1b2c3d4", "project")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "a1b2c3d4" {
		t.Errorf("got %q, want %q", got, "a1b2c3d4")
	}
}

func TestResolveByRefMatchesName(t *testing.T) {
	for _, ref := range []string{"Habits", "habits", "HABITS"} {
		got, err := resolveByRef(projects, ref, "project")
		if err != nil {
			t.Fatalf("resolveByRef(%q): %v", ref, err)
		}
		if got != "a1b2c3d4" {
			t.Errorf("resolveByRef(%q) = %q, want a1b2c3d4", ref, got)
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
	for _, want := range []string{"Habits", "a1b2c3d4"} {
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
