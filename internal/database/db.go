package database

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"fmt"
	"math"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

type Command struct {
	ID            int64
	Command       string
	Timestamp     time.Time
	ExitCode      int
	Duration      int64
	WorkingDir    string
	OutputPreview string
	Embedding     []byte
}

func NewDB(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := &DB{conn: conn}

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
		embedding BLOB
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
	`

	_, err := db.conn.Exec(schema)
	return err
}

func (db *DB) SaveCommand(cmd Command) error {
	query := `
		INSERT INTO commands (command, timestamp, exit_code, duration_ms, working_dir, output_preview, embedding)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.conn.Exec(
		query,
		cmd.Command,
		cmd.Timestamp,
		cmd.ExitCode,
		cmd.Duration,
		cmd.WorkingDir,
		cmd.OutputPreview,
		cmd.Embedding,
	)

	return err
}

func (db *DB) SearchCommands(query string, limit int) ([]Command, error) {
	sqlQuery := `
		SELECT c.id, c.command, c.timestamp, c.exit_code, c.duration_ms, c.working_dir, c.output_preview
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
		SELECT id, command, timestamp, exit_code, duration_ms, working_dir, output_preview
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
			SELECT id, command, timestamp, exit_code, duration_ms, working_dir, output_preview
			FROM commands
			WHERE exit_code = 0
			ORDER BY timestamp DESC
			LIMIT ?
		`
	} else {
		query = `
			SELECT id, command, timestamp, exit_code, duration_ms, working_dir, output_preview
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

func (db *DB) SearchCommandsSemantic(queryEmbedding []byte, limit int, threshold float32) ([]Command, error) {
	query := `
		SELECT id, command, timestamp, exit_code, duration_ms, working_dir, output_preview, embedding
		FROM commands
		WHERE embedding IS NOT NULL
		ORDER BY timestamp DESC
		LIMIT 100
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type scoredCommand struct {
		cmd   Command
		score float32
	}

	var candidates []scoredCommand

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

		similarity := calculateSimilarity(queryEmbedding, embeddingBytes)

		if similarity >= threshold {
			candidates = append(candidates, scoredCommand{
				cmd:   cmd,
				score: similarity,
			})
		}
	}

	for i := 0; i < len(candidates); i++ {
		for j := i + 1; j < len(candidates); j++ {
			if candidates[j].score > candidates[i].score {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			}
		}
	}

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

func (db *DB) Close() error {
	return db.conn.Close()
}
