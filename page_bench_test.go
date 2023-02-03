package gonertia

import (
	"context"
	"net/http"
	"testing"
)

func BenchmarkInertia_buildPage(b *testing.B) {
	inertia := Inertia{}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/foo", nil)
	if err != nil {
		b.Fatalf("unexpected error: %#v", err)
	}

	props := Props{
		"foo": "bar",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = inertia.buildPage(req, "foobar", props)
	}
}
