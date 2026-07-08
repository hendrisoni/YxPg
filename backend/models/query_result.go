package models

// QueryResult holds the result of a database query
type QueryResult struct {
	Columns      []ColumnMeta    `json:"columns"`
	Rows         [][]interface{} `json:"rows"`
	RowCount     int             `json:"row_count"`
	TotalCount   int64           `json:"total_count"`
	Duration     int64           `json:"duration_ms"`
	Error        string          `json:"error,omitempty"`
	QueryType    string          `json:"query_type"` // select, insert, update, delete, ddl
	RowsAffected int64           `json:"rows_affected,omitempty"`
	RawSQL       string          `json:"raw_sql,omitempty"`
}

// ColumnMeta holds metadata about a result column
type ColumnMeta struct {
	Name     string `json:"name"`
	DataType string `json:"data_type"`
	TableOID int    `json:"table_oid,omitempty"`
}

// BrowseOptions holds options for browsing a table
type BrowseOptions struct {
	Page      int               `json:"page"`
	PageSize  int               `json:"page_size"`
	SortBy    string            `json:"sort_by,omitempty"`
	SortOrder string            `json:"sort_order,omitempty"` // asc, desc
	Filters   []FilterCondition `json:"filters,omitempty"`
	Columns   []string          `json:"columns,omitempty"` // specific columns to select
}

// FilterCondition holds a filter for browsing
type FilterCondition struct {
	Column   string `json:"column"`
	Operator string `json:"operator"` // =, !=, >, <, >=, <=, LIKE, ILIKE, IS NULL, IS NOT NULL, IN
	Value    string `json:"value"`
}

// ExplainResult holds the result of EXPLAIN ANALYZE
type ExplainResult struct {
	RawText string       `json:"raw_text"`
	Plan    *ExplainNode `json:"plan,omitempty"`
}

// ExplainNode represents a node in the EXPLAIN tree
type ExplainNode struct {
	NodeType            string         `json:"node_type"`
	RelationName        string         `json:"relation_name,omitempty"`
	Alias               string         `json:"alias,omitempty"`
	_startupCost        float64        `json:"startup_cost"`
	TotalCost           float64        `json:"total_cost"`
	PlanRows            int            `json:"plan_rows"`
	PlanWidth           int            `json:"plan_width"`
	ActualTime          float64        `json:"actual_time,omitempty"`
	ActualRows          int            `json:"actual_rows,omitempty"`
	ActualLoops         int            `json:"actual_loops,omitempty"`
	Filter              string         `json:"filter,omitempty"`
	RowsRemovedByFilter int            `json:"rows_removed_by_filter,omitempty"`
	SharedHitBlocks     int            `json:"shared_hit_blocks,omitempty"`
	SharedReadBlocks    int            `json:"shared_read_blocks,omitempty"`
	Plans               []*ExplainNode `json:"plans,omitempty"`
}
