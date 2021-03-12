package nmap

type Hash func(key string) uint32

const (
	kFNVPrime = uint32(16777619)
	kBKDRSeed = uint32(131)
)

func FNV1(key string) uint32 {
	var hash = uint32(2166136261)
	for i := 0; i < len(key); i++ {
		hash *= kFNVPrime
		hash ^= uint32(key[i])
	}
	return hash
}

func BKDR(key string) uint32 {
	var hash uint32 = 0
	for i := 0; i < len(key); i++ {
		hash *= kBKDRSeed
		hash += uint32(key[i])
	}
	return hash
}
