package nmap

// 参考 https://github.com/orcaman/concurrent-map

import (
	"sync"
)

const (
	kShardCount = uint32(32)
)

type Option func(opt *options)

func WithShardCount(count uint32) Option {
	return func(opt *options) {
		opt.shard = count
	}
}

func WithFNVHash() Option {
	return func(opts *options) {
		//opts.hash = FNV1()
	}
}

func WithBKDRHash() Option {
	return func(opts *options) {
		//opts.hash = BKDR()
	}
}

func WithDJBHash() Option {
	return func(opts *options) {
		//opts.hash = DJB()
	}
}

type options struct {
	shard uint32
}

type Key interface {
	~string | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Map[K Key, V any] struct {
	*options
	hash   func(key K) uint32
	shards []*shardMap[K, V]
}

type shardMap[K Key, V any] struct {
	elements map[K]V
	sync.RWMutex
}

func New[V any](opts ...Option) *Map[string, V] {
	return NewMap[string, V](DJB(), opts...)
}

func NewMap[K Key, V any](hash func(key K) uint32, opts ...Option) *Map[K, V] {
	var m = &Map[K, V]{}
	m.options = &options{}

	for _, opt := range opts {
		if opt != nil {
			opt(m.options)
		}
	}
	m.hash = hash
	//if m.hash == nil {
	//	WithDJBHash()(m.options)
	//}
	if m.shard == 0 {
		m.shard = kShardCount
	}

	m.shards = make([]*shardMap[K, V], m.shard)
	for i := uint32(0); i < m.shard; i++ {
		m.shards[i] = &shardMap[K, V]{elements: make(map[K]V)}
	}
	return m
}

func (this *Map[K, V]) getShard(key K) *shardMap[K, V] {
	var index = this.hash(key) % this.shard
	return this.shards[index]
}

func (this *Map[K, V]) Set(key K, value V) {
	var shard = this.getShard(key)
	shard.Lock()
	shard.elements[key] = value
	shard.Unlock()
}

func (this *Map[K, V]) SetNx(key K, value V) bool {
	var shard = this.getShard(key)
	shard.Lock()
	var _, ok = shard.elements[key]
	if ok == false {
		shard.elements[key] = value
		shard.Unlock()
		return true
	}
	shard.Unlock()
	return false
}

func (this *Map[K, V]) Get(key K) (V, bool) {
	var shard = this.getShard(key)
	shard.RLock()
	var value, ok = shard.elements[key]
	shard.RUnlock()
	return value, ok
}

func (this *Map[K, V]) Exists(key K) bool {
	var shard = this.getShard(key)
	shard.RLock()
	var _, ok = shard.elements[key]
	shard.RUnlock()
	return ok
}

func (this *Map[K, V]) RemoveAll() {
	for i := uint32(0); i < this.shard; i++ {
		var shard = this.shards[i]
		shard.Lock()
		shard.elements = make(map[K]V)
		shard.Unlock()
	}
}

func (this *Map[K, V]) Remove(key K) {
	var shard = this.getShard(key)
	shard.Lock()
	delete(shard.elements, key)
	shard.Unlock()
}

func (this *Map[K, V]) Pop(key K) (V, bool) {
	var shard = this.getShard(key)
	shard.Lock()
	var value, ok = shard.elements[key]
	delete(shard.elements, key)
	shard.Unlock()
	return value, ok
}

func (this *Map[K, V]) Len() int {
	var count = 0
	for i := uint32(0); i < this.shard; i++ {
		var shard = this.shards[i]
		shard.RLock()
		count += len(shard.elements)
		shard.RUnlock()
	}
	return count
}

func (this *Map[K, V]) Range(f func(key K, value V) bool) {
	if f == nil {
		return
	}
	for i := uint32(0); i < this.shard; i++ {
		var shard = this.shards[i]
		shard.RLock()
		for k, v := range shard.elements {
			shard.RUnlock()
			if f(k, v) == false {
				return
			}
			shard.RLock()
		}
		shard.RUnlock()
	}
}

func (this *Map[K, V]) Elements() map[K]V {
	var nMap = make(map[K]V, this.Len())
	for i := uint32(0); i < this.shard; i++ {
		var shard = this.shards[i]
		shard.RLock()
		for k, v := range shard.elements {
			nMap[k] = v
		}
		shard.RUnlock()
	}
	return nMap
}

func (this *Map[K, V]) Keys() []K {
	var nKeys = make([]K, 0, this.Len())
	for i := uint32(0); i < this.shard; i++ {
		var shard = this.shards[i]
		shard.RLock()
		for k := range shard.elements {
			nKeys = append(nKeys, k)
		}
		shard.RUnlock()
	}
	return nKeys
}
