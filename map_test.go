package nmap_test

import (
	"fmt"
	"github.com/smartwalle/nmap"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
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

func BenchmarkMapFNV_GetSet(b *testing.B) {
	var m = nmap.New(nmap.WithFNVHash())
	b.ResetTimer()
	go set(m, b)
	get(m, b)
}

func BenchmarkMap_GetSet(b *testing.B) {
	var m = make(map[string]string)
	var mu = sync.RWMutex{}
	b.ResetTimer()
	go func() {
		for i := 0; i < b.N; i++ {
			mu.Lock()
			m["sss"+strconv.Itoa(i)] = "hello"
			mu.Unlock()
		}
	}()

	for i := 0; i < b.N; i++ {
		mu.RLock()
		_ = m["sss"+strconv.Itoa(i)]
		mu.RUnlock()
	}
}

func BenchmarkSyncMap_GetSet(b *testing.B) {
	var m = &sync.Map{}
	go func() {
		for i := 0; i < b.N; i++ {
			m.Store("sss"+strconv.Itoa(i), "hello")
		}
	}()
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

const N = 30000000

func TimeGC() time.Duration {
	start := time.Now()
	runtime.GC()
	return time.Since(start)
}

func TestMap_GC_SV(t *testing.T) {
	var m = make(map[string]int32)
	for i := 0; i < N; i++ {
		n := int32(i)
		m[fmt.Sprintf("%d", n)] = n
	}
	fmt.Printf("With map[string]int32, GC took %s\n", TimeGC())
	_ = m["0"]
}

func TestMap_GC_SP(t *testing.T) {
	var m = make(map[string]*int32)
	for i := 0; i < N; i++ {
		n := int32(i)
		m[fmt.Sprintf("%d", n)] = &n
	}
	fmt.Printf("With map[string]*int32, GC took %s\n", TimeGC())
	_ = m["0"]
}

func TestMap_GC_IV(t *testing.T) {
	var m = make(map[int32]int32)
	for i := 0; i < N; i++ {
		n := int32(i)
		m[n] = n
	}
	fmt.Printf("With map[int32]int32, GC took %s\n", TimeGC())
	_ = m[0]
}

func TestMap_GC_IP(t *testing.T) {
	var m = make(map[int32]*int32)
	for i := 0; i < N; i++ {
		n := int32(i)
		m[n] = &n
	}
	fmt.Printf("With map[int32]*int32, GC took %s\n", TimeGC())
	_ = m[0]
}

func TestNMap_GC_SV(t *testing.T) {
	var m = nmap.New()
	for i := 0; i < N; i++ {
		n := int32(i)
		m.Set(fmt.Sprintf("%d", n), n)
	}
	fmt.Printf("With nmap[string]int32, GC took %s\n", TimeGC())
	_, _ = m.Get("0")
}

func TestNMap_GC_SP(t *testing.T) {
	var m = nmap.New()
	for i := 0; i < N; i++ {
		n := int32(i)
		m.Set(fmt.Sprintf("%d", n), &n)
	}
	fmt.Printf("With nmap[string]*int32, GC took %s\n", TimeGC())
	_, _ = m.Get("0")
}
