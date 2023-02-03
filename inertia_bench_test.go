package gonertia

import (
	"encoding/json"
	"net/http"
	"testing"
)

func BenchmarkInertia_inertiaContainer(b *testing.B) {
	inertia := Inertia{
		containerID: "foobar",
	}

	page, err := json.Marshal(map[string]any{
		"foo": "bar",
	})
	if err != nil {
		b.Fatalf("unexpected error: %#v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		inertia.inertiaContainerHTML(page)
	}
}

func BenchmarkInertia_buildPage(b *testing.B) {
	inertia := Inertia{}

	req, err := http.NewRequest("GET", "/foo", nil)
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
