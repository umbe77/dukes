package cache

import (
	"fmt"
	"sync"

	"github.com/umbe77/dukes/datatypes"
)

type CacheValue struct {
	Kind  datatypes.DataType
	Value any
}

type Cache interface {
	Set(key string, value *CacheValue) error
	Get(key string) (*CacheValue, error)
	Has(key string) bool
	Del(key string) error
	Dump() chan string
}

type MemoryCache struct {
	sync.RWMutex
	c map[string]*CacheValue
}

func NewCache() *MemoryCache {
	return &MemoryCache{
		c: make(map[string]*CacheValue),
	}
}

func (c *MemoryCache) Set(key string, value *CacheValue) error {
	c.Lock()
	c.c[key] = value
	c.Unlock()
	return nil
}

func (c *MemoryCache) Has(key string) bool {
	var result bool
	c.RLock()
	_, result = c.c[key]
	c.RUnlock()
	return result
}

func (c *MemoryCache) Get(key string) (*CacheValue, error) {
	var (
		v  *CacheValue
		ok bool
	)

	c.RLock()
	if v, ok = c.c[key]; !ok {
		return nil, fmt.Errorf("key %s not present in cache", key)
	}
	c.RUnlock()
	return v, nil
}

func (c *MemoryCache) Del(key string) error {

	c.Lock()
	delete(c.c, key)
	c.Unlock()
	return nil
}

func (c *MemoryCache) Dump() <-chan string {
	keysCh := make(chan string)

	go func(mc *MemoryCache) {

		for k, _ := range c.c {
			keysCh <- k
		}

		close(keysCh)
	}(c)

	return keysCh
}
