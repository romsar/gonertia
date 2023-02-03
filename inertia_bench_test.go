package gonertia

import (
	"encoding/json"
	"testing"
)

func BenchmarkInertia_inertiaContainerHTML(b *testing.B) {
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
