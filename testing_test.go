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
			page: &page{URL: "https://foobar.com"},
		}

		i.AssertURL("https://foobar.com")

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

func TestAssertInertiaFromString(t *testing.T) {
	t.Parallel()

	t.Run("without inertia container", func(t *testing.T) {
		t.Parallel()

		mock := new(tMock)

		AssertInertiaFromString(mock, `<html>
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

		AssertInertiaFromString(mock, `<html>
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

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mock := new(tMock)

		assertable := AssertInertiaFromString(mock, stubHTML)

		testStubSuccess(t, mock, assertable)
	})
}

func TestAssertInertiaFromBytes(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mock := new(tMock)

		assertable := AssertInertiaFromBytes(mock, []byte(stubHTML))

		testStubSuccess(t, mock, assertable)
	})
}

func TestAssertInertia(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mock := new(tMock)

		assertable := AssertInertia(mock, strings.NewReader(stubHTML))

		testStubSuccess(t, mock, assertable)
	})
}

func testStubSuccess(t *testing.T, mock *tMock, assertable AssertableInertia) {
	if !mock.helperInvoked {
		t.Fatal("expected Helper() to be invoked")
	}

	if mock.isFailed {
		t.Fatal("unexpected assertion failure")
	}

	wantComponent := "Foo/Bar"
	wantVersion := "foobar"
	wantURL := "https://example.com"
	wantProps := Props{"foo": "bar"}

	if assertable.Component != wantComponent {
		t.Fatalf("Component=%s, want=%s", assertable.Component, wantComponent)
	}

	if assertable.Version != wantVersion {
		t.Fatalf("Version=%s, want=%s", assertable.Version, wantVersion)
	}

	if assertable.URL != wantURL {
		t.Fatalf("URL=%s, want=%s", assertable.URL, wantURL)
	}

	if !reflect.DeepEqual(assertable.Props, wantProps) {
		t.Fatalf("Props=%#v, want=%#v", assertable.Props, wantProps)
	}
}