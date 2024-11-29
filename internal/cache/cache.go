package cache

import "sync"

type Cache struct {
	msg sync.Map
}

func NewCache() *Cache {
	return &Cache{}
}
