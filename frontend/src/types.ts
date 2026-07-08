// Wails runtime bindings type definitions
// These types match the Go backend models

export interface Connection {
  id: string
  name: string
  host: string
  port: number
  database: string
  username: string
  password: string
  ssl_mode: string
  color: string
  created_at: string
}

export interface ConnectionTestResult {
  ok: boolean
  latency_ms: number
  message: string
}

export interface SchemaInfo {
  name: string
  owner: string
}

export interface TableInfo {
  schema: string
  name: string
  type: string // table, view, materialized_view
  row_count: number
  comment?: string
}

export interface ColumnInfo {
  table_schema: string
  table_name: string
  column_name: string
  data_type: string
  is_nullable: boolean
  default_value?: any
  is_primary_key: boolean
  ordinal_position: number
  character_maximum_length?: number
  comment?: string
}

export interface IndexInfo {
  index_name: string
  table_name: string
  columns: string[]
  is_unique: boolean
  index_type: string
  definition: string
}

export interface FKInfo {
  constraint_name: string
  source_schema: string
  source_table: string
  source_column: string
  target_schema: string
  target_table: string
  target_column: string
  on_update: string
  on_delete: string
}

export interface ViewInfo {
  schema: string
  name: string
  read_only: boolean
  definition: string
}

export interface FunctionInfo {
  schema: string
  name: string
  arguments: string
  return_type: string
  language: string
}

export interface SequenceInfo {
  schema: string
  name: string
  data_type: string
  min_value: number
  max_value: number
  increment: number
  cache_size: number
}

export interface TriggerInfo {
  schema: string
  table_name: string
  trigger_name: string
  event: string
  timing: string
  definition: string
}

export interface TypeInfo {
  schema: string
  name: string
  type: string
  values?: string[]
}

export interface QueryResult {
  columns: ColumnMeta[]
  rows: any[][]
  row_count: number
  total_count: number
  duration_ms: number
  error?: string
  query_type: string
  rows_affected?: number
  raw_sql?: string
}

export interface ColumnMeta {
  name: string
  data_type: string
  table_oid?: number
}

export interface BrowseOptions {
  page: number
  page_size: number
  sort_by?: string
  sort_order?: string
  filters?: FilterCondition[]
  columns?: string[]
}

export interface FilterCondition {
  column: string
  operator: string
  value: string
}

export interface ExplainResult {
  raw_text: string
  plan?: ExplainNode
}

export interface ExplainNode {
  node_type: string
  relation_name?: string
  alias?: string
  startup_cost: number
  total_cost: number
  plan_rows: number
  plan_width: number
  actual_time?: number
  actual_rows?: number
  actual_loops?: number
  filter?: string
  rows_removed_by_filter?: number
  shared_hit_blocks?: number
  shared_read_blocks?: number
  plans?: ExplainNode[]
}

export interface QueryHistoryEntry {
  id: number
  connection_id: string
  database: string
  sql: string
  duration_ms: number
  rows_returned: number
  executed_at: string
  error?: string
  bookmarked: boolean
}

export interface SavedQuery {
  id: number
  name: string
  sql: string
  folder?: string
  created_at: string
  updated_at: string
}

export interface TableDefinition {
  schema: string
  table_name: string
  columns: ColumnDefinition[]
  indexes?: IndexDefinition[]
  foreign_keys?: FKDefinition[]
}

export interface ColumnDefinition {
  name: string
  data_type: string
  length?: number
  precision?: number
  scale?: number
  is_nullable: boolean
  default_value?: string
  is_primary_key: boolean
  is_unique: boolean
  is_auto_increment: boolean
  references?: FKReference
  comment?: string
}

export interface IndexDefinition {
  name: string
  table_name: string
  columns: string[]
  is_unique: boolean
  index_type: string
  where?: string
}

export interface FKDefinition {
  name: string
  source_columns: string[]
  target_schema: string
  target_table: string
  target_columns: string[]
  on_update: string
  on_delete: string
}

export interface FKReference {
  target_schema: string
  target_table: string
  target_column: string
  on_update: string
  on_delete: string
}

// UI Types
export interface Tab {
  id: string
  title: string
  type: 'query' | 'table' | 'builder' | 'ddl' | 'home' | 'log'
  connectionId?: string
  schema?: string
  table?: string
  sql?: string
  modified?: boolean
  data?: any
}

export interface TreeNode {
  id: string
  label: string
  type: 'connection' | 'schema' | 'table' | 'view' | 'function' | 'sequence' | 'column' | 'index' | 'category'
  icon?: string
  children?: TreeNode[]
  data?: any
  expanded?: boolean
}

export interface Notification {
  id: string
  type: 'success' | 'error' | 'warning' | 'info'
  title: string
  message: string
  duration?: number
}
