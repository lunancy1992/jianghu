package cache

import (
	"time"

	"github.com/dgraph-io/ristretto"
)

type Cache struct {
	store *ristretto.Cache
}

func New(numCounters, maxCost int64) (*Cache, error) {
	c, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: numCounters,
		MaxCost:     maxCost,
		BufferItems: 64,
	})
	if err != nil {
		return nil, err
	}
	return &Cache{store: c}, nil
}

func (c *Cache) Get(key string) (interface{}, bool) {
	return c.store.Get(key)
}

func (c *Cache) Set(key string, value interface{}, cost int64) {
	c.store.SetWithTTL(key, value, cost, 5*time.Minute)
}

func (c *Cache) SetWithTTL(key string, value interface{}, cost int64, ttl time.Duration) {
	c.store.SetWithTTL(key, value, cost, ttl)
}

func (c *Cache) Del(key string) {
	c.store.Del(key)
}

func (c *Cache) Close() {
	c.store.Close()
}
