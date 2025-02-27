# Testing the weak package

Investigating how to use the [weak](https://pkg.go.dev/weak) package with a cache.

The `lru` package implements a vanilla Least Recently Used cache utilizing a map and a doubly linked list.

The `wlru` package implements an LRU cache with weak pointers to hold values in the linked list.

## Using Weak Pointers

There are 2 primary changes needed to the standard LRU cache to use weak references:

1. The value field in the node struct now has the `weak.Pointer[T]` type. 
```go{3}
type node[K comparable, V any] struct {
	key   K
	value weak.Pointer[V]
	next  *node[K, V]
	prev  *node[K, V]
}
```

2. The signatures for `Get` and `Put` have changed, note the `*V` below.

```go
func (c *Cache[K, V]) Put(key K, value *V) {...}
func (c *Cache[K, V]) Get(key K) (*V, bool) {...}
```

## The Case for Weak Pointers
The canonical use case for weak pointers in Go are illustrated in the test `TestWeakPtr`.

```go
	v := &value{1, 2, 3}
	cache := NewWeakLRU[string, value](2)

	cache.Put("foo", v)
	v = nil
	runtime.GC() // run the garbage collector
	
    actual, ok := cache.Get("foo")
    if actual != nil {
        t.Fatal("the weak pointer should have been cleaned up.")
    }
```

The GC is allowed to clean up the memory when the last strong pointer is removed, regardless of how many weak pointers are referencing the memory.
When the GC runs and cleans up memory the pointer value in `weak.Pointer` will be set to nil.

## An Alternative

Implementing a Delete method on the cache would have the same effect, if not better.  Leaving a nil value in the cache
until eviction feels weird.
