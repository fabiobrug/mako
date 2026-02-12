package database

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// EmbeddingService defines the interface for generating embeddings
type EmbeddingService interface {
	GenerateEmbedding(text string) ([]byte, error)
}

// EmbeddingWorker manages background embedding generation
type EmbeddingWorker struct {
	db              *DB
	embedService    EmbeddingService
	numWorkers      int
	queue           chan int64
	ctx             context.Context
	cancel          context.CancelFunc
	wg              sync.WaitGroup
	retryDelay      time.Duration
	maxRetries      int
	processedCount  int64
	failedCount     int64
	mu              sync.Mutex
}

// NewEmbeddingWorker creates a new embedding worker pool
func NewEmbeddingWorker(db *DB, embedService EmbeddingService, numWorkers int) *EmbeddingWorker {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &EmbeddingWorker{
		db:           db,
		embedService: embedService,
		numWorkers:   numWorkers,
		queue:        make(chan int64, 1000), // Buffer up to 1000 pending commands
		ctx:          ctx,
		cancel:       cancel,
		retryDelay:   time.Second * 5,
		maxRetries:   3,
	}
}

// Start begins the worker pool
func (w *EmbeddingWorker) Start() {
	// Start worker goroutines
	for i := 0; i < w.numWorkers; i++ {
		w.wg.Add(1)
		go w.worker(i)
	}

	// Start queue feeder
	w.wg.Add(1)
	go w.feedQueue()
}

// Stop gracefully shuts down the worker pool
func (w *EmbeddingWorker) Stop() {
	w.cancel()
	close(w.queue)
	w.wg.Wait()
}

// Enqueue adds a command ID to the embedding queue
func (w *EmbeddingWorker) Enqueue(cmdID int64) {
	select {
	case w.queue <- cmdID:
	case <-w.ctx.Done():
	default:
		// Queue full, will be picked up by feedQueue later
	}
}

// worker processes commands from the queue
func (w *EmbeddingWorker) worker(id int) {
	defer w.wg.Done()

	for {
		select {
		case <-w.ctx.Done():
			return
		case cmdID, ok := <-w.queue:
			if !ok {
				return
			}
			w.processCommand(cmdID)
		}
	}
}

// feedQueue continuously feeds pending commands into the queue
func (w *EmbeddingWorker) feedQueue() {
	defer w.wg.Done()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-w.ctx.Done():
			return
		case <-ticker.C:
			w.loadPendingCommands()
		}
	}
}

// loadPendingCommands loads pending commands from database
func (w *EmbeddingWorker) loadPendingCommands() {
	cmds, err := w.db.GetPendingEmbeddings(100)
	if err != nil {
		return
	}

	for _, cmd := range cmds {
		select {
		case w.queue <- cmd.ID:
		case <-w.ctx.Done():
			return
		default:
			// Queue full, try again next cycle
			return
		}
	}
}

// processCommand generates and saves embedding for a command
func (w *EmbeddingWorker) processCommand(cmdID int64) {
	// Mark as processing
	if err := w.db.UpdateEmbeddingStatus(cmdID, "processing", nil); err != nil {
		return
	}

	// Get command details
	var command string
	err := w.db.GetConn().QueryRow("SELECT command FROM commands WHERE id = ?", cmdID).Scan(&command)
	if err != nil {
		w.db.UpdateEmbeddingStatus(cmdID, "failed", nil)
		w.incrementFailed()
		return
	}

	// Generate embedding with retries
	var embedding []byte
	var lastErr error
	
	for attempt := 0; attempt < w.maxRetries; attempt++ {
		embedding, lastErr = w.embedService.GenerateEmbedding(command)
		if lastErr == nil {
			break
		}

		// Exponential backoff
		backoff := w.retryDelay * time.Duration(1<<uint(attempt))
		select {
		case <-time.After(backoff):
		case <-w.ctx.Done():
			return
		}
	}

	if lastErr != nil {
		w.db.UpdateEmbeddingStatus(cmdID, "failed", nil)
		w.incrementFailed()
		log.Printf("Failed to generate embedding for command %d after %d attempts: %v", cmdID, w.maxRetries, lastErr)
		return
	}

	// Save embedding
	if err := w.db.UpdateEmbeddingStatus(cmdID, "completed", embedding); err != nil {
		log.Printf("Failed to save embedding for command %d: %v", cmdID, err)
		w.incrementFailed()
		return
	}

	w.incrementProcessed()
}

// Stats returns worker statistics
type WorkerStats struct {
	QueueSize      int
	ProcessedCount int64
	FailedCount    int64
	NumWorkers     int
}

func (w *EmbeddingWorker) Stats() WorkerStats {
	w.mu.Lock()
	defer w.mu.Unlock()

	return WorkerStats{
		QueueSize:      len(w.queue),
		ProcessedCount: w.processedCount,
		FailedCount:    w.failedCount,
		NumWorkers:     w.numWorkers,
	}
}

func (w *EmbeddingWorker) incrementProcessed() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.processedCount++
}

func (w *EmbeddingWorker) incrementFailed() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.failedCount++
}

// RetryFailed re-queues all failed embeddings
func (w *EmbeddingWorker) RetryFailed() error {
	rows, err := w.db.GetConn().Query(`
		SELECT id FROM commands 
		WHERE embedding_status = 'failed' 
		ORDER BY timestamp DESC
	`)
	if err != nil {
		return fmt.Errorf("failed to query failed embeddings: %w", err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		var cmdID int64
		if err := rows.Scan(&cmdID); err != nil {
			continue
		}

		// Reset status to pending
		w.db.UpdateEmbeddingStatus(cmdID, "pending", nil)
		
		// Enqueue
		select {
		case w.queue <- cmdID:
			count++
		default:
			// Queue full
			return fmt.Errorf("queue full after retrying %d commands", count)
		}
	}

	return nil
}
