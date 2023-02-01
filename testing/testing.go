package testing

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func AssertInertia(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()

	got := w.Result().Header.Get("X-Inertia")

	if got != "true" {
		t.Fatalf("not intertia request")
	}
}

func AssertInertiaLocation(t *testing.T, w *httptest.ResponseRecorder, want string) {
	t.Helper()

	if w.Result().StatusCode != http.StatusConflict {
		t.Fatalf("got status=%#v, want status=%#v", w.Result().StatusCode, http.StatusConflict)
	}

	got := w.Result().Header.Get("X-Inertia-Location")
	if got != want {
		t.Fatalf("got url=%#v, want url=%#v", got, want)
	}
}
