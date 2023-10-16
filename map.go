package nmap

// 参考 https://github.com/orcaman/concurrent-map

import (
	"sync"
)

const (
	ShardCount = 32 // 分片数量
)

type Key interface {
	~string | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Map[K Key, V any] struct {
	sharding   func(key K) uint32
	shardCount uint32
	shards     []*shardMap[K, V]
}

type shardMap[K Key, V any] struct {
	elements map[K]V
	sync.RWMutex
}

func New[V any]() *Map[string, V] {
	return NewMap[string, V](DJBSharding())
}

func NewMap[K Key, V any](sharding func(key K) uint32) *Map[K, V] {
	var m = &Map[K, V]{}
	m.sharding = sharding
	m.shardCount = ShardCount
	m.shards = make([]*shardMap[K, V], m.shardCount)
	for i := uint32(0); i < m.shardCount; i++ {
		m.shards[i] = &shardMap[K, V]{elements: make(map[K]V)}
	}
	return m
}

func (m *Map[K, V]) getShard(key K) *shardMap[K, V] {
	var index = m.sharding(key)
	return m.shards[index]
}

func (m *Map[K, V]) Set(key K, value V) {
	var shard = m.getShard(key)
	shard.Lock()
	shard.elements[key] = value
	shard.Unlock()
}

func (m *Map[K, V]) SetNx(key K, value V) bool {
	var shard = m.getShard(key)
	shard.Lock()
	var _, found = shard.elements[key]
	if !found {
		shard.elements[key] = value
		shard.Unlock()
		return true
	}
	shard.Unlock()
	return false
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	var shard = m.getShard(key)
	shard.RLock()
	var value, ok = shard.elements[key]
	shard.RUnlock()
	return value, ok
}

func (m *Map[K, V]) Exists(key K) bool {
	var shard = m.getShard(key)
	shard.RLock()
	var _, ok = shard.elements[key]
	shard.RUnlock()
	return ok
}

func (m *Map[K, V]) RemoveAll() {
	for i := uint32(0); i < m.shardCount; i++ {
		var shard = m.shards[i]
		shard.Lock()
		shard.elements = make(map[K]V)
		shard.Unlock()
	}
}

func (m *Map[K, V]) Remove(key K) {
	m.Pop(key)
}

func (m *Map[K, V]) Pop(key K) (V, bool) {
	var shard = m.getShard(key)
	shard.Lock()
	var value, found = shard.elements[key]
	if found {
		delete(shard.elements, key)
	}
	shard.Unlock()
	return value, found
}

func (m *Map[K, V]) Len() int {
	var count = 0
	for i := uint32(0); i < m.shardCount; i++ {
		var shard = m.shards[i]
		shard.RLock()
		count += len(shard.elements)
		shard.RUnlock()
	}
	return count
}

func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	if f == nil {
		return
	}
	for i := uint32(0); i < m.shardCount; i++ {
		var shard = m.shards[i]
		shard.RLock()
		for k, v := range shard.elements {
			shard.RUnlock()
			if !f(k, v) {
				return
			}
			shard.RLock()
		}
		shard.RUnlock()
	}
}

func (m *Map[K, V]) Elements() map[K]V {
	var nMap = make(map[K]V, m.Len())
	for i := uint32(0); i < m.shardCount; i++ {
		var shard = m.shards[i]
		shard.RLock()
		for k, v := range shard.elements {
			nMap[k] = v
		}
		shard.RUnlock()
	}
	return nMap
}

func (m *Map[K, V]) Keys() []K {
	var nKeys = make([]K, 0, m.Len())
	for i := uint32(0); i < m.shardCount; i++ {
		var shard = m.shards[i]
		shard.RLock()
		for k := range shard.elements {
			nKeys = append(nKeys, k)
		}
		shard.RUnlock()
	}
	return nKeys
}
