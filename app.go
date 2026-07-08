package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"yxpg/backend/connection"
	"yxpg/backend/ddl"
	"yxpg/backend/export"
	"yxpg/backend/models"
	"yxpg/backend/query"
	"yxpg/backend/schema"

	"github.com/jackc/pgx/v5"
	_ "github.com/mattn/go-sqlite3"
)

// App struct
type App struct {
	ctx            context.Context
	manager        *connection.Manager
	inspector      *schema.Inspector
	executor       *query.Executor
	history        *query.History
	ddlExec        *ddl.Executor
	builder        *ddl.Builder
	workspaceStore *connection.WorkspaceStore
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize connection manager
	var err error
	a.manager, err = connection.NewManager()
	if err != nil {
		fmt.Printf("Failed to create connection manager: %v\n", err)
		return
	}

	// Initialize schema inspector
	a.inspector = schema.NewInspector(a.manager)

	// Initialize query history (optional - gracefully skip if CGO unavailable)
	a.history, err = query.NewHistory()
	if err != nil {
		fmt.Printf("Warning: Query history disabled: %v\n", err)
		a.history = nil
	}

	// Initialize query executor
	a.executor = query.NewExecutor(a.manager, a.history)

	// Initialize DDL executor
	a.ddlExec = ddl.NewExecutor(a.manager)
	a.builder = ddl.NewBuilder()

	// Initialize workspace store
	a.workspaceStore, err = connection.NewWorkspaceStore()
	if err != nil {
		fmt.Printf("Failed to create workspace store: %v\n", err)
	}

	// Sync and connect on startup in background
	go func() {
		time.Sleep(500 * time.Millisecond)
		_, err := a.SyncServerConnections()
		if err != nil {
			fmt.Printf("[Startup] SyncServerConnections failed: %v. Connecting saved local connections.\n", err)
		}
		a.manager.ConnectAll()
	}()
}

// shutdown is called when the app stops
func (a *App) shutdown(ctx context.Context) {
	if a.manager != nil {
		a.manager.Cleanup()
	}
	if a.history != nil {
		a.history.Close()
	}
}

// ==================== CONNECTION METHODS ====================

// AddConnection saves a new connection
func (a *App) AddConnection(conn models.Connection) error {
	return a.manager.AddConnection(conn)
}

// UpdateConnection updates an existing connection
func (a *App) UpdateConnection(conn models.Connection) error {
	return a.manager.UpdateConnection(conn)
}

// TestConnection tests if a connection can be established
func (a *App) TestConnection(conn models.Connection) (models.ConnectionTestResult, error) {
	result, err := a.manager.TestConnection(conn)
	return result, err
}

// ListConnections returns all saved connections
func (a *App) ListConnections() []models.Connection {
	return a.manager.ListConnections()
}

// DeleteConnection removes a connection
func (a *App) DeleteConnection(id string) error {
	return a.manager.DeleteConnection(id)
}

// Connect opens a connection pool
func (a *App) Connect(id string) error {
	return a.manager.Connect(id)
}

// Disconnect closes a connection pool
func (a *App) Disconnect(id string) error {
	return a.manager.Disconnect(id)
}

// GetActiveConnections returns IDs of active connections
func (a *App) GetActiveConnections() []string {
	return a.manager.GetActiveConnections()
}

// ==================== SCHEMA METHODS ====================

// GetSchemas returns all schemas
func (a *App) GetSchemas(connID string) ([]models.SchemaInfo, error) {
	return a.inspector.GetSchemas(a.ctx, connID)
}

// GetTables returns tables in a schema
func (a *App) GetTables(connID, schemaName string) ([]models.TableInfo, error) {
	return a.inspector.GetTables(a.ctx, connID, schemaName)
}

// GetColumns returns columns for a table
func (a *App) GetColumns(connID, schemaName, table string) ([]models.ColumnInfo, error) {
	return a.inspector.GetColumns(a.ctx, connID, schemaName, table)
}

// GetIndexes returns indexes for a table
func (a *App) GetIndexes(connID, schemaName, table string) ([]models.IndexInfo, error) {
	return a.inspector.GetIndexes(a.ctx, connID, schemaName, table)
}

// GetForeignKeys returns foreign keys for a table
func (a *App) GetForeignKeys(connID, schemaName, table string) ([]models.FKInfo, error) {
	return a.inspector.GetForeignKeys(a.ctx, connID, schemaName, table)
}

// GetViews returns views in a schema
func (a *App) GetViews(connID, schemaName string) ([]models.ViewInfo, error) {
	return a.inspector.GetViews(a.ctx, connID, schemaName)
}

// GetFunctions returns functions in a schema
func (a *App) GetFunctions(connID, schemaName string) ([]models.FunctionInfo, error) {
	return a.inspector.GetFunctions(a.ctx, connID, schemaName)
}

// GetSequences returns sequences in a schema
func (a *App) GetSequences(connID, schemaName string) ([]models.SequenceInfo, error) {
	return a.inspector.GetSequences(a.ctx, connID, schemaName)
}

// GetTriggers returns triggers for a table
func (a *App) GetTriggers(connID, schemaName, table string) ([]models.TriggerInfo, error) {
	return a.inspector.GetTriggers(a.ctx, connID, schemaName, table)
}

// GetTypes returns custom types in a schema
func (a *App) GetTypes(connID, schemaName string) ([]models.TypeInfo, error) {
	return a.inspector.GetTypes(a.ctx, connID, schemaName)
}

// GetTableDDL returns DDL for a table
func (a *App) GetTableDDL(connID, schemaName, table string) (string, error) {
	return a.inspector.GetTableDDL(a.ctx, connID, schemaName, table)
}

// RefreshSchema refreshes schema cache
func (a *App) RefreshSchema(connID string) error {
	return a.inspector.RefreshSchema(a.ctx, connID)
}

// GetFullSchema returns all schema objects
func (a *App) GetFullSchema(connID string) (map[string]interface{}, error) {
	return a.inspector.GetFullSchema(a.ctx, connID)
}

// ==================== QUERY METHODS ====================

// ExecuteQuery executes a SQL query
func (a *App) ExecuteQuery(connID, sql string, timeout int) models.QueryResult {
	return a.executor.ExecuteWithConnection(a.ctx, connID, sql, timeout)
}

// CancelQuery cancels an active query
func (a *App) CancelQuery(connID string) error {
	return a.executor.CancelQuery(connID)
}

// ExplainQuery runs EXPLAIN ANALYZE
func (a *App) ExplainQuery(connID, sql string) models.ExplainResult {
	return a.executor.ExplainQuery(a.ctx, connID, sql)
}

// GetQueryHistory returns query history
func (a *App) GetQueryHistory(connID string, limit int) ([]models.QueryHistoryEntry, error) {
	return a.executor.GetHistory(connID, limit)
}

// SaveQuery saves a named query
func (a *App) SaveQuery(name, sql string) error {
	return a.executor.SaveQuery(name, sql)
}

// GetSavedQueries returns all saved queries
func (a *App) GetSavedQueries() ([]models.SavedQuery, error) {
	return a.executor.GetSavedQueries()
}

// BrowseTable retrieves paginated data from a table
func (a *App) BrowseTable(connID, schemaName, table string, opts models.BrowseOptions) models.QueryResult {
	return query.BrowseTable(a.ctx, a.manager, connID, schemaName, table, opts)
}

// ==================== DDL METHODS ====================

// CreateTable creates a new table
func (a *App) CreateTable(connID string, def models.TableDefinition) (string, error) {
	return a.ddlExec.CreateTable(a.ctx, connID, def)
}

// AlterTable alters an existing table
func (a *App) AlterTable(connID string, original, modified models.TableDefinition) ([]string, error) {
	return a.ddlExec.AlterTable(a.ctx, connID, original, modified)
}

// DropTable drops a table
func (a *App) DropTable(connID, schemaName, table string, cascade bool) error {
	return a.ddlExec.DropTable(a.ctx, connID, schemaName, table, cascade)
}

// RenameTable renames a table
func (a *App) RenameTable(connID, schemaName, oldName, newName string) error {
	return a.ddlExec.RenameTable(a.ctx, connID, schemaName, oldName, newName)
}

// CreateIndex creates a new index
func (a *App) CreateIndex(connID string, def models.IndexDefinition) error {
	return a.ddlExec.CreateIndex(a.ctx, connID, def)
}

// DropIndex drops an index
func (a *App) DropIndex(connID, indexName string) error {
	return a.ddlExec.DropIndex(a.ctx, connID, indexName)
}

// AddForeignKey adds a foreign key
func (a *App) AddForeignKey(connID string, def models.FKDefinition) error {
	return a.ddlExec.AddForeignKey(a.ctx, connID, def)
}

// DropConstraint drops a constraint
func (a *App) DropConstraint(connID, schemaName, table, constraint string) error {
	return a.ddlExec.DropConstraint(a.ctx, connID, schemaName, table, constraint)
}

// ExecuteRawDDL executes a raw DDL statement
func (a *App) ExecuteRawDDL(connID, sql string) error {
	return a.ddlExec.ExecuteRaw(a.ctx, connID, sql)
}

// ==================== EXPORT METHODS ====================

// ExportData exports query result data
func (a *App) ExportData(result models.QueryResult, format string, schemaName, table string) (string, error) {
	switch format {
	case "csv":
		return export.ExportCSV(result, ",")
	case "json":
		return export.ExportJSON(result)
	case "sql":
		return export.ExportSQL(result, schemaName, table)
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

// ==================== UTILITY METHODS ====================

// GetConnInfo returns connection info for display
func (a *App) GetConnInfo(connID string) (host, user, db string, ok bool) {
	return a.manager.GetConnInfo(connID)
}

// GetDatabaseName returns database name for a connection
func (a *App) GetDatabaseName(connID string) string {
	return a.manager.GetDatabaseName(connID)
}

// IsConnected checks if a connection is active
func (a *App) IsConnected(connID string) bool {
	return a.manager.IsConnected(connID)
}

// GetTimestamp returns current timestamp
func (a *App) GetTimestamp() string {
	return time.Now().Format("15:04:05")
}

// loadConfigMap searches for and loads configuration from yxpg.conf.
// It checks current working directory, executable directory, and home directory config folder.
// If none exists, it writes a default one to the home folder.
func (a *App) loadConfigMap() map[string]string {
	filename := "yxpg.conf"
	config := make(map[string]string)

	// Paths to check
	var paths []string

	// 1. Current working directory
	if cwd, err := os.Getwd(); err == nil {
		paths = append(paths, filepath.Join(cwd, filename))
	}

	// 2. Executable directory
	if exePath, err := os.Executable(); err == nil {
		paths = append(paths, filepath.Join(filepath.Dir(exePath), filename))
	}

	// 3. User home directory ~/.yxpg/yxpg.conf
	if homeDir, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(homeDir, ".yxpg", filename))
	}

	var foundPath string
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			foundPath = p
			break
		}
	}

	// If not found, create default config in user home ~/.yxpg/yxpg.conf
	if foundPath == "" {
		if homeDir, err := os.UserHomeDir(); err == nil {
			configDir := filepath.Join(homeDir, ".yxpg")
			_ = os.MkdirAll(configDir, 0755)
			defaultPath := filepath.Join(configDir, filename)
			
			content := `# YxPg Configuration File

# ==========================================
# Central Database Connection Settings
# Configure connection to the central database
# ==========================================
host=localhost
port=5432
user=postgres
password=
dbname=postgres
sslmode=disable

# ==========================================
# Default Settings for New Connections
# Configure default parameters for the New Connection dialog
# ==========================================
default_host=localhost
default_port=5432
default_username=postgres
default_password=
default_database=postgres
`
			if err := os.WriteFile(defaultPath, []byte(content), 0644); err == nil {
				foundPath = defaultPath
				fmt.Printf("[Config] Created default config file at: %s\n", defaultPath)
			}
		}
	}

	if foundPath == "" {
		return config
	}

	file, err := os.Open(foundPath)
	if err != nil {
		return config
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines or comments
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		config[key] = val
	}

	return config
}

// loadDSNFromConfig loads the DSN connection string from yxpg.conf.
func (a *App) loadDSNFromConfig() string {
	defaultDSN := "host=localhost port=5432 user=postgres password= dbname=postgres sslmode=disable"
	config := a.loadConfigMap()

	if len(config) == 0 {
		return defaultDSN
	}

	// If user defined a full DSN key, use it
	if dsnVal, ok := config["dsn"]; ok && dsnVal != "" {
		return dsnVal
	}

	// Build connection string from individual params if available
	host := config["host"]
	port := config["port"]
	user := config["user"]
	password := config["password"]
	dbname := config["dbname"]
	sslmode := config["sslmode"]

	// Fallback to defaults for missing fields
	if host == "" { host = "localhost" }
	if port == "" { port = "5432" }
	if user == "" { user = "postgres" }
	if password == "" { password = "" }
	if dbname == "" { dbname = "postgres" }
	if sslmode == "" { sslmode = "disable" }

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", 
		host, port, user, password, dbname, sslmode)
}

// DefaultConnectionConfig represents the default parameters for a new connection
type DefaultConnectionConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// GetDefaultConnectionConfig returns the configured default settings for new connections from yxpg.conf
func (a *App) GetDefaultConnectionConfig() DefaultConnectionConfig {
	config := a.loadConfigMap()
	
	// Parse port
	port := 5432
	if portStr, ok := config["default_port"]; ok && portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}
	
	dHost := "localhost"
	if h, ok := config["default_host"]; ok && h != "" {
		dHost = h
	}
	
	dUser := "postgres"
	if u, ok := config["default_username"]; ok && u != "" {
		dUser = u
	}
	
	dPass := ""
	if p, ok := config["default_password"]; ok {
		dPass = p
	}
	
	dDB := "postgres"
	if d, ok := config["default_database"]; ok && d != "" {
		dDB = d
	}
	
	return DefaultConnectionConfig{
		Host:     dHost,
		Port:     port,
		Username: dUser,
		Password: dPass,
		Database: dDB,
	}
}

// SyncServerConnections connects to the yxz database and fetches servers from public._server table
func (a *App) SyncServerConnections() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	dsn := a.loadDSNFromConfig()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to yxz db: %w", err)
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, `
		SELECT
			pkserver,
			server_key,
			host,
			port,
			username,
			password,
			db_name,
			is_active,
			created_at,
			_default
		FROM
			public._server
	`)
	if err != nil {
		return 0, fmt.Errorf("failed to query _server table: %w", err)
	}
	defer rows.Close()

	var syncedCount int
	existingConns := a.manager.ListConnections()
	existingMap := make(map[string]models.Connection)
	for _, c := range existingConns {
		existingMap[c.ID] = c
	}

	for rows.Next() {
		var pkserver interface{}
		var serverKey string
		var host string
		var port int
		var username string
		var password string
		var dbName string
		var isActive bool
		var createdAt time.Time
		var isDefault interface{}

		err := rows.Scan(
			&pkserver,
			&serverKey,
			&host,
			&port,
			&username,
			&password,
			&dbName,
			&isActive,
			&createdAt,
			&isDefault,
		)
		if err != nil {
			return syncedCount, fmt.Errorf("failed to scan row: %w", err)
		}

		if !isActive {
			continue
		}

		connID := serverKey
		if connID == "" {
			connID = fmt.Sprintf("%v", pkserver)
		}

		newConn := models.Connection{
			ID:        connID,
			Name:      serverKey,
			Host:      host,
			Port:      port,
			Database:  dbName,
			Username:  username,
			Password:  password,
			SSLMode:   "disable",
			Color:     getRandomColor(),
			CreatedAt: createdAt,
		}

		// Save to store
		if _, exists := existingMap[connID]; exists {
			err = a.manager.UpdateConnection(newConn)
		} else {
			err = a.manager.AddConnection(newConn)
		}

		if err != nil {
			return syncedCount, fmt.Errorf("failed to save connection %s: %w", serverKey, err)
		}
		syncedCount++
	}

	// Connect all databases concurrently
	a.manager.ConnectAll()

	return syncedCount, nil
}

// SyncPgAdminConnections imports server configs from pgAdmin 4 (pgadmin4.db SQLite database)
func (a *App) SyncPgAdminConnections() (int, error) {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		return 0, fmt.Errorf("APPDATA environment variable not set")
	}
	dbPath := filepath.Join(appData, "pgAdmin", "pgadmin4.db")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return 0, fmt.Errorf("pgAdmin 4 database not found at: %s", dbPath)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return 0, fmt.Errorf("failed to open pgadmin database: %w", err)
	}
	defer db.Close()

	// Query the server table from pgAdmin 4 SQLite db
	rows, err := db.Query("SELECT name, host, port, maintenance_db, username FROM server")
	if err != nil {
		return 0, fmt.Errorf("failed to query servers from pgadmin database: %w", err)
	}
	defer rows.Close()

	var syncedCount int
	existingConns := a.manager.ListConnections()
	existingMap := make(map[string]models.Connection)
	for _, c := range existingConns {
		existingMap[c.ID] = c
	}

	for rows.Next() {
		var name sql.NullString
		var host sql.NullString
		var port sql.NullInt64
		var dbName sql.NullString
		var username sql.NullString

		err := rows.Scan(&name, &host, &port, &dbName, &username)
		if err != nil {
			continue // Skip unreadable rows
		}

		if !host.Valid || host.String == "" {
			continue
		}

		p := 5432
		if port.Valid {
			p = int(port.Int64)
		}

		dName := "postgres"
		if dbName.Valid && dbName.String != "" {
			dName = dbName.String
		}

		uName := "postgres"
		if username.Valid && username.String != "" {
			uName = username.String
		}

		connName := name.String
		if connName == "" {
			connName = host.String
		}

		// Generate a unique ID for this connection
		connID := fmt.Sprintf("pgadmin_%s_%s_%d_%s", connName, host.String, p, dName)

		newConn := models.Connection{
			ID:        connID,
			Name:      connName,
			Host:      host.String,
			Port:      p,
			Database:  dName,
			Username:  uName,
			Password:  "",
			SSLMode:   "disable",
			Color:     getRandomColor(),
			CreatedAt: time.Now(),
		}

		// If a connection with this ID or endpoints already exists, skip it
		exists := false
		for _, c := range existingMap {
			if c.ID == connID || (c.Host == newConn.Host && c.Port == newConn.Port && c.Database == newConn.Database && c.Username == newConn.Username) {
				exists = true
				break
			}
		}

		if !exists {
			err = a.manager.AddConnection(newConn)
			if err == nil {
				syncedCount++
			}
		}
	}

	return syncedCount, nil
}

// CatalogItem represents a searchable database object
type CatalogItem struct {
	ConnectionID   string `json:"connection_id"`
	ConnectionName string `json:"connection_name"`
	DatabaseName   string `json:"database_name"`
	Schema         string `json:"schema"`
	Name           string `json:"name"`
	Type           string `json:"type"` // "table", "view", "function"
}

// GetSearchCatalog returns all tables, views, and functions for all active connections
func (a *App) GetSearchCatalog() ([]CatalogItem, error) {
	activeIDs := a.manager.GetActiveConnections()
	fmt.Printf("[GetSearchCatalog] activeIDs: %v\n", activeIDs)
	var result []CatalogItem

	for _, connID := range activeIDs {
		pool, err := a.manager.GetPool(connID)
		if err != nil {
			fmt.Printf("[GetSearchCatalog] GetPool error for %s: %v\n", connID, err)
			continue
		}

		_, _, dbName, ok := a.manager.GetConnInfo(connID)
		if !ok {
			fmt.Printf("[GetSearchCatalog] GetConnInfo failed for %s\n", connID)
			continue
		}

		connConfig, err := a.manager.GetConn(connID)
		connName := connID
		if err == nil && connConfig != nil {
			connName = connConfig.Name
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		query := `
			SELECT table_schema AS schema_name, table_name AS item_name, CASE WHEN table_type = 'VIEW' THEN 'view' ELSE 'table' END AS item_type
			FROM information_schema.tables
			WHERE table_schema NOT IN ('pg_catalog', 'information_schema', 'pg_toast')
			UNION ALL
			SELECT n.nspname AS schema_name, p.proname AS item_name, 'function' AS item_type
			FROM pg_proc p
			JOIN pg_namespace n ON p.pronamespace = n.oid
			WHERE n.nspname NOT IN ('pg_catalog', 'information_schema', 'pg_toast')
			ORDER BY schema_name, item_name;
		`

		rows, err := pool.Query(ctx, query)
		if err != nil {
			fmt.Printf("[GetSearchCatalog] Query error for %s: %v\n", connID, err)
			cancel()
			continue
		}

		count := 0
		for rows.Next() {
			var item CatalogItem
			item.ConnectionID = connID
			item.ConnectionName = connName
			item.DatabaseName = dbName
			if err := rows.Scan(&item.Schema, &item.Name, &item.Type); err == nil {
				result = append(result, item)
				count++
			} else {
				fmt.Printf("[GetSearchCatalog] Scan error: %v\n", err)
			}
		}
		rows.Close()
		cancel()
		fmt.Printf("[GetSearchCatalog] Found %d items for connection %s (%s)\n", count, connID, connName)
	}

	if result == nil {
		result = []CatalogItem{}
	}

	fmt.Printf("[GetSearchCatalog] Total items: %d\n", len(result))
	return result, nil
}

// LoadWorkspace returns the saved workspace tree JSON
func (a *App) LoadWorkspace() (string, error) {
	if a.workspaceStore == nil {
		return "[]", fmt.Errorf("workspace store not initialized")
	}
	err := a.workspaceStore.Load()
	return a.workspaceStore.Data, err
}

// SaveWorkspace saves the workspace tree JSON
func (a *App) SaveWorkspace(data string) error {
	if a.workspaceStore == nil {
		return fmt.Errorf("workspace store not initialized")
	}
	return a.workspaceStore.Save(data)
}

func getRandomColor() string {
	colors := []string{
		"#00C9A7", // Teal
		"#3B82F6", // Blue
		"#F59E0B", // Amber
		"#EF4444", // Red
		"#10B981", // Green
		"#8B5CF6", // Purple
		"#EC4899", // Pink
		"#06B6D4", // Cyan
	}
	index := time.Now().UnixNano() % int64(len(colors))
	return colors[index]
}


