package gonertia

import (
	"net/http"
	"testing"
)

func TestInertia_Middleware(t *testing.T) {
	t.Parallel()

	t.Run("plain request", func(t *testing.T) {
		t.Parallel()

		t.Run("do nothing, call next handler", func(t *testing.T) {
			t.Parallel()

			w, r := requestMock(http.MethodGet, "/")

			I().Middleware(assertNextHandlerServed(t)).ServeHTTP(w, r)

			assertInertiaVary(t, w)
			assertResponseStatusCode(t, w, http.StatusOK)
		})
	})

	t.Run("inertia request", func(t *testing.T) {
		t.Parallel()

		t.Run("assert versioning", func(t *testing.T) {
			t.Parallel()

			t.Run("diff version with GET, should change location with 409", func(t *testing.T) {
				t.Parallel()

				i := I(func(i *Inertia) {
					i.version = "foo"
				})

				w, r := requestMock(http.MethodGet, "https://example.com/home")
				asInertiaRequest(r)
				withInertiaVersion(r, "bar")

				i.Middleware(assertNextHandlerServed(t, successJSONHandler)).ServeHTTP(w, r)

				assertInertiaVary(t, w)
				assertResponseStatusCode(t, w, http.StatusConflict)
				assertInertiaLocation(t, w, "https://example.com/home")
			})

			t.Run("diff version with POST, do nothing", func(t *testing.T) {
				t.Parallel()

				i := I(func(i *Inertia) {
					i.version = "foo"
				})

				w, r := requestMock(http.MethodPost, "https://example.com/home")
				asInertiaRequest(r)
				withInertiaVersion(r, "bar")

				i.Middleware(assertNextHandlerServed(t, successJSONHandler)).ServeHTTP(w, r)

				assertInertiaVary(t, w)
				assertResponseStatusCode(t, w, http.StatusOK)
			})
		})

		t.Run("redirect back if empty response body", func(t *testing.T) {
			t.Parallel()

			t.Run("redirect back if empty request and status ok", func(t *testing.T) {
				t.Parallel()

				w, r := requestMock(http.MethodGet, "/")
				asInertiaRequest(r)
				withReferer(r, "/foo")

				I().Middleware(assertNextHandlerServed(t)).ServeHTTP(w, r)

				assertInertiaVary(t, w)
				assertResponseStatusCode(t, w, http.StatusConflict)
				assertInertiaLocation(t, w, "/foo")
			})

			t.Run("don't redirect back if empty request and status not ok", func(t *testing.T) {
				t.Parallel()

				w, r := requestMock(http.MethodGet, "/")
				asInertiaRequest(r)
				withReferer(r, "/foo")

				I().Middleware(assertNextHandlerServed(t, errorJSONHandler)).ServeHTTP(w, r)

				assertInertiaVary(t, w)
				assertResponseStatusCode(t, w, http.StatusBadRequest)
				assertInertiaLocation(t, w, "")
			})
		})

		t.Run("POST, PUT and PATCH requests cannot have the status 302", func(t *testing.T) {
			t.Parallel()

			t.Run("GET can have 302 status", func(t *testing.T) {
				t.Parallel()

				w, r := requestMock(http.MethodGet, "/")
				asInertiaRequest(r)

				I().Middleware(assertNextHandlerServed(t, redirectHandler(http.StatusFound))).ServeHTTP(w, r)

				assertInertiaVary(t, w)
				assertResponseStatusCode(t, w, http.StatusFound)
			})

			for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch} {
				method := method

				t.Run(method+" cannot have 302 status", func(t *testing.T) {
					t.Parallel()

					w, r := requestMock(method, "/")
					asInertiaRequest(r)

					I().Middleware(assertNextHandlerServed(t, redirectHandler(http.StatusFound))).ServeHTTP(w, r)

					assertInertiaVary(t, w)
					assertResponseStatusCode(t, w, http.StatusSeeOther)
				})
			}
		})
	})
}

var (
	successJSON = `{"success": true}`
	errorJSON   = `{"success": false}`
)

func successJSONHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte(successJSON))
}

func errorJSONHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte(errorJSON))
	w.WriteHeader(http.StatusBadRequest)
}

func redirectHandler(status int) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(status)
	}
}
