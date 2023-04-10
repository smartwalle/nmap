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

func BenchmarkStdSyncMap_SetStringString(b *testing.B) {
	var m = &sync.Map{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Store(strconv.Itoa(i), "hello")
	}
}

func BenchmarkStdMap_SetStringString(b *testing.B) {
	var m = make(map[string]string)
	var mu = sync.RWMutex{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mu.Lock()
		m[strconv.Itoa(i)] = "hello"
		mu.Unlock()
	}
}

func BenchmarkStdMap_SetIntString(b *testing.B) {
	var m = make(map[int]string)
	var mu = sync.RWMutex{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mu.Lock()
		m[i] = "hello"
		mu.Unlock()
	}
}

func BenchmarkMap_SetStringString(b *testing.B) {
	var m = nmap.New[string]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set(strconv.Itoa(i), "hello")
	}
}

func BenchmarkMap_SetIntInt(b *testing.B) {
	var m = nmap.NewMap[int, int](func(key int) uint32 {
		return uint32(key % nmap.ShardCount)
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set(i, i)
	}
}

func BenchmarkMapDJB_SetIntString(b *testing.B) {
	var m = nmap.NewMap[int, string](func(key int) uint32 {
		return uint32(key % nmap.ShardCount)
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set(i, "hello")
	}
}

type Human struct {
	Name string
}

func BenchmarkMap_SetIntPointer(b *testing.B) {
	var m = nmap.NewMap[int, *Human](func(key int) uint32 {
		return uint32(key % nmap.ShardCount)
	})
	var h = &Human{
		Name: "name",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set(i, h)
	}
}

func BenchmarkStdSyncMap_GetStringString(b *testing.B) {
	var m = &sync.Map{}
	for i := 0; i < b.N; i++ {
		m.Store(strconv.Itoa(i), "hello")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Load(strconv.Itoa(i))
	}
}

func BenchmarkStdSyncMap_GetIntInt(b *testing.B) {
	var m = &sync.Map{}
	for i := 0; i < b.N; i++ {
		m.Store(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Load(i)
	}
}

func BenchmarkMap_GetStringString(b *testing.B) {
	var m = nmap.New[string]()
	for i := 0; i < b.N; i++ {
		m.Set(strconv.Itoa(i), "hello")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get(strconv.Itoa(i))
	}
}

func BenchmarkMap_GetIntInt(b *testing.B) {
	var m = nmap.NewMap[int, int](func(key int) uint32 {
		return uint32(key % nmap.ShardCount)
	})
	for i := 0; i < b.N; i++ {
		m.Set(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get(i)
	}
}

func BenchmarkStdMap_GetSetStringString(b *testing.B) {
	var m = make(map[string]string)
	var mu = sync.RWMutex{}
	b.ResetTimer()
	go func() {
		for i := 0; i < b.N; i++ {
			mu.Lock()
			m[strconv.Itoa(i)] = "hello"
			mu.Unlock()
		}
	}()

	for i := 0; i < b.N; i++ {
		mu.RLock()
		_ = m[strconv.Itoa(i)]
		mu.RUnlock()
	}
}

func BenchmarkStdSyncMap_GetSetStringString(b *testing.B) {
	var m = &sync.Map{}
	b.ResetTimer()
	go func() {
		for i := 0; i < b.N; i++ {
			m.Store(strconv.Itoa(i), "hello")
		}
	}()
	for i := 0; i < b.N; i++ {
		m.Load(strconv.Itoa(i))
	}
}

func BenchmarkMap_GetSetStringString(b *testing.B) {
	var m = nmap.New[string]()
	b.ResetTimer()
	go func() {
		for i := 0; i < b.N; i++ {
			m.Set(strconv.Itoa(i), "hello")
		}
	}()
	for i := 0; i < b.N; i++ {
		m.Get(strconv.Itoa(i))
	}
}

func BenchmarkStdSyncMap_Range(b *testing.B) {
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

func BenchmarkMap_Range(b *testing.B) {
	var m = nmap.New[string]()
	for i := 0; i < 1000; i++ {
		m.Set(strconv.Itoa(i), "hello")
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Range(func(key string, value string) bool {
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

func TestStdMap_GC_SV(t *testing.T) {
	var m = make(map[string]int32)
	for i := 0; i < N; i++ {
		n := int32(i)
		m[fmt.Sprintf("%d", n)] = n
	}
	runtime.GC()
	t.Logf("With %T, GC took %s\n", m, TimeGC())
	_ = m["0"]
}

func TestStdMap_GC_SP(t *testing.T) {
	var m = make(map[string]*int32)
	for i := 0; i < N; i++ {
		n := int32(i)
		m[fmt.Sprintf("%d", n)] = &n
	}
	runtime.GC()
	t.Logf("With %T, GC took %s\n", m, TimeGC())
	_ = m["0"]
}

func TestStdMap_GC_IV(t *testing.T) {
	var m = make(map[int32]int32)
	for i := 0; i < N; i++ {
		n := int32(i)
		m[n] = n
	}
	runtime.GC()
	t.Logf("With %T, GC took %s\n", m, TimeGC())
	_ = m[0]
}

func TestStdMap_GC_IP(t *testing.T) {
	var m = make(map[int32]*int32)
	for i := 0; i < N; i++ {
		n := int32(i)
		m[n] = &n
	}
	runtime.GC()
	t.Logf("With %T, GC took %s\n", m, TimeGC())
	_ = m[0]
}

func TestMap_GC_SV(t *testing.T) {
	var m = nmap.New[int32]()
	for i := 0; i < N; i++ {
		n := int32(i)
		m.Set(fmt.Sprintf("%d", n), n)
	}
	runtime.GC()
	t.Logf("With %T, GC took %s\n", m, TimeGC())
	_, _ = m.Get("0")
}

func TestMap_GC_SP(t *testing.T) {
	var m = nmap.New[*int32]()
	for i := 0; i < N; i++ {
		n := int32(i)
		m.Set(fmt.Sprintf("%d", n), &n)
	}
	runtime.GC()
	t.Logf("With %T, GC took %s  \n", m, TimeGC())
	_, _ = m.Get("0")
}

func TestMap_GC_IV(t *testing.T) {
	var m = nmap.NewMap[int32, int32](func(key int32) uint32 {
		return uint32(key % nmap.ShardCount)
	})
	for i := 0; i < N; i++ {
		n := int32(i)
		m.Set(n, n)
	}
	runtime.GC()
	t.Logf("With %T, GC took %s\n", m, TimeGC())
	_, _ = m.Get(0)
}

func TestMap_GC_IP(t *testing.T) {
	var m = nmap.NewMap[int32, *int32](func(key int32) uint32 {
		return uint32(key % nmap.ShardCount)
	})
	for i := 0; i < N; i++ {
		n := int32(i)
		m.Set(n, &n)
	}
	runtime.GC()
	t.Logf("With %T, GC took %s\n", m, TimeGC())
	_, _ = m.Get(0)
}
