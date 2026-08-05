package connection

import (
	"database/sql"
	"os"
	"path/filepath"
	"sync"

	_ "modernc.org/sqlite"
)

// WorkspaceStore manages workspace tree persistence in SQLite
type WorkspaceStore struct {
	mu   sync.RWMutex
	db   *sql.DB
	Data string // JSON string representing the workspace tree
}

// NewWorkspaceStore creates a new workspace store backed by SQLite
func NewWorkspaceStore(db *sql.DB) (*WorkspaceStore, error) {
	s := &WorkspaceStore{db: db}
	if err := s.initAndMigrate(); err != nil {
		return nil, err
	}
	return s, nil
}

// SetDB updates the database instance for the workspace store
func (s *WorkspaceStore) SetDB(db *sql.DB) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.db = db
	s.Data = "[]"
	return s.initAndMigrateLocked()
}

func (s *WorkspaceStore) initAndMigrate() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.initAndMigrateLocked()
}

func (s *WorkspaceStore) initAndMigrateLocked() error {
	if s.db == nil {
		s.Data = "[]"
		return nil
	}

	tableStmt := `CREATE TABLE IF NOT EXISTS workspace (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		data TEXT NOT NULL
	);`

	if _, err := s.db.Exec(tableStmt); err != nil {
		return err
	}

	var count int
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM workspace`).Scan(&count)
	if count == 0 {
		s.importLegacyJSON()
	}

	return s.loadLocked()
}

func (s *WorkspaceStore) importLegacyJSON() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	legacyPath := filepath.Join(homeDir, ".yxpg", "workspace.json")
	data, err := os.ReadFile(legacyPath)
	if err != nil || len(data) == 0 {
		return
	}

	_, _ = s.db.Exec(`INSERT OR REPLACE INTO workspace (id, data) VALUES (1, ?)`, string(data))
}

// Load reads workspace from the SQLite database
func (s *WorkspaceStore) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.loadLocked()
}

func (s *WorkspaceStore) loadLocked() error {
	if s.db == nil {
		s.Data = "[]"
		return nil
	}

	var data string
	err := s.db.QueryRow(`SELECT data FROM workspace WHERE id = 1`).Scan(&data)
	if err != nil {
		s.Data = "[]"
		return nil
	}

	s.Data = data
	return nil
}

// Save writes workspace JSON to the SQLite database
func (s *WorkspaceStore) Save(data string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Data = data
	if s.db == nil {
		return nil
	}

	_, err := s.db.Exec(`INSERT OR REPLACE INTO workspace (id, data) VALUES (1, ?)`, data)
	return err
}
