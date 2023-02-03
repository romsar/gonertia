package gonertia

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func asInertiaRequest(r *http.Request) {
	r.Header.Set("X-Inertia", "true")
}

func withInertiaVersion(r *http.Request, ver string) {
	r.Header.Set("X-Inertia-Version", ver)
}

func withReferer(r *http.Request, referer string) {
	r.Header.Set("Referer", referer)
}

func requestMock(method, target string) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, target, nil)

	return w, r
}

func assertResponseStatusCode(t *testing.T, w *httptest.ResponseRecorder, want int) {
	t.Helper()

	if w.Result().StatusCode != want {
		t.Fatalf("status=%d, want=%d", w.Result().StatusCode, want)
	}
}

func assertHeader(t *testing.T, w *httptest.ResponseRecorder, key, want string) {
	t.Helper()

	got := w.Result().Header.Get(key)

	if got != want {
		t.Fatalf("header=%s, want=%s", got, want)
	}
}

func assertLocation(t *testing.T, w *httptest.ResponseRecorder, want string) {
	t.Helper()

	assertHeader(t, w, "Location", want)
}

func assertInertiaLocation(t *testing.T, w *httptest.ResponseRecorder, want string) {
	t.Helper()

	assertHeader(t, w, "X-Inertia-Location", want)
}

func assertInertiaVary(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()

	gotVary := w.Result().Header.Get("Vary")
	wantVary := "X-Inertia"

	if gotVary != wantVary {
		t.Fatalf("vary=%s, want=%s", gotVary, wantVary)
	}
}

func assertNextHandlerServed(t *testing.T, handlers ...http.HandlerFunc) http.HandlerFunc {
	t.Helper()

	called := false

	t.Cleanup(func() {
		if !called {
			t.Fatal("next handler was not called")
		}
	})

	return func(w http.ResponseWriter, r *http.Request) {
		for _, handler := range handlers {
			handler(w, r)
		}

		called = true
	}
}
