package connection

import (
	"context"
	"fmt"
	"sync"
	"time"

	"yxpg/backend/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Manager manages multiple PostgreSQL connection pools
type Manager struct {
	mu          sync.RWMutex
	pools       map[string]*pgxpool.Pool
	configs     map[string]*pgxpool.Config
	store       *Store
	cancelFuncs map[string]context.CancelFunc
}

// NewManager creates a new connection manager
func NewManager() (*Manager, error) {
	store, err := NewStore()
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	return &Manager{
		pools:       make(map[string]*pgxpool.Pool),
		configs:     make(map[string]*pgxpool.Config),
		store:       store,
		cancelFuncs: make(map[string]context.CancelFunc),
	}, nil
}

// buildDSN constructs a PostgreSQL connection string
func buildDSN(conn models.Connection) string {
	sslMode := conn.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		conn.Host, conn.Port, conn.Username, conn.Password, conn.Database, sslMode,
	)
}

// ListConnections returns all saved connections
func (m *Manager) ListConnections() []models.Connection {
	return m.store.List()
}

// AddConnection saves a new connection
func (m *Manager) AddConnection(conn models.Connection) error {
	return m.store.Add(conn)
}

// UpdateConnection updates an existing connection
func (m *Manager) UpdateConnection(conn models.Connection) error {
	return m.store.Update(conn)
}

// DeleteConnection removes a connection and disconnects if active
func (m *Manager) DeleteConnection(id string) error {
	m.Disconnect(id)
	return m.store.Delete(id)
}

// TestConnection tests if a PostgreSQL connection can be established
func (m *Manager) TestConnection(conn models.Connection) (models.ConnectionTestResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dsn := buildDSN(conn)
	start := time.Now()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return models.ConnectionTestResult{
			OK:      false,
			Message: fmt.Sprintf("Failed to create pool: %v", err),
		}, nil
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		return models.ConnectionTestResult{
			OK:      false,
			Latency: latency,
			Message: fmt.Sprintf("Ping failed: %v", err),
		}, nil
	}

	return models.ConnectionTestResult{
		OK:      true,
		Latency: latency,
		Message: "Connection successful",
	}, nil
}

// Connect opens a connection pool for the given connection ID
func (m *Manager) Connect(id string) error {
	m.mu.RLock()
	if _, exists := m.pools[id]; exists {
		m.mu.RUnlock()
		return nil
	}
	m.mu.RUnlock()

	conn, err := m.store.Get(id)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	dsn := buildDSN(*conn)
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// Set pool limits
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 1
	poolConfig.MaxConnLifetime = 30 * time.Minute
	poolConfig.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return fmt.Errorf("failed to ping: %w", err)
	}

	m.mu.Lock()
	if existingPool, exists := m.pools[id]; exists {
		m.mu.Unlock()
		pool.Close()
		_ = existingPool
		return nil
	}
	m.pools[id] = pool
	m.configs[id] = poolConfig
	m.mu.Unlock()

	return nil
}

// ConnectAll opens connection pools for all saved connections concurrently
func (m *Manager) ConnectAll() map[string]error {
	conns := m.store.List()
	results := make(map[string]error)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, c := range conns {
		if m.IsConnected(c.ID) {
			continue
		}
		wg.Add(1)
		go func(connID string) {
			defer wg.Done()
			err := m.Connect(connID)
			mu.Lock()
			results[connID] = err
			mu.Unlock()
		}(c.ID)
	}
	wg.Wait()
	return results
}

// Disconnect closes a connection pool
func (m *Manager) Disconnect(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	pool, exists := m.pools[id]
	if !exists {
		return nil
	}

	pool.Close()
	delete(m.pools, id)
	delete(m.configs, id)

	return nil
}

// GetPool returns the connection pool for a given connection ID
func (m *Manager) GetPool(id string) (*pgxpool.Pool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	pool, exists := m.pools[id]
	if !exists {
		return nil, ErrConnectionClosed
	}

	return pool, nil
}

// GetConn returns the connection config for a given connection ID
func (m *Manager) GetConn(id string) (*models.Connection, error) {
	return m.store.Get(id)
}

// GetActiveConnections returns IDs of all active connections
func (m *Manager) GetActiveConnections() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ids := make([]string, 0, len(m.pools))
	for id := range m.pools {
		ids = append(ids, id)
	}

	return ids
}

// IsConnected checks if a connection is active
func (m *Manager) IsConnected(id string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.pools[id]
	return exists
}

// GetDatabaseName returns the database name for an active connection
func (m *Manager) GetDatabaseName(connID string) string {
	conn, err := m.store.Get(connID)
	if err != nil {
		return ""
	}
	return conn.Database
}

// GetConnInfo returns connection info (host, user, db) for display
func (m *Manager) GetConnInfo(connID string) (host, user, db string, ok bool) {
	conn, err := m.store.Get(connID)
	if err != nil {
		return "", "", "", false
	}
	return conn.Host, conn.Username, conn.Database, true
}

// Execute executes a query on a specific connection
func (m *Manager) Execute(ctx context.Context, connID string, sql string) (pgx.Rows, error) {
	pool, err := m.GetPool(connID)
	if err != nil {
		return nil, err
	}
	return pool.Query(ctx, sql)
}

// Exec executes a query that doesn't return rows
func (m *Manager) Exec(ctx context.Context, connID string, sql string) error {
	pool, err := m.GetPool(connID)
	if err != nil {
		return err
	}
	_, err = pool.Exec(ctx, sql)
	return err
}

// Cleanup closes all pools
func (m *Manager) Cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, pool := range m.pools {
		pool.Close()
		delete(m.pools, id)
	}
}
