package gonertia

import (
	"embed"
	"io"
	"log"
	"reflect"
	"testing"
)

func TestWithTemplateFS(t *testing.T) {
	t.Parallel()

	i := new(Inertia)
	fs := embed.FS{}

	option := WithTemplateFS(fs)

	if err := option(i); err != nil {
		t.Fatalf("unexpected error: %#v", err)
	}

	if !reflect.DeepEqual(i.templateFS, fs) {
		t.Fatalf("templateFS=%#v, want=%#v", i.templateFS, fs)
	}
}

func TestWithVersion(t *testing.T) {
	t.Parallel()

	i := new(Inertia)

	want := "foo bar"

	option := WithVersion(want)

	if err := option(i); err != nil {
		t.Fatalf("unexpected error: %#v", err)
	}

	if i.version != want {
		t.Fatalf("version=%s, want=%s", i.version, want)
	}
}

func TestWithAssetURL(t *testing.T) {
	t.Parallel()

	i := new(Inertia)

	url := "https://example.com/foo/bar"
	want := md5(url)

	option := WithAssetURL(url)

	if err := option(i); err != nil {
		t.Fatalf("unexpected error: %#v", err)
	}

	if i.version != want {
		t.Fatalf("version=%s, want=%s", i.version, want)
	}
}

func TestWithManifestFile(t *testing.T) {
	t.Parallel()

	i := new(Inertia)

	f := tmpFile(t, "foo")

	option := WithManifestFile(f.Name())

	want := "acbd18db4cc2f85cedef654fccc4a4d8"

	if err := option(i); err != nil {
		t.Fatalf("unexpected error: %#v", err)
	}

	if i.version != want {
		t.Fatalf("version=%s, want=%s", i.version, want)
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
		t.Fatalf("unexpected error: %#v", err)
	}

	got, err := i.marshallJSON([]byte{})
	if err != nil {
		t.Fatalf("unexpected error: %#v", err)
	}

	if string(got) != want {
		t.Fatalf("marshallJSON()=%s, want=%s", string(got), want)
	}
}

func TestWithoutLogger(t *testing.T) {
	t.Parallel()

	i := new(Inertia)

	want := log.New(io.Discard, "", 0)

	option := WithoutLogger()

	if err := option(i); err != nil {
		t.Fatalf("unexpected error: %#v", err)
	}

	if !reflect.DeepEqual(i.logger, want) {
		t.Fatalf("logger=%#v, want=%#v", i.logger, want)
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
			t.Fatalf("unexpected error: %#v", err)
		}

		if !reflect.DeepEqual(i.logger, want) {
			t.Fatalf("logger=%#v, want=%#v", i.logger, want)
		}
	})

	t.Run("with logger", func(t *testing.T) {
		t.Parallel()

		i := new(Inertia)

		want := log.New(io.Discard, "foo bar", 0)

		option := WithLogger(want)

		if err := option(i); err != nil {
			t.Fatalf("unexpected error: %#v", err)
		}

		if !reflect.DeepEqual(i.logger, want) {
			t.Fatalf("logger=%#v, want=%#v", i.logger, want)
		}
	})
}

func TestWithContainerID(t *testing.T) {
	t.Parallel()

	i := new(Inertia)

	want := "foo"

	option := WithContainerID(want)

	if err := option(i); err != nil {
		t.Fatalf("unexpected error: %#v", err)
	}

	if i.containerID != want {
		t.Fatalf("containerID=%s, want=%s", i.containerID, want)
	}
}
