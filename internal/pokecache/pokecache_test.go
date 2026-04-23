package pokecache

import (
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	cache := NewCache(5 * time.Second)

	cases := []struct {
		key string
		val []byte
	}{
		{"https://pokeapi.co/api/v2/location-area?offset=0", []byte("test data 1")},
		{"https://pokeapi.co/api/v2/location-area?offset=20", []byte("test data 2")},
	}

	for _, c := range cases {
		cache.Add(c.key, c.val)
		got, ok := cache.Get(c.key)
		if !ok {
			t.Errorf("expected to find key %q in cache", c.key)
			continue
		}
		if string(got) != string(c.val) {
			t.Errorf("expected %q, got %q", c.val, got)
		}
	}
}

func TestReap(t *testing.T) {
	interval := 10 * time.Millisecond
	cache := NewCache(interval)

	key := "https://pokeapi.co/api/v2/location-area?offset=0"
	cache.Add(key, []byte("test data"))

	// Entry should exist immediately
	if _, ok := cache.Get(key); !ok {
		t.Error("expected entry to exist before reap interval")
	}

	// Wait for reap to kick in
	time.Sleep(interval * 3)

	// Entry should be gone now
	if _, ok := cache.Get(key); ok {
		t.Error("expected entry to have been reaped")
	}
}

func TestGetMiss(t *testing.T) {
	cache := NewCache(5 * time.Second)
	if _, ok := cache.Get("nonexistent-key"); ok {
		t.Error("expected cache miss for nonexistent key")
	}
}
