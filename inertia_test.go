package gonertia

import (
	"net/http"
	"testing"
	"testing/fstest"
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

var rootTemplate = `<html>
	<head>{{ .inertiaHead }}</head>
	<body>{{ .inertia }}</body>
</html>`

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
	})

	t.Run("inertia request", func(t *testing.T) {
		t.Parallel()

		f := tmpFile(t, rootTemplate)

		i := I(func(i *Inertia) {
			i.rootTemplatePath = f.Name()
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

		assertable := AssertInertiaFromString(t, w.Body.String())
		assertable.AssertComponent("Some/Component")
		assertable.AssertProps(Props{"foo": "bar"})
		assertable.AssertVersion("f8v01xv4h4")
		assertable.AssertURL("/home")

		assertInertiaResponse(t, w)
		assertJSONResponse(t, w)
		assertResponseStatusCode(t, w, http.StatusOK)
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

	assertable := AssertInertia(t, w.Body)
	assertable.AssertComponent("Some/Component")
	assertable.AssertProps(Props{"foo": "bar"})
	assertable.AssertVersion("f8v01xv4h4")
	assertable.AssertURL("/home")

	assertNotInertiaResponse(t, w)
	assertHTMLResponse(t, w)
	assertResponseStatusCode(t, w, http.StatusOK)
}
