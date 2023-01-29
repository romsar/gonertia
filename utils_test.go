package gonertia

import (
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
			t.Fatalf("got=%#v, want %#v", got, want)
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		got := set[string]([]string{})
		var want map[string]struct{}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got=%#v, want %#v", got, want)
		}
	})

	t.Run("duplicates", func(t *testing.T) {
		t.Parallel()

		got := set[string]([]string{"foo", "foo"})
		want := map[string]struct{}{
			"foo": {},
		}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got=%#v, want %#v", got, want)
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
			t.Fatalf("got=%#v, want %#v", got, want)
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
			t.Fatalf("got=%#v, want %#v", got, want)
		}
	})

	t.Run("floats", func(t *testing.T) {
		t.Parallel()

		got := set[float64]([]float64{123.45, 456.78})
		want := map[float64]struct{}{
			123.45: {},
			456.78: {},
		}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got=%#v, want %#v", got, want)
		}
	})
}
