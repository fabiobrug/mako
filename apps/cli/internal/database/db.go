package database

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
}

type Command struct {
	ID              int64
	Command         string
	Timestamp       time.Time
	ExitCode        int
	Duration        int64
	WorkingDir      string
	OutputPreview   string
	Embedding       []byte
	CommandHash     string
	LastUsed        time.Time
	EmbeddingStatus string // "pending", "processing", "completed", "failed"
}

func NewDB(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := &DB{conn: conn}

	// Enable WAL mode for better concurrent access
	_, err = conn.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	if err := db.initTables(); err != nil {
		conn.Close()
		return nil, err
	}

	return db, nil
}

func (db *DB) initTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS commands (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		command TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		exit_code INTEGER DEFAULT 0,
		duration_ms INTEGER DEFAULT 0,
		working_dir TEXT,
		output_preview TEXT,
		embedding BLOB,
		command_hash TEXT,
		last_used DATETIME,
		embedding_status TEXT DEFAULT 'pending'
	);

	CREATE VIRTUAL TABLE IF NOT EXISTS commands_fts USING fts5(
		command,
		output_preview,
		content='commands',
		content_rowid='id'
	);
	
	CREATE TRIGGER IF NOT EXISTS commands_ai AFTER INSERT ON commands BEGIN
		INSERT INTO commands_fts(rowid, command, output_preview)
		VALUES (new.id, new.command, new.output_preview);
	END;
	
	CREATE TRIGGER IF NOT EXISTS commands_ad AFTER DELETE ON commands BEGIN
		DELETE FROM commands_fts WHERE rowid = old.id;
	END;
	
	CREATE TRIGGER IF NOT EXISTS commands_au AFTER UPDATE ON commands BEGIN
		UPDATE commands_fts 
		SET command = new.command, output_preview = new.output_preview
		WHERE rowid = new.id;
	END;
	
	CREATE INDEX IF NOT EXISTS idx_timestamp ON commands(timestamp DESC);
	CREATE INDEX IF NOT EXISTS idx_working_dir ON commands(working_dir);
	CREATE INDEX IF NOT EXISTS idx_has_embedding ON commands(embedding) WHERE embedding IS NOT NULL;

	-- Embedding cache table
	CREATE TABLE IF NOT EXISTS embedding_cache (
		command_text TEXT PRIMARY KEY,
		embedding BLOB NOT NULL,
		hit_count INTEGER DEFAULT 0,
		last_accessed DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Sync metadata table
	CREATE TABLE IF NOT EXISTS sync_metadata (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.conn.Exec(schema); err != nil {
		return err
	}

	// Run migrations for existing databases
	return db.runMigrations()
}

func (db *DB) runMigrations() error {
	// Check if command_hash column exists
	var columnExists bool
	err := db.conn.QueryRow(`
		SELECT COUNT(*) > 0 
		FROM pragma_table_info('commands') 
		WHERE name='command_hash'
	`).Scan(&columnExists)
	
	if err != nil {
		return fmt.Errorf("failed to check schema: %w", err)
	}

	if !columnExists {
		migrations := []string{
			"ALTER TABLE commands ADD COLUMN command_hash TEXT",
			"ALTER TABLE commands ADD COLUMN last_used DATETIME",
			"ALTER TABLE commands ADD COLUMN embedding_status TEXT DEFAULT 'completed'",
		}

		for _, migration := range migrations {
			if _, err := db.conn.Exec(migration); err != nil {
				// Ignore errors for columns that might already exist
				continue
			}
		}

		// Update embedding_status for existing rows
		_, _ = db.conn.Exec(`
			UPDATE commands 
			SET embedding_status = CASE 
				WHEN embedding IS NOT NULL THEN 'completed'
				ELSE 'pending'
			END
			WHERE embedding_status IS NULL
		`)
	}

	// Create indexes after ensuring columns exist (safe to run multiple times)
	indexCreations := []string{
		"CREATE INDEX IF NOT EXISTS idx_embedding_status ON commands(embedding_status)",
		// Note: Don't create unique index on command_hash for existing databases
		// as it may have NULL values. New databases get it from schema above.
	}

	for _, indexSQL := range indexCreations {
		_, _ = db.conn.Exec(indexSQL)
	}

	return nil
}

func (db *DB) SaveCommand(cmd Command) error {
	hash := hashCommand(cmd.Command)
	
	query := `
		INSERT INTO commands (command, timestamp, exit_code, duration_ms, working_dir, output_preview, embedding, command_hash, embedding_status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, COALESCE(?, 'pending'))
	`

	embeddingStatus := cmd.EmbeddingStatus
	if embeddingStatus == "" {
		embeddingStatus = "pending"
	}

	_, err := db.conn.Exec(
		query,
		cmd.Command,
		cmd.Timestamp,
		cmd.ExitCode,
		cmd.Duration,
		cmd.WorkingDir,
		cmd.OutputPreview,
		cmd.Embedding,
		hash,
		embeddingStatus,
	)

	return err
}

func (db *DB) SearchCommands(query string, limit int) ([]Command, error) {
	sqlQuery := `
		SELECT c.id, c.command, c.timestamp, c.exit_code, c.duration_ms, c.working_dir, c.output_preview,
		       c.embedding, COALESCE(c.embedding_status, 'pending') as embedding_status
		FROM commands c
		JOIN commands_fts fts ON c.id = fts.rowid
		WHERE commands_fts MATCH ?
		ORDER BY c.timestamp DESC
		LIMIT ?
	`

	rows, err := db.conn.Query(sqlQuery, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []Command
	for rows.Next() {
		var cmd Command
		err := rows.Scan(
			&cmd.ID,
			&cmd.Command,
			&cmd.Timestamp,
			&cmd.ExitCode,
			&cmd.Duration,
			&cmd.WorkingDir,
			&cmd.OutputPreview,
			&cmd.Embedding,
			&cmd.EmbeddingStatus,
		)
		if err != nil {
			return nil, err
		}
		commands = append(commands, cmd)
	}

	return commands, nil
}

func (db *DB) GetRecentCommands(limit int) ([]Command, error) {
	query := `
		SELECT id, command, timestamp, exit_code, duration_ms, working_dir, output_preview, 
		       embedding, COALESCE(embedding_status, 'pending') as embedding_status
		FROM commands
		ORDER BY timestamp DESC
		LIMIT ?
	`

	rows, err := db.conn.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []Command
	for rows.Next() {
		var cmd Command
		err := rows.Scan(
			&cmd.ID,
			&cmd.Command,
			&cmd.Timestamp,
			&cmd.ExitCode,
			&cmd.Duration,
			&cmd.WorkingDir,
			&cmd.OutputPreview,
			&cmd.Embedding,
			&cmd.EmbeddingStatus,
		)
		if err != nil {
			return nil, err
		}
		commands = append(commands, cmd)
	}

	return commands, nil
}

// GetCommandsByExitCode returns commands filtered by success/failure
func (db *DB) GetCommandsByExitCode(successful bool, limit int) ([]Command, error) {
	var query string
	if successful {
		query = `
			SELECT id, command, timestamp, exit_code, duration_ms, working_dir, output_preview,
			       embedding, COALESCE(embedding_status, 'pending') as embedding_status
			FROM commands
			WHERE exit_code = 0
			ORDER BY timestamp DESC
			LIMIT ?
		`
	} else {
		query = `
			SELECT id, command, timestamp, exit_code, duration_ms, working_dir, output_preview,
			       embedding, COALESCE(embedding_status, 'pending') as embedding_status
			FROM commands
			WHERE exit_code != 0
			ORDER BY timestamp DESC
			LIMIT ?
		`
	}

	rows, err := db.conn.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []Command
	for rows.Next() {
		var cmd Command
		err := rows.Scan(
			&cmd.ID,
			&cmd.Command,
			&cmd.Timestamp,
			&cmd.ExitCode,
			&cmd.Duration,
			&cmd.WorkingDir,
			&cmd.OutputPreview,
			&cmd.Embedding,
			&cmd.EmbeddingStatus,
		)
		if err != nil {
			return nil, err
		}
		commands = append(commands, cmd)
	}

	return commands, nil
}

func (db *DB) GetCommandsByDirectory(dir string, limit int) ([]Command, error) {
	query := `
		SELECT id, command, timestamp, exit_code, duration_ms, working_dir, output_preview
		FROM commands
		WHERE working_dir = ?
		ORDER BY timestamp DESC
		LIMIT ?
	`

	rows, err := db.conn.Query(query, dir, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []Command
	for rows.Next() {
		var cmd Command
		err := rows.Scan(
			&cmd.ID,
			&cmd.Command,
			&cmd.Timestamp,
			&cmd.ExitCode,
			&cmd.Duration,
			&cmd.WorkingDir,
			&cmd.OutputPreview,
		)
		if err != nil {
			return nil, err
		}
		commands = append(commands, cmd)
	}

	return commands, nil
}

func (db *DB) GetStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var total int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM commands").Scan(&total)
	if err != nil {
		return nil, err
	}
	stats["total_commands"] = total

	var today int
	err = db.conn.QueryRow(`
		SELECT COUNT(*) FROM commands 
		WHERE DATE(timestamp) = DATE('now')
	`).Scan(&today)
	if err != nil {
		return nil, err
	}
	stats["commands_today"] = today

	var avgDuration sql.NullFloat64
	err = db.conn.QueryRow(`
		SELECT COALESCE(AVG(duration_ms), 0) FROM commands WHERE duration_ms > 0
	`).Scan(&avgDuration)
	if err != nil {
		return nil, err
	}

	if avgDuration.Valid {
		stats["avg_duration_ms"] = avgDuration.Float64
	} else {
		stats["avg_duration_ms"] = 0.0
	}

	return stats, nil
}

// SearchCommandsSemantic performs two-phase semantic search
// Phase 1: Use FTS5 for keyword filtering (fast)
// Phase 2: Rank by vector similarity (accurate)
func (db *DB) SearchCommandsSemantic(query string, queryEmbedding []byte, limit int, threshold float32) ([]Command, error) {
	type scoredCommand struct {
		cmd   Command
		score float32
	}

	var candidates []scoredCommand

	// Phase 1: FTS5 keyword search to narrow down candidates
	ftsQuery := `
		SELECT c.id, c.command, c.timestamp, c.exit_code, c.duration_ms, c.working_dir, c.output_preview, c.embedding
		FROM commands c
		JOIN commands_fts fts ON c.id = fts.rowid
		WHERE commands_fts MATCH ? AND c.embedding IS NOT NULL
		ORDER BY c.timestamp DESC
		LIMIT 1000
	`

	rows, err := db.conn.Query(ftsQuery, query)
	if err != nil {
		// FTS might fail on certain queries, fallback to recent commands
		rows, err = db.conn.Query(`
			SELECT id, command, timestamp, exit_code, duration_ms, working_dir, output_preview, embedding
			FROM commands
			WHERE embedding IS NOT NULL
			ORDER BY timestamp DESC
			LIMIT 1000
		`)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	// Collect FTS results
	var ftsResults []Command
	for rows.Next() {
		var cmd Command
		var embeddingBytes []byte

		err := rows.Scan(
			&cmd.ID,
			&cmd.Command,
			&cmd.Timestamp,
			&cmd.ExitCode,
			&cmd.Duration,
			&cmd.WorkingDir,
			&cmd.OutputPreview,
			&embeddingBytes,
		)
		if err != nil {
			continue
		}

		cmd.Embedding = embeddingBytes
		ftsResults = append(ftsResults, cmd)
	}

	// If FTS returned < 50 results, expand to recent 1000
	if len(ftsResults) < 50 {
		fallbackQuery := `
			SELECT id, command, timestamp, exit_code, duration_ms, working_dir, output_preview, embedding
			FROM commands
			WHERE embedding IS NOT NULL
			ORDER BY timestamp DESC
			LIMIT 1000
		`

		rows, err := db.conn.Query(fallbackQuery)
		if err != nil {
			// Use what we have from FTS
		} else {
			defer rows.Close()
			ftsResults = nil // Clear and rebuild

			for rows.Next() {
				var cmd Command
				var embeddingBytes []byte

				err := rows.Scan(
					&cmd.ID,
					&cmd.Command,
					&cmd.Timestamp,
					&cmd.ExitCode,
					&cmd.Duration,
					&cmd.WorkingDir,
					&cmd.OutputPreview,
					&embeddingBytes,
				)
				if err != nil {
					continue
				}

				cmd.Embedding = embeddingBytes
				ftsResults = append(ftsResults, cmd)
			}
		}
	}

	// Phase 2: Rank by vector similarity
	for _, cmd := range ftsResults {
		similarity := calculateSimilarity(queryEmbedding, cmd.Embedding)

		if similarity >= threshold {
			candidates = append(candidates, scoredCommand{
				cmd:   cmd,
				score: similarity,
			})
		}
	}

	// Sort by similarity score (bubble sort for simplicity, could use sort.Slice)
	for i := 0; i < len(candidates); i++ {
		for j := i + 1; j < len(candidates); j++ {
			if candidates[j].score > candidates[i].score {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			}
		}
	}

	// Return top results
	var results []Command
	for i := 0; i < len(candidates) && i < limit; i++ {
		results = append(results, candidates[i].cmd)
	}

	return results, nil
}

func calculateSimilarity(a, b []byte) float32 {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}

	vecA := make([]float32, len(a)/4)
	vecB := make([]float32, len(b)/4)

	bufA := bytes.NewReader(a)
	bufB := bytes.NewReader(b)

	binary.Read(bufA, binary.LittleEndian, &vecA)
	binary.Read(bufB, binary.LittleEndian, &vecB)

	var dotProduct, normA, normB float64

	for i := range vecA {
		if i >= len(vecB) {
			break
		}
		dotProduct += float64(vecA[i]) * float64(vecB[i])
		normA += float64(vecA[i]) * float64(vecA[i])
		normB += float64(vecB[i]) * float64(vecB[i])
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return float32(dotProduct / (math.Sqrt(normA) * math.Sqrt(normB)))
}

// SaveCommandAsync saves a command without blocking on embedding generation
func (db *DB) SaveCommandAsync(cmd Command) (int64, error) {
	hash := hashCommand(cmd.Command)
	
	query := `
		INSERT INTO commands (command, timestamp, exit_code, duration_ms, working_dir, output_preview, command_hash, last_used, embedding_status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'pending')
	`

	result, err := db.conn.Exec(
		query,
		cmd.Command,
		cmd.Timestamp,
		cmd.ExitCode,
		cmd.Duration,
		cmd.WorkingDir,
		cmd.OutputPreview,
		hash,
		cmd.Timestamp,
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// SaveCommandDeduplicated saves a command with deduplication
// If command exists, updates last_used timestamp
// Returns (isNew, commandID, error)
func (db *DB) SaveCommandDeduplicated(cmd Command) (bool, int64, error) {
	hash := hashCommand(cmd.Command)
	
	// Check if command exists
	var existingID int64
	err := db.conn.QueryRow(`
		SELECT id FROM commands WHERE command_hash = ?
	`, hash).Scan(&existingID)

	if err == sql.ErrNoRows {
		// New command
		id, err := db.SaveCommandAsync(cmd)
		return true, id, err
	} else if err != nil {
		return false, 0, err
	}

	// Update existing command's last_used timestamp
	_, err = db.conn.Exec(`
		UPDATE commands SET last_used = ? WHERE id = ?
	`, cmd.Timestamp, existingID)

	return false, existingID, err
}

// GetCommandByHash retrieves a command by its hash
func (db *DB) GetCommandByHash(hash string) (*Command, error) {
	query := `
		SELECT id, command, timestamp, exit_code, duration_ms, working_dir, output_preview, embedding, embedding_status
		FROM commands
		WHERE command_hash = ?
	`

	var cmd Command
	var embeddingStatus sql.NullString
	
	err := db.conn.QueryRow(query, hash).Scan(
		&cmd.ID,
		&cmd.Command,
		&cmd.Timestamp,
		&cmd.ExitCode,
		&cmd.Duration,
		&cmd.WorkingDir,
		&cmd.OutputPreview,
		&cmd.Embedding,
		&embeddingStatus,
	)

	if err != nil {
		return nil, err
	}

	if embeddingStatus.Valid {
		cmd.EmbeddingStatus = embeddingStatus.String
	}

	return &cmd, nil
}

// GetEmbeddingStatus returns the embedding status for a command
func (db *DB) GetEmbeddingStatus(cmdID int64) (string, error) {
	var status string
	err := db.conn.QueryRow(`
		SELECT embedding_status FROM commands WHERE id = ?
	`, cmdID).Scan(&status)
	
	return status, err
}

// UpdateEmbeddingStatus updates the embedding generation status
func (db *DB) UpdateEmbeddingStatus(cmdID int64, status string, embedding []byte) error {
	_, err := db.conn.Exec(`
		UPDATE commands 
		SET embedding_status = ?, embedding = ?
		WHERE id = ?
	`, status, embedding, cmdID)
	
	return err
}

// GetPendingEmbeddings returns commands that need embeddings generated
func (db *DB) GetPendingEmbeddings(limit int) ([]Command, error) {
	query := `
		SELECT id, command, timestamp, exit_code, duration_ms, working_dir, output_preview,
		       embedding, COALESCE(embedding_status, 'pending') as embedding_status
		FROM commands
		WHERE COALESCE(embedding_status, 'pending') = 'pending'
		ORDER BY timestamp DESC
		LIMIT ?
	`

	rows, err := db.conn.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []Command
	for rows.Next() {
		var cmd Command
		err := rows.Scan(
			&cmd.ID,
			&cmd.Command,
			&cmd.Timestamp,
			&cmd.ExitCode,
			&cmd.Duration,
			&cmd.WorkingDir,
			&cmd.OutputPreview,
			&cmd.Embedding,
			&cmd.EmbeddingStatus,
		)
		if err != nil {
			return nil, err
		}
		commands = append(commands, cmd)
	}

	return commands, nil
}

// BulkInsertCommands inserts multiple commands in a transaction
func (db *DB) BulkInsertCommands(cmds []Command) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO commands (command, timestamp, exit_code, duration_ms, working_dir, output_preview, command_hash, last_used, embedding_status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'pending')
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, cmd := range cmds {
		hash := hashCommand(cmd.Command)
		_, err = stmt.Exec(
			cmd.Command,
			cmd.Timestamp,
			cmd.ExitCode,
			cmd.Duration,
			cmd.WorkingDir,
			cmd.OutputPreview,
			hash,
			cmd.Timestamp,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetLastSyncTime retrieves the last sync timestamp
func (db *DB) GetLastSyncTime() (time.Time, error) {
	var timeStr string
	err := db.conn.QueryRow(`
		SELECT value FROM sync_metadata WHERE key = 'last_sync'
	`).Scan(&timeStr)

	if err == sql.ErrNoRows {
		return time.Time{}, nil
	} else if err != nil {
		return time.Time{}, err
	}

	return time.Parse(time.RFC3339, timeStr)
}

// SetLastSyncTime updates the last sync timestamp
func (db *DB) SetLastSyncTime(t time.Time) error {
	_, err := db.conn.Exec(`
		INSERT OR REPLACE INTO sync_metadata (key, value, updated_at)
		VALUES ('last_sync', ?, ?)
	`, t.Format(time.RFC3339), time.Now())

	return err
}

// GetDatabaseSize returns the database file size in bytes
func (db *DB) GetDatabaseSize() (int64, error) {
	var pageCount, pageSize int64
	
	err := db.conn.QueryRow("PRAGMA page_count").Scan(&pageCount)
	if err != nil {
		return 0, err
	}
	
	err = db.conn.QueryRow("PRAGMA page_size").Scan(&pageSize)
	if err != nil {
		return 0, err
	}
	
	return pageCount * pageSize, nil
}

// GetCommandCount returns total number of commands
func (db *DB) GetCommandCount() (int64, error) {
	var count int64
	err := db.conn.QueryRow("SELECT COUNT(*) FROM commands").Scan(&count)
	return count, err
}

// hashCommand generates a SHA256 hash of a command
func hashCommand(command string) string {
	h := sha256.Sum256([]byte(command))
	return hex.EncodeToString(h[:])
}

// GetConn returns the underlying database connection
// This is needed for advanced operations like cache loading/saving
func (db *DB) GetConn() *sql.DB {
	return db.conn
}

func (db *DB) Close() error {
	return db.conn.Close()
}
