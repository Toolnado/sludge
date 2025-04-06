package token

import "testing"

func BenchmarkNew(b *testing.B) {
	pos := Position{
		Filename: "test.go",
		Line:     1,
		Column:   1,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		New(pos, IDENTIFIER, "test", "test")
	}
}
