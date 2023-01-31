package gonertia

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsInertiaRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		header http.Header
		want   bool
	}{
		{
			"positive",
			http.Header{"X-Inertia": []string{"foo"}},
			true,
		},
		{
			"negative",
			http.Header{},
			false,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest("GET", "/", nil)
			r.Header = tt.header

			got := IsInertiaRequest(r)

			if got != tt.want {
				t.Fatalf("got=%#v, want=%#v", got, tt.want)
			}
		})
	}
}

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
		t.Fatalf("got=%#v, want=%#v", w.Result().StatusCode, want)
	}
}

func assertHeader(t *testing.T, w *httptest.ResponseRecorder, key, want string) {
	t.Helper()

	got := w.Result().Header.Get(key)

	if got != want {
		t.Fatalf("got=%#v, want=%#v", got, want)
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
		t.Fatalf("got=%#v, want=%#v", gotVary, wantVary)
	}
}
