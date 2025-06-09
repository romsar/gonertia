package gonertia

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"
)

//nolint:gocognit,gocyclo
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

		t.Run("success with pre-parsed template", func(t *testing.T) {
			t.Parallel()

			tmpl, err := template.New("root").
				Funcs(template.FuncMap(make(TemplateFuncs))).
				Parse(rootTemplate)
			if err != nil {
				t.Fatalf("parse root template: %v", err)
			}

			i := I(func(i *Inertia) {
				i.rootTemplate = tmpl
				i.version = "f8v01xv4h4"
			})

			assertRootTemplateSuccess(t, i)
		})

		t.Run("ssr", func(t *testing.T) {
			t.Parallel()

			newTestServerSSR := func(t *testing.T) *httptest.Server {
				t.Helper()
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			}

			successRunner := func(t *testing.T, i *Inertia) {
				t.Helper()
				w, r := requestMock(http.MethodGet, "/home")

				err := i.Render(w, r, "Some/Component", Props{"foo": "bar"})
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}

				var buf bytes.Buffer

				assertable := AssertFromReader(t, io.TeeReader(w.Body, &buf))
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
			}

			errorRunner := func(t *testing.T, i *Inertia) {
				t.Helper()
				w, r := requestMock(http.MethodGet, "/home")

				err := i.Render(w, r, "Some/Component", Props{"foo": "bar"})
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}

				var buf bytes.Buffer

				assertable := AssertFromReader(t, io.TeeReader(w.Body, &buf))
				assertable.AssertComponent("Some/Component")
				assertable.AssertProps(Props{"foo": "bar", "errors": map[string]any{}})
				assertable.AssertVersion("f8v01xv4h4")
				assertable.AssertURL("/home")
			}

			t.Run("success", func(t *testing.T) {
				t.Parallel()

				ts := newTestServerSSR(t)

				defer ts.Close()

				i := I(func(i *Inertia) {
					i.rootTemplateHTML = rootTemplate
					i.version = "f8v01xv4h4"
					i.ssrURL = ts.URL
					i.ssrHTTPClient = ts.Client()
				})

				successRunner(t, i)
			})

			t.Run("success with pre-parsed root template", func(t *testing.T) {
				t.Parallel()

				ts := newTestServerSSR(t)

				defer ts.Close()

				tmpl, err := template.New("root").
					Funcs(template.FuncMap(make(TemplateFuncs))).
					Parse(rootTemplate)
				if err != nil {
					t.Fatalf("parse root template: %v", err)
				}

				i := I(func(i *Inertia) {
					i.rootTemplate = tmpl
					i.version = "f8v01xv4h4"
					i.ssrURL = ts.URL
					i.ssrHTTPClient = ts.Client()
				})

				successRunner(t, i)
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

				errorRunner(t, i)
			})

			t.Run("error with fallback and pre-parsed root template", func(t *testing.T) {
				t.Parallel()

				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}))
				defer ts.Close()

				tmpl, err := template.New("root").
					Funcs(template.FuncMap(make(TemplateFuncs))).
					Parse(rootTemplate)
				if err != nil {
					t.Fatalf("parse root template: %v", err)
				}

				i := I(func(i *Inertia) {
					i.rootTemplate = tmpl
					i.version = "f8v01xv4h4"
					i.ssrURL = ts.URL
					i.ssrHTTPClient = ts.Client()
				})

				errorRunner(t, i)
			})
		})

		t.Run("shared funcs", func(t *testing.T) {
			t.Parallel()

			runner := func(t *testing.T, i *Inertia) {
				t.Helper()
				w, r := requestMock(http.MethodGet, "/")

				err := i.Render(w, r, "Some/Component")
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}

				got := w.Body.String()
				want := "foo bar"

				if got != want {
					t.Fatalf("got=%s, want=%s", got, want)
				}
			}

			t.Run("success", func(t *testing.T) {
				i := I(func(i *Inertia) {
					i.rootTemplateHTML = `{{ trim " foo bar " }}`
					i.sharedTemplateFuncs = TemplateFuncs{
						"trim": strings.TrimSpace,
					}
				})

				runner(t, i)
			})

			t.Run("success with pre-parsed root template", func(t *testing.T) {
				tFuncs := make(TemplateFuncs)
				tFuncs["trim"] = strings.TrimSpace

				tmpl, err := template.New("root").
					Funcs(template.FuncMap(tFuncs)).
					Parse(`{{ trim " foo bar " }}`)
				if err != nil {
					t.Fatalf("parse root template: %v", err)
				}

				i := I(func(i *Inertia) {
					i.rootTemplate = tmpl
				})

				runner(t, i)
			})
		})

		t.Run("shared template data", func(t *testing.T) {
			t.Parallel()

			runner := func(t *testing.T, i *Inertia) {
				t.Helper()
				w, r := requestMock(http.MethodGet, "/")

				err := i.Render(w, r, "Some/Component")
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}

				got := w.Body.String()
				want := "Hello, world!"

				if got != want {
					t.Fatalf("got=%s, want=%s", got, want)
				}
			}

			t.Run("success", func(t *testing.T) {
				i := I(func(i *Inertia) {
					i.rootTemplateHTML = `Hello, {{ .text }}!`
					i.sharedTemplateData = TemplateData{
						"text": "world",
					}
				})

				runner(t, i)
			})

			t.Run("success with pre-parsed root template", func(t *testing.T) {
				tmpl, err := template.New("root").
					Funcs(template.FuncMap(make(TemplateFuncs))).
					Parse(`Hello, {{ .text }}!`)
				if err != nil {
					t.Fatalf("parse root template: %v", err)
				}

				i := I(func(i *Inertia) {
					i.rootTemplate = tmpl
					i.sharedTemplateData = TemplateData{
						"text": "world",
					}
				})

				runner(t, i)
			})
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

			err := i.Render(w, r, "Some/Component", Props{"foo": "bar"})
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			assertable := AssertFromString(t, w.Body.String())
			assertable.AssertComponent("Some/Component")
			assertable.AssertProps(Props{"foo": "bar", "errors": map[string]any{}})
			assertable.AssertVersion("f8v01xv4h4")
			assertable.AssertURL("/home")
			assertable.AssertEncryptHistory(false)
			assertable.AssertEncryptHistory(false)
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

			ctx := SetProps(r.Context(), Props{"foo": "baz", "abc": "456", "ctx": "prop"})

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

			ctx := SetValidationErrors(r.Context(), ValidationErrors{"foo": "bar"})

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

		t.Run("history encryption", func(t *testing.T) {
			t.Parallel()

			w, r := requestMock(http.MethodGet, "/home")
			asInertiaRequest(r)

			ctx := r.Context()
			ctx = SetEncryptHistory(ctx, true)
			ctx = ClearHistory(ctx)

			err := I().Render(w, r.WithContext(ctx), "Some/Component")
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			assertable := AssertFromString(t, w.Body.String())
			assertable.AssertEncryptHistory(true)
			assertable.AssertClearHistory(true)
		})

		t.Run("props value resolving", func(t *testing.T) {
			t.Parallel()

			t.Run("reject ignoreFirstLoad props", func(t *testing.T) {
				t.Parallel()

				w, r := requestMock(http.MethodGet, "/home")
				asInertiaRequest(r)

				err := I().Render(w, r, "Some/Component", Props{
					"foo":                       "bar",
					"closure":                   func() any { return "prop" },
					"closure_with_ctx":          func(_ context.Context) any { return "prop" },
					"closure_with_err":          func() (any, error) { return "prop", nil },
					"closure_with_ctx_with_err": func(_ context.Context) (any, error) { return "prop", nil },
					"optional":                  Optional(func() (any, error) { return "prop", nil }),
					"defer":                     Defer(func() (any, error) { return "prop", nil }),
				})
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}

				assertable := AssertFromString(t, w.Body.String())
				assertable.AssertProps(Props{
					"foo":                       "bar",
					"closure":                   "prop",
					"closure_with_ctx":          "prop",
					"closure_with_err":          "prop",
					"closure_with_ctx_with_err": "prop",
					"errors":                    map[string]any{},
				})
			})

			t.Run("partial resolving", func(t *testing.T) {
				t.Parallel()

				w, r := requestMock(http.MethodGet, "/home")
				asInertiaRequest(r)
				withPartialComponent(r, "Some/Component")

				err := I().Render(w, r, "Some/Component", Props{
					"foo": "bar",
					"baz": "quz",
				})
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}

				assertable := AssertFromString(t, w.Body.String())
				assertable.AssertProps(Props{
					"foo":    "bar",
					"baz":    "quz",
					"errors": map[string]any{},
				})
			})

			t.Run("only", func(t *testing.T) {
				t.Parallel()

				t.Run("partial", func(t *testing.T) {
					t.Parallel()

					w, r := requestMock(http.MethodGet, "/home")
					asInertiaRequest(r)
					withOnly(r, []string{"foo", "closure", "optional"})
					withPartialComponent(r, "Some/Component")

					err := I().Render(w, r, "Some/Component", Props{
						"foo":      "bar",
						"abc":      "123",
						"closure":  func() (any, error) { return "prop", nil },
						"optional": Optional(func() (any, error) { return "prop", nil }),
						"always":   Always("prop"),
					})
					if err != nil {
						t.Fatalf("unexpected error: %s", err)
					}

					assertable := AssertFromString(t, w.Body.String())
					assertable.AssertProps(Props{
						"foo":      "bar",
						"closure":  "prop",
						"optional": "prop",
						"always":   "prop",
						"errors":   map[string]any{},
					})
				})

				t.Run("not partial", func(t *testing.T) {
					t.Parallel()

					w, r := requestMock(http.MethodGet, "/home")
					asInertiaRequest(r)
					withOnly(r, []string{"foo", "closure", "optional"})
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

				t.Run("partial", func(t *testing.T) {
					t.Parallel()

					w, r := requestMock(http.MethodGet, "/home")
					asInertiaRequest(r)
					withOnly(r, []string{"foo", "baz"})
					withExcept(r, []string{"foo", "abc", "optional", "always"})
					withPartialComponent(r, "Some/Component")

					err := I().Render(w, r, "Some/Component", Props{
						"foo":      "bar",
						"baz":      "quz",
						"bez":      "bee",
						"optional": Optional(func() (any, error) { return "prop", nil }),
						"always":   Always("prop"),
					})
					if err != nil {
						t.Fatalf("unexpected error: %s", err)
					}

					assertable := AssertFromString(t, w.Body.String())
					assertable.AssertProps(Props{
						"baz":    "quz",
						"always": "prop",
						"errors": map[string]any{},
					})
				})

				t.Run("not partial", func(t *testing.T) {
					t.Parallel()

					w, r := requestMock(http.MethodGet, "/home")
					asInertiaRequest(r)
					withExcept(r, []string{"baz"})
					withPartialComponent(r, "Other/Component")

					err := I().Render(w, r, "Some/Component", Props{
						"foo": "bar",
						"baz": "quz",
					})
					if err != nil {
						t.Fatalf("unexpected error: %s", err)
					}

					assertable := AssertFromString(t, w.Body.String())
					assertable.AssertProps(Props{
						"foo":    "bar",
						"baz":    "quz",
						"errors": map[string]any{},
					})
				})
			})

			t.Run("deferred props", func(t *testing.T) {
				t.Parallel()

				t.Run("partial", func(t *testing.T) {
					t.Parallel()

					w, r := requestMock(http.MethodGet, "/home")
					asInertiaRequest(r)
					withPartialComponent(r, "Some/Component")

					i := I()
					i.sharedProps = Props{
						"defer_shared_prop": Defer("prop_defer", "shared"),
					}
					err := i.Render(w, r, "Some/Component", Props{
						"defer_with_default_group1": Defer(func() (any, error) { return "prop1", nil }),
						"defer_with_default_group2": Defer("prop2", "default"),
						"defer_with_custom_group":   Defer("prop3", "foobar"),
					})
					if err != nil {
						t.Fatalf("unexpected error: %s", err)
					}

					assertable := AssertFromString(t, w.Body.String())
					assertable.AssertProps(Props{
						"defer_shared_prop":         "prop_defer",
						"defer_with_default_group1": "prop1",
						"defer_with_default_group2": "prop2",
						"defer_with_custom_group":   "prop3",
						"errors":                    map[string]any{},
					})
					assertable.AssertDeferredProps(nil)
				})

				t.Run("not partial", func(t *testing.T) {
					t.Parallel()

					w, r := requestMock(http.MethodGet, "/home")
					asInertiaRequest(r)
					withPartialComponent(r, "Other/Component")

					i := I()
					i.sharedProps = Props{
						"defer_shared_prop": Defer("prop_defer", "shared"),
					}
					err := i.Render(w, r, "Some/Component", Props{
						"defer_with_default_group1": Defer(func() (any, error) { return "prop1", nil }),
						"defer_with_default_group2": Defer("prop2", "default"),
						"defer_with_custom_group":   Defer("prop3", "foobar"),
					})
					if err != nil {
						t.Fatalf("unexpected error: %s", err)
					}

					assertable := AssertFromString(t, w.Body.String())
					assertable.AssertProps(Props{"errors": map[string]any{}})

					sort.Strings(assertable.DeferredProps["default"]) // fix flacks
					assertable.AssertDeferredProps(map[string][]string{
						"default": {"defer_with_default_group1", "defer_with_default_group2"},
						"foobar":  {"defer_with_custom_group"},
						"shared":  {"defer_shared_prop"},
					})
				})
			})

			t.Run("merge props", func(t *testing.T) {
				t.Parallel()

				t.Run("array", func(t *testing.T) {
					t.Parallel()

					w, r := requestMock(http.MethodGet, "/home")
					asInertiaRequest(r)

					err := I().Render(w, r, "Some/Component", Props{
						"ids": Merge([]int{1, 2, 3}),
						"foo": "bar",
					})
					if err != nil {
						t.Fatalf("unexpected error: %s", err)
					}

					assertable := AssertFromString(t, w.Body.String())
					assertable.AssertProps(Props{
						"ids":    []any{float64(1), float64(2), float64(3)},
						"foo":    "bar",
						"errors": map[string]any{},
					})
					assertable.AssertMergeProps([]string{"ids"})
				})

				t.Run("scalar", func(t *testing.T) {
					t.Parallel()

					w, r := requestMock(http.MethodGet, "/home")
					asInertiaRequest(r)

					err := I().Render(w, r, "Some/Component", Props{
						"foo": Merge("bar"),
					})
					if err != nil {
						t.Fatalf("unexpected error: %s", err)
					}

					assertable := AssertFromString(t, w.Body.String())
					assertable.AssertProps(Props{
						"foo":    "bar",
						"errors": map[string]any{},
					})
					assertable.AssertMergeProps([]string{"foo"})
				})

				t.Run("reset", func(t *testing.T) {
					t.Parallel()

					w, r := requestMock(http.MethodGet, "/home")
					asInertiaRequest(r)
					withReset(r, []string{"foo"})

					err := I().Render(w, r, "Some/Component", Props{
						"foo": Merge([]int{1, 2}),
						"bar": Merge([]int{3, 4}),
						"baz": "quz",
					})
					if err != nil {
						t.Fatalf("unexpected error: %s", err)
					}

					assertable := AssertFromString(t, w.Body.String())
					assertable.AssertProps(Props{
						"foo":    []any{float64(1), float64(2)},
						"bar":    []any{float64(3), float64(4)},
						"baz":    "quz",
						"errors": map[string]any{},
					})
					assertable.AssertMergeProps([]string{"bar"})
				})

				t.Run("deferred props", func(t *testing.T) {
					t.Parallel()

					w, r := requestMock(http.MethodGet, "/home")
					asInertiaRequest(r)
					withPartialComponent(r, "Some/Component")

					err := I().Render(w, r, "Some/Component", Props{
						"foo": Defer([]int{1, 2, 3}).Merge(),
					})
					if err != nil {
						t.Fatalf("unexpected error: %s", err)
					}

					assertable := AssertFromString(t, w.Body.String())
					assertable.AssertProps(Props{
						"foo":    []any{float64(1), float64(2), float64(3)},
						"errors": map[string]any{},
					})
					assertable.AssertDeferredProps(nil)
					assertable.AssertMergeProps([]string{"foo"})
				})

				t.Run("shared props", func(t *testing.T) {
					t.Parallel()

					w, r := requestMock(http.MethodGet, "/home")
					asInertiaRequest(r)

					i := I()
					i.sharedProps = Props{"foo": Merge("bar")}

					err := i.Render(w, r, "Some/Component")
					if err != nil {
						t.Fatalf("unexpected error: %s", err)
					}

					assertable := AssertFromString(t, w.Body.String())
					assertable.AssertProps(Props{
						"foo":    "bar",
						"errors": map[string]any{},
					})
					assertable.AssertMergeProps([]string{"foo"})
				})
			})

			t.Run("proper interfaces", func(t *testing.T) {
				t.Parallel()

				w, r := requestMock(http.MethodGet, "/home")
				asInertiaRequest(r)

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

	t.Run("flash", func(t *testing.T) {
		t.Parallel()

		t.Run("validation errors", func(t *testing.T) {
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

		t.Run("clear history", func(t *testing.T) {
			t.Parallel()

			t.Run("plain redirect", func(t *testing.T) {
				t.Parallel()

				w, r := requestMock(http.MethodGet, "/")

				flashProvider := &flashProviderMock{}

				i := I(func(i *Inertia) {
					i.flash = flashProvider
				})

				withClearHistory(r)
				i.Location(w, r, "/foo")

				if !flashProvider.clearHistory {
					t.Fatalf("got clear history=%v, want=true", flashProvider.clearHistory)
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

				withClearHistory(r)
				i.Location(w, r, "/foo", http.StatusMovedPermanently)

				if !flashProvider.clearHistory {
					t.Fatalf("got clear history=%v, want=true", flashProvider.clearHistory)
				}
			})
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

	t.Run("flash", func(t *testing.T) {
		t.Parallel()

		t.Run("validation errors", func(t *testing.T) {
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

		t.Run("clear history", func(t *testing.T) {
			t.Parallel()

			w, r := requestMock(http.MethodGet, "/")

			flashProvider := &flashProviderMock{}

			i := I(func(i *Inertia) {
				i.flash = flashProvider
			})

			withClearHistory(r)
			i.Redirect(w, r, "https://example.com/foo")

			if !flashProvider.clearHistory {
				t.Fatalf("got clear history=%v, want=true", flashProvider.clearHistory)
			}
		})
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

	t.Run("flash", func(t *testing.T) {
		t.Parallel()

		t.Run("validation errors", func(t *testing.T) {
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

		t.Run("clear history", func(t *testing.T) {
			t.Parallel()

			w, r := requestMock(http.MethodGet, "/")
			r.Header.Set("Referer", "https://example.com/foo")

			flashProvider := &flashProviderMock{}

			i := I(func(i *Inertia) {
				i.flash = flashProvider
			})

			withClearHistory(r)
			i.Back(w, r)

			if !flashProvider.clearHistory {
				t.Fatalf("got clear history=%v, want=true", flashProvider.clearHistory)
			}
		})
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

	assertable := AssertFromReader(t, w.Body)
	assertable.AssertComponent("Some/Component")
	assertable.AssertProps(Props{"foo": "bar", "errors": map[string]any{}})
	assertable.AssertVersion("f8v01xv4h4")
	assertable.AssertURL("/home")

	assertNotInertiaResponse(t, w)
	assertHTMLResponse(t, w)
	assertResponseStatusCode(t, w, http.StatusOK)
}
