package internalapi

import (
	"reflect"
	"testing"
)

// The whole point of --metadata-merge: a bare --metadata PATCH replaces the
// document (verified live — it took an 18-key object down to 1).
func TestMergeMetadataPreservesUntouchedKeys(t *testing.T) {
	current := map[string]interface{}{
		"title":    "Unlock Recur Pro",
		"subtitle": "Unlimited Subscriptions",
		"features": []interface{}{"a", "b"},
		"reviews":  []interface{}{map[string]interface{}{"stars": 5}},
	}
	got := mergeMetadata(current, map[string]interface{}{"title": "MERGED"})

	if got["title"] != "MERGED" {
		t.Errorf("title = %v, want MERGED", got["title"])
	}
	if len(got) != len(current) {
		t.Errorf("key count = %d, want %d — merge must not drop keys", len(got), len(current))
	}
	if got["subtitle"] != "Unlimited Subscriptions" {
		t.Errorf("subtitle was clobbered: %v", got["subtitle"])
	}
	if !reflect.DeepEqual(got["reviews"], current["reviews"]) {
		t.Errorf("reviews was clobbered: %v", got["reviews"])
	}
}

func TestMergeMetadataNullDeletesKey(t *testing.T) {
	current := map[string]interface{}{"title": "T", "footer": "F"}
	got := mergeMetadata(current, map[string]interface{}{"footer": nil})

	if _, ok := got["footer"]; ok {
		t.Error("a null value must delete the key")
	}
	if got["title"] != "T" {
		t.Error("deleting one key must not affect others")
	}
}

func TestMergeMetadataAddsNewKeys(t *testing.T) {
	got := mergeMetadata(
		map[string]interface{}{"title": "T"},
		map[string]interface{}{"headerImageURL": "https://x"},
	)
	if got["headerImageURL"] != "https://x" {
		t.Error("merge must add keys that are not already present")
	}
	if got["title"] != "T" {
		t.Error("merge must keep existing keys")
	}
}

func TestMergeMetadataDoesNotMutateInput(t *testing.T) {
	current := map[string]interface{}{"title": "original"}
	mergeMetadata(current, map[string]interface{}{"title": "changed", "extra": 1})

	if current["title"] != "original" {
		t.Error("merge must not mutate the caller's map")
	}
	if _, ok := current["extra"]; ok {
		t.Error("merge must not add keys to the caller's map")
	}
}

func TestMergeMetadataNestedObjectReplacesWholesale(t *testing.T) {
	// Documented behaviour: the merge is one level deep.
	current := map[string]interface{}{
		"colorScheme": map[string]interface{}{"dark": "#000", "light": "#fff"},
	}
	got := mergeMetadata(current, map[string]interface{}{
		"colorScheme": map[string]interface{}{"dark": "#111"},
	})
	cs := got["colorScheme"].(map[string]interface{})
	if len(cs) != 1 {
		t.Errorf("nested object should be replaced wholesale, got %v", cs)
	}
}

func TestMergeMetadataEmptyCurrent(t *testing.T) {
	got := mergeMetadata(map[string]interface{}{}, map[string]interface{}{"title": "T"})
	if got["title"] != "T" {
		t.Error("merging into empty metadata should just set the key")
	}
}

// If every key is deleted we must still send `metadata: {}` rather than omitting
// the field, which the caller distinguishes by nil-ness.
func TestMergeMetadataNeverReturnsNil(t *testing.T) {
	got := mergeMetadata(map[string]interface{}{"a": 1}, map[string]interface{}{"a": nil})
	if got == nil {
		t.Fatal("merge must return a non-nil map so the field is still sent")
	}
	if len(got) != 0 {
		t.Errorf("expected an empty map, got %v", got)
	}
}
