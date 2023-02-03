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

// AssertableInertia is an Inertia response struct with assert methods.
type AssertableInertia struct {
	t t
	*page
}

// AssertComponent verifies that component from Inertia
// response and the passed component are the same.
func (i AssertableInertia) AssertComponent(component string) {
	i.t.Helper()

	if i.Component != component {
		i.t.Fatalf("inertia: Component=%s, want=%s", i.Component, component)
	}
}

// AssertVersion verifies that version from Inertia
// response and the passed version are the same.
func (i AssertableInertia) AssertVersion(version string) {
	i.t.Helper()

	if i.Version != version {
		i.t.Fatalf("inertia: Version=%s, want=%s", i.Version, version)
	}
}

// AssertURL verifies that url from Inertia
// response and the passed url are the same.
func (i AssertableInertia) AssertURL(url string) {
	i.t.Helper()

	if i.URL != url {
		i.t.Fatalf("inertia: URL=%s, want=%s", i.URL, url)
	}
}

// AssertProps verifies that props from Inertia
// response and the passed props are the same.
func (i AssertableInertia) AssertProps(props Props) {
	i.t.Helper()

	if !reflect.DeepEqual(i.Props, props) {
		i.t.Fatalf("inertia: Props=%#v, want=%#v", i.Props, props)
	}
}

// AssertInertia creates AssertableInertia from the io.Reader body.
func AssertInertia(t t, body io.Reader) AssertableInertia {
	t.Helper()

	bodyBs, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("read response bodyBs error: %#v", err)
	}

	return AssertInertiaFromBytes(t, bodyBs)
}

// AssertInertiaFromBytes creates AssertableInertia from the bytes body.
func AssertInertiaFromBytes(t t, body []byte) AssertableInertia {
	t.Helper()

	return AssertInertiaFromString(t, string(body))
}

var pageRe = regexp.MustCompile(` data-page="(.*?)"`)

// AssertInertiaFromString creates AssertableInertia from the string body.
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
