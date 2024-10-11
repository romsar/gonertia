package gonertia

import (
	"io"
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

const stubHTML = `<!DOCTYPE html>
<html lang="en">
	<head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <script type="module" src="/build/assets/main.js"></script>
            <link rel="stylesheet" href="/build/assets/index.css">
	</head>
	<body>
		<div data-page="foo bar">
			<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla a metus condimentum, pulvinar arcu in, lacinia urna.</p>
			<p>Proin tincidunt, leo ut consectetur tincidunt, sem ex fermentum ipsum, a sollicitudin odio magna et dui.</p>
			<p>Aliquam efficitur, purus quis porttitor placerat, massa mi hendrerit nulla, id convallis eros tortor non augue. Duis id varius arcu.</p>
		</div>
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

func TestAssertableInertia_EncryptHistory(t *testing.T) {
	t.Parallel()

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		mock := new(tMock)

		i := AssertableInertia{
			t:    mock,
			page: &page{EncryptHistory: true},
		}

		i.AssertEncryptHistory(true)

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
			page: &page{EncryptHistory: true},
		}

		i.AssertEncryptHistory(false)

		if !mock.helperInvoked {
			t.Fatal("expected Helper() to be invoked")
		}

		if !mock.isFailed {
			t.Fatal("expected assertion failure")
		}
	})
}

func TestAssertableInertia_ClearHistory(t *testing.T) {
	t.Parallel()

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		mock := new(tMock)

		i := AssertableInertia{
			t:    mock,
			page: &page{ClearHistory: true},
		}

		i.AssertClearHistory(true)

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
			page: &page{ClearHistory: true},
		}

		i.AssertClearHistory(false)

		if !mock.helperInvoked {
			t.Fatal("expected Helper() to be invoked")
		}

		if !mock.isFailed {
			t.Fatal("expected assertion failure")
		}
	})
}

func TestAssertableInertia_DeferredProps(t *testing.T) {
	t.Parallel()

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		mock := new(tMock)

		i := AssertableInertia{
			t:    mock,
			page: &page{DeferredProps: map[string][]string{"foo": {"bar", "baz"}}},
		}

		i.AssertDeferredProps(map[string][]string{"foo": {"bar", "baz"}})

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
			page: &page{DeferredProps: map[string][]string{"foo": {"bar", "baz"}}},
		}

		i.AssertDeferredProps(map[string][]string{"foo": {"bar", "quz"}})

		if !mock.helperInvoked {
			t.Fatal("expected Helper() to be invoked")
		}

		if !mock.isFailed {
			t.Fatal("expected assertion failure")
		}
	})
}

func TestAssertableInertia_MergeProps(t *testing.T) {
	t.Parallel()

	t.Run("positive", func(t *testing.T) {
		t.Parallel()

		mock := new(tMock)

		i := AssertableInertia{
			t:    mock,
			page: &page{MergeProps: []string{"foo", "bar"}},
		}

		i.AssertMergeProps([]string{"foo", "bar"})

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
			page: &page{MergeProps: []string{"foo", "bar"}},
		}

		i.AssertMergeProps([]string{"foo", "baz"})

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

		assertStubSuccess(t, mock, stubJSON, assertable)
	})

	t.Run("success with html", func(t *testing.T) {
		t.Parallel()

		mock := new(tMock)

		assertable := AssertFromString(mock, stubHTML)

		assertStubSuccess(t, mock, stubHTML, assertable)
	})
}

func TestAssertFromBytes(t *testing.T) {
	t.Parallel()

	mock := new(tMock)

	assertable := AssertFromBytes(mock, []byte(stubHTML))

	assertStubSuccess(t, mock, stubHTML, assertable)
}

func TestAssertFromReader(t *testing.T) {
	t.Parallel()

	mock := new(tMock)

	assertable := AssertFromReader(mock, strings.NewReader(stubHTML))

	assertStubSuccess(t, mock, stubHTML, assertable)
}

func assertStubSuccess(t *testing.T, mock *tMock, wantBody string, assertable AssertableInertia) {
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

	gotBody, err := io.ReadAll(assertable.Body)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if string(gotBody) != wantBody {
		t.Fatalf("got body=%s, want=%s", gotBody, wantBody)
	}
}
