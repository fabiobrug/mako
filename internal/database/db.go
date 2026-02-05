package database

import (
	"database/sql"
	"fmt"
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
		output_preview TEXT
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
	`

	_, err := db.conn.Exec(schema)
	return err
}

func (db *DB) SaveCommand(cmd Command) error {
	query := `
		INSERT INTO commands (command, timestamp, exit_code, duration_ms, working_dir, output_preview)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := db.conn.Exec(
		query,
		cmd.Command,
		cmd.Timestamp,
		cmd.ExitCode,
		cmd.Duration,
		cmd.WorkingDir,
		cmd.OutputPreview,
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

func (db *DB) Close() error {
	return db.conn.Close()
}
