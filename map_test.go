package nmap

import (
	"strconv"
	"testing"
)

func BenchmarkMapFNV_Set(b *testing.B) {
	var m = New(WithFNVHash())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set("sss"+strconv.Itoa(i), "hello")
	}
}

func BenchmarkMapBKDR_Set(b *testing.B) {
	var m = New(WithBKDRHash())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set("sss"+strconv.Itoa(i), "hello")
	}
}
