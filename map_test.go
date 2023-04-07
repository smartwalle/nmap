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

func BenchmarkSyncMap_SetStringString(b *testing.B) {
	var m = &sync.Map{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Store(strconv.Itoa(i), "hello")
	}
}

func BenchmarkMap_SetStringString(b *testing.B) {
	var m = make(map[string]string)
	var mu = sync.RWMutex{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mu.Lock()
		m[strconv.Itoa(i)] = "hello"
		mu.Unlock()
	}
}

func BenchmarkMap_SetIntString(b *testing.B) {
	var m = make(map[int]string)
	var mu = sync.RWMutex{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mu.Lock()
		m[i] = "hello"
		mu.Unlock()
	}
}

// setStringString
// key - string
// value - string
func setStringString(m *nmap.Map[string, string], b *testing.B) {
	for i := 0; i < b.N; i++ {
		m.Set(strconv.Itoa(i), "hello")
	}
}

func BenchmarMap_SetStringString(b *testing.B, opts ...nmap.Option) {
	var m = nmap.New[string, string](opts...)
	b.ResetTimer()
	setStringString(m, b)
}

func BenchmarkMapFNV_SetStringString(b *testing.B) {
	BenchmarMap_SetStringString(b, nmap.WithFNVHash())
}

func BenchmarkMapBKDR_SetStringString(b *testing.B) {
	BenchmarMap_SetStringString(b, nmap.WithBKDRHash())
}

func BenchmarkMapDJB_SetStringString(b *testing.B) {
	BenchmarMap_SetStringString(b, nmap.WithDJBHash())
}

// setIntInt
// key - int
// value - int
func setIntInt(m *nmap.Map[int, int], b *testing.B) {
	for i := 0; i < b.N; i++ {
		m.Set(i, i)
	}
}

func BenchmarMap_SetIntInt(b *testing.B, opts ...nmap.Option) {
	//var m = nmap.New[int, int](opts...)
	//b.ResetTimer()
	//setIntInt(m, b)
}

func BenchmarkMapDJB_SetIntInt(b *testing.B) {
	BenchmarMap_SetIntInt(b, nmap.WithDJBHash())
}

// setIntString
// key - int
// value - string
func setIntString(m *nmap.Map[int, string], b *testing.B) {
	for i := 0; i < b.N; i++ {
		m.Set(i, "hello")
	}
}

func BenchmarMap_SetIntString(b *testing.B, opts ...nmap.Option) {
	//var m = nmap.New[int, string](opts...)
	//b.ResetTimer()
	//setIntString(m, b)
}

func BenchmarkMapDJB_SetIntString(b *testing.B) {
	BenchmarMap_SetIntString(b, nmap.WithDJBHash())
}

// setIntString
// key - int
// value - pointer
type Human struct {
	Name string
}

func setIntPointer(m *nmap.Map[int, *Human], b *testing.B) {
	var h = &Human{
		Name: "haha",
	}
	for i := 0; i < b.N; i++ {
		m.Set(i, h)
	}
}

func BenchmarMap_SetIntPointer(b *testing.B, opts ...nmap.Option) {
	//var m = nmap.New[int, *Human](opts...)
	//b.ResetTimer()
	//setIntPointer(m, b)
}

func BenchmarkMapDJB_SetIntPointer(b *testing.B) {
	BenchmarMap_SetIntPointer(b, nmap.WithDJBHash())
}

//func get(m *nmap.Map[string], b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		m.Get("sss" + strconv.Itoa(i))
//	}
//}
//
//func BenchmarkMapFNV_Get(b *testing.B) {
//	var m = nmap.New[string](nmap.WithFNVHash())
//	set(m, b)
//	b.ResetTimer()
//	get(m, b)
//}
//
//func BenchmarkMapBKDR_Get(b *testing.B) {
//	var m = nmap.New[string](nmap.WithBKDRHash())
//	set(m, b)
//	b.ResetTimer()
//	get(m, b)
//}
//
//func BenchmarkMapDJB_Get(b *testing.B) {
//	var m = nmap.New[string](nmap.WithDJBHash())
//	set(m, b)
//	b.ResetTimer()
//	get(m, b)
//}
//
//func BenchmarkSyncMap_Get(b *testing.B) {
//	var m = &sync.Map{}
//	for i := 0; i < b.N; i++ {
//		m.Store("sss"+strconv.Itoa(i), "hello")
//	}
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		m.Load("sss" + strconv.Itoa(i))
//	}
//}
//
//func BenchmarkMapFNV_GetSet(b *testing.B) {
//	var m = nmap.New[string](nmap.WithFNVHash())
//	b.ResetTimer()
//	go set(m, b)
//	get(m, b)
//}
//
//func BenchmarkMapDJB_GetSet(b *testing.B) {
//	var m = nmap.New[string](nmap.WithDJBHash())
//	b.ResetTimer()
//	go set(m, b)
//	get(m, b)
//}
//
//func BenchmarkMap_GetSet(b *testing.B) {
//	var m = make(map[string]string)
//	var mu = sync.RWMutex{}
//	b.ResetTimer()
//	go func() {
//		for i := 0; i < b.N; i++ {
//			mu.Lock()
//			m["sss"+strconv.Itoa(i)] = "hello"
//			mu.Unlock()
//		}
//	}()
//
//	for i := 0; i < b.N; i++ {
//		mu.RLock()
//		_ = m["sss"+strconv.Itoa(i)]
//		mu.RUnlock()
//	}
//}
//
//func BenchmarkSyncMap_GetSet(b *testing.B) {
//	var m = &sync.Map{}
//	go func() {
//		for i := 0; i < b.N; i++ {
//			m.Store("sss"+strconv.Itoa(i), "hello")
//		}
//	}()
//	for i := 0; i < b.N; i++ {
//		m.Load("sss" + strconv.Itoa(i))
//	}
//}
//
//func BenchmarkDJB_Range(b *testing.B) {
//	var m = nmap.New[string](nmap.WithDJBHash())
//	for i := 0; i < 1000; i++ {
//		m.Set(strconv.Itoa(i), "hello")
//	}
//	b.ResetTimer()
//
//	for i := 0; i < b.N; i++ {
//		m.Range(func(key string, value string) bool {
//			return true
//		})
//	}
//}
//
//func BenchmarkSyncMap_Range(b *testing.B) {
//	var m = &sync.Map{}
//	for i := 0; i < 1000; i++ {
//		m.Store(strconv.Itoa(i), "hello")
//	}
//	b.ResetTimer()
//
//	for i := 0; i < b.N; i++ {
//		m.Range(func(key, value interface{}) bool {
//			return true
//		})
//	}
//}

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
	runtime.GC()
	t.Logf("With %T, GC took %s\n", m, TimeGC())
	_ = m["0"]
}

func TestMap_GC_SP(t *testing.T) {
	var m = make(map[string]*int32)
	for i := 0; i < N; i++ {
		n := int32(i)
		m[fmt.Sprintf("%d", n)] = &n
	}
	runtime.GC()
	t.Logf("With %T, GC took %s\n", m, TimeGC())
	_ = m["0"]
}

func TestMap_GC_IV(t *testing.T) {
	var m = make(map[int32]int32)
	for i := 0; i < N; i++ {
		n := int32(i)
		m[n] = n
	}
	runtime.GC()
	t.Logf("With %T, GC took %s\n", m, TimeGC())
	_ = m[0]
}

func TestMap_GC_IP(t *testing.T) {
	var m = make(map[int32]*int32)
	for i := 0; i < N; i++ {
		n := int32(i)
		m[n] = &n
	}
	runtime.GC()
	t.Logf("With %T, GC took %s\n", m, TimeGC())
	_ = m[0]
}

//func TestMap_GC_NMap_SV(t *testing.T) {
//	var m = nmap.New[int32, int32]()
//	for i := 0; i < N; i++ {
//		n := int32(i)
//		m.Set(n, n)
//	}
//	runtime.GC()
//	t.Logf("With %T, GC took %s\n", m, TimeGC())
//	_, _ = m.Get(0)
//}
//
//func TestMap_GC_NMap_SP(t *testing.T) {
//	var m = nmap.New[int32, *int32]()
//	for i := 0; i < N; i++ {
//		n := int32(i)
//		m.Set(n, &n)
//	}
//	runtime.GC()
//	t.Logf("With %T, GC took %s  \n", m, TimeGC())
//	_, _ = m.Get(0)
//}
