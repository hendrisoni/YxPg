package query

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"yxpg/backend/models"

	_ "modernc.org/sqlite"
)

// History manages query history using SQLite
type History struct {
	mu sync.RWMutex
	db *sql.DB
}

// NewHistory creates a new history manager using the provided SQLite connection
func NewHistory(db *sql.DB) (*History, error) {
	h := &History{db: db}
	if err := h.migrateAndImport(); err != nil {
		return nil, err
	}
	return h, nil
}

// SetDB updates the database connection for history manager
func (h *History) SetDB(db *sql.DB) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.db = db
	return h.migrateAndImportLocked()
}

func (h *History) migrateAndImport() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.migrateAndImportLocked()
}

func (h *History) migrateAndImportLocked() error {
	if h.db == nil {
		return fmt.Errorf("history database connection is nil")
	}

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
			return fmt.Errorf("history migration failed: %w", err)
		}
	}

	var count int
	_ = h.db.QueryRow(`SELECT COUNT(*) FROM query_history`).Scan(&count)
	if count == 0 {
		h.importLegacyHistory()
	}

	return nil
}

func (h *History) importLegacyHistory() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	legacyPath := filepath.Join(homeDir, ".yxpg", "history.db")
	if _, err := os.Stat(legacyPath); os.IsNotExist(err) {
		return
	}

	cleanPath := filepath.ToSlash(legacyPath)
	attachQuery := fmt.Sprintf("ATTACH DATABASE '%s' AS legacy_hist", cleanPath)
	if _, err := h.db.Exec(attachQuery); err != nil {
		return
	}
	defer func() {
		_, _ = h.db.Exec("DETACH DATABASE legacy_hist")
	}()

	_, _ = h.db.Exec(`INSERT OR IGNORE INTO query_history (id, connection_id, database, sql, duration_ms, rows_returned, executed_at, error, bookmarked)
		SELECT id, connection_id, database, sql, duration_ms, rows_returned, executed_at, error, bookmarked FROM legacy_hist.query_history`)

	_, _ = h.db.Exec(`INSERT OR IGNORE INTO saved_queries (id, name, sql, folder, created_at, updated_at)
		SELECT id, name, sql, folder, created_at, updated_at FROM legacy_hist.saved_queries`)
}

// Save saves a query history entry
func (h *History) Save(entry models.QueryHistoryEntry) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.db == nil {
		return nil
	}

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
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.db == nil {
		return []models.QueryHistoryEntry{}, nil
	}

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

	if entries == nil {
		entries = []models.QueryHistoryEntry{}
	}

	return entries, rows.Err()
}

// Search searches history by SQL text
func (h *History) Search(query string, limit int) ([]models.QueryHistoryEntry, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.db == nil {
		return []models.QueryHistoryEntry{}, nil
	}

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

	if entries == nil {
		entries = []models.QueryHistoryEntry{}
	}

	return entries, rows.Err()
}

// ToggleBookmark toggles the bookmark status of a history entry
func (h *History) ToggleBookmark(id int64) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.db == nil {
		return nil
	}

	_, err := h.db.Exec(`UPDATE query_history SET bookmarked = NOT bookmarked WHERE id = ?`, id)
	return err
}

// SaveQuery saves a named query
func (h *History) SaveQuery(name, sql string) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.db == nil {
		return nil
	}

	now := time.Now()
	_, err := h.db.Exec(
		`INSERT INTO saved_queries (name, sql, created_at, updated_at) VALUES (?, ?, ?, ?)`,
		name, sql, now, now,
	)
	return err
}

// ListSavedQueries returns all saved queries
func (h *History) ListSavedQueries() ([]models.SavedQuery, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.db == nil {
		return []models.SavedQuery{}, nil
	}

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

	if queries == nil {
		queries = []models.SavedQuery{}
	}

	return queries, rows.Err()
}

// Close closes the database connection (handled by storage)
func (h *History) Close() error {
	return nil
}
