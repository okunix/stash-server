package cache

import (
	"crypto/sha1"
	"errors"
	"sync"
)

type Shard struct {
	store map[string]any
	mu    sync.RWMutex
}

type Cache struct {
	shards []*Shard
}

func NewCache(nshards uint) (*Cache, error) {
	if nshards < 1 {
		return nil, errors.New("cache should contain at least one shard")
	}
	shards := make([]*Shard, nshards)
	for i := range len(shards) {
		shards[i] = &Shard{store: make(map[string]any)}
	}
	return &Cache{shards: shards}, nil
}

func (c *Cache) getShard(key string) *Shard {
	shardIndex := int(sha1.Sum([]byte(key))[0]) % len(c.shards)
	return c.shards[shardIndex]
}

func (c *Cache) Get(key string) (any, bool) {
	shard := c.getShard(key)
	shard.mu.RLock()
	defer shard.mu.RUnlock()
	value, ok := shard.store[key]
	return value, ok
}

func (c *Cache) Set(key string, value any) {
	shard := c.getShard(key)

	shard.mu.Lock()
	defer shard.mu.Unlock()

	shard.store[key] = value
}

func (c *Cache) Delete(key string) {
	shard := c.getShard(key)

	shard.mu.Lock()
	defer shard.mu.Unlock()
	delete(shard.store, key)
}

func (c *Cache) Len() int {
	count := 0
	for _, v := range c.shards {
		v.mu.RLock()
		count += len(v.store)
		v.mu.RUnlock()
	}
	return count
}

func (c *Cache) Contains(key string) bool {
	shard := c.getShard(key)
	shard.mu.RLock()
	defer shard.mu.RUnlock()
	return shard.store[key] != nil
}

func (c *Cache) Keys() []string {
	keys := make([]string, 0)
	mu := &sync.Mutex{}
	wg := sync.WaitGroup{}

	for _, shard := range c.shards {
		wg.Add(1)
		go func() {
			shard.mu.RLock()

			for k := range shard.store {
				mu.Lock()
				keys = append(keys, k)
				mu.Unlock()
			}

			wg.Done()
			shard.mu.RUnlock()
		}()
	}
	wg.Wait()
	return keys
}
