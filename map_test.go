package nmap_test

import (
	"github.com/smartwalle/nmap"
	"strconv"
	"sync"
	"testing"
)

func set(m *nmap.Map, b *testing.B) {
	for i := 0; i < b.N; i++ {
		m.Set("sss"+strconv.Itoa(i), "hello")
	}
}

func BenchmarkMapFNV_Set(b *testing.B) {
	var m = nmap.New(nmap.WithFNVHash())
	b.ResetTimer()
	set(m, b)
}

func BenchmarkMapBKDR_Set(b *testing.B) {
	var m = nmap.New(nmap.WithBKDRHash())
	b.ResetTimer()
	set(m, b)
}

func BenchmarkMapDJB_Set(b *testing.B) {
	var m = nmap.New(nmap.WithDJBHash())
	b.ResetTimer()
	set(m, b)
}

func BenchmarkSyncMap_Set(b *testing.B) {
	var m = &sync.Map{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Store("sss"+strconv.Itoa(i), "hello")
	}
}

func get(m *nmap.Map, b *testing.B) {
	for i := 0; i < b.N; i++ {
		m.Get("sss" + strconv.Itoa(i))
	}
}

func BenchmarkMapFNV_Get(b *testing.B) {
	var m = nmap.New(nmap.WithFNVHash())
	set(m, b)
	b.ResetTimer()
	get(m, b)
}

func BenchmarkMapBKDR_Get(b *testing.B) {
	var m = nmap.New(nmap.WithBKDRHash())
	set(m, b)
	b.ResetTimer()
	get(m, b)
}

func BenchmarkMapDJB_Get(b *testing.B) {
	var m = nmap.New(nmap.WithDJBHash())
	set(m, b)
	b.ResetTimer()
	get(m, b)
}

func BenchmarkSyncMap_Get(b *testing.B) {
	var m = &sync.Map{}
	for i := 0; i < b.N; i++ {
		m.Store("sss"+strconv.Itoa(i), "hello")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Load("sss" + strconv.Itoa(i))
	}
}

func BenchmarkMap_Range(b *testing.B) {
	var m = nmap.New()
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
