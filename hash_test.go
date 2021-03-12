package nmap

import (
	"strconv"
	"testing"
)

func BenchmarkFNV1(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FNV1("有中文" + strconv.Itoa(i))
	}
}

func BenchmarkBKDR(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BKDR("有中文" + strconv.Itoa(i))
	}
}
