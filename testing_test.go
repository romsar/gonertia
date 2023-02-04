package gonertia

import (
	"reflect"
	"strings"
	"testing"
)

type tMock struct {
	helperInvoked bool
	isFailed      bool
}

func (m *tMock) Helper() {
	m.helperInvoked = true
}

func (m *tMock) Fatal(args ...any) {
	m.isFailed = true
}

func (m *tMock) Fatalf(format string, args ...any) {
	m.isFailed = true
}

const stubHTML = `<html>
	<head></head>
	<body>
		<div data-page="foo bar"></div>
		<div id="app" data-page="{&#34;component&#34;:&#34;Foo/Bar&#34;,&#34;props&#34;:{&#34;foo&#34;: &#34;bar&#34;},&#34;url&#34;:&#34;https://example.com&#34;,&#34;version&#34;:&#34;foobar&#34;}"></div>
	</body>
</html>`

const stubJSON = `{"component":"Foo/Bar","props":{"foo": "bar"},"url":"https://example.com","version":"foobar"}`

func TestAssertableInertia_AssertComponent(t *testing.T) {
	t.Parallel()

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		mock := new(tMock)

		i := AssertableInertia{
			t:    mock,
			page: &page{Component: "Foo/Bar"},
		}

		i.AssertComponent("Foo/Bar")

		if !mock.helperInvoked {
			t.Fatal("expected Helper() to be invoked")
		}

		if mock.isFailed {
			t.Fatal("unexpected assertion failure")
		}
	})

	t.Run("negative", func(t *testing.T) {
		t.Parallel()

		mock := new(tMock)

		i := AssertableInertia{
			t:    mock,
			page: &page{Component: "Foo/Bar"},
		}

		i.AssertComponent("Some/Component")

		if !mock.helperInvoked {
			t.Fatal("expected Helper() to be invoked")
		}

		if !mock.isFailed {
			t.Fatal("expected assertion failure")
		}
	})
}

func TestAssertableInertia_AssertVersion(t *testing.T) {
	t.Parallel()

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		mock := new(tMock)

		i := AssertableInertia{
			t:    mock,
			page: &page{Version: "foo bar"},
		}

		i.AssertVersion("foo bar")

		if !mock.helperInvoked {
			t.Fatal("expected Helper() to be invoked")
		}

		if mock.isFailed {
			t.Fatal("unexpected assertion failure")
		}
	})

	t.Run("negative", func(t *testing.T) {
		t.Parallel()

		mock := new(tMock)

		i := AssertableInertia{
			t:    mock,
			page: &page{Version: "foo bar"},
		}

		i.AssertVersion("foo barrr")

		if !mock.helperInvoked {
			t.Fatal("expected Helper() to be invoked")
		}

		if !mock.isFailed {
			t.Fatal("expected assertion failure")
		}
	})
}

func TestAssertableInertia_AssertURL(t *testing.T) {
	t.Parallel()

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		mock := new(tMock)

		i := AssertableInertia{
			t:    mock,
			page: &page{URL: "/home"},
		}

		i.AssertURL("/home")

		if !mock.helperInvoked {
			t.Fatal("expected Helper() to be invoked")
		}

		if mock.isFailed {
			t.Fatal("unexpected assertion failure")
		}
	})

	t.Run("negative", func(t *testing.T) {
		t.Parallel()

		mock := new(tMock)

		i := AssertableInertia{
			t:    mock,
			page: &page{URL: "https://foobar.com"},
		}

		i.AssertURL("https://foobarrr.com")

		if !mock.helperInvoked {
			t.Fatal("expected Helper() to be invoked")
		}

		if !mock.isFailed {
			t.Fatal("expected assertion failure")
		}
	})
}

func TestAssertableInertia_AssertProps(t *testing.T) {
	t.Parallel()

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		mock := new(tMock)

		i := AssertableInertia{
			t:    mock,
			page: &page{Props: Props{"foo": "bar"}},
		}

		i.AssertProps(Props{"foo": "bar"})

		if !mock.helperInvoked {
			t.Fatal("expected Helper() to be invoked")
		}

		if mock.isFailed {
			t.Fatal("unexpected assertion failure")
		}
	})

	t.Run("negative", func(t *testing.T) {
		t.Parallel()

		mock := new(tMock)

		i := AssertableInertia{
			t:    mock,
			page: &page{Props: Props{"foo": "bar"}},
		}

		i.AssertProps(Props{"foo": "barrr"})

		if !mock.helperInvoked {
			t.Fatal("expected Helper() to be invoked")
		}

		if !mock.isFailed {
			t.Fatal("expected assertion failure")
		}
	})
}

func TestAssertFromString(t *testing.T) {
	t.Parallel()

	t.Run("without inertia container", func(t *testing.T) {
		t.Parallel()

		mock := new(tMock)

		AssertFromString(mock, `<html>
	<head></head>
	<body><div id="app" data-foo="bar"></div></body>
</html>`)

		if !mock.helperInvoked {
			t.Fatal("expected Helper() to be invoked")
		}

		if !mock.isFailed {
			t.Fatal("expected assertion failure")
		}
	})

	t.Run("with invalid inertia data", func(t *testing.T) {
		t.Parallel()

		mock := new(tMock)

		AssertFromString(mock, `<html>
	<head></head>
	<body><div id="app" data-page="foo bar"></div></body>
</html>`)

		if !mock.helperInvoked {
			t.Fatal("expected Helper() to be invoked")
		}

		if !mock.isFailed {
			t.Fatal("expected assertion failure")
		}
	})

	t.Run("success with json", func(t *testing.T) {
		t.Parallel()

		mock := new(tMock)

		assertable := AssertFromString(mock, stubJSON)

		assertStubSuccess(t, mock, assertable)
	})

	t.Run("success with html", func(t *testing.T) {
		t.Parallel()

		mock := new(tMock)

		assertable := AssertFromString(mock, stubHTML)

		assertStubSuccess(t, mock, assertable)
	})
}

func TestAssertFromBytes(t *testing.T) {
	t.Parallel()

	mock := new(tMock)

	assertable := AssertFromBytes(mock, []byte(stubHTML))

	assertStubSuccess(t, mock, assertable)
}

func TestAssert(t *testing.T) {
	t.Parallel()

	mock := new(tMock)

	assertable := Assert(mock, strings.NewReader(stubHTML))

	assertStubSuccess(t, mock, assertable)
}

func assertStubSuccess(t *testing.T, mock *tMock, assertable AssertableInertia) {
	t.Helper()

	if !mock.helperInvoked {
		t.Fatal("expected Helper() to be invoked")
	}

	if mock.isFailed {
		t.Fatal("unexpected assertion failure")
	}

	if wantComponent := "Foo/Bar"; assertable.Component != wantComponent {
		t.Fatalf("Component=%s, want=%s", assertable.Component, wantComponent)
	}

	if wantVersion := "foobar"; assertable.Version != wantVersion {
		t.Fatalf("Version=%s, want=%s", assertable.Version, wantVersion)
	}

	if wantURL := "https://example.com"; assertable.URL != wantURL {
		t.Fatalf("URL=%s, want=%s", assertable.URL, wantURL)
	}

	wantProps := Props{"foo": "bar"}
	if !reflect.DeepEqual(assertable.Props, wantProps) {
		t.Fatalf("Props=%#v, want=%#v", assertable.Props, wantProps)
	}
}
