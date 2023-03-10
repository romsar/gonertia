package gonertia

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"
)

func TestInertia_Middleware(t *testing.T) {
	t.Parallel()

	t.Run("plain request", func(t *testing.T) {
		t.Parallel()

		t.Run("do nothing, call next handler", func(t *testing.T) {
			t.Parallel()

			w, r := requestMock(http.MethodGet, "/")

			I().Middleware(assertHandlerServed(t)).ServeHTTP(w, r)

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

				i.Middleware(assertHandlerServed(t, successJSONHandler)).ServeHTTP(w, r)

				assertInertiaVary(t, w)
				assertResponseStatusCode(t, w, http.StatusConflict)
				assertInertiaLocation(t, w, "/home")
			})

			t.Run("diff version with POST, do nothing", func(t *testing.T) {
				t.Parallel()

				i := I(func(i *Inertia) {
					i.version = "foo"
				})

				w, r := requestMock(http.MethodPost, "/home")
				asInertiaRequest(r)
				withInertiaVersion(r, "bar")

				i.Middleware(assertHandlerServed(t, successJSONHandler)).ServeHTTP(w, r)

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

				I().Middleware(assertHandlerServed(t)).ServeHTTP(w, r)

				assertInertiaVary(t, w)
				assertResponseStatusCode(t, w, http.StatusConflict)
				assertInertiaLocation(t, w, "/foo")
			})

			t.Run("don't redirect back if empty request and status not ok", func(t *testing.T) {
				t.Parallel()

				w, r := requestMock(http.MethodGet, "/")
				asInertiaRequest(r)
				withReferer(r, "/foo")

				I().Middleware(assertHandlerServed(t, errorJSONHandler)).ServeHTTP(w, r)

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

				I().Middleware(assertHandlerServed(t, setStatusHandler(http.StatusFound))).ServeHTTP(w, r)

				assertInertiaVary(t, w)
				assertResponseStatusCode(t, w, http.StatusFound)
			})

			for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch} {
				method := method

				t.Run(method+" cannot have 302 status", func(t *testing.T) {
					t.Parallel()

					w, r := requestMock(method, "/")
					asInertiaRequest(r)

					I().Middleware(assertHandlerServed(t, setStatusHandler(http.StatusFound))).ServeHTTP(w, r)

					assertInertiaVary(t, w)
					assertResponseStatusCode(t, w, http.StatusSeeOther)
				})
			}
		})

		t.Run("success", func(t *testing.T) {
			t.Parallel()

			t.Run("with new response writer", func(t *testing.T) {
				t.Parallel()

				w, r := requestMock(http.MethodGet, "/")
				asInertiaRequest(r)

				handlers := []http.HandlerFunc{
					successJSONHandler,
					setHeadersHandler(map[string]string{
						"foo": "bar",
					}),
				}

				I().Middleware(assertHandlerServed(t, handlers...)).ServeHTTP(w, r)

				assertInertiaVary(t, w)
				assertResponseStatusCode(t, w, http.StatusOK)

				if !reflect.DeepEqual(w.Body.String(), successJSON) {
					t.Fatalf("JSON=%#v, want=%#v", w.Body.String(), successJSON)
				}

				gotHeader := w.Header().Get("foo")
				wantHeader := "bar"

				if gotHeader != wantHeader {
					t.Fatalf("header=%#v, want=%#v", gotHeader, wantHeader)
				}
			})

			t.Run("with passed response writer", func(t *testing.T) {
				t.Parallel()

				w, r := requestMock(http.MethodGet, "/")
				asInertiaRequest(r)

				buf := bytes.NewBufferString(successJSON)

				i := I()

				wrap := &inertiaResponseWrapper{
					statusCode: http.StatusNotFound,
					buf:        buf,
					header:     http.Header{"foo": []string{"bar"}},
				}

				I().Middleware(assertHandlerServed(t, successJSONHandler)).ServeHTTP(wrap, r)
				i.copyWrapperResponse(w, wrap)

				assertInertiaVary(t, w)
				assertResponseStatusCode(t, w, http.StatusNotFound)

				if !reflect.DeepEqual(w.Body.String(), successJSON+successJSON) {
					t.Fatalf("JSON=%#v, want=%#v", w.Body.String(), successJSON)
				}

				gotHeader := w.Header().Get("foo")
				wantHeader := "bar"

				if gotHeader != wantHeader {
					t.Fatalf("header=%#v, want=%#v", gotHeader, wantHeader)
				}
			})
		})
	})
}

var (
	successJSON = `{"success": true}`
	errorJSON   = `{"success": false}`
)

func successJSONHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(successJSON))
}

func errorJSONHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(errorJSON))
	w.WriteHeader(http.StatusBadRequest)
}

func setStatusHandler(status int) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(status)
	}
}

func setHeadersHandler(headers map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		for key, val := range headers {
			w.Header().Set(key, val)
		}
	}
}
