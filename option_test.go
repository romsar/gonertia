package gonertia

import (
	"io"
	"log"
	"net/http"
	"reflect"
	"testing"
	"testing/fstest"
	"time"
)

func TestWithVersion(t *testing.T) {
	t.Parallel()

	i := I()

	want := "327b6f07435811239bc47e1544353273"

	option := WithVersion("foo bar")

	if err := option(i); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if i.version != want {
		t.Fatalf("version=%s, want=%s", i.version, want)
	}
}

func TestWithVersionFromFile(t *testing.T) {
	t.Parallel()

	i := I()

	f := tmpFile(t, "foo")

	option := WithVersionFromFile(f.Name())

	want := "acbd18db4cc2f85cedef654fccc4a4d8"

	if err := option(i); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if i.version != want {
		t.Fatalf("version=%s, want=%s", i.version, want)
	}
}

func TestWithVersionFromFileFS(t *testing.T) {
	t.Parallel()

	i := I()

	mapFS := fstest.MapFS{
		"foo": {
			Data: []byte("foo"),
		},
	}

	option := WithVersionFromFileFS(mapFS, "foo")

	want := "acbd18db4cc2f85cedef654fccc4a4d8"

	if err := option(i); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if i.version != want {
		t.Fatalf("version=%s, want=%s", i.version, want)
	}
}

func TestWithJSONMarshaller(t *testing.T) {
	t.Parallel()

	t.Run("marshal", func(t *testing.T) {
		t.Parallel()

		i := I()

		want := "foo bar"

		option := WithJSONMarshaller(jsonTestMarshaller{val: want})

		if err := option(i); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		got, err := i.jsonMarshaller.Marshal([]byte{})
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if string(got) != want {
			t.Fatalf("JSONMarshaller.Marshal()=%s, want=%s", string(got), want)
		}
	})

	t.Run("decode", func(t *testing.T) {
		t.Parallel()

		i := I()

		want := "foo bar"

		option := WithJSONMarshaller(jsonTestMarshaller{val: want})

		if err := option(i); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		var got string

		err := i.jsonMarshaller.Decode(nil, &got)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if got != want {
			t.Fatalf("JSONMarshaller.Decode()=%s, want=%s", got, want)
		}
	})
}

type jsonTestMarshaller struct {
	val string
}

func (j jsonTestMarshaller) Decode(_ io.Reader, v any) error {
	if ptr, ok := v.(*string); ok {
		*ptr = j.val
	}
	return nil
}

func (j jsonTestMarshaller) Marshal(v any) ([]byte, error) {
	return []byte(j.val), nil
}

func TestWithLogger(t *testing.T) {
	t.Parallel()

	t.Run("with nil", func(t *testing.T) {
		t.Parallel()

		i := I()

		option := WithLogger(nil)

		if err := option(i); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if i.logger == nil {
			t.Fatal("Logger is nil")
		}
	})

	t.Run("with default", func(t *testing.T) {
		t.Parallel()

		i := I()

		option := WithLogger()

		if err := option(i); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if i.logger == nil {
			t.Fatal("Logger is nil")
		}
	})

	t.Run("with custom", func(t *testing.T) {
		t.Parallel()

		i := I()

		want := log.New(io.Discard, "foo bar", 0)

		option := WithLogger(want)

		if err := option(i); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if !reflect.DeepEqual(i.logger, want) {
			t.Fatalf("Logger=%#v, want=%#v", i.logger, want)
		}
	})
}

func TestWithContainerID(t *testing.T) {
	t.Parallel()

	i := I()

	want := "foo"

	option := WithContainerID(want)

	if err := option(i); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if i.containerID != want {
		t.Fatalf("containerID=%s, want=%s", i.containerID, want)
	}
}

func TestWithSSR(t *testing.T) {
	t.Parallel()

	t.Run("with default url", func(t *testing.T) {
		t.Parallel()

		i := I()

		wantURL := "http://127.0.0.1:13714"

		option := WithSSR()

		if err := option(i); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if i.ssrURL != wantURL {
			t.Fatalf("ssrURL=%s, want=%s", i.containerID, wantURL)
		}
	})

	t.Run("with specified url", func(t *testing.T) {
		t.Parallel()

		i := I()

		wantURL := "https://foo.bar/baz/quz"

		option := WithSSR(wantURL)

		if err := option(i); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if i.ssrURL != wantURL {
			t.Fatalf("ssrURL=%s, want=%s", i.containerID, wantURL)
		}
	})
}

func TestWithSSRHTTPClient(t *testing.T) {
	t.Parallel()

	i := I()

	want := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       5 * time.Second,
	}

	option := WithSSRHTTPClient(want)

	if err := option(i); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if i.ssrHTTPClient != want {
		t.Fatal("ssr http client was not set")
	}
}

func TestWithFlashProvider(t *testing.T) {
	t.Parallel()

	i := I()

	want := &flashProviderMock{}

	option := WithFlashProvider(want)

	if err := option(i); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if i.flash != want {
		t.Fatalf("flash provider=%v, want=%v", i.flash, want)
	}
}

func TestWithEncryptHistory(t *testing.T) {
	t.Parallel()

	i := I()

	option := WithEncryptHistory(true)

	if err := option(i); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !i.encryptHistory {
		t.Fatalf("encryptHistory=%t, want=%t", i.encryptHistory, true)
	}
}
