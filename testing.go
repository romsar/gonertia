package gonertia

import (
	"encoding/json"
	"html"
	"io"
	"reflect"
	"regexp"
)

type t interface {
	Helper()
	Fatal(args ...any)
	Fatalf(format string, args ...any)
}

type AssertableInertia struct {
	t t
	*page
}

func (i AssertableInertia) AssertComponent(component string) {
	i.t.Helper()

	if i.Component != component {
		i.t.Fatalf("inertia: Component=%s, want=%s", i.Component, component)
	}
}

func (i AssertableInertia) AssertVersion(version string) {
	i.t.Helper()

	if i.Version != version {
		i.t.Fatalf("inertia: Version=%s, want=%s", i.Version, version)
	}
}

func (i AssertableInertia) AssertURL(url string) {
	i.t.Helper()

	if i.URL != url {
		i.t.Fatalf("inertia: URL=%s, want=%s", i.URL, url)
	}
}

func (i AssertableInertia) AssertProps(props Props) {
	i.t.Helper()

	if !reflect.DeepEqual(i.Props, props) {
		i.t.Fatalf("inertia: Props=%#v, want=%#v", i.Props, props)
	}
}

func AssertInertia(t t, body io.Reader) AssertableInertia {
	t.Helper()

	bodyBs, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("read response bodyBs error: %#v", err)
	}

	return AssertInertiaFromBytes(t, bodyBs)
}

func AssertInertiaFromBytes(t t, body []byte) AssertableInertia {
	t.Helper()

	return AssertInertiaFromString(t, string(body))
}

var pageRe = regexp.MustCompile(` data-page="(.*?)"`)

func AssertInertiaFromString(t t, body string) AssertableInertia {
	t.Helper()

	matched := pageRe.FindAllStringSubmatch(body, -1)
	if len(matched) == 0 {
		invalidInertiaResponse(t)
	}

	assertable := AssertableInertia{t: t}

	for _, m := range matched {
		if len(m) <= 1 {
			invalidInertiaResponse(t)
		}

		pageJSON := []byte(html.UnescapeString(m[1]))
		if err := json.Unmarshal(pageJSON, &assertable.page); err == nil {
			break
		}
	}

	if assertable.page == nil {
		invalidInertiaResponse(t)
	}

	return assertable
}

func invalidInertiaResponse(t t) {
	t.Fatal("invalid inertia response")
}
