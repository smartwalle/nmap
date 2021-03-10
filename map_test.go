package nmap_test

import (
	"github.com/smartwalle/nmap"
	"strconv"
	"testing"
)

func BenchmarkMap_Set(b *testing.B) {
	var m = nmap.New()

	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), strconv.Itoa(i))
	}
	for i := 0; i < b.N; i++ {
		m.Items()
	}
}

func BenchmarkMap_Keys(b *testing.B) {
	m := nmap.New()

	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), strconv.Itoa(i))
	}
	for i := 0; i < b.N; i++ {
		m.Keys()
	}
}
