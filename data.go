package lru

// ListData lru data
type ListData interface {
	Make(size uint64)
	KeyEqual(idx uint64, key interface{}) bool
	Get(idx uint64) interface{}
	Clone(idx uint64) interface{}
	Set(idx uint64, k, v interface{})
	Hash(key interface{}) uint64
	MoveTop(idx uint64)
}
