package piimask

import (
	"testing"
)

func BenchmarkMaskEmail(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = MaskEmail("alexander@example.com")
	}
}

func BenchmarkMaskIP(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = MaskIP("192.168.1.100")
	}
}

func BenchmarkMaskToken(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = MaskToken("dev_tok_9f823a789bc123")
	}
}