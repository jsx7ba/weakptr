package wlru

import (
	"math/rand/v2"
	"runtime"
	"testing"
)

type value struct {
	a int
	b int
	c int
}

func TestPut(t *testing.T) {
	cache := NewWeakLRU[int, value](10)

	for i := 0; i != 10; i++ {
		cache.Put(i, &value{i, i, i})
	}

	for i := 10; i != 20; i++ {
		cache.Put(i, &value{i, i, i})
	}

	// check that 0-9 have been evicted from the cache
	for i := 0; i != 10; i++ {
		_, exists := cache.Get(i)
		if exists {
			t.Errorf("key %d  should have been evicted", i)
		}
	}
}

func TestGetEmpty(t *testing.T) {
	cache := NewWeakLRU[string, string](3)
	_, exists := cache.Get("foo")
	if exists {
		t.Fatal("get on an empty list returned a value")
	}
}

func TestGet(t *testing.T) {
	v := value{1, 2, 3}
	cache := NewWeakLRU[int, value](2)

	cache.Put(1, &v)
	_, exists := cache.Get(0)
	if exists {
		t.Fatal("key 0 should not exist")
	}

	p2, exists := cache.Get(1)
	if !exists {
		t.Fatal("key 1 should exist in the cache")
	}
	if p2 != &v {
		t.Fatalf("pointers should be equal, but are not %v != %v", p2, &v)
	}
}

func TestWeakPtr(t *testing.T) {
	v := &value{1, 2, 3}
	cache := NewWeakLRU[string, value](2)

	cache.Put("foo", v)
	v = nil
	runtime.GC() // run the garbage collector

	actual, ok := cache.Get("foo")
	if !ok {
		t.Fatal("cache lost the key")
	}

	if actual != nil {
		t.Fatal("the weak pointer should have been cleaned up.")
	}
}

func TestDelete(t *testing.T) {
	cache := NewWeakLRU[int, value](10)

	for i := 0; i != 10; i++ {
		cache.Put(i, &value{i, i, i})
	}

	key := rand.IntN(10)
	_, ok := cache.Get(key)
	if !ok {
		t.Fatalf("couldn't find key %d", key)
	}

	if !cache.Delete(key) {
		t.Fatalf("key %d couldn't be deleted", key)
	}

	_, ok = cache.Get(key)
	if ok {
		t.Fatalf("key %d should not exist in the cache", key)
	}
}

func TestDeleteAll(t *testing.T) {
	cache := NewWeakLRU[int, value](10)

	for i := 0; i != 10; i++ {
		cache.Put(i, &value{i, i, i})
	}

	for i := 0; i != 10; i++ {
		cache.Delete(i)
	}

	// now try adding things back in to the cache.
	for i := 0; i != 10; i++ {
		cache.Put(i, &value{i, i, i})
	}

	for i := 0; i != 10; i++ {
		v, ok := cache.Get(i)
		if !ok {
			t.Fatalf("unable to retrieve key %d after deletion", i)
		}
		if v.a != i {
			t.Fatalf("cache return the wrong value for key %d", i)
		}
	}
}
