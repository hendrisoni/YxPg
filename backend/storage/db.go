package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	_ "modernc.org/sqlite"
)

// Storage manages the SQLite central database for YxPg
type Storage struct {
	mu     sync.RWMutex
	db     *sql.DB
	dbPath string
}

// GetDefaultDbPath resolves the default location for yxpg.sys alongside the executable
func GetDefaultDbPath() string {
	if exePath, err := os.Executable(); err == nil {
		dir := filepath.Dir(exePath)
		if !strings.Contains(strings.ToLower(dir), "go-build") {
			return filepath.Join(dir, "yxpg.sys")
		}
	}
	if cwd, err := os.Getwd(); err == nil {
		return filepath.Join(cwd, "yxpg.sys")
	}
	return `D:\Yx\YxPg\yxpg.sys`
}

// NewStorage initializes SQLite storage at the specified or default path
func NewStorage(targetPath string) (*Storage, error) {
	s := &Storage{}
	if err := s.Open(targetPath); err != nil {
		return nil, err
	}
	return s, nil
}

// Open opens/re-opens the SQLite database at targetPath
func (s *Storage) Open(targetPath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	finalPath := targetPath
	if strings.TrimSpace(finalPath) == "" {
		finalPath = GetDefaultDbPath()
	}

	absPath, err := filepath.Abs(finalPath)
	if err == nil {
		finalPath = absPath
	}

	if err := os.MkdirAll(filepath.Dir(finalPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory for database: %w", err)
	}

	if s.db != nil {
		_ = s.db.Close()
	}

	db, err := sql.Open("sqlite", finalPath)
	if err != nil {
		return fmt.Errorf("failed to open sqlite database at %s: %w", finalPath, err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("failed to ping sqlite database at %s: %w", finalPath, err)
	}

	s.db = db
	s.dbPath = finalPath

	if err := s.initTablesLocked(); err != nil {
		fmt.Printf("[Storage] Warning: table init error: %v\n", err)
	}

	return nil
}

// initTablesLocked creates base tables if they don't exist
func (s *Storage) initTablesLocked() error {
	if s.db == nil {
		return fmt.Errorf("database connection is nil")
	}

	tables := []string{
		`CREATE TABLE IF NOT EXISTS connections (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			host TEXT NOT NULL,
			port INTEGER NOT NULL,
			database TEXT NOT NULL,
			username TEXT NOT NULL,
			password TEXT NOT NULL,
			ssl_mode TEXT DEFAULT 'disable',
			color TEXT DEFAULT '#00C9A7',
			created_at DATETIME NOT NULL
		);`,
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
		);`,
		`CREATE TABLE IF NOT EXISTS saved_queries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			sql TEXT NOT NULL,
			folder TEXT DEFAULT '',
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS workspace (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			data TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS app_settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		);`,
		`CREATE INDEX IF NOT EXISTS idx_history_connection ON query_history(connection_id);`,
		`CREATE INDEX IF NOT EXISTS idx_history_executed ON query_history(executed_at DESC);`,
	}

	for _, stmt := range tables {
		if _, err := s.db.Exec(stmt); err != nil {
			return fmt.Errorf("migration failed on stmt: %s: %w", stmt, err)
		}
	}

	return nil
}

// DB returns the underlying *sql.DB
func (s *Storage) DB() *sql.DB {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.db
}

// Path returns the current database file path
func (s *Storage) Path() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.dbPath
}

// Close closes the database connection
func (s *Storage) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db != nil {
		err := s.db.Close()
		s.db = nil
		return err
	}
	return nil
}
