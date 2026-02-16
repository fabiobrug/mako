package cache

import (
	"bytes"
	"database/sql"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func TestNewEmbeddingCache(t *testing.T) {
	cache := NewEmbeddingCache(100)
	
	if cache == nil {
		t.Fatal("Expected non-nil cache")
	}
	if cache.maxSize != 100 {
		t.Errorf("Expected maxSize 100, got %d", cache.maxSize)
	}
	if cache.items == nil {
		t.Error("Expected initialized items map")
	}
	if cache.lru == nil {
		t.Error("Expected initialized LRU list")
	}
}

func TestCacheGetSet(t *testing.T) {
	cache := NewEmbeddingCache(10)
	
	// Test miss
	_, hit := cache.Get("command1")
	if hit {
		t.Error("Expected cache miss for non-existent command")
	}
	
	// Set value
	embedding := []byte{1, 2, 3, 4, 5}
	cache.Set("command1", embedding)
	
	// Test hit
	retrieved, hit := cache.Get("command1")
	if !hit {
		t.Error("Expected cache hit")
	}
	if !bytes.Equal(retrieved, embedding) {
		t.Errorf("Expected %v, got %v", embedding, retrieved)
	}
}

func TestCacheUpdate(t *testing.T) {
	cache := NewEmbeddingCache(10)
	
	// Set initial value
	embedding1 := []byte{1, 2, 3}
	cache.Set("command1", embedding1)
	
	// Update value
	embedding2 := []byte{4, 5, 6}
	cache.Set("command1", embedding2)
	
	// Verify update
	retrieved, hit := cache.Get("command1")
	if !hit {
		t.Error("Expected cache hit")
	}
	if !bytes.Equal(retrieved, embedding2) {
		t.Errorf("Expected %v, got %v", embedding2, retrieved)
	}
}

func TestCacheLRUEviction(t *testing.T) {
	cache := NewEmbeddingCache(3)
	
	// Fill cache
	cache.Set("cmd1", []byte{1})
	cache.Set("cmd2", []byte{2})
	cache.Set("cmd3", []byte{3})
	
	// Add one more to trigger eviction
	cache.Set("cmd4", []byte{4})
	
	// cmd1 should be evicted (least recently used)
	_, hit := cache.Get("cmd1")
	if hit {
		t.Error("Expected cmd1 to be evicted")
	}
	
	// Others should still exist
	_, hit = cache.Get("cmd2")
	if !hit {
		t.Error("Expected cmd2 to exist")
	}
	_, hit = cache.Get("cmd4")
	if !hit {
		t.Error("Expected cmd4 to exist")
	}
}

func TestCacheLRUOrdering(t *testing.T) {
	cache := NewEmbeddingCache(3)
	
	// Fill cache
	cache.Set("cmd1", []byte{1})
	cache.Set("cmd2", []byte{2})
	cache.Set("cmd3", []byte{3})
	
	// Access cmd1 to make it most recently used
	cache.Get("cmd1")
	
	// Add new item - cmd2 should be evicted (now LRU)
	cache.Set("cmd4", []byte{4})
	
	_, hit := cache.Get("cmd2")
	if hit {
		t.Error("Expected cmd2 to be evicted")
	}
	
	_, hit = cache.Get("cmd1")
	if !hit {
		t.Error("Expected cmd1 to still exist (was accessed)")
	}
}

func TestCacheStats(t *testing.T) {
	cache := NewEmbeddingCache(10)
	
	// Initial stats
	stats := cache.Stats()
	if stats.Size != 0 {
		t.Errorf("Expected size 0, got %d", stats.Size)
	}
	if stats.Hits != 0 {
		t.Errorf("Expected 0 hits, got %d", stats.Hits)
	}
	if stats.Misses != 0 {
		t.Errorf("Expected 0 misses, got %d", stats.Misses)
	}
	
	// Add items
	cache.Set("cmd1", []byte{1})
	cache.Set("cmd2", []byte{2})
	
	// Generate hits and misses
	cache.Get("cmd1") // hit
	cache.Get("cmd1") // hit
	cache.Get("cmd3") // miss
	
	stats = cache.Stats()
	if stats.Size != 2 {
		t.Errorf("Expected size 2, got %d", stats.Size)
	}
	if stats.Hits != 2 {
		t.Errorf("Expected 2 hits, got %d", stats.Hits)
	}
	if stats.Misses != 1 {
		t.Errorf("Expected 1 miss, got %d", stats.Misses)
	}
	if stats.HitRate != 2.0/3.0 {
		t.Errorf("Expected hit rate 0.666, got %.3f", stats.HitRate)
	}
}

func TestCacheStatsHitRate(t *testing.T) {
	cache := NewEmbeddingCache(10)
	
	// Hit rate should be 0 with no operations
	stats := cache.Stats()
	if stats.HitRate != 0.0 {
		t.Errorf("Expected hit rate 0.0, got %.3f", stats.HitRate)
	}
	
	cache.Set("cmd1", []byte{1})
	
	// All hits
	cache.Get("cmd1")
	cache.Get("cmd1")
	cache.Get("cmd1")
	
	stats = cache.Stats()
	if stats.HitRate != 1.0 {
		t.Errorf("Expected hit rate 1.0, got %.3f", stats.HitRate)
	}
	
	// Add misses
	cache.Get("cmd2")
	cache.Get("cmd3")
	
	stats = cache.Stats()
	expectedRate := 3.0 / 5.0 // 3 hits, 2 misses
	if stats.HitRate != expectedRate {
		t.Errorf("Expected hit rate %.3f, got %.3f", expectedRate, stats.HitRate)
	}
}

func TestCacheClear(t *testing.T) {
	cache := NewEmbeddingCache(10)
	
	// Add items
	cache.Set("cmd1", []byte{1})
	cache.Set("cmd2", []byte{2})
	cache.Get("cmd1")
	
	// Clear cache
	cache.Clear()
	
	// Verify empty
	stats := cache.Stats()
	if stats.Size != 0 {
		t.Errorf("Expected size 0 after clear, got %d", stats.Size)
	}
	if stats.Hits != 0 {
		t.Errorf("Expected 0 hits after clear, got %d", stats.Hits)
	}
	if stats.Misses != 0 {
		t.Errorf("Expected 0 misses after clear, got %d", stats.Misses)
	}
	
	_, hit := cache.Get("cmd1")
	if hit {
		t.Error("Expected cache miss after clear")
	}
}

func TestCacheLoadSave(t *testing.T) {
	// Create in-memory database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	
	// Create cache table
	_, err = db.Exec(`
		CREATE TABLE embedding_cache (
			command_text TEXT PRIMARY KEY,
			embedding BLOB NOT NULL,
			hit_count INTEGER DEFAULT 0,
			last_accessed DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	
	// Create and populate cache
	cache1 := NewEmbeddingCache(10)
	cache1.Set("cmd1", []byte{1, 2, 3})
	cache1.Set("cmd2", []byte{4, 5, 6})
	
	// Save to database
	err = cache1.Save(db)
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}
	
	// Load into new cache
	cache2 := NewEmbeddingCache(10)
	err = cache2.Load(db)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	
	// Verify loaded data
	retrieved, hit := cache2.Get("cmd1")
	if !hit {
		t.Error("Expected cmd1 to be loaded")
	}
	expected := []byte{1, 2, 3}
	if !bytes.Equal(retrieved, expected) {
		t.Errorf("Expected %v, got %v", expected, retrieved)
	}
}

func TestCacheConcurrency(t *testing.T) {
	cache := NewEmbeddingCache(100)
	
	// Run concurrent operations
	done := make(chan bool)
	
	// Writer goroutine
	go func() {
		for i := 0; i < 100; i++ {
			cache.Set("cmd", []byte{byte(i)})
		}
		done <- true
	}()
	
	// Reader goroutine
	go func() {
		for i := 0; i < 100; i++ {
			cache.Get("cmd")
		}
		done <- true
	}()
	
	// Wait for completion
	<-done
	<-done
	
	// Should not panic or deadlock
}

func TestCacheTimestamp(t *testing.T) {
	cache := NewEmbeddingCache(10)
	
	before := time.Now()
	cache.Set("cmd1", []byte{1})
	after := time.Now()
	
	// Access to update timestamp
	time.Sleep(10 * time.Millisecond)
	cache.Get("cmd1")
	
	// Verify timestamp is reasonable (can't directly check, but ensure no panic)
	stats := cache.Stats()
	if stats.Size != 1 {
		t.Error("Cache should have 1 item")
	}
	
	_ = before
	_ = after
}

func TestCacheMemoryEstimate(t *testing.T) {
	cache := NewEmbeddingCache(10)
	
	// Add items
	cache.Set("cmd1", make([]byte, 768*4)) // ~3KB per entry
	cache.Set("cmd2", make([]byte, 768*4))
	
	stats := cache.Stats()
	
	// Should estimate ~6KB (2 entries * 3KB)
	expectedMemory := int64(2 * 3072)
	if stats.MemoryUsed != expectedMemory {
		t.Errorf("Expected memory %d, got %d", expectedMemory, stats.MemoryUsed)
	}
}

func TestCacheMaxSize(t *testing.T) {
	maxSize := 5
	cache := NewEmbeddingCache(maxSize)
	
	// Add more items than max size
	for i := 0; i < 10; i++ {
		cache.Set(string(rune('a'+i)), []byte{byte(i)})
	}
	
	stats := cache.Stats()
	if stats.Size > maxSize {
		t.Errorf("Cache size %d exceeds max size %d", stats.Size, maxSize)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	cache := NewEmbeddingCache(1000)
	embedding := make([]byte, 768*4)
	cache.Set("command", embedding)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get("command")
	}
}

func BenchmarkCacheSet(b *testing.B) {
	cache := NewEmbeddingCache(1000)
	embedding := make([]byte, 768*4)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("command", embedding)
	}
}

func BenchmarkCacheLRUEviction(b *testing.B) {
	cache := NewEmbeddingCache(100)
	embedding := make([]byte, 768*4)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(string(rune(i%200)), embedding)
	}
}
