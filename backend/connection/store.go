package connection

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"yxpg/backend/models"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

// Store manages persistence of connection configurations in SQLite
type Store struct {
	mu sync.RWMutex
	db *sql.DB
}

// NewStore creates a new connection store backed by SQLite
func NewStore(db *sql.DB) (*Store, error) {
	store := &Store{db: db}
	if err := store.initAndMigrate(); err != nil {
		return nil, err
	}
	return store, nil
}

// SetDB updates the database instance for the store
func (s *Store) SetDB(db *sql.DB) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.db = db
	return s.initAndMigrateLocked()
}

func (s *Store) initAndMigrate() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.initAndMigrateLocked()
}

func (s *Store) initAndMigrateLocked() error {
	if s.db == nil {
		return fmt.Errorf("store db is nil")
	}

	tableStmt := `CREATE TABLE IF NOT EXISTS connections (
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
	);`

	if _, err := s.db.Exec(tableStmt); err != nil {
		return fmt.Errorf("failed to create connections table: %w", err)
	}

	// Check if table is empty to attempt legacy JSON import
	var count int
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM connections`).Scan(&count)
	if count == 0 {
		s.importLegacyJSON()
	}

	return nil
}

func (s *Store) importLegacyJSON() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	legacyPath := filepath.Join(homeDir, ".yxpg", "connections.json")
	data, err := os.ReadFile(legacyPath)
	if err != nil {
		return
	}

	var conns []models.Connection
	if err := json.Unmarshal(data, &conns); err != nil {
		return
	}

	for _, conn := range conns {
		if conn.ID == "" {
			conn.ID = uuid.New().String()
		}
		if conn.CreatedAt.IsZero() {
			conn.CreatedAt = time.Now()
		}
		if conn.SSLMode == "" {
			conn.SSLMode = "disable"
		}
		_, _ = s.db.Exec(
			`INSERT OR IGNORE INTO connections (id, name, host, port, database, username, password, ssl_mode, color, created_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			conn.ID, conn.Name, conn.Host, conn.Port, conn.Database, conn.Username, conn.Password, conn.SSLMode, conn.Color, conn.CreatedAt,
		)
	}
}

// List returns all saved connections
func (s *Store) List() []models.Connection {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []models.Connection
	if s.db == nil {
		return result
	}

	rows, err := s.db.Query(`SELECT id, name, host, port, database, username, password, ssl_mode, color, created_at FROM connections ORDER BY name`)
	if err != nil {
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Connection
		if err := rows.Scan(&c.ID, &c.Name, &c.Host, &c.Port, &c.Database, &c.Username, &c.Password, &c.SSLMode, &c.Color, &c.CreatedAt); err == nil {
			result = append(result, c)
		}
	}

	if result == nil {
		result = []models.Connection{}
	}
	return result
}

// Get returns a connection by ID
func (s *Store) Get(id string) (*models.Connection, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.db == nil {
		return nil, ErrConnectionNotFound
	}

	var c models.Connection
	err := s.db.QueryRow(
		`SELECT id, name, host, port, database, username, password, ssl_mode, color, created_at FROM connections WHERE id = ?`,
		id,
	).Scan(&c.ID, &c.Name, &c.Host, &c.Port, &c.Database, &c.Username, &c.Password, &c.SSLMode, &c.Color, &c.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrConnectionNotFound
	} else if err != nil {
		return nil, err
	}

	return &c, nil
}

// Add saves a new connection
func (s *Store) Add(conn models.Connection) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db == nil {
		return fmt.Errorf("store database connection is closed")
	}

	// Check for duplicates
	conns := s.listLocked()
	for _, c := range conns {
		if c.Name == conn.Name {
			return fmt.Errorf("connection name '%s' already exists", conn.Name)
		}
		if c.Host == conn.Host && c.Port == conn.Port && c.Database == conn.Database && c.Username == conn.Username {
			return fmt.Errorf("connection parameters already exist (see connection '%s')", c.Name)
		}
	}

	if conn.ID == "" {
		conn.ID = uuid.New().String()
	}
	if conn.CreatedAt.IsZero() {
		conn.CreatedAt = time.Now()
	}
	if conn.SSLMode == "" {
		conn.SSLMode = "disable"
	}

	_, err := s.db.Exec(
		`INSERT INTO connections (id, name, host, port, database, username, password, ssl_mode, color, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		conn.ID, conn.Name, conn.Host, conn.Port, conn.Database, conn.Username, conn.Password, conn.SSLMode, conn.Color, conn.CreatedAt,
	)

	return err
}

// Update modifies an existing connection
func (s *Store) Update(conn models.Connection) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db == nil {
		return fmt.Errorf("store database connection is closed")
	}

	// Check duplicates excluding itself
	conns := s.listLocked()
	for _, c := range conns {
		if c.ID == conn.ID {
			continue
		}
		if c.Name == conn.Name {
			return fmt.Errorf("connection name '%s' already exists", conn.Name)
		}
		if c.Host == conn.Host && c.Port == conn.Port && c.Database == conn.Database && c.Username == conn.Username {
			return fmt.Errorf("connection parameters already exist (see connection '%s')", c.Name)
		}
	}

	res, err := s.db.Exec(
		`UPDATE connections SET name = ?, host = ?, port = ?, database = ?, username = ?, password = ?, ssl_mode = ?, color = ? WHERE id = ?`,
		conn.Name, conn.Host, conn.Port, conn.Database, conn.Username, conn.Password, conn.SSLMode, conn.Color, conn.ID,
	)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrConnectionNotFound
	}

	return nil
}

// Delete removes a connection by ID
func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db == nil {
		return fmt.Errorf("store database connection is closed")
	}

	res, err := s.db.Exec(`DELETE FROM connections WHERE id = ?`, id)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrConnectionNotFound
	}

	return nil
}

func (s *Store) listLocked() []models.Connection {
	var result []models.Connection
	if s.db == nil {
		return result
	}

	rows, err := s.db.Query(`SELECT id, name, host, port, database, username, password, ssl_mode, color, created_at FROM connections`)
	if err != nil {
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Connection
		if err := rows.Scan(&c.ID, &c.Name, &c.Host, &c.Port, &c.Database, &c.Username, &c.Password, &c.SSLMode, &c.Color, &c.CreatedAt); err == nil {
			result = append(result, c)
		}
	}
	return result
}
