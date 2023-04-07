package nmap

import "math/rand"

func FNV1() func(key string) uint32 {
	var seed = uint32(16777619)
	return func(key string) uint32 {
		var hash = uint32(2166136261)
		var keyLen = len(key)
		for i := 0; i < keyLen; i++ {
			hash *= seed
			hash ^= uint32(key[i])
		}
		return hash
	}
}

func BKDR() func(key string) uint32 {
	var seed = uint32(131)
	return func(key string) uint32 {
		var hash uint32 = 0
		var keyLen = len(key)
		for i := 0; i < keyLen; i++ {
			hash *= seed
			hash += uint32(key[i])
		}
		return hash
	}
}

func DJB() func(key string) uint32 {
	var seed = rand.Uint32() + 5381
	return func(key string) uint32 {
		var (
			l = uint32(len(key))
			d = seed + l
			i = uint32(0)
		)
		if l >= 4 {
			for i < l-4 {
				d = (d * 33) ^ uint32(key[i])
				d = (d * 33) ^ uint32(key[i+1])
				d = (d * 33) ^ uint32(key[i+2])
				d = (d * 33) ^ uint32(key[i+3])
				i += 4
			}
		}
		switch l - i {
		case 1:
		case 2:
			d = (d * 33) ^ uint32(key[i])
		case 3:
			d = (d * 33) ^ uint32(key[i])
			d = (d * 33) ^ uint32(key[i+1])
		case 4:
			d = (d * 33) ^ uint32(key[i])
			d = (d * 33) ^ uint32(key[i+1])
			d = (d * 33) ^ uint32(key[i+2])
		}
		return d ^ (d >> 16)
	}
}
