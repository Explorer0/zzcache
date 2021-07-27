package zzcache

import "testing"

func TestGet(t *testing.T) {
	lru := NewLRU(10)
	lru.Set("key1", "1234")
	if v, ok := lru.Get("key1"); !ok || v.(string) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestRemoveoldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"

	lru := NewLRU(2)
	lru.Set(k1, v1)
	lru.Set(k2, v2)
	lru.Set(k3, v3)

	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatalf("Removeoldest key1 failed")
	}
}
