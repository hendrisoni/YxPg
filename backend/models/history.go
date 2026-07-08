package models

import "time"

// QueryHistoryEntry represents a saved query history entry
type QueryHistoryEntry struct {
	ID           int64     `json:"id"`
	ConnectionID string    `json:"connection_id"`
	Database     string    `json:"database"`
	SQL          string    `json:"sql"`
	DurationMs   int64     `json:"duration_ms"`
	RowsReturned int       `json:"rows_returned"`
	ExecutedAt   time.Time `json:"executed_at"`
	Error        string    `json:"error,omitempty"`
	Bookmarked   bool      `json:"bookmarked"`
}

// SavedQuery represents a user-saved query
type SavedQuery struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	SQL       string    `json:"sql"`
	Folder    string    `json:"folder,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ExportOpts holds options for data export
type ExportOpts struct {
	Format    string `json:"format"`               // csv, json, sql
	Delimiter string `json:"delimiter,omitempty"`  // for CSV
	TableOnly bool   `json:"table_only,omitempty"` // for SQL: INSERT only (no CREATE)
}

// TableDefinition represents a full table definition for DDL
type TableDefinition struct {
	Schema      string             `json:"schema"`
	TableName   string             `json:"table_name"`
	Columns     []ColumnDefinition `json:"columns"`
	Indexes     []IndexDefinition  `json:"indexes,omitempty"`
	ForeignKeys []FKDefinition     `json:"foreign_keys,omitempty"`
}

// ColumnDefinition represents a column in DDL context
type ColumnDefinition struct {
	Name            string       `json:"name"`
	DataType        string       `json:"data_type"`
	Length          *int         `json:"length,omitempty"`
	Precision       *int         `json:"precision,omitempty"`
	Scale           *int         `json:"scale,omitempty"`
	IsNullable      bool         `json:"is_nullable"`
	DefaultValue    string       `json:"default_value,omitempty"`
	IsPrimaryKey    bool         `json:"is_primary_key"`
	IsUnique        bool         `json:"is_unique"`
	IsAutoIncrement bool         `json:"is_auto_increment"`
	References      *FKReference `json:"references,omitempty"`
	Comment         string       `json:"comment,omitempty"`
}

// IndexDefinition represents an index definition for DDL
type IndexDefinition struct {
	Name      string   `json:"name"`
	TableName string   `json:"table_name"`
	Columns   []string `json:"columns"`
	IsUnique  bool     `json:"is_unique"`
	IndexType string   `json:"index_type"`      // btree, hash, gin, gist, brin
	Where     string   `json:"where,omitempty"` // partial index condition
}

// FKDefinition represents a foreign key definition for DDL
type FKDefinition struct {
	Name          string   `json:"name"`
	SourceColumns []string `json:"source_columns"`
	TargetSchema  string   `json:"target_schema"`
	TargetTable   string   `json:"target_table"`
	TargetColumns []string `json:"target_columns"`
	OnUpdate      string   `json:"on_update"`
	OnDelete      string   `json:"on_delete"`
}

// FKReference represents a column-level foreign key reference
type FKReference struct {
	TargetSchema string `json:"target_schema"`
	TargetTable  string `json:"target_table"`
	TargetColumn string `json:"target_column"`
	OnUpdate     string `json:"on_update"`
	OnDelete     string `json:"on_delete"`
}
