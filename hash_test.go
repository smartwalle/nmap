package nmap_test

import (
	"github.com/smartwalle/nmap"
	"math/rand"
	"strconv"
	"testing"
)

const (
	kFNVSeed  = uint32(16777619)
	kBKDRSeed = uint32(131)
)

func BenchmarkFNV1(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nmap.FNV1(kFNVSeed, "有中文"+strconv.Itoa(i))
	}
}

func BenchmarkBKDR(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nmap.BKDR(kBKDRSeed, "有中文"+strconv.Itoa(i))
	}
}

func BenchmarkDJB(b *testing.B) {
	b.ResetTimer()
	var seed = rand.Uint32()
	for i := 0; i < b.N; i++ {
		nmap.DJB(seed, "有中文"+strconv.Itoa(i))
	}
}
