package gonertia

import (
	"net/http"
	"testing"
)

func TestIsInertiaRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		r    *http.Request
		want bool
	}{
		{
			"positive",
			&http.Request{
				Header: http.Header{"X-Inertia": []string{"foo"}},
			},
			true,
		},
		{
			"negative",
			&http.Request{},
			false,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := IsInertiaRequest(tt.r)

			if got != tt.want {
				t.Fatalf("got=%#v, want=%#v", got, tt.want)
			}
		})
	}
}
