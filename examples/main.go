package main

import (
	"fmt"
	"github.com/smartwalle/nmap"
)

func main() {
	var m1 = nmap.New[string]()
	m1.Set("s1", "s1")
	m1.Set("s2", "s2")
	m1.Set("s3", "s3")
	m1.Set("s4", "s4")
	m1.Set("s5", "s5")
	m1.Range(func(key string, value string) bool {
		fmt.Println(key, value)
		return true
	})

	var m2 = nmap.NewMap[int, string](func(key int) uint32 {
		return uint32(key % nmap.ShardCount)
	})
	m2.Set(1, "1")
	m2.Set(2, "2")
	m2.Set(3, "3")
	m2.Set(4, "4")
	m2.Set(5, "5")
	m2.Range(func(key int, value string) bool {
		fmt.Println(key, value)
		return true
	})
}
