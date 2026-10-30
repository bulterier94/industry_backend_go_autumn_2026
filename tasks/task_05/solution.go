package main

type Cache[K comparable, V any] struct {
	capacity int
	items    map[K]V
}

func NewCache[K comparable, V any](capacity int) *Cache[K, V] {
	return &Cache[K, V]{
		items: make(map[K]V),
		capacity: capacity,
	}
}
func (c *Cache[K, V]) Get(k K) (v V, ok bool) {
	v, ok = c.items[k]

	return v, ok
}
func (c *Cache[K, V]) Set(k K, v V) bool {
	if c.capacity <= 0 {
		return false
	}

	if _, exists := c.Get(k); exists {
		c.items[k] = v
		return true
	}

	if len(c.items) >= c.capacity {
		return false
	}

	c.items[k] = v
	return true
}
