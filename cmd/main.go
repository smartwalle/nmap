package main

import (
	"fmt"
	"github.com/smartwalle/nmap"
)

func main() {
	var m = nmap.New()
	m.Set("h1", "hh1")
	m.Set("h2", "hh2")
	m.Set("h3", "hh3")
	m.Set("h4", "hh4")
	m.Set("h5", "hh5")

	m.Range(func(key string, value interface{}) bool {
		fmt.Println(key, value)
		return true
	})
}
