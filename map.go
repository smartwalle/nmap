package nmap

// 参考 https://github.com/orcaman/concurrent-map

import (
	"math/rand"
	"sync"
)

const (
	kShardCount = uint32(32)
)

type Option func(opt *option)

func WithShardCount(count uint32) Option {
	return func(opt *option) {
		opt.shard = count
	}
}

func WithFNVHash() Option {
	return func(opt *option) {
		opt.hashSeed = kFNVSeed
		opt.hash = FNV1
	}
}

func WithBKDRHash() Option {
	return func(opt *option) {
		opt.hashSeed = kBKDRSeed
		opt.hash = BKDR
	}
}

func WithDJBHash() Option {
	return func(opt *option) {
		opt.hashSeed = rand.Uint32()
		opt.hash = DJB
	}
}

type option struct {
	hashSeed uint32
	hash     Hash
	shard    uint32
}

type Map[T any] struct {
	*option
	shards []*shardMap[T]
}

type shardMap[T any] struct {
	*sync.RWMutex
	items map[string]T
}

func New[T any](opts ...Option) *Map[T] {
	var m = &Map[T]{}
	m.option = &option{}

	for _, opt := range opts {
		opt(m.option)
	}
	if m.hash == nil {
		WithDJBHash()(m.option)
	}
	if m.shard == 0 {
		m.shard = kShardCount
	}

	m.shards = make([]*shardMap[T], m.shard)
	for i := uint32(0); i < m.shard; i++ {
		m.shards[i] = &shardMap[T]{RWMutex: &sync.RWMutex{}, items: make(map[string]T)}
	}
	return m
}

func (this *Map[T]) getShard(key string) *shardMap[T] {
	var index = this.hash(this.hashSeed, key) % this.shard
	return this.shards[index]
}

func (this *Map[T]) Set(key string, value T) {
	var shard = this.getShard(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}

func (this *Map[T]) SetNx(key string, value T) bool {
	var shard = this.getShard(key)
	shard.Lock()
	var _, ok = shard.items[key]
	if ok == false {
		shard.items[key] = value
		shard.Unlock()
		return true
	}
	shard.Unlock()
	return false
}

func (this *Map[T]) Get(key string) (T, bool) {
	var shard = this.getShard(key)
	shard.RLock()
	var value, ok = shard.items[key]
	shard.RUnlock()
	return value, ok
}

func (this *Map[T]) Exists(key string) bool {
	var shard = this.getShard(key)
	shard.RLock()
	var _, ok = shard.items[key]
	shard.RUnlock()
	return ok
}

func (this *Map[T]) RemoveAll() {
	for i := uint32(0); i < this.shard; i++ {
		var shard = this.shards[i]
		shard.Lock()
		shard.items = make(map[string]T)
		shard.Unlock()
	}
}

func (this *Map[T]) Remove(key string) {
	var shard = this.getShard(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

func (this *Map[T]) Pop(key string) (T, bool) {
	var shard = this.getShard(key)
	shard.Lock()
	var value, ok = shard.items[key]
	delete(shard.items, key)
	shard.Unlock()
	return value, ok
}

func (this *Map[T]) Len() int {
	var count = 0
	for i := uint32(0); i < this.shard; i++ {
		var shard = this.shards[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

func (this *Map[T]) Range(f func(key string, value T) bool) {
	if f == nil {
		return
	}
	for i := uint32(0); i < this.shard; i++ {
		var shard = this.shards[i]
		shard.RLock()
		for k, v := range shard.items {
			shard.RUnlock()
			if f(k, v) == false {
				return
			}
			shard.RLock()
		}
		shard.RUnlock()
	}
}

func (this *Map[T]) Items() map[string]T {
	var nMap = make(map[string]T, this.Len())
	for i := uint32(0); i < this.shard; i++ {
		var shard = this.shards[i]
		shard.RLock()
		for k, v := range shard.items {
			nMap[k] = v
		}
		shard.RUnlock()
	}
	return nMap
}

func (this *Map[T]) Keys() []string {
	var nKeys = make([]string, 0, this.Len())
	for i := uint32(0); i < this.shard; i++ {
		var shard = this.shards[i]
		shard.RLock()
		for k := range shard.items {
			nKeys = append(nKeys, k)
		}
		shard.RUnlock()
	}
	return nKeys
}
