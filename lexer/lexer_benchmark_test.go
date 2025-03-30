package lexer

import (
	"strings"
	"testing"
)

func BenchmarkScanTokens(b *testing.B) {
	input := `let x = 42;
if (x > 0) {
    return x;
} else {
    return -x;
}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := New(strings.NewReader(input))
		l.ScanTokens()
	}
}
