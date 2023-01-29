package gonertia

import (
	"testing"
)

func Benchmark_md5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		md5("foo bar")
	}
}
