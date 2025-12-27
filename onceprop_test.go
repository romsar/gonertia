package gonertia

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOncePropFirstLoad(t *testing.T) {
	t.Parallel()
	i, err := New(`<div id="app"></div>`)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	props := Props{
		"foo": Once("bar"),
		"baz": "qux",
	}
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("X-Inertia", "true")
	w := httptest.NewRecorder()
	err = i.Render(w, r, "TestComponent", props)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !containsJSON(w.Body.String(), `"foo":"bar"`) {
		t.Errorf("expected OnceProp to be present on first load")
	}
	if !containsJSON(w.Body.String(), `"baz":"qux"`) {
		t.Errorf("expected normal prop to be present on first load")
	}
}

func TestOncePropOmittedOnPartialReloadUnlessRequested(t *testing.T) {
	t.Parallel()
	i, err := New(`<div id="app"></div>`)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	props := Props{
		"foo": Once("bar"),
		"baz": "qux",
	}
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("X-Inertia", "true")
	r.Header.Set("X-Inertia-Partial-Component", "TestComponent")
	r.Header.Set("X-Inertia-Partial-Data", "baz")
	w := httptest.NewRecorder()
	err = i.Render(w, r, "TestComponent", props)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if containsJSON(w.Body.String(), `"foo":"bar"`) {
		t.Errorf("expected OnceProp to be omitted on partial reload if not requested")
	}
	if !containsJSON(w.Body.String(), `"baz":"qux"`) {
		t.Errorf("expected normal prop to be present on partial reload")
	}
}

func TestOncePropIncludedOnPartialReloadIfRequested(t *testing.T) {
	t.Parallel()
	i, err := New(`<div id="app"></div>`)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	props := Props{
		"foo": Once("bar"),
		"baz": "qux",
	}
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("X-Inertia", "true")
	r.Header.Set("X-Inertia-Partial-Component", "TestComponent")
	r.Header.Set("X-Inertia-Partial-Data", "foo,baz")
	w := httptest.NewRecorder()
	err = i.Render(w, r, "TestComponent", props)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !containsJSON(w.Body.String(), `"foo":"bar"`) {
		t.Errorf("expected OnceProp to be present when explicitly requested on partial reload")
	}
	if !containsJSON(w.Body.String(), `"baz":"qux"`) {
		t.Errorf("expected normal prop to be present on partial reload")
	}
}

func TestLazyOptionalPropBehavior(t *testing.T) {
	t.Parallel()
	i, err := New(`<div id="app"></div>`)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	props := Props{
		"foo": Once(func() (any, error) {
			return "bar", nil
		}),
		"baz": "qux",
	}
	// First load: should Include both props
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("X-Inertia", "true")
	w := httptest.NewRecorder()
	err = i.Render(w, r, "TestComponent", props)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !containsJSON(w.Body.String(), `"foo":"bar"`) {
		t.Errorf("expected OptionalProp to be present on first load")
	}
	if !containsJSON(w.Body.String(), `"baz":"qux"`) {
		t.Errorf("expected normal prop to be present on first load")
	}
	// partial reload: should NOT include foo if not requested
	r1 := httptest.NewRequest("GET", "/", nil)
	r1.Header.Set("X-Inertia", "true")
	r1.Header.Set("X-Inertia-Partial-Component", "TestComponent")
	r1.Header.Set("X-Inertia-Partial-Data", "baz")
	w1 := httptest.NewRecorder()
	err = i.Render(w1, r1, "TestComponent", props)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if containsJSON(w1.Body.String(), `"foo":"bar"`) {
		t.Errorf("expected OptionalProp to be omitted on partial reload when not requested")
	}
	if !containsJSON(w1.Body.String(), `"baz":"qux"`) {
		t.Errorf("expected normal prop to be present on partial reload")
	}
	// Partial reload: should include foo if requested
	r2 := httptest.NewRequest("GET", "/", nil)
	r2.Header.Set("X-Inertia", "true")
	r2.Header.Set("X-Inertia-Partial-Component", "TestComponent")
	r2.Header.Set("X-Inertia-Partial-Data", "foo,baz")
	w2 := httptest.NewRecorder()
	err = i.Render(w2, r2, "TestComponent", props)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !containsJSON(w2.Body.String(), `"foo":"bar"`) {
		t.Errorf("expected OptionalProp to be present on partial reload when requested")
	}
	if !containsJSON(w2.Body.String(), `"baz":"qux"`) {
		t.Errorf("expected normal prop to be present on partial reload")
	}
}

// containsJSON is a helper to check if a substring is present in a JSON string (not strict, but good enough for this test).
func containsJSON(s, sub string) bool {
	return strings.Contains(s, sub)
}
