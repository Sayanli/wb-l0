package cache

import (
	"sync"
	"time"
	"wb-l0/internal/models"
)

type Item struct {
	Key     string
	Value   models.Order
	Created time.Time
}

type Cache struct {
	sync.RWMutex
	items    map[string]Item
	TTL      time.Duration
	TTLCheck time.Duration
}

func NewCache(ttl, ttlCheck time.Duration) *Cache {
	items := make(map[string]Item)

	cache := &Cache{
		items:    items,
		TTL:      ttl,
		TTLCheck: ttlCheck,
	}

	return cache
}

func (c *Cache) RestoreCache(ords []models.Order) {
	for _, ord := range ords {
		c.Set(ord.Order_uid, ord)
	}
}

func (c *Cache) Set(key string, value models.Order) {
	c.Lock()
	defer c.Unlock()
	c.items[key] = Item{
		Key:     key,
		Value:   value,
		Created: time.Now(),
	}
}

func (c *Cache) Get(key string) (models.Order, bool) {
	c.RLock()
	defer c.RUnlock()
	item, ok := c.items[key]
	return item.Value, ok
}

func (c *Cache) Delete(key string) {
	c.Lock()
	defer c.Unlock()
	delete(c.items, key)
}

func (c *Cache) Clear() {
	c.Lock()
	defer c.Unlock()
	c.items = make(map[string]Item)
}

func (c *Cache) GC() {
	c.Lock()
	defer c.Unlock()
	for key, item := range c.items {
		if time.Since(item.Created) > c.TTLCheck {
			delete(c.items, key)
		}
	}
}
