package cache

import (
	"testing"
)

func TestPut(t *testing.T) {
	cache, err := NewLRU[int, int](10)
	if err != nil {
		t.Fatal(err)
	}

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
	cache, err := NewLRU[string, string](3)
	if err != nil {
		t.Fatalf("failed to make the cache: %+v", err)
	}
	_, exists := cache.Get("foo")
	if exists {
		t.Fatal("get on an empty list returned a value")
	}
}

func TestGet(t *testing.T) {
	cache, err := NewLRU[int, int](2)
	if err != nil {
		t.Fatalf("unable to create cache: %+v", err)
	}

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
