package query

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"yxpg/backend/models"

	_ "github.com/mattn/go-sqlite3"
)

// History manages query history using SQLite
type History struct {
	db *sql.DB
}

// NewHistory creates a new history manager
func NewHistory() (*History, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(homeDir, ".yxpg")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(configDir, "history.db")
	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("failed to open history database: %w", err)
	}

	h := &History{db: db}
	if err := h.migrate(); err != nil {
		db.Close()
		return nil, err
	}

	return h, nil
}

func (h *History) migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS query_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			connection_id TEXT NOT NULL,
			database TEXT DEFAULT '',
			sql TEXT NOT NULL,
			duration_ms INTEGER DEFAULT 0,
			rows_returned INTEGER DEFAULT 0,
			executed_at DATETIME NOT NULL,
			error TEXT DEFAULT '',
			bookmarked BOOLEAN DEFAULT FALSE
		)`,
		`CREATE TABLE IF NOT EXISTS saved_queries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			sql TEXT NOT NULL,
			folder TEXT DEFAULT '',
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_history_connection ON query_history(connection_id)`,
		`CREATE INDEX IF NOT EXISTS idx_history_executed ON query_history(executed_at DESC)`,
	}

	for _, q := range queries {
		if _, err := h.db.Exec(q); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}

// Save saves a query history entry
func (h *History) Save(entry models.QueryHistoryEntry) error {
	_, err := h.db.Exec(
		`INSERT INTO query_history (connection_id, database, sql, duration_ms, rows_returned, executed_at, error)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		entry.ConnectionID, entry.Database, entry.SQL, entry.DurationMs,
		entry.RowsReturned, entry.ExecutedAt, entry.Error,
	)

	if err != nil {
		return err
	}

	// Prune old entries (keep max 10,000)
	h.db.Exec(`DELETE FROM query_history WHERE id NOT IN (SELECT id FROM query_history ORDER BY executed_at DESC LIMIT 10000)`)

	return nil
}

// ListByConnection returns history entries for a connection
func (h *History) ListByConnection(connID string, limit int) ([]models.QueryHistoryEntry, error) {
	if limit <= 0 {
		limit = 100
	}

	rows, err := h.db.Query(
		`SELECT id, connection_id, database, sql, duration_ms, rows_returned, executed_at, error, bookmarked
		 FROM query_history
		 WHERE connection_id = ?
		 ORDER BY executed_at DESC
		 LIMIT ?`,
		connID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.QueryHistoryEntry
	for rows.Next() {
		var e models.QueryHistoryEntry
		if err := rows.Scan(
			&e.ID, &e.ConnectionID, &e.Database, &e.SQL,
			&e.DurationMs, &e.RowsReturned, &e.ExecutedAt,
			&e.Error, &e.Bookmarked,
		); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	return entries, rows.Err()
}

// Search searches history by SQL text
func (h *History) Search(query string, limit int) ([]models.QueryHistoryEntry, error) {
	if limit <= 0 {
		limit = 100
	}

	rows, err := h.db.Query(
		`SELECT id, connection_id, database, sql, duration_ms, rows_returned, executed_at, error, bookmarked
		 FROM query_history
		 WHERE sql LIKE ?
		 ORDER BY executed_at DESC
		 LIMIT ?`,
		"%"+query+"%", limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.QueryHistoryEntry
	for rows.Next() {
		var e models.QueryHistoryEntry
		if err := rows.Scan(
			&e.ID, &e.ConnectionID, &e.Database, &e.SQL,
			&e.DurationMs, &e.RowsReturned, &e.ExecutedAt,
			&e.Error, &e.Bookmarked,
		); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	return entries, rows.Err()
}

// ToggleBookmark toggles the bookmark status of a history entry
func (h *History) ToggleBookmark(id int64) error {
	_, err := h.db.Exec(`UPDATE query_history SET bookmarked = NOT bookmarked WHERE id = ?`, id)
	return err
}

// SaveQuery saves a named query
func (h *History) SaveQuery(name, sql string) error {
	now := time.Now()
	_, err := h.db.Exec(
		`INSERT INTO saved_queries (name, sql, created_at, updated_at) VALUES (?, ?, ?, ?)`,
		name, sql, now, now,
	)
	return err
}

// ListSavedQueries returns all saved queries
func (h *History) ListSavedQueries() ([]models.SavedQuery, error) {
	query := "SELECT id, name, sql, folder, created_at, updated_at FROM saved_queries ORDER BY name"
	rows, err := h.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var queries []models.SavedQuery
	for rows.Next() {
		var q models.SavedQuery
		if err := rows.Scan(&q.ID, &q.Name, &q.SQL, &q.Folder, &q.CreatedAt, &q.UpdatedAt); err != nil {
			return nil, err
		}
		queries = append(queries, q)
	}

	return queries, rows.Err()
}

// Close closes the database connection
func (h *History) Close() error {
	return h.db.Close()
}
