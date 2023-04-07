package nmap_test

import (
	"github.com/smartwalle/nmap"
	"strconv"
	"testing"
)

func BenchmarkFNV1(b *testing.B) {
	var hash = nmap.FNV1()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash("有中文" + strconv.Itoa(i))
	}
}

func BenchmarkBKDR(b *testing.B) {
	var hash = nmap.BKDR()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash("有中文" + strconv.Itoa(i))
	}
}

func BenchmarkDJB(b *testing.B) {
	var hash = nmap.DJB()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash("有中文" + strconv.Itoa(i))
	}
}
