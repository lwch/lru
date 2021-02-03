package main

import (
	"container/list"
	"crypto/sha1"
	"encoding/binary"
	"fmt"

	"github.com/lwch/lru"
)

type stringList struct {
	list *list.List
}

func (l *stringList) Make(uint64) {
	l.list = list.New()
}

func (l *stringList) KeyEqual(idx uint64, k interface{}) bool {
	node := l.list.Front()
	for i := uint64(0); i < idx; i++ {
		node = node.Next()
	}
	return node.Value == k
}

func (l *stringList) Get(idx uint64) interface{} {
	node := l.list.Front()
	for i := uint64(0); i < idx; i++ {
		node = node.Next()
	}
	return node.Value
}

func (l *stringList) Clone(idx uint64) interface{} {
	node := l.list.Front()
	for i := uint64(0); i < idx; i++ {
		node = node.Next()
	}
	return node.Value
}

func (l *stringList) Set(idx uint64, k, v interface{}) {
	if idx >= uint64(l.list.Len()) {
		l.list.PushBack(v)
		return
	}
	node := l.list.Front()
	for i := uint64(0); i < idx; i++ {
		node = node.Next()
	}
	node.Value = v
}

func (l *stringList) Hash(k interface{}) uint64 {
	enc := sha1.Sum([]byte(k.(string)))
	var id uint64
	for i := 0; i < 5; i++ {
		id += uint64(binary.LittleEndian.Uint32(enc[i*4:]))
	}
	return id
}

func (l *stringList) MoveTop(idx uint64) {
	node := l.list.Front()
	for i := uint64(0); i < idx; i++ {
		node = node.Next()
	}
	l.list.Remove(node)
	l.list.PushFront(node.Value)
}

func (l *stringList) print() {
	list := make([]string, 0, l.list.Len())
	for node := l.list.Front(); node != nil; node = node.Next() {
		list = append(list, node.Value.(string))
	}
	fmt.Println(list)
}

func main() {
	list := &stringList{}
	cache := lru.New(list, 10)
	for i := 0; i < 10; i++ {
		str := fmt.Sprintf("%c", '0'+i)
		cache.Set(str, str)
	}
	list.print()
	cache.Get("9")
	list.print()
	cache.Set("10", "10")
	list.print()
}
