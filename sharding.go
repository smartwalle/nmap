package nmap

import "math/rand"

func FNV1Sharding() func(key string) uint32 {
	var seed = uint32(16777619)
	return func(key string) uint32 {
		var hash = uint32(2166136261)
		var keyLen = len(key)
		for i := 0; i < keyLen; i++ {
			hash *= seed
			hash ^= uint32(key[i])
		}
		return hash % ShardCount
	}
}

func BKDRSharding() func(key string) uint32 {
	var seed = uint32(131)
	return func(key string) uint32 {
		var hash uint32 = 0
		var keyLen = len(key)
		for i := 0; i < keyLen; i++ {
			hash *= seed
			hash += uint32(key[i])
		}
		return hash % ShardCount
	}
}

func DJBSharding() func(key string) uint32 {
	var seed = rand.Uint32() + 5381
	return func(key string) uint32 {
		var (
			l    = uint32(len(key))
			hash = seed + l
			i    = uint32(0)
		)
		if l >= 4 {
			for i < l-4 {
				hash = (hash * 33) ^ uint32(key[i])
				hash = (hash * 33) ^ uint32(key[i+1])
				hash = (hash * 33) ^ uint32(key[i+2])
				hash = (hash * 33) ^ uint32(key[i+3])
				i += 4
			}
		}
		switch l - i {
		case 1:
		case 2:
			hash = (hash * 33) ^ uint32(key[i])
		case 3:
			hash = (hash * 33) ^ uint32(key[i])
			hash = (hash * 33) ^ uint32(key[i+1])
		case 4:
			hash = (hash * 33) ^ uint32(key[i])
			hash = (hash * 33) ^ uint32(key[i+1])
			hash = (hash * 33) ^ uint32(key[i+2])
		}
		return (hash ^ (hash >> 16)) % ShardCount
	}
}
