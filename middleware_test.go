package gonertia

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"
)

//nolint:gocognit
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

		t.Run("flash", func(t *testing.T) {
			t.Parallel()

			t.Run("validation errors", func(t *testing.T) {
				t.Parallel()

				w, r := requestMock(http.MethodGet, "/")

				want := ValidationErrors{
					"foo": "baz",
					"baz": "quz",
				}

				flashProvider := &flashProviderMock{
					errors: want,
				}

				i := I(func(i *Inertia) {
					i.flash = flashProvider
				})

				var got ValidationErrors
				i.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					got = ValidationErrorsFromContext(r.Context())
				})).ServeHTTP(w, r)

				if !reflect.DeepEqual(got, want) {
					t.Fatalf("validation errors=%#v, want=%#v", got, want)
				}
			})

			t.Run("clear history", func(t *testing.T) {
				t.Parallel()

				w, r := requestMock(http.MethodGet, "/")

				flashProvider := &flashProviderMock{
					clearHistory: true,
				}

				i := I(func(i *Inertia) {
					i.flash = flashProvider
				})

				var got bool
				i.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					got = ClearHistoryFromContext(r.Context())
				})).ServeHTTP(w, r)

				if !got {
					t.Fatalf("clear history=%v, want=true", got)
				}
			})
		})
	})

	t.Run("inertia request", func(t *testing.T) {
		t.Parallel()

		t.Run("assert versioning", func(t *testing.T) {
			t.Parallel()

			t.Run("diff version with GET, should change location with 409 and flash errors", func(t *testing.T) {
				t.Parallel()

				errors := ValidationErrors{
					"foo": "baz",
					"baz": "quz",
				}

				flashProvider := &flashProviderMock{
					errors: errors,
				}

				i := I(func(i *Inertia) {
					i.version = "foo"
					i.flash = flashProvider
				})

				w, r := requestMock(http.MethodGet, "https://example.com/home")
				asInertiaRequest(r)
				withInertiaVersion(r, "bar")

				i.Middleware(assertHandlerServed(t, setInertiaResponseHandler, successJSONHandler)).ServeHTTP(w, r)

				assertInertiaNotVary(t, w)
				assertNotInertiaResponse(t, w)
				assertResponseStatusCode(t, w, http.StatusConflict)
				assertInertiaLocation(t, w, "/home")

				if !reflect.DeepEqual(flashProvider.errors, errors) {
					t.Fatalf("got validation errors=%#v, want=%#v", flashProvider.errors, errors)
				}
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
				assertResponseStatusCode(t, w, http.StatusFound)
				assertLocation(t, w, "/foo")
				assertInertiaLocation(t, w, "")
			})

			t.Run("don't redirect back if empty request and status not ok", func(t *testing.T) {
				t.Parallel()

				w, r := requestMock(http.MethodGet, "/")
				asInertiaRequest(r)
				withReferer(r, "/foo")

				I().Middleware(assertHandlerServed(t, errorJSONHandler)).ServeHTTP(w, r)

				assertInertiaVary(t, w)
				assertResponseStatusCode(t, w, http.StatusBadRequest)
				assertLocation(t, w, "")
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

			for _, method := range []string{http.MethodPut, http.MethodPatch, http.MethodDelete} {
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

func setInertiaResponseHandler(w http.ResponseWriter, _ *http.Request) {
	setInertiaInResponse(w)
}
