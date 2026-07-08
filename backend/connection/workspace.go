package connection

import (
	"os"
	"path/filepath"
	"sync"
)

// WorkspaceStore manages workspace tree persistence
type WorkspaceStore struct {
	mu       sync.RWMutex
	filePath string
	Data     string // JSON string representing the workspace tree
}

// NewWorkspaceStore creates a new workspace store
func NewWorkspaceStore() (*WorkspaceStore, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(homeDir, ".yxpg")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	s := &WorkspaceStore{
		filePath: filepath.Join(configDir, "workspace.json"),
	}

	if err := s.Load(); err != nil {
		s.Data = "[]" // Empty tree by default
	}

	return s, nil
}

// Load reads workspace from the JSON file
func (s *WorkspaceStore) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	s.Data = string(data)
	return nil
}

// Save writes workspace JSON to the file
func (s *WorkspaceStore) Save(data string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Data = data
	return os.WriteFile(s.filePath, []byte(data), 0644)
}
