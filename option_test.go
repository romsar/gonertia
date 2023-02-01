package gonertia

import (
	"embed"
	"io"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestWithTemplateFS(t *testing.T) {
	t.Parallel()

	i := new(Inertia)
	fs := embed.FS{}

	option := WithTemplateFS(fs)

	if err := option(i); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(i.templateFS, fs) {
		t.Fatalf("got=%#v, want=%#v", i.templateFS, fs)
	}
}

func TestWithVersion(t *testing.T) {
	t.Parallel()

	i := new(Inertia)

	want := "foo bar"

	option := WithVersion(want)

	if err := option(i); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if i.version != want {
		t.Fatalf("got=%#v, want=%#v", i.version, want)
	}
}

func TestWithAssetURL(t *testing.T) {
	t.Parallel()

	i := new(Inertia)

	url := "https://example.com/foo/bar"
	want := md5(url)

	option := WithAssetURL(url)

	if err := option(i); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if i.version != want {
		t.Fatalf("got=%#v, want=%#v", i.version, want)
	}
}

func TestWithManifestFile(t *testing.T) {
	t.Parallel()

	i := new(Inertia)

	f, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	closed := false

	t.Cleanup(func() {
		if !closed {
			if err := f.Close(); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}

		if err := os.Remove(f.Name()); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	if _, err := f.WriteString("foo"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := f.Close(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	closed = true

	option := WithManifestFile(f.Name())

	want := "acbd18db4cc2f85cedef654fccc4a4d8"

	if err := option(i); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if i.version != want {
		t.Fatalf("got=%#v, want=%#v", i.version, want)
	}
}

func TestWithMarshalJSON(t *testing.T) {
	t.Parallel()

	i := new(Inertia)

	want := "bar"

	option := WithMarshalJSON(func(v any) ([]byte, error) {
		return []byte(want), nil
	})

	if err := option(i); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := i.marshallJSON([]byte{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(got) != want {
		t.Fatalf("got=%#v, want=%#v", string(got), want)
	}
}

func TestWithoutLogger(t *testing.T) {
	t.Parallel()

	i := new(Inertia)

	want := log.New(io.Discard, "", 0)

	option := WithoutLogger()

	if err := option(i); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(i.logger, want) {
		t.Fatalf("got=%#v, want=%#v", i.logger, want)
	}
}

func TestWithLogger(t *testing.T) {
	t.Parallel()

	t.Run("with nil", func(t *testing.T) {
		t.Parallel()

		i := new(Inertia)

		want := log.New(io.Discard, "", 0)

		option := WithLogger(nil)

		if err := option(i); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(i.logger, want) {
			t.Fatalf("got=%#v, want=%#v", i.logger, want)
		}
	})

	t.Run("with logger", func(t *testing.T) {
		t.Parallel()

		i := new(Inertia)

		want := log.New(io.Discard, "foo bar", 0)

		option := WithLogger(want)

		if err := option(i); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(i.logger, want) {
			t.Fatalf("got=%#v, want=%#v", i.logger, want)
		}
	})
}

func TestWithContainerID(t *testing.T) {
	t.Parallel()

	i := new(Inertia)

	want := "foo"

	option := WithContainerID(want)

	if err := option(i); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if i.containerID != want {
		t.Fatalf("got=%#v, want=%#v", i.containerID, want)
	}
}
