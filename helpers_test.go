package gonertia

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func I(opts ...func(i *Inertia)) *Inertia {
	i := &Inertia{
		containerID:        "app",
		marshallJSON:       json.Marshal,
		sharedProps:        make(Props),
		sharedTemplateData: make(TemplateData),
	}

	for _, opt := range opts {
		opt(i)
	}

	return i
}

func requestMock(method, target string) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, target, nil)

	return w, r
}

func asInertiaRequest(r *http.Request) {
	r.Header.Set("X-Inertia", "true")
}

func withPartialData(r *http.Request, data []string) {
	r.Header.Set("X-Inertia-Partial-Data", strings.Join(data, ","))
}

func withPartialComponent(r *http.Request, component string) {
	r.Header.Set("X-Inertia-Partial-Component", component)
}

func withInertiaVersion(r *http.Request, ver string) {
	r.Header.Set("X-Inertia-Version", ver)
}

func withReferer(r *http.Request, referer string) {
	r.Header.Set("Referer", referer)
}

func assertResponseStatusCode(t *testing.T, w *httptest.ResponseRecorder, want int) {
	t.Helper()

	if w.Code != want {
		t.Fatalf("status=%d, want=%d", w.Code, want)
	}
}

func assertHeader(t *testing.T, w *httptest.ResponseRecorder, key, want string) {
	t.Helper()

	if got := w.Header().Get(key); got != want {
		t.Fatalf("header=%s, want=%s", got, want)
	}
}

func assertHeaderMissing(t *testing.T, w *httptest.ResponseRecorder, key string) {
	t.Helper()

	if got := w.Header().Get(key); got != "" {
		t.Fatalf("header=%s, want=%s", got, "")
	}
}

func assertLocation(t *testing.T, w *httptest.ResponseRecorder, want string) {
	t.Helper()

	assertHeader(t, w, "Location", want)
}

func assertInertiaResponse(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()

	assertHeader(t, w, "X-Inertia", "true")
}

func assertNotInertiaResponse(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()

	assertHeaderMissing(t, w, "X-Inertia")
}

func assertInertiaLocation(t *testing.T, w *httptest.ResponseRecorder, want string) {
	t.Helper()

	assertHeader(t, w, "X-Inertia-Location", want)
}

func assertJSONResponse(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()

	assertHeader(t, w, "Content-Type", "application/json")
}

func assertHTMLResponse(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()

	assertHeader(t, w, "Content-Type", "text/html")
}

func assertInertiaVary(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()

	gotVary := w.Header().Get("Vary")
	wantVary := "X-Inertia"

	if gotVary != wantVary {
		t.Fatalf("vary=%s, want=%s", gotVary, wantVary)
	}
}

func assertHandlerServed(t *testing.T, handlers ...http.HandlerFunc) http.HandlerFunc {
	t.Helper()

	called := false

	t.Cleanup(func() {
		if !called {
			t.Fatal("handler was not called")
		}
	})

	return func(w http.ResponseWriter, r *http.Request) {
		for _, handler := range handlers {
			handler(w, r)
		}

		called = true
	}
}

func tmpFile(t *testing.T, content string) *os.File {
	t.Helper()

	f, err := os.CreateTemp("", "gonertia")
	if err != nil {
		t.Fatalf("unexpected error: %#v", err)
	}

	closed := false

	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("unexpected error: %#v", err)
	}

	if err := f.Close(); err != nil {
		t.Fatalf("unexpected error: %#v", err)
	}

	closed = true

	t.Cleanup(func() {
		if !closed {
			if err := f.Close(); err != nil {
				t.Fatalf("unexpected error: %#v", err)
			}
		}

		if err := os.Remove(f.Name()); err != nil {
			t.Fatalf("unexpected error: %#v", err)
		}
	})

	return f
}
