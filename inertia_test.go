package gonertia

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestInertia_Location(t *testing.T) {
	t.Parallel()

	t.Run("plain redirect with default status", func(t *testing.T) {
		t.Parallel()

		w, r := requestMock(http.MethodGet, "/")

		wantStatus := http.StatusFound
		wantLocation := "/foo"

		I().Location(w, r, wantLocation)

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

func I(opts ...func(i *Inertia)) *Inertia {
	i := &Inertia{
		containerID:  "app",
		marshallJSON: json.Marshal,
	}

	for _, opt := range opts {
		opt(i)
	}

	return i
}
