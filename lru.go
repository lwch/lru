package lru

import (
	"sync"
)

// LRU lru structure
type LRU struct {
	sync.Mutex
	data    ListData
	filter  uint64
	size    uint64
	maxSize uint64
}

// New new lru cache
func New(data ListData, size uint64) *LRU {
	lru := &LRU{data: data, maxSize: size}
	lru.data.Make(size)
	return lru
}

// Set set data
func (lru *LRU) Set(k, v interface{}) interface{} {
	hash := lru.data.Hash(k)
	lru.Lock()
	defer lru.Unlock()
	if lru.filter&hash == 0 {
		return lru.push(k, v, hash)
	}
	return lru.set(k, v, hash)
}

// Get get data
func (lru *LRU) Get(k interface{}) interface{} {
	lru.Lock()
	defer lru.Unlock()
	if lru.filter&lru.data.Hash(k) == 0 {
		return nil
	}
	idx := lru.search(k)
	if idx < 0 {
		return nil
	}
	value := lru.data.Get(uint64(idx))
	lru.data.MoveTop(uint64(idx))
	return value
}

func (lru *LRU) push(k, v interface{}, hash uint64) interface{} {
	if lru.size >= lru.maxSize {
		lru.filter |= hash
		expire := lru.data.Clone(lru.maxSize - 1)
		lru.data.Set(lru.maxSize-1, k, v)
		lru.size++
		return expire
	}
	lru.filter |= hash
	lru.data.Set(lru.size, k, v)
	lru.size++
	return nil
}

func (lru *LRU) search(k interface{}) int64 {
	for i := uint64(0); i < uint64(lru.size); i++ {
		if lru.data.KeyEqual(i, k) {
			return int64(i)
		}
	}
	return -1
}

func (lru *LRU) set(k, v interface{}, hash uint64) interface{} {
	idx := lru.search(k)
	if idx < 0 {
		return lru.push(k, v, hash)
	}
	lru.data.Set(uint64(idx), k, v)
	return nil
}
