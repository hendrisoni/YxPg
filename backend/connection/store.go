package connection

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"yxpg/backend/models"

	"github.com/google/uuid"
)

// Store manages persistence of connection configurations to JSON file
type Store struct {
	mu          sync.RWMutex
	filePath    string
	Connections []models.Connection
}

// NewStore creates a new connection store
func NewStore() (*Store, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(homeDir, ".yxpg")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	store := &Store{
		filePath: filepath.Join(configDir, "connections.json"),
	}

	if err := store.load(); err != nil {
		// Start with empty connections if file doesn't exist
		store.Connections = []models.Connection{}
	}

	return store, nil
}

// load reads connections from the JSON file
func (s *Store) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &s.Connections)
}

// save writes connections to the JSON file
func (s *Store) save() error {
	data, err := json.MarshalIndent(s.Connections, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0644)
}

// List returns all saved connections
func (s *Store) List() []models.Connection {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Connections
}

// Get returns a connection by ID
func (s *Store) Get(id string) (*models.Connection, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := range s.Connections {
		if s.Connections[i].ID == id {
			return &s.Connections[i], nil
		}
	}

	return nil, ErrConnectionNotFound
}

// Add saves a new connection
func (s *Store) Add(conn models.Connection) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check for duplicates
	for _, c := range s.Connections {
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

	s.Connections = append(s.Connections, conn)
	return s.save()
}

// Update modifies an existing connection
func (s *Store) Update(conn models.Connection) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check for duplicates (excluding itself)
	for _, c := range s.Connections {
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

	for i := range s.Connections {
		if s.Connections[i].ID == conn.ID {
			s.Connections[i] = conn
			return s.save()
		}
	}

	return ErrConnectionNotFound
}

// Delete removes a connection by ID
func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.Connections {
		if s.Connections[i].ID == id {
			s.Connections = append(s.Connections[:i], s.Connections[i+1:]...)
			return s.save()
		}
	}

	return ErrConnectionNotFound
}
