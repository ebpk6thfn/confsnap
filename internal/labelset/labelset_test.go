package labelset

import (
	"testing"
)

func TestNew_NotNil(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("expected non-nil Manager")
	}
}

func TestSet_And_Get(t *testing.T) {
	m := New()
	if err := m.Set("host1", "env", "prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := m.Get("host1", "env")
	if !ok || v != "prod" {
		t.Errorf("expected prod, got %q ok=%v", v, ok)
	}
}

func TestSet_EmptyKey_ReturnsError(t *testing.T) {
	m := New()
	if err := m.Set("host1", "", "value"); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestGet_MissingEntity_ReturnsFalse(t *testing.T) {
	m := New()
	_, ok := m.Get("ghost", "env")
	if ok {
		t.Error("expected not found for missing entity")
	}
}

func TestGet_MissingKey_ReturnsFalse(t *testing.T) {
	m := New()
	_ = m.Set("host1", "env", "prod")
	_, ok := m.Get("host1", "region")
	if ok {
		t.Error("expected not found for missing key")
	}
}

func TestLabels_ReturnsCopy(t *testing.T) {
	m := New()
	_ = m.Set("host1", "env", "prod")
	_ = m.Set("host1", "region", "us-east")
	ls := m.Labels("host1")
	if ls["env"] != "prod" || ls["region"] != "us-east" {
		t.Errorf("unexpected labels: %v", ls)
	}
	// mutating copy should not affect manager
	ls["env"] = "staging"
	v, _ := m.Get("host1", "env")
	if v != "prod" {
		t.Error("copy mutation affected manager state")
	}
}

func TestRemove_DeletesKey(t *testing.T) {
	m := New()
	_ = m.Set("host1", "env", "prod")
	m.Remove("host1", "env")
	_, ok := m.Get("host1", "env")
	if ok {
		t.Error("expected key to be removed")
	}
}

func TestMatch_ReturnsMatchingEntities(t *testing.T) {
	m := New()
	_ = m.Set("host1", "env", "prod")
	_ = m.Set("host1", "region", "us-east")
	_ = m.Set("host2", "env", "prod")
	_ = m.Set("host2", "region", "eu-west")
	_ = m.Set("host3", "env", "staging")

	results := m.Match(LabelSet{"env": "prod"})
	if len(results) != 2 {
		t.Fatalf("expected 2 matches, got %d: %v", len(results), results)
	}
	if results[0] != "host1" || results[1] != "host2" {
		t.Errorf("unexpected match order: %v", results)
	}
}

func TestMatch_MultipleSelectors_NarrowsResults(t *testing.T) {
	m := New()
	_ = m.Set("host1", "env", "prod")
	_ = m.Set("host1", "region", "us-east")
	_ = m.Set("host2", "env", "prod")
	_ = m.Set("host2", "region", "eu-west")

	results := m.Match(LabelSet{"env": "prod", "region": "eu-west"})
	if len(results) != 1 || results[0] != "host2" {
		t.Errorf("expected [host2], got %v", results)
	}
}

func TestMatch_NoSelector_ReturnsAll(t *testing.T) {
	m := New()
	_ = m.Set("a", "x", "1")
	_ = m.Set("b", "y", "2")
	results := m.Match(LabelSet{})
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}
