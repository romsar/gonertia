package gonertia

import (
	"net/http"
	"testing"
)

func TestInertia_Location(t *testing.T) {
	t.Parallel()

	t.Run("plain redirect with default status", func(t *testing.T) {
		t.Parallel()

		w, r := requestMock(http.MethodGet, "/")

		i := I()

		wantStatus := http.StatusFound
		wantLocation := "/foo"

		i.Location(w, r, wantLocation)

		assertResponseStatusCode(t, w, wantStatus)
		assertLocation(t, w, wantLocation)
	})

	t.Run("plain redirect with specified status", func(t *testing.T) {
		t.Parallel()

		w, r := requestMock(http.MethodGet, "/")

		wantStatus := http.StatusMovedPermanently
		wantLocation := "/foo"

		I().Location(w, r, wantLocation, wantStatus)

		assertResponseStatusCode(t, w, wantStatus)
		assertLocation(t, w, wantLocation)
	})

	t.Run("inertia location", func(t *testing.T) {
		t.Parallel()

		w, r := requestMock(http.MethodGet, "/")
		asInertiaRequest(r)

		wantLocation := ""
		wantInertiaLocation := "/foo"

		I().Location(w, r, wantInertiaLocation, http.StatusMovedPermanently)

		assertLocation(t, w, wantLocation)
		assertResponseStatusCode(t, w, http.StatusConflict)
		assertInertiaLocation(t, w, wantInertiaLocation)
	})
}

func TestInertia_Back(t *testing.T) {
	t.Parallel()

	t.Run("plain redirect with default status", func(t *testing.T) {
		t.Parallel()

		wantStatus := http.StatusFound
		wantLocation := "https://example.com/foo"

		w, r := requestMock(http.MethodGet, "/")
		r.Header.Set("Referer", wantLocation)

		i := I()

		i.Back(w, r)

		assertResponseStatusCode(t, w, wantStatus)
		assertLocation(t, w, wantLocation)
	})

	t.Run("plain redirect with specified status", func(t *testing.T) {
		t.Parallel()

		wantStatus := http.StatusMovedPermanently
		wantLocation := "https://example.com/foo"

		w, r := requestMock(http.MethodGet, "/")
		r.Header.Set("Referer", wantLocation)

		I().Location(w, r, wantLocation, wantStatus)

		assertResponseStatusCode(t, w, wantStatus)
		assertLocation(t, w, wantLocation)
	})

	t.Run("inertia location", func(t *testing.T) {
		t.Parallel()

		wantLocation := ""
		wantInertiaLocation := "https://example.com/foo"

		w, r := requestMock(http.MethodGet, "/")
		r.Header.Set("Referer", wantLocation)
		asInertiaRequest(r)

		I().Location(w, r, wantInertiaLocation, http.StatusMovedPermanently)

		assertLocation(t, w, wantLocation)
		assertResponseStatusCode(t, w, http.StatusConflict)
		assertInertiaLocation(t, w, wantInertiaLocation)
	})
}
