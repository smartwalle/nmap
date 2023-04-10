package nmap_test

import (
	"github.com/smartwalle/nmap"
	"strconv"
	"testing"
)

func BenchmarkFNV1(b *testing.B) {
	var sharding = nmap.FNV1Sharding()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sharding("有中文" + strconv.Itoa(i))
	}
}

func BenchmarkBKDR(b *testing.B) {
	var sharding = nmap.BKDRSharding()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sharding("有中文" + strconv.Itoa(i))
	}
}

func BenchmarkDJB(b *testing.B) {
	var sharding = nmap.DJBSharding()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sharding("有中文" + strconv.Itoa(i))
	}
}
