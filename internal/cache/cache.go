package cache

import "container/list"

type Cache struct {
	list       *list.List
	elemByKeys map[uint64]*list.Element
	elemByExp  map[uint64]*list.Element
}

// lru
// expiration
