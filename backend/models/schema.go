package models

// SchemaInfo represents a database schema
type SchemaInfo struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

// TableInfo represents a table in the database
type TableInfo struct {
	Schema   string `json:"schema"`
	Name     string `json:"name"`
	Type     string `json:"type"` // table, view, materialized_view
	RowCount int64  `json:"row_count"`
	Comment  string `json:"comment,omitempty"`
}

// ColumnInfo represents a column in a table
type ColumnInfo struct {
	TableSchema            string      `json:"table_schema"`
	TableName              string      `json:"table_name"`
	ColumnName             string      `json:"column_name"`
	DataType               string      `json:"data_type"`
	IsNullable             bool        `json:"is_nullable"`
	DefaultValue           interface{} `json:"default_value,omitempty"`
	IsPrimaryKey           bool        `json:"is_primary_key"`
	OrdinalPosition        int         `json:"ordinal_position"`
	CharacterMaximumLength *int        `json:"character_maximum_length,omitempty"`
	Comment                string      `json:"comment,omitempty"`
}

// IndexInfo represents an index on a table
type IndexInfo struct {
	IndexName  string   `json:"index_name"`
	TableName  string   `json:"table_name"`
	Columns    []string `json:"columns"`
	IsUnique   bool     `json:"is_unique"`
	IndexType  string   `json:"index_type"`
	Definition string   `json:"definition"`
}

// FKInfo represents a foreign key relationship
type FKInfo struct {
	ConstraintName string `json:"constraint_name"`
	SourceSchema   string `json:"source_schema"`
	SourceTable    string `json:"source_table"`
	SourceColumn   string `json:"source_column"`
	TargetSchema   string `json:"target_schema"`
	TargetTable    string `json:"target_table"`
	TargetColumn   string `json:"target_column"`
	OnUpdate       string `json:"on_update"`
	OnDelete       string `json:"on_delete"`
}

// ViewInfo represents a database view
type ViewInfo struct {
	Schema     string `json:"schema"`
	Name       string `json:"name"`
	ReadOnly   bool   `json:"read_only"`
	Definition string `json:"definition"`
}

// FunctionInfo represents a stored function/procedure
type FunctionInfo struct {
	Schema     string `json:"schema"`
	Name       string `json:"name"`
	Arguments  string `json:"arguments"`
	ReturnType string `json:"return_type"`
	Language   string `json:"language"`
}

// SequenceInfo represents a database sequence
type SequenceInfo struct {
	Schema    string `json:"schema"`
	Name      string `json:"name"`
	DataType  string `json:"data_type"`
	MinValue  int64  `json:"min_value"`
	MaxValue  int64  `json:"max_value"`
	Increment int    `json:"increment"`
	CacheSize int    `json:"cache_size"`
}

// TriggerInfo represents a database trigger
type TriggerInfo struct {
	Schema      string `json:"schema"`
	TableName   string `json:"table_name"`
	TriggerName string `json:"trigger_name"`
	Event       string `json:"event"`
	Timing      string `json:"timing"`
	Definition  string `json:"definition"`
}

// TypeInfo represents a custom/enum type
type TypeInfo struct {
	Schema string   `json:"schema"`
	Name   string   `json:"name"`
	Type   string   `json:"type"`             // enum, composite, range, etc.
	Values []string `json:"values,omitempty"` // for enums
}
