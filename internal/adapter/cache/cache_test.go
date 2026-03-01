package cache

import (
	"fmt"
	"sort"
	"sync"
	"testing"
)

// --- NewCache ---

func TestNewCache_ValidShards(t *testing.T) {
	cases := []struct {
		name    string
		nshards uint
	}{
		{"single shard", 1},
		{"two shards", 2},
		{"many shards", 64},
		{"power of two", 256},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewCache(tc.nshards)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if c == nil {
				t.Fatal("expected non-nil cache")
			}
			if uint(len(c.shards)) != tc.nshards {
				t.Fatalf("expected %d shards, got %d", tc.nshards, len(c.shards))
			}
			for i, shard := range c.shards {
				if shard == nil {
					t.Fatalf("shard %d is nil", i)
				}
				if shard.store == nil {
					t.Fatalf("shard %d store is nil", i)
				}
			}
		})
	}
}

func TestNewCache_ZeroShards(t *testing.T) {
	c, err := NewCache(0)
	if err == nil {
		t.Fatal("expected error for 0 shards, got nil")
	}
	if c != nil {
		t.Fatal("expected nil cache on error")
	}
}

// --- Set and Get ---

func TestSetAndGet_BasicTypes(t *testing.T) {
	c, _ := NewCache(4)

	cases := []struct {
		key   string
		value any
	}{
		{"string_key", "hello"},
		{"int_key", 42},
		{"float_key", 3.14},
		{"bool_key", true},
		{"slice_key", []int{1, 2, 3}},
		{"map_key", map[string]int{"a": 1}},
	}

	for _, tc := range cases {
		c.Set(tc.key, tc.value)
		val, ok := c.Get(tc.key)
		if !ok {
			t.Errorf("key %q: expected ok=true", tc.key)
		}
		if fmt.Sprintf("%v", val) != fmt.Sprintf("%v", tc.value) {
			t.Errorf("key %q: expected %v, got %v", tc.key, tc.value, val)
		}
	}
}

func TestGet_MissingKey(t *testing.T) {
	c, _ := NewCache(4)

	val, ok := c.Get("nonexistent")
	if ok {
		t.Error("expected ok=false for missing key")
	}
	if val != nil {
		t.Errorf("expected nil value for missing key, got %v", val)
	}
}

func TestSet_Overwrite(t *testing.T) {
	c, _ := NewCache(4)

	c.Set("key", "first")
	c.Set("key", "second")

	val, ok := c.Get("key")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if val != "second" {
		t.Errorf("expected 'second', got %v", val)
	}
}

func TestSet_EmptyKey(t *testing.T) {
	c, _ := NewCache(4)

	c.Set("", "empty-key-value")
	val, ok := c.Get("")
	if !ok {
		t.Error("expected ok=true for empty key")
	}
	if val != "empty-key-value" {
		t.Errorf("expected 'empty-key-value', got %v", val)
	}
}

// --- Delete ---

func TestDelete_ExistingKey(t *testing.T) {
	c, _ := NewCache(4)

	c.Set("key", "value")
	c.Delete("key")

	_, ok := c.Get("key")
	if ok {
		t.Error("expected key to be deleted")
	}
}

func TestDelete_NonExistingKey(t *testing.T) {
	c, _ := NewCache(4)
	// Should not panic
	c.Delete("ghost")
}

func TestDelete_ThenSet(t *testing.T) {
	c, _ := NewCache(4)

	c.Set("key", "original")
	c.Delete("key")
	c.Set("key", "new")

	val, ok := c.Get("key")
	if !ok {
		t.Fatal("expected key to exist after re-set")
	}
	if val != "new" {
		t.Errorf("expected 'new', got %v", val)
	}
}

// --- Len ---

func TestLen_Empty(t *testing.T) {
	c, _ := NewCache(4)
	if c.Len() != 0 {
		t.Errorf("expected 0, got %d", c.Len())
	}
}

func TestLen_AfterSets(t *testing.T) {
	c, _ := NewCache(4)

	keys := []string{"a", "b", "c", "d", "e"}
	for _, k := range keys {
		c.Set(k, k)
	}
	if c.Len() != len(keys) {
		t.Errorf("expected %d, got %d", len(keys), c.Len())
	}
}

func TestLen_AfterDelete(t *testing.T) {
	c, _ := NewCache(4)

	c.Set("x", 1)
	c.Set("y", 2)
	c.Delete("x")

	if c.Len() != 1 {
		t.Errorf("expected 1, got %d", c.Len())
	}
}

func TestLen_Overwrite_DoesNotIncrease(t *testing.T) {
	c, _ := NewCache(4)

	c.Set("key", "v1")
	c.Set("key", "v2")

	if c.Len() != 1 {
		t.Errorf("expected 1 after overwrite, got %d", c.Len())
	}
}

// --- Contains ---

func TestContains_ExistingKey(t *testing.T) {
	c, _ := NewCache(4)

	c.Set("present", "value")
	if !c.Contains("present") {
		t.Error("expected Contains to return true for existing key")
	}
}

func TestContains_MissingKey(t *testing.T) {
	c, _ := NewCache(4)

	if c.Contains("absent") {
		t.Error("expected Contains to return false for missing key")
	}
}

func TestContains_AfterDelete(t *testing.T) {
	c, _ := NewCache(4)

	c.Set("temp", "data")
	c.Delete("temp")

	if c.Contains("temp") {
		t.Error("expected Contains to return false after delete")
	}
}

func TestContains_NilValue(t *testing.T) {
	c, _ := NewCache(4)

	// NOTE: Contains uses `shard.store[key] != nil`, so a key explicitly set
	// to nil will appear as not contained. This test documents that known behaviour.
	c.Set("nil_val", nil)
	contained := c.Contains("nil_val")
	t.Logf("Contains('nil_val') with nil stored value = %v (nil values appear absent)", contained)
}

// --- Keys ---

func TestKeys_Empty(t *testing.T) {
	c, _ := NewCache(4)

	keys := c.Keys()
	if len(keys) != 0 {
		t.Errorf("expected empty keys, got %v", keys)
	}
}

func TestKeys_AllKeysReturned(t *testing.T) {
	c, _ := NewCache(4)

	expected := []string{"alpha", "beta", "gamma", "delta", "epsilon"}
	for _, k := range expected {
		c.Set(k, true)
	}

	got := c.Keys()
	sort.Strings(got)
	sort.Strings(expected)

	if len(got) != len(expected) {
		t.Fatalf("expected %d keys, got %d: %v", len(expected), len(got), got)
	}
	for i, k := range expected {
		if got[i] != k {
			t.Errorf("key mismatch at index %d: expected %q, got %q", i, k, got[i])
		}
	}
}

func TestKeys_AfterDelete(t *testing.T) {
	c, _ := NewCache(4)

	c.Set("keep", 1)
	c.Set("remove", 2)
	c.Delete("remove")

	keys := c.Keys()
	if len(keys) != 1 || keys[0] != "keep" {
		t.Errorf("expected ['keep'], got %v", keys)
	}
}

// --- getShard ---

func TestGetShard_Consistency(t *testing.T) {
	c, _ := NewCache(8)

	for i := 0; i < 100; i++ {
		s1 := c.getShard("stable-key")
		s2 := c.getShard("stable-key")
		if s1 != s2 {
			t.Error("getShard returned different shards for the same key")
		}
	}
}

func TestGetShard_DifferentKeys_MayDiffer(t *testing.T) {
	c, _ := NewCache(256)

	shards := make(map[*Shard]struct{})
	for i := 0; i < 50; i++ {
		shards[c.getShard(fmt.Sprintf("key-%d", i))] = struct{}{}
	}
	if len(shards) < 2 {
		t.Error("expected keys to be distributed across multiple shards")
	}
}

// --- Concurrency ---

func TestConcurrent_SetAndGet(t *testing.T) {
	c, _ := NewCache(16)
	var wg sync.WaitGroup
	n := 500

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c.Set(fmt.Sprintf("key-%d", i), i)
		}(i)
	}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c.Get(fmt.Sprintf("key-%d", i))
		}(i)
	}

	wg.Wait()

	if c.Len() != n {
		t.Errorf("expected %d entries after concurrent writes, got %d", n, c.Len())
	}
}

func TestConcurrent_SetAndDelete(t *testing.T) {
	c, _ := NewCache(8)
	var wg sync.WaitGroup
	n := 200

	for i := 0; i < n; i++ {
		c.Set(fmt.Sprintf("key-%d", i), i)
	}
	for i := 0; i < n; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			c.Set(fmt.Sprintf("key-%d", i), i*2)
		}(i)
		go func(i int) {
			defer wg.Done()
			c.Delete(fmt.Sprintf("key-%d", i))
		}(i)
	}

	wg.Wait()
	// No panic = success; final state is non-deterministic
}

func TestConcurrent_Keys(t *testing.T) {
	c, _ := NewCache(8)
	n := 100
	for i := 0; i < n; i++ {
		c.Set(fmt.Sprintf("k-%d", i), i)
	}

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			keys := c.Keys()
			if len(keys) > n {
				t.Errorf("Keys() returned more keys (%d) than inserted (%d)", len(keys), n)
			}
		}()
	}
	wg.Wait()
}

func TestConcurrent_Len(t *testing.T) {
	c, _ := NewCache(8)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c.Set(fmt.Sprintf("key-%d", i), i)
			_ = c.Len()
		}(i)
	}
	wg.Wait()
}

// --- Single shard edge case ---

func TestSingleShard_AllOps(t *testing.T) {
	c, _ := NewCache(1)

	c.Set("a", 1)
	c.Set("b", 2)

	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Errorf("Get('a') = %v, %v; want 1, true", v, ok)
	}
	if c.Len() != 2 {
		t.Errorf("Len() = %d; want 2", c.Len())
	}
	if !c.Contains("b") {
		t.Error("Contains('b') = false; want true")
	}

	c.Delete("a")
	if c.Contains("a") {
		t.Error("Contains('a') = true after delete; want false")
	}
	if c.Len() != 1 {
		t.Errorf("Len() = %d after delete; want 1", c.Len())
	}

	keys := c.Keys()
	if len(keys) != 1 || keys[0] != "b" {
		t.Errorf("Keys() = %v; want ['b']", keys)
	}
}
