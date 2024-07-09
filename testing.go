package gonertia

import (
	"bytes"
	"encoding/json"
	"html"
	"io"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

type t interface {
	Helper()
	Fatal(args ...any)
	Fatalf(format string, args ...any)
}

var _ t = (*testing.T)(nil)

// AssertableInertia is an Inertia response struct with assert methods.
type AssertableInertia struct {
	t t
	*page
	Body *bytes.Buffer
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

var containerRe = regexp.MustCompile(` data-page="(.*?)"`)

// Assert creates AssertableInertia from the io.Reader body.
func Assert(t t, body io.Reader) AssertableInertia {
	t.Helper()

	assertable := AssertableInertia{t: t}

	var buf bytes.Buffer

	// Might be body is a json? Let's try to unmarshal first.
	if err := json.NewDecoder(io.TeeReader(body, &buf)).Decode(&assertable.page); err == nil {
		assertable.Body = &buf
		return assertable
	}

	matched := containerRe.FindAllStringSubmatch(buf.String(), -1)
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

	assertable.Body = &buf
	return assertable
}

// AssertFromBytes creates AssertableInertia from the bytes body.
func AssertFromBytes(t t, body []byte) AssertableInertia {
	t.Helper()

	return Assert(t, bytes.NewReader(body))
}

// AssertFromString creates AssertableInertia from the string body.
func AssertFromString(t t, body string) AssertableInertia {
	t.Helper()

	return Assert(t, strings.NewReader(body))
}

func invalidInertiaResponse(t t) {
	t.Fatal("invalid inertia response")
}
