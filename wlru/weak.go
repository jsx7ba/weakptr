package wlru

import (
	"weak"
)

type Cache[K comparable, V any] struct {
	cap     int
	hashMap map[K]*node[K, V]
	head    *node[K, V]
	tail    *node[K, V]
}

// NewWeakLRU creates a new Cache with the capacity of cap.
func NewWeakLRU[K comparable, V any](cap int) *Cache[K, V] {
	if cap <= 1 {
		panic("capacity must be greater than 1")
	}

	return &Cache[K, V]{
		cap:     cap,
		hashMap: make(map[K]*node[K, V]),
		head:    nil,
		tail:    nil,
	}
}

type node[K comparable, V any] struct {
	key   K
	value weak.Pointer[V]
	next  *node[K, V]
	prev  *node[K, V]
}

// Get - returns the value associated with the key, or false if no value exists.  If the value
// exists in the cache, that value is now the most recently used value.
func (c *Cache[K, V]) Get(key K) (*V, bool) {
	n, exists := c.hashMap[key]
	if !exists {
		return nil, false
	}
	moveToTop(c, n)
	return n.value.Value(), true
}

// Delete an entry from the cache.
func (c *Cache[K, V]) Delete(key K) bool {
	n, ok := c.hashMap[key]
	if ok {
		delete(c.hashMap, key)
		if c.head == n {
			c.head = n.next
		}
		if n.next != nil {
			n.next.prev = n.prev
		}

		if c.tail == n {
			c.tail = n.prev
		}
		if n.prev != nil {
			n.prev.next = n.next
		}
	}

	return ok
}

// Put - adds a new value or updates an existing value. If the capacity is exceeded when a new
// value is added then the oldest value in the cache will be evicted.
func (c *Cache[K, V]) Put(key K, value *V) {
	n, exists := c.hashMap[key]
	if exists {
		moveToTop(c, n)
		return
	}

	n = &node[K, V]{
		key:   key,
		value: weak.Make(value),
		next:  c.head,
		prev:  nil,
	}

	if c.head != nil {
		c.head.prev = n
	}

	c.head = n
	c.hashMap[key] = n
	if c.tail == nil {
		c.tail = n
	}

	if len(c.hashMap) > c.cap {
		delete(c.hashMap, c.tail.key)
		c.tail.next = nil
		parent := c.tail.prev
		c.tail.prev = nil

		c.tail = parent
		parent.next = nil
	}
}

// shuffle node to the top of the head
// A given node can never be both the head and the tail due to the
// requirement that the capacity must be 2.
func moveToTop[K comparable, V any](c *Cache[K, V], n *node[K, V]) {
	if n != c.head {
		if n == c.tail {
			c.tail = n.prev
		}
		if n.prev != nil {
			n.prev.next = n.next
			n.prev = nil
		}

		n.next = c.head
		c.head = n
	}
}
