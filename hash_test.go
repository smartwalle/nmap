package nmap

import (
	"math/rand"
	"strconv"
	"testing"
)

func BenchmarkFNV1(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FNV1(kFNVSeed, "有中文"+strconv.Itoa(i))
	}
}

func BenchmarkBKDR(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BKDR(kBKDRSeed, "有中文"+strconv.Itoa(i))
	}
}

func BenchmarkDJB(b *testing.B) {
	b.ResetTimer()
	var seed = rand.Uint32()
	for i := 0; i < b.N; i++ {
		DJB(seed, "有中文"+strconv.Itoa(i))
	}
}
