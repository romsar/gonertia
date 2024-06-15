package gonertia

import (
	"net/http"
	"strings"
	"testing"
	"testing/fstest"
)

var rootTemplate = `<html>
	<head>{{ .inertiaHead }}</head>
	<body>{{ .inertia }}</body>
</html>`

//nolint:gocognit
func TestInertia_Render(t *testing.T) {
	t.Parallel()

	t.Run("plain request", func(t *testing.T) {
		t.Parallel()

		t.Run("file template", func(t *testing.T) {
			t.Parallel()

			f := tmpFile(t, rootTemplate)

			i := I(func(i *Inertia) {
				i.rootTemplatePath = f.Name()
				i.version = "f8v01xv4h4"
			})

			assertRootTemplateSuccess(t, i)
		})

		t.Run("embed fs template", func(t *testing.T) {
			t.Parallel()

			fs := fstest.MapFS{
				"app.html": {
					Data: []byte(rootTemplate),
				},
			}

			i := I(func(i *Inertia) {
				i.rootTemplatePath = "app.html"
				i.version = "f8v01xv4h4"
				i.templateFS = fs
			})

			assertRootTemplateSuccess(t, i)
		})

		t.Run("shared funcs", func(t *testing.T) {
			t.Parallel()

			f := tmpFile(t, `{{ trim " foo bar " }}`)
			w, r := requestMock(http.MethodGet, "/")

			i := I(func(i *Inertia) {
				i.rootTemplatePath = f.Name()
				i.sharedTemplateFuncs = TemplateFuncs{
					"trim": strings.TrimSpace,
				}
			})

			err := i.Render(w, r, "Some/Component")
			if err != nil {
				t.Fatalf("unexpected error: %#v", err)
			}

			got := w.Body.String()
			want := "foo bar"

			if got != want {
				t.Fatalf("got=%s, want=%s", got, want)
			}
		})

		t.Run("shared template data", func(t *testing.T) {
			t.Parallel()

			f := tmpFile(t, `Hello, {{ .text }}!`)
			w, r := requestMock(http.MethodGet, "/")

			i := I(func(i *Inertia) {
				i.rootTemplatePath = f.Name()
				i.sharedTemplateData = TemplateData{
					"text": "world",
				}
			})

			err := i.Render(w, r, "Some/Component")
			if err != nil {
				t.Fatalf("unexpected error: %#v", err)
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
				t.Fatalf("unexpected error: %#v", err)
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
				t.Fatalf("unexpected error: %#v", err)
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

			i := I()

			w, r := requestMock(http.MethodGet, "/home")
			asInertiaRequest(r)

			ctx := WithValidationErrors(r.Context(), ValidationErrors{"foo": "bar"})

			err := i.Render(w, r.WithContext(ctx), "Some/Component", Props{
				"abc": "123",
			})
			if err != nil {
				t.Fatalf("unexpected error: %#v", err)
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
					"lazy":             LazyProp(func() (any, error) { return "prop", nil }),
				})
				if err != nil {
					t.Fatalf("unexpected error: %#v", err)
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
						"lazy":    LazyProp(func() (any, error) { return "prop", nil }),
						"always":  AlwaysProp(func() any { return "prop" }),
					})
					if err != nil {
						t.Fatalf("unexpected error: %#v", err)
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
						t.Fatalf("unexpected error: %#v", err)
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
					"lazy":   LazyProp(func() (any, error) { return "prop", nil }),
					"always": AlwaysProp(func() any { return "prop" }),
				})
				if err != nil {
					t.Fatalf("unexpected error: %#v", err)
				}

				assertable := AssertFromString(t, w.Body.String())
				assertable.AssertProps(Props{
					"baz":    "quz",
					"errors": map[string]any{},
				})
			})
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
		t.Fatalf("unexpected error: %#v", err)
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
