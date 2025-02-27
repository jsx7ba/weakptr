package lru

import (
	"testing"
)

func TestPut(t *testing.T) {
	cache := NewLRU[int, int](10)
	for i := 0; i != 10; i++ {
		cache.Put(i, i)
	}

	for i := 10; i != 20; i++ {
		cache.Put(i, i)
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
	cache := NewLRU[string, string](3)
	_, exists := cache.Get("foo")
	if exists {
		t.Fatal("get on an empty list returned a value")
	}
}

func TestGet(t *testing.T) {
	cache := NewLRU[int, int](2)
	cache.Put(1, 1)
	_, exists := cache.Get(0)
	if exists {
		t.Fatal("key 0 should not exist")
	}

	val, exists := cache.Get(1)
	if !exists {
		t.Fatal("key 1 should exist in the cache")
	}
	if val != 1 {
		t.Fatalf("expected 1, got %d", val)
	}
}
