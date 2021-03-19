package nmap

import (
	"strconv"
	"sync"
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

func BenchmarkSyncMap_Set(b *testing.B) {
	var m = &sync.Map{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Store("sss"+strconv.Itoa(i), "hello")
	}
}

func BenchmarkMap_Range(b *testing.B) {
	var m = New()
	for i := 0; i < 1000; i++ {
		m.Set(strconv.Itoa(i), "hello")
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Range(func(key string, value interface{}) bool {
			return true
		})
	}
}

func BenchmarkSyncMap_Range(b *testing.B) {
	var m = &sync.Map{}
	for i := 0; i < 1000; i++ {
		m.Store(strconv.Itoa(i), "hello")
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Range(func(key, value interface{}) bool {
			return true
		})
	}
}
