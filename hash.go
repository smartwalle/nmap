package nmap

type Hash func(seed uint32, key string) uint32

const (
	kFNVSeed  = uint32(16777619)
	kBKDRSeed = uint32(131)
)

func FNV1(seed uint32, key string) uint32 {
	var hash = uint32(2166136261)
	var keyLen = len(key)
	for i := 0; i < keyLen; i++ {
		hash *= seed
		hash ^= uint32(key[i])
	}
	return hash
}

func BKDR(seed uint32, key string) uint32 {
	var hash uint32 = 0
	var keyLen = len(key)
	for i := 0; i < keyLen; i++ {
		hash *= seed
		hash += uint32(key[i])
	}
	return hash
}

func DJB(seed uint32, k string) uint32 {
	var (
		l = uint32(len(k))
		d = 5381 + seed + l
		i = uint32(0)
	)
	if l >= 4 {
		for i < l-4 {
			d = (d * 33) ^ uint32(k[i])
			d = (d * 33) ^ uint32(k[i+1])
			d = (d * 33) ^ uint32(k[i+2])
			d = (d * 33) ^ uint32(k[i+3])
			i += 4
		}
	}
	switch l - i {
	case 1:
	case 2:
		d = (d * 33) ^ uint32(k[i])
	case 3:
		d = (d * 33) ^ uint32(k[i])
		d = (d * 33) ^ uint32(k[i+1])
	case 4:
		d = (d * 33) ^ uint32(k[i])
		d = (d * 33) ^ uint32(k[i+1])
		d = (d * 33) ^ uint32(k[i+2])
	}
	return d ^ (d >> 16)
}
