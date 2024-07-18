package gonertia

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

//nolint:gocognit
func TestInertia_Render(t *testing.T) {
	t.Parallel()

	t.Run("plain request", func(t *testing.T) {
		t.Parallel()

		t.Run("success", func(t *testing.T) {
			t.Parallel()

			i := I(func(i *Inertia) {
				i.rootTemplateHTML = rootTemplate
				i.version = "f8v01xv4h4"
			})

			assertRootTemplateSuccess(t, i)
		})

		t.Run("ssr", func(t *testing.T) {
			t.Parallel()

			t.Run("success", func(t *testing.T) {
				t.Parallel()

				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					reqContentType := r.Header.Get("Content-Type")
					wantContentType := "application/json"
					if reqContentType != wantContentType {
						t.Fatalf("reqest content type=%s, want=%s", reqContentType, wantContentType)
					}

					pageJSON, err := io.ReadAll(r.Body)
					if err != nil {
						t.Fatalf("unexpected error: %s", err)
					}

					assertable := AssertFromBytes(t, pageJSON)
					assertable.AssertComponent("Some/Component")
					assertable.AssertProps(Props{"foo": "bar", "errors": map[string]any{}})
					assertable.AssertVersion("f8v01xv4h4")
					assertable.AssertURL("/home")

					setJSONResponse(w)

					ssr := map[string]any{
						"head": []string{`<title inertia>foo</title>`, `<meta charset="UTF-8">`},
						"body": `<div id="app" data-page="` + template.HTMLEscapeString(string(pageJSON)) + `">foo bar</div>`,
					}

					pageJSON, err = json.Marshal(ssr)
					if err != nil {
						t.Fatalf("unexpected error: %s", err)
					}

					_, err = w.Write(pageJSON)
					if err != nil {
						t.Fatalf("unexpected error: %s", err)
					}
				}))
				defer ts.Close()

				i := I(func(i *Inertia) {
					i.rootTemplateHTML = rootTemplate
					i.version = "f8v01xv4h4"
					i.ssrURL = ts.URL
					i.ssrHTTPClient = ts.Client()
				})

				w, r := requestMock(http.MethodGet, "/home")

				err := i.Render(w, r, "Some/Component", Props{"foo": "bar"})
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}

				var buf bytes.Buffer

				assertable := Assert(t, io.TeeReader(w.Body, &buf))
				assertable.AssertComponent("Some/Component")
				assertable.AssertProps(Props{"foo": "bar", "errors": map[string]any{}})
				assertable.AssertVersion("f8v01xv4h4")
				assertable.AssertURL("/home")

				re := regexp.MustCompile(`<div\sid="app"\sdata-page="[^"]+">([^<]+)</div>`)

				got := re.FindStringSubmatch(buf.String())[1]
				want := "foo bar"
				if got != want {
					t.Fatalf("got content=%s, want=%s", got, want)
				}
			})

			t.Run("error with fallback", func(t *testing.T) {
				t.Parallel()

				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}))
				defer ts.Close()

				i := I(func(i *Inertia) {
					i.rootTemplateHTML = rootTemplate
					i.version = "f8v01xv4h4"
					i.ssrURL = ts.URL
					i.ssrHTTPClient = ts.Client()
				})

				w, r := requestMock(http.MethodGet, "/home")

				err := i.Render(w, r, "Some/Component", Props{"foo": "bar"})
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}

				var buf bytes.Buffer

				assertable := Assert(t, io.TeeReader(w.Body, &buf))
				assertable.AssertComponent("Some/Component")
				assertable.AssertProps(Props{"foo": "bar", "errors": map[string]any{}})
				assertable.AssertVersion("f8v01xv4h4")
				assertable.AssertURL("/home")
			})
		})

		t.Run("shared funcs", func(t *testing.T) {
			t.Parallel()

			w, r := requestMock(http.MethodGet, "/")

			i := I(func(i *Inertia) {
				i.rootTemplateHTML = `{{ trim " foo bar " }}`
				i.sharedTemplateFuncs = TemplateFuncs{
					"trim": strings.TrimSpace,
				}
			})

			err := i.Render(w, r, "Some/Component")
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			got := w.Body.String()
			want := "foo bar"

			if got != want {
				t.Fatalf("got=%s, want=%s", got, want)
			}
		})

		t.Run("shared template data", func(t *testing.T) {
			t.Parallel()

			w, r := requestMock(http.MethodGet, "/")

			i := I(func(i *Inertia) {
				i.rootTemplateHTML = `Hello, {{ .text }}!`
				i.sharedTemplateData = TemplateData{
					"text": "world",
				}
			})

			err := i.Render(w, r, "Some/Component")
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			got := w.Body.String()
			want := "Hello, world!"

			if got != want {
				t.Fatalf("got=%s, want=%s", got, want)
			}
		})
	})

	t.Run("inertia request", func(t *testing.T) {
		t.Parallel()

		t.Run("success", func(t *testing.T) {
			t.Parallel()

			i := I(func(i *Inertia) {
				i.version = "f8v01xv4h4"
			})

			w, r := requestMock(http.MethodGet, "/home")
			asInertiaRequest(r)

			err := i.Render(w, r, "Some/Component", Props{
				"foo": "bar",
			})
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			assertable := AssertFromString(t, w.Body.String())
			assertable.AssertComponent("Some/Component")
			assertable.AssertProps(Props{"foo": "bar", "errors": map[string]any{}})
			assertable.AssertVersion("f8v01xv4h4")
			assertable.AssertURL("/home")

			assertInertiaResponse(t, w)
			assertJSONResponse(t, w)
			assertResponseStatusCode(t, w, http.StatusOK)
		})

		t.Run("props priority", func(t *testing.T) {
			t.Parallel()

			i := I(func(i *Inertia) {
				i.sharedProps = Props{"foo": "bar", "abc": "123", "shared": "prop"}
			})

			w, r := requestMock(http.MethodGet, "/home")
			asInertiaRequest(r)

			ctx := WithProps(r.Context(), Props{"foo": "baz", "abc": "456", "ctx": "prop"})

			err := i.Render(w, r.WithContext(ctx), "Some/Component", Props{
				"foo": "zzz",
			})
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			assertable := AssertFromString(t, w.Body.String())
			assertable.AssertProps(Props{
				"abc":    "456",
				"ctx":    "prop",
				"foo":    "zzz",
				"shared": "prop",
				"errors": map[string]any{},
			})
		})

		t.Run("validation errors", func(t *testing.T) {
			t.Parallel()

			w, r := requestMock(http.MethodGet, "/home")
			asInertiaRequest(r)

			ctx := WithValidationErrors(r.Context(), ValidationErrors{"foo": "bar"})

			err := I().Render(w, r.WithContext(ctx), "Some/Component", Props{
				"abc": "123",
			})
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			assertable := AssertFromString(t, w.Body.String())
			assertable.AssertProps(Props{
				"abc": "123",
				"errors": map[string]any{
					"foo": "bar",
				},
			})
		})

		t.Run("props value resolving", func(t *testing.T) {
			t.Parallel()

			t.Run("reject lazy props", func(t *testing.T) {
				t.Parallel()

				w, r := requestMock(http.MethodGet, "/home")
				asInertiaRequest(r)

				err := I().Render(w, r, "Some/Component", Props{
					"foo":              "bar",
					"closure":          func() any { return "prop" },
					"closure_with_err": func() (any, error) { return "prop", nil },
					"lazy":             LazyProp{func() (any, error) { return "prop", nil }},
				})
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}

				assertable := AssertFromString(t, w.Body.String())
				assertable.AssertProps(Props{
					"foo":              "bar",
					"closure":          "prop",
					"closure_with_err": "prop",
					"errors":           map[string]any{},
				})
			})

			t.Run("only", func(t *testing.T) {
				t.Parallel()

				t.Run("resolve lazy props, same component", func(t *testing.T) {
					t.Parallel()

					w, r := requestMock(http.MethodGet, "/home")
					asInertiaRequest(r)
					withOnly(r, []string{"foo", "closure", "lazy"})
					withPartialComponent(r, "Some/Component")

					err := I().Render(w, r, "Some/Component", Props{
						"foo":     "bar",
						"abc":     "123",
						"closure": func() (any, error) { return "prop", nil },
						"lazy":    LazyProp{func() (any, error) { return "prop", nil }},
						"always":  AlwaysProp{"prop"},
					})
					if err != nil {
						t.Fatalf("unexpected error: %s", err)
					}

					assertable := AssertFromString(t, w.Body.String())
					assertable.AssertProps(Props{
						"foo":     "bar",
						"closure": "prop",
						"lazy":    "prop",
						"always":  "prop",
						"errors":  map[string]interface{}{},
					})
				})

				t.Run("resolve lazy props, other component", func(t *testing.T) {
					t.Parallel()

					w, r := requestMock(http.MethodGet, "/home")
					asInertiaRequest(r)
					withOnly(r, []string{"foo", "closure", "lazy"})
					withPartialComponent(r, "Other/Component")

					err := I().Render(w, r, "Some/Component", Props{
						"foo":     "bar",
						"abc":     "123",
						"closure": func() (any, error) { return "prop", nil },
					})
					if err != nil {
						t.Fatalf("unexpected error: %s", err)
					}

					assertable := AssertFromString(t, w.Body.String())
					assertable.AssertProps(Props{
						"foo":     "bar",
						"abc":     "123",
						"closure": "prop",
						"errors":  map[string]any{},
					})
				})
			})

			t.Run("except", func(t *testing.T) {
				t.Parallel()

				w, r := requestMock(http.MethodGet, "/home")
				asInertiaRequest(r)
				withOnly(r, []string{"foo", "baz"})
				withExcept(r, []string{"foo", "abc", "lazy", "always"})
				withPartialComponent(r, "Some/Component")

				err := I().Render(w, r, "Some/Component", Props{
					"foo":    "bar",
					"baz":    "quz",
					"bez":    "bee",
					"lazy":   LazyProp{func() (any, error) { return "prop", nil }},
					"always": AlwaysProp{"prop"},
				})
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}

				assertable := AssertFromString(t, w.Body.String())
				assertable.AssertProps(Props{
					"baz":    "quz",
					"errors": map[string]any{},
				})
			})

			t.Run("proper interfaces", func(t *testing.T) {
				t.Parallel()

				w, r := requestMock(http.MethodGet, "/home")
				asInertiaRequest(r)
				withPartialComponent(r, "Some/Component")

				err := I().Render(w, r, "Some/Component", Props{
					"proper":     testProper{"prop1"},
					"try_proper": testTryProper{"prop2"},
				})
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}

				assertable := AssertFromString(t, w.Body.String())
				if assertable.Props["proper"] != "prop1" {
					t.Fatalf("resolved prop=%v, want=%v", assertable.Props["proper"], "prop1")
				}
				if assertable.Props["try_proper"] != "prop2" {
					t.Fatalf("try resolved prop=%v, want=%v", assertable.Props["try_proper"], "prop2")
				}
			})
		})
	})
}

type testProper struct {
	Value any
}

func (p testProper) Prop() any {
	return p.Value
}

type testTryProper struct {
	Value any
}

func (p testTryProper) TryProp() (any, error) {
	return p.Value, nil
}

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

	t.Run("flash validation errors", func(t *testing.T) {
		t.Parallel()

		t.Run("plain redirect", func(t *testing.T) {
			t.Parallel()

			w, r := requestMock(http.MethodGet, "/")

			flashProvider := &flashProviderMock{}

			i := I(func(i *Inertia) {
				i.flash = flashProvider
			})

			errors := ValidationErrors{
				"foo": "bar",
				"baz": "quz",
			}

			withValidationErrors(r, errors)
			i.Location(w, r, "/foo")

			if !reflect.DeepEqual(flashProvider.errors, errors) {
				t.Fatalf("got validation errors=%#v, want=%#v", flashProvider.errors, errors)
			}
		})

		t.Run("inertia location", func(t *testing.T) {
			t.Parallel()

			w, r := requestMock(http.MethodGet, "/")
			asInertiaRequest(r)

			flashProvider := &flashProviderMock{}

			i := I(func(i *Inertia) {
				i.flash = flashProvider
			})

			errors := ValidationErrors{
				"foo": "bar",
				"baz": "quz",
			}

			withValidationErrors(r, errors)
			i.Location(w, r, "/foo", http.StatusMovedPermanently)

			if !reflect.DeepEqual(flashProvider.errors, errors) {
				t.Fatalf("got validation errors=%#v, want=%#v", flashProvider.errors, errors)
			}
		})
	})
}

func TestInertia_Redirect(t *testing.T) {
	t.Parallel()

	t.Run("with default status", func(t *testing.T) {
		t.Parallel()

		wantStatus := http.StatusFound
		wantLocation := "https://example.com/foo"

		w, r := requestMock(http.MethodGet, "/")

		i := I()

		i.Redirect(w, r, wantLocation)

		assertResponseStatusCode(t, w, wantStatus)
		assertLocation(t, w, wantLocation)
	})

	t.Run("with specified status", func(t *testing.T) {
		t.Parallel()

		wantStatus := http.StatusMovedPermanently
		wantLocation := "https://example.com/foo"

		w, r := requestMock(http.MethodGet, "/")

		I().Redirect(w, r, wantLocation, wantStatus)

		assertResponseStatusCode(t, w, wantStatus)
		assertLocation(t, w, wantLocation)
	})

	t.Run("inertia request", func(t *testing.T) {
		t.Parallel()

		wantLocation := "https://example.com/foo"
		wantInertiaLocation := ""

		w, r := requestMock(http.MethodGet, "/")
		asInertiaRequest(r)

		I().Redirect(w, r, wantLocation, http.StatusMovedPermanently)

		assertLocation(t, w, wantLocation)
		assertResponseStatusCode(t, w, http.StatusMovedPermanently)
		assertInertiaLocation(t, w, wantInertiaLocation)
	})

	t.Run("flash validation errors", func(t *testing.T) {
		t.Parallel()

		w, r := requestMock(http.MethodGet, "/")

		flashProvider := &flashProviderMock{}

		i := I(func(i *Inertia) {
			i.flash = flashProvider
		})

		errors := ValidationErrors{
			"foo": "bar",
			"baz": "quz",
		}

		withValidationErrors(r, errors)
		i.Redirect(w, r, "https://example.com/foo")

		if !reflect.DeepEqual(flashProvider.errors, errors) {
			t.Fatalf("got validation errors=%#v, want=%#v", flashProvider.errors, errors)
		}
	})
}

func TestInertia_Back(t *testing.T) {
	t.Parallel()

	t.Run("with default status", func(t *testing.T) {
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

	t.Run("with specified status", func(t *testing.T) {
		t.Parallel()

		wantStatus := http.StatusMovedPermanently
		wantLocation := "https://example.com/foo"

		w, r := requestMock(http.MethodGet, "/")
		r.Header.Set("Referer", wantLocation)

		I().Back(w, r, wantStatus)

		assertResponseStatusCode(t, w, wantStatus)
		assertLocation(t, w, wantLocation)
	})

	t.Run("inertia request", func(t *testing.T) {
		t.Parallel()

		wantLocation := "https://example.com/foo"
		wantInertiaLocation := ""

		w, r := requestMock(http.MethodGet, "/")
		r.Header.Set("Referer", wantLocation)
		asInertiaRequest(r)

		I().Back(w, r, http.StatusMovedPermanently)

		assertLocation(t, w, wantLocation)
		assertResponseStatusCode(t, w, http.StatusMovedPermanently)
		assertInertiaLocation(t, w, wantInertiaLocation)
	})

	t.Run("flash validation errors", func(t *testing.T) {
		t.Parallel()

		w, r := requestMock(http.MethodGet, "/")
		r.Header.Set("Referer", "https://example.com/foo")

		flashProvider := &flashProviderMock{}

		i := I(func(i *Inertia) {
			i.flash = flashProvider
		})

		errors := ValidationErrors{
			"foo": "bar",
			"baz": "quz",
		}

		withValidationErrors(r, errors)
		i.Back(w, r)

		if !reflect.DeepEqual(flashProvider.errors, errors) {
			t.Fatalf("got validation errors=%#v, want=%#v", flashProvider.errors, errors)
		}
	})
}

func assertRootTemplateSuccess(t *testing.T, i *Inertia) {
	t.Helper()

	w, r := requestMock(http.MethodGet, "/home")

	err := i.Render(w, r, "Some/Component", Props{
		"foo": "bar",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	assertable := Assert(t, w.Body)
	assertable.AssertComponent("Some/Component")
	assertable.AssertProps(Props{"foo": "bar", "errors": map[string]any{}})
	assertable.AssertVersion("f8v01xv4h4")
	assertable.AssertURL("/home")

	assertNotInertiaResponse(t, w)
	assertHTMLResponse(t, w)
	assertResponseStatusCode(t, w, http.StatusOK)
}
