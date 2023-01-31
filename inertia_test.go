package gonertia

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInertia_Location(t *testing.T) {
	t.Parallel()

	t.Run("plain redirect with default status", func(t *testing.T) {
		t.Parallel()

		i := new(Inertia)

		w, r := requestMock(http.MethodGet, "/")

		wantStatus := http.StatusFound
		wantLocation := "/foo"

		i.Location(w, r, wantLocation)

		resp := w.Result()

		if resp.StatusCode != wantStatus {
			t.Fatalf("got=%#v, props=%#v", resp.StatusCode, wantStatus)
		}

		gotLocation := mustLocation(t, resp)

		if gotLocation != wantLocation {
			t.Fatalf("got=%#v, props=%#v", gotLocation, wantLocation)
		}
	})

	t.Run("plain redirect with specified status", func(t *testing.T) {
		t.Parallel()

		i := new(Inertia)

		w, r := requestMock(http.MethodGet, "/")

		wantStatus := http.StatusMovedPermanently
		wantLocation := "/foo"

		i.Location(w, r, wantLocation, wantStatus)

		resp := w.Result()

		if resp.StatusCode != wantStatus {
			t.Fatalf("got=%#v, props=%#v", resp.StatusCode, wantStatus)
		}

		gotLocation := mustLocation(t, resp)

		if gotLocation != wantLocation {
			t.Fatalf("got=%#v, props=%#v", gotLocation, wantLocation)
		}
	})

	t.Run("inertia location", func(t *testing.T) {
		t.Parallel()

		i := new(Inertia)

		w, r := requestMock(http.MethodGet, "/")
		asInertiaRequest(r)

		wantStatus := http.StatusConflict
		wantLocation := ""
		wantInertiaLocation := "/foo"

		i.Location(w, r, wantInertiaLocation, http.StatusMovedPermanently)

		resp := w.Result()

		if resp.StatusCode != wantStatus {
			t.Fatalf("got=%#v, props=%#v", resp.StatusCode, wantStatus)
		}

		gotLocation := resp.Header.Get("Location")

		if gotLocation != wantLocation {
			t.Fatalf("got=%#v, props=%#v", gotLocation, wantLocation)
		}

		gotInertiaLocation := resp.Header.Get("X-Inertia-Location")

		if gotInertiaLocation != wantInertiaLocation {
			t.Fatalf("got=%#v, props=%#v", gotInertiaLocation, wantInertiaLocation)
		}
	})
}

func requestMock(method, target string) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, target, nil)

	return w, r
}

func asInertiaRequest(r *http.Request) {
	r.Header.Set("X-Inertia", "true")
}

func mustLocation(t *testing.T, resp *http.Response) string {
	t.Helper()

	location, err := resp.Location()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	return location.String()
}
