package cache

import (
	"container/list"
	"database/sql"
	"sync"
	"time"
)

// EmbeddingCache implements an LRU cache for command embeddings
type EmbeddingCache struct {
	maxSize int
	mu      sync.RWMutex
	items   map[string]*list.Element
	lru     *list.List
	hits    int64
	misses  int64
}

type cacheEntry struct {
	key       string
	embedding []byte
	timestamp time.Time
}

// NewEmbeddingCache creates a new LRU embedding cache
func NewEmbeddingCache(maxSize int) *EmbeddingCache {
	return &EmbeddingCache{
		maxSize: maxSize,
		items:   make(map[string]*list.Element),
		lru:     list.New(),
	}
}

// Get retrieves an embedding from cache
// Returns (embedding, hit) where hit is true if found
func (c *EmbeddingCache) Get(command string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, exists := c.items[command]
	if !exists {
		c.misses++
		return nil, false
	}

	// Move to front (most recently used)
	c.lru.MoveToFront(elem)
	entry := elem.Value.(*cacheEntry)
	entry.timestamp = time.Now()
	
	c.hits++
	return entry.embedding, true
}

// Set stores an embedding in the cache
func (c *EmbeddingCache) Set(command string, embedding []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if already exists
	if elem, exists := c.items[command]; exists {
		c.lru.MoveToFront(elem)
		entry := elem.Value.(*cacheEntry)
		entry.embedding = embedding
		entry.timestamp = time.Now()
		return
	}

	// Add new entry
	entry := &cacheEntry{
		key:       command,
		embedding: embedding,
		timestamp: time.Now(),
	}
	elem := c.lru.PushFront(entry)
	c.items[command] = elem

	// Evict LRU if over capacity
	if c.lru.Len() > c.maxSize {
		oldest := c.lru.Back()
		if oldest != nil {
			c.lru.Remove(oldest)
			oldEntry := oldest.Value.(*cacheEntry)
			delete(c.items, oldEntry.key)
		}
	}
}

// Stats returns cache statistics
type CacheStats struct {
	Size       int
	MaxSize    int
	Hits       int64
	Misses     int64
	HitRate    float64
	MemoryUsed int64 // Approximate memory usage in bytes
}

func (c *EmbeddingCache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := c.hits + c.misses
	hitRate := 0.0
	if total > 0 {
		hitRate = float64(c.hits) / float64(total)
	}

	// Estimate memory: each embedding is ~768 dims * 4 bytes = 3KB
	memoryUsed := int64(c.lru.Len()) * 3072

	return CacheStats{
		Size:       c.lru.Len(),
		MaxSize:    c.maxSize,
		Hits:       c.hits,
		Misses:     c.misses,
		HitRate:    hitRate,
		MemoryUsed: memoryUsed,
	}
}

// Load loads cache entries from database
func (c *EmbeddingCache) Load(db *sql.DB) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	query := `
		SELECT command_text, embedding, last_accessed
		FROM embedding_cache
		ORDER BY hit_count DESC, last_accessed DESC
		LIMIT ?
	`

	rows, err := db.Query(query, c.maxSize)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var commandText string
		var embedding []byte
		var lastAccessed time.Time

		if err := rows.Scan(&commandText, &embedding, &lastAccessed); err != nil {
			continue
		}

		entry := &cacheEntry{
			key:       commandText,
			embedding: embedding,
			timestamp: lastAccessed,
		}
		elem := c.lru.PushBack(entry)
		c.items[commandText] = elem
	}

	return nil
}

// Save persists cache to database
func (c *EmbeddingCache) Save(db *sql.DB) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clear existing cache
	if _, err := tx.Exec("DELETE FROM embedding_cache"); err != nil {
		return err
	}

	// Insert current cache entries
	stmt, err := tx.Prepare(`
		INSERT INTO embedding_cache (command_text, embedding, hit_count, last_accessed)
		VALUES (?, ?, 1, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for elem := c.lru.Front(); elem != nil; elem = elem.Next() {
		entry := elem.Value.(*cacheEntry)
		if _, err := stmt.Exec(entry.key, entry.embedding, entry.timestamp); err != nil {
			// Continue on error to save as much as possible
			continue
		}
	}

	return tx.Commit()
}

// Clear removes all entries from the cache
func (c *EmbeddingCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*list.Element)
	c.lru = list.New()
	c.hits = 0
	c.misses = 0
}
