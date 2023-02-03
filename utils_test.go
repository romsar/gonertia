package gonertia

import (
	"os"
	"reflect"
	"testing"
)

func Test_set(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		got := set[string](nil)
		var want map[string]struct{}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("set()=%#v, want=%#v", got, want)
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		got := set[string]([]string{})
		var want map[string]struct{}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("set()=%#v, want=%#v", got, want)
		}
	})

	t.Run("duplicates", func(t *testing.T) {
		t.Parallel()

		got := set[string]([]string{"foo", "foo"})
		want := map[string]struct{}{
			"foo": {},
		}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("set()=%#v, want=%#v", got, want)
		}
	})

	t.Run("strings", func(t *testing.T) {
		t.Parallel()

		got := set[string]([]string{"foo", "bar"})
		want := map[string]struct{}{
			"foo": {},
			"bar": {},
		}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("set()=%#v, want=%#v", got, want)
		}
	})

	t.Run("integers", func(t *testing.T) {
		t.Parallel()

		got := set[int]([]int{123, 456})
		want := map[int]struct{}{
			123: {},
			456: {},
		}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("set()=%#v, want=%#v", got, want)
		}
	})
}

func Test_firstOr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		items    []string
		fallback string
		want     string
	}{
		{
			"nil",
			nil,
			"zoo",
			"zoo",
		},
		{
			"empty",
			[]string{},
			"zoo",
			"zoo",
		},
		{
			"not empty",
			[]string{"foo", "bar"},
			"zoo",
			"foo",
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := firstOr[string](tt.items, tt.fallback)
			if got != tt.want {
				t.Fatalf("firstOr()=%s, want=%s", got, tt.want)
			}
		})
	}
}

func Test_md5(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		str  string
		want string
	}{
		{
			"empty",
			"",
			"d41d8cd98f00b204e9800998ecf8427e",
		},
		{
			"not empty",
			"foo",
			"acbd18db4cc2f85cedef654fccc4a4d8",
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			got := md5(tt.str)
			if got != tt.want {
				t.Fatalf("md5()=%s, want=%s", got, tt.want)
			}
		})
	}
}

func Test_md5File(t *testing.T) {
	t.Parallel()

	f, err := os.CreateTemp("", "gonertia")
	if err != nil {
		t.Fatalf("unexpected error: %#v", err)
	}

	closed := false

	t.Cleanup(func() {
		if !closed {
			if err := f.Close(); err != nil {
				t.Fatalf("unexpected error: %#v", err)
			}
		}

		if err := os.Remove(f.Name()); err != nil {
			t.Fatalf("unexpected error: %#v", err)
		}
	})

	if _, err := f.WriteString("foo"); err != nil {
		t.Fatalf("unexpected error: %#v", err)
	}

	if err := f.Close(); err != nil {
		t.Fatalf("unexpected error: %#v", err)
	}

	closed = true

	got, err := md5File(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %#v", err)
	}

	want := md5("foo")
	if got != want {
		t.Fatalf("md5File()=%s, want=%s", got, want)
	}
}
