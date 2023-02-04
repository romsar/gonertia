package gonertia

import (
	"encoding/json"
	"html"
	"io"
	"reflect"
	"regexp"
)

// t is the contract of testing.T.
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
func (i AssertableInertia) AssertComponent(want string) {
	i.t.Helper()

	if i.Component != want {
		i.t.Fatalf("inertia: Component=%s, want=%s", i.Component, want)
	}
}

// AssertVersion verifies that version from Inertia
// response and the passed version are the same.
func (i AssertableInertia) AssertVersion(want string) {
	i.t.Helper()

	if i.Version != want {
		i.t.Fatalf("inertia: Version=%s, want=%s", i.Version, want)
	}
}

// AssertURL verifies that url from Inertia
// response and the passed url are the same.
func (i AssertableInertia) AssertURL(want string) {
	i.t.Helper()

	if i.URL != want {
		i.t.Fatalf("inertia: URL=%s, want=%s", i.URL, want)
	}
}

// AssertProps verifies that props from Inertia
// response and the passed props are the same.
func (i AssertableInertia) AssertProps(want Props) {
	i.t.Helper()

	if !reflect.DeepEqual(i.Props, want) {
		i.t.Fatalf("inertia: Props=%#v, want=%#v", i.Props, want)
	}
}

// Assert creates AssertableInertia from the io.Reader body.
func Assert(t t, body io.Reader) AssertableInertia {
	t.Helper()

	bodyBs, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("read response bodyBs error: %#v", err)
	}

	return AssertFromBytes(t, bodyBs)
}

// AssertFromBytes creates AssertableInertia from the bytes body.
func AssertFromBytes(t t, body []byte) AssertableInertia {
	t.Helper()

	return AssertFromString(t, string(body))
}

var pageRe = regexp.MustCompile(` data-page="(.*?)"`)

// AssertFromString creates AssertableInertia from the string body.
func AssertFromString(t t, body string) AssertableInertia {
	t.Helper()

	assertable := AssertableInertia{t: t}

	// Might be body is a json? Let's try to unmarshall first.
	if err := json.Unmarshal([]byte(body), &assertable.page); err == nil {
		return assertable
	}

	matched := pageRe.FindAllStringSubmatch(body, -1)
	if len(matched) == 0 {
		invalidInertiaResponse(t)
	}

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

// invalidInertiaResponse fail test with invalid inertia response message.
func invalidInertiaResponse(t t) {
	t.Fatal("invalid inertia response")
}
