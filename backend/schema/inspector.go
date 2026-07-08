package schema

import (
	"context"
	"fmt"

	"yxpg/backend/connection"
	"yxpg/backend/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Inspector handles database schema introspection
type Inspector struct {
	manager *connection.Manager
}

// NewInspector creates a new schema inspector
func NewInspector(manager *connection.Manager) *Inspector {
	return &Inspector{manager: manager}
}

// GetSchemas returns all schemas in the database
func (i *Inspector) GetSchemas(ctx context.Context, connID string) ([]models.SchemaInfo, error) {
	pool, err := i.manager.GetPool(connID)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT schema_name, schema_owner
		FROM information_schema.schemata
		WHERE schema_name NOT IN ('pg_catalog', 'information_schema', 'pg_toast')
		ORDER BY schema_name
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query schemas: %w", err)
	}
	defer rows.Close()

	var schemas []models.SchemaInfo
	for rows.Next() {
		var s models.SchemaInfo
		if err := rows.Scan(&s.Name, &s.Owner); err != nil {
			return nil, err
		}
		schemas = append(schemas, s)
	}

	return schemas, rows.Err()
}

// GetTables returns all tables in a schema
func (i *Inspector) GetTables(ctx context.Context, connID, schema string) ([]models.TableInfo, error) {
	pool, err := i.manager.GetPool(connID)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT 
			t.table_schema,
			t.table_name,
			CASE 
			 WHEN t.table_type = 'VIEW' THEN 'view'
			 ELSE 'table'
			END AS type,
			COALESCE(pg_catalog.obj_description(
				(c.table_schema || '.' || c.table_name)::regclass
			), '') AS comment,
			COALESCE(
				(SELECT reltuples::bigint FROM pg_class 
				 WHERE oid = (c.table_schema || '.' || c.table_name)::regclass),
				0
			) AS row_count
		FROM information_schema.tables t
		LEFT JOIN information_schema.tables c ON t.table_schema = c.table_schema AND t.table_name = c.table_name
		WHERE t.table_schema = $1
		ORDER BY t.table_name
	`

	if schema == "" {
		schema = "public"
	}

	rows, err := pool.Query(ctx, query, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to query tables: %w", err)
	}
	defer rows.Close()

	var tables []models.TableInfo
	for rows.Next() {
		var t models.TableInfo
		if err := rows.Scan(&t.Schema, &t.Name, &t.Type, &t.Comment, &t.RowCount); err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}

	return tables, rows.Err()
}

// GetColumns returns all columns in a table
func (i *Inspector) GetColumns(ctx context.Context, connID, schema, table string) ([]models.ColumnInfo, error) {
	pool, err := i.manager.GetPool(connID)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT 
			c.table_schema,
			c.table_name,
			c.column_name,
			c.data_type,
			c.is_nullable = 'YES' AS is_nullable,
			c.column_default,
			CASE WHEN pk.column_name IS NOT NULL THEN true ELSE false END AS is_primary_key,
			c.ordinal_position,
			c.character_maximum_length,
			COALESCE(pgd.description, '') AS comment
		FROM information_schema.columns c
		LEFT JOIN (
			SELECT ku.column_name, ku.table_schema, ku.table_name
			FROM information_schema.table_constraints tc
			JOIN information_schema.key_column_usage ku
				ON tc.constraint_name = ku.constraint_name
				AND tc.table_schema = ku.table_schema
			WHERE tc.constraint_type = 'PRIMARY KEY'
		) pk ON c.column_name = pk.column_name 
			AND c.table_schema = pk.table_schema 
			AND c.table_name = pk.table_name
		LEFT JOIN pg_catalog.pg_statio_all_tables st
			ON c.table_schema = st.schemaname AND c.table_name = st.relname
		LEFT JOIN pg_catalog.pg_description pgd
			ON pgd.objoid = st.relid AND pgd.objsubid = c.ordinal_position
		WHERE c.table_schema = $1 AND c.table_name = $2
		ORDER BY c.ordinal_position
	`

	rows, err := pool.Query(ctx, query, schema, table)
	if err != nil {
		return nil, fmt.Errorf("failed to query columns: %w", err)
	}
	defer rows.Close()

	var columns []models.ColumnInfo
	for rows.Next() {
		var col models.ColumnInfo
		var defaultVal interface{}
		if err := rows.Scan(
			&col.TableSchema,
			&col.TableName,
			&col.ColumnName,
			&col.DataType,
			&col.IsNullable,
			&defaultVal,
			&col.IsPrimaryKey,
			&col.OrdinalPosition,
			&col.CharacterMaximumLength,
			&col.Comment,
		); err != nil {
			return nil, err
		}
		col.DefaultValue = defaultVal
		columns = append(columns, col)
	}

	return columns, rows.Err()
}

// GetIndexes returns all indexes for a table
func (i *Inspector) GetIndexes(ctx context.Context, connID, schema, table string) ([]models.IndexInfo, error) {
	pool, err := i.manager.GetPool(connID)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT
			i.relname AS index_name,
			t.relname AS table_name,
			ARRAY_AGG(a.attname ORDER BY array_position(ix.indkey, a.attnum)) AS columns,
			ix.indisunique AS is_unique,
			am.amname AS index_type,
			pg_get_indexdef(ix.indexrelid) AS definition
		FROM pg_class t
		JOIN pg_index ix ON t.oid = ix.indrelid
		JOIN pg_class i ON ix.indexrelid = i.oid
		JOIN pg_am am ON i.relam = am.oid
		JOIN pg_namespace n ON t.relnamespace = n.oid
		JOIN pg_attribute a ON a.attrelid = t.oid AND a.attnum = ANY(ix.indkey)
		WHERE n.nspname = $1 AND t.relname = $2
		GROUP BY i.relname, t.relname, ix.indisunique, am.amname, ix.indexrelid
		ORDER BY i.relname
	`

	rows, err := pool.Query(ctx, query, schema, table)
	if err != nil {
		return nil, fmt.Errorf("failed to query indexes: %w", err)
	}
	defer rows.Close()

	var indexes []models.IndexInfo
	for rows.Next() {
		var idx models.IndexInfo
		if err := rows.Scan(
			&idx.IndexName,
			&idx.TableName,
			&idx.Columns,
			&idx.IsUnique,
			&idx.IndexType,
			&idx.Definition,
		); err != nil {
			return nil, err
		}
		indexes = append(indexes, idx)
	}

	return indexes, rows.Err()
}

// GetForeignKeys returns all foreign keys for a table
func (i *Inspector) GetForeignKeys(ctx context.Context, connID, schema, table string) ([]models.FKInfo, error) {
	pool, err := i.manager.GetPool(connID)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT
			tc.constraint_name,
			kcu.table_schema AS source_schema,
			kcu.table_name AS source_table,
			kcu.column_name AS source_column,
			ccu.table_schema AS target_schema,
			ccu.table_name AS target_table,
			ccu.column_name AS target_column,
			COALESCE(rc.update_rule, 'NO ACTION') AS on_update,
			COALESCE(rc.delete_rule, 'NO ACTION') AS on_delete
		FROM information_schema.table_constraints tc
		JOIN information_schema.key_column_usage kcu
			ON tc.constraint_name = kcu.constraint_name
			AND tc.table_schema = kcu.table_schema
		JOIN information_schema.constraint_column_usage ccu
			ON tc.constraint_name = ccu.constraint_name
			AND tc.table_schema = ccu.table_schema
		LEFT JOIN information_schema.referential_constraints rc
			ON tc.constraint_name = rc.constraint_name
			AND tc.table_schema = rc.constraint_schema
		WHERE tc.constraint_type = 'FOREIGN KEY'
			AND tc.table_schema = $1
			AND tc.table_name = $2
		ORDER BY tc.constraint_name, kcu.ordinal_position
	`

	rows, err := pool.Query(ctx, query, schema, table)
	if err != nil {
		return nil, fmt.Errorf("failed to query foreign keys: %w", err)
	}
	defer rows.Close()

	var fks []models.FKInfo
	for rows.Next() {
		var fk models.FKInfo
		if err := rows.Scan(
			&fk.ConstraintName,
			&fk.SourceSchema,
			&fk.SourceTable,
			&fk.SourceColumn,
			&fk.TargetSchema,
			&fk.TargetTable,
			&fk.TargetColumn,
			&fk.OnUpdate,
			&fk.OnDelete,
		); err != nil {
			return nil, err
		}
		fks = append(fks, fk)
	}

	return fks, rows.Err()
}

// GetViews returns all views in a schema
func (i *Inspector) GetViews(ctx context.Context, connID, schema string) ([]models.ViewInfo, error) {
	pool, err := i.manager.GetPool(connID)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT 
			table_schema,
			table_name,
			is_insertable_into = 'NO' AS read_only,
			view_definition
		FROM information_schema.views
		WHERE table_schema = $1
		ORDER BY table_name
	`

	rows, err := pool.Query(ctx, query, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to query views: %w", err)
	}
	defer rows.Close()

	var views []models.ViewInfo
	for rows.Next() {
		var v models.ViewInfo
		if err := rows.Scan(&v.Schema, &v.Name, &v.ReadOnly, &v.Definition); err != nil {
			return nil, err
		}
		views = append(views, v)
	}

	return views, rows.Err()
}

// GetFunctions returns all functions/procedures in a schema
func (i *Inspector) GetFunctions(ctx context.Context, connID, schema string) ([]models.FunctionInfo, error) {
	pool, err := i.manager.GetPool(connID)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT 
			n.nspname AS schema,
			p.proname AS name,
			pg_get_function_arguments(p.oid) AS arguments,
			pg_get_function_result(p.oid) AS return_type,
			l.lanname AS language
		FROM pg_proc p
		JOIN pg_namespace n ON p.pronamespace = n.oid
		JOIN pg_language l ON p.prolang = l.oid
		WHERE n.nspname = $1
		ORDER BY p.proname
	`

	rows, err := pool.Query(ctx, query, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to query functions: %w", err)
	}
	defer rows.Close()

	var funcs []models.FunctionInfo
	for rows.Next() {
		var f models.FunctionInfo
		if err := rows.Scan(&f.Schema, &f.Name, &f.Arguments, &f.ReturnType, &f.Language); err != nil {
			return nil, err
		}
		funcs = append(funcs, f)
	}

	return funcs, rows.Err()
}

// GetSequences returns all sequences in a schema
func (i *Inspector) GetSequences(ctx context.Context, connID, schema string) ([]models.SequenceInfo, error) {
	pool, err := i.manager.GetPool(connID)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT 
		 schemaname AS schema,
		 sequencename AS name,
		 data_type,
		 min_value,
		 max_value,
		 increment_by AS increment,
		 cache_size
		FROM pg_sequences
		WHERE schemaname = $1
		ORDER BY sequencename
	`

	rows, err := pool.Query(ctx, query, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to query sequences: %w", err)
	}
	defer rows.Close()

	var seqs []models.SequenceInfo
	for rows.Next() {
		var s models.SequenceInfo
		if err := rows.Scan(
			&s.Schema,
			&s.Name,
			&s.DataType,
			&s.MinValue,
			&s.MaxValue,
			&s.Increment,
			&s.CacheSize,
		); err != nil {
			return nil, err
		}
		seqs = append(seqs, s)
	}

	return seqs, rows.Err()
}

// GetTriggers returns all triggers for a table
func (i *Inspector) GetTriggers(ctx context.Context, connID, schema, table string) ([]models.TriggerInfo, error) {
	pool, err := i.manager.GetPool(connID)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT
			tg.tgname AS trigger_name,
			n.nspname AS schema,
			c.relname AS table_name,
			CASE WHEN tg.tgtype & 2 = 2 THEN 'BEFORE' ELSE 'AFTER' END AS timing,
			CASE
				WHEN tg.tgtype & 4 = 4 THEN 'INSERT'
				WHEN tg.tgtype & 8 = 8 THEN 'DELETE'
				WHEN tg.tgtype & 16 = 16 THEN 'UPDATE'
				WHEN tg.tgtype & 20 = 20 THEN 'INSERT OR UPDATE'
				ELSE 'UNKNOWN'
			END AS event,
			pg_get_triggerdef(tg.oid) AS definition
		FROM pg_trigger tg
		JOIN pg_class c ON tg.tgrelid = c.oid
		JOIN pg_namespace n ON c.relnamespace = n.oid
		WHERE NOT tg.tgisinternal
			AND n.nspname = $1
			AND c.relname = $2
		ORDER BY tg.tgname
	`

	rows, err := pool.Query(ctx, query, schema, table)
	if err != nil {
		return nil, fmt.Errorf("failed to query triggers: %w", err)
	}
	defer rows.Close()

	var triggers []models.TriggerInfo
	for rows.Next() {
		var t models.TriggerInfo
		if err := rows.Scan(&t.TriggerName, &t.Schema, &t.TableName, &t.Timing, &t.Event, &t.Definition); err != nil {
			return nil, err
		}
		triggers = append(triggers, t)
	}

	return triggers, rows.Err()
}

// GetTypes returns custom types (enums, etc.) in a schema
func (i *Inspector) GetTypes(ctx context.Context, connID, schema string) ([]models.TypeInfo, error) {
	pool, err := i.manager.GetPool(connID)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT 
			n.nspname AS schema,
			t.typname AS name,
			CASE 
				WHEN t.typtype = 'e' THEN 'enum'
				WHEN t.typtype = 'c' THEN 'composite'
				WHEN t.typtype = 'r' THEN 'range'
				ELSE t.typtype
			END AS type,
			CASE 
				WHEN t.typtype = 'e' THEN
					ARRAY_AGG(e.enumlabel ORDER BY e.enumsortorder)
				ELSE NULL
			END AS values
		FROM pg_type t
		JOIN pg_namespace n ON t.typnamespace = n.oid
		LEFT JOIN pg_enum e ON t.oid = e.enumtypid
		WHERE n.nspname = $1
			AND t.typtype IN ('e', 'c', 'r')
		GROUP BY n.nspname, t.typname, t.typtype
		ORDER BY t.typname
	`

	rows, err := pool.Query(ctx, query, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to query types: %w", err)
	}
	defer rows.Close()

	var types []models.TypeInfo
	for rows.Next() {
		var t models.TypeInfo
		if err := rows.Scan(&t.Schema, &t.Name, &t.Type, &t.Values); err != nil {
			return nil, err
		}
		types = append(types, t)
	}

	return types, rows.Err()
}

// GetTableDDL generates the CREATE TABLE DDL for a table
func (i *Inspector) GetTableDDL(ctx context.Context, connID, schema, table string) (string, error) {
	// Verify connection exists
	_, err := i.manager.GetPool(connID)
	if err != nil {
		return "", err
	}

	// Get column definitions
	columns, err := i.GetColumns(ctx, connID, schema, table)
	if err != nil {
		return "", err
	}

	// Get primary key
	pkCols := make([]string, 0)
	for _, col := range columns {
		if col.IsPrimaryKey {
			pkCols = append(pkCols, col.ColumnName)
		}
	}

	ddl := fmt.Sprintf("CREATE TABLE %s.%s (\n", schema, table)
	for idx, col := range columns {
		ddl += fmt.Sprintf("    %s %s", col.ColumnName, formatDataType(col))
		if col.IsNullable {
			ddl += " NULL"
		} else {
			ddl += " NOT NULL"
		}
		if col.DefaultValue != nil {
			ddl += fmt.Sprintf(" DEFAULT %v", col.DefaultValue)
		}
		if idx < len(columns)-1 || len(pkCols) > 0 {
			ddl += ","
		}
		ddl += "\n"
	}

	if len(pkCols) > 0 {
		ddl += fmt.Sprintf("    CONSTRAINT %s_pkey PRIMARY KEY (%s)\n", table, joinColumns(pkCols))
	}

	ddl += ");"

	// Add indexes
	indexes, err := i.GetIndexes(ctx, connID, schema, table)
	if err == nil {
		for _, idx := range indexes {
			ddl += "\n\n" + idx.Definition + ";"
		}
	}

	return ddl, nil
}

func formatDataType(col models.ColumnInfo) string {
	dt := col.DataType
	if col.CharacterMaximumLength != nil && *col.CharacterMaximumLength > 0 {
		return fmt.Sprintf("%s(%d)", dt, *col.CharacterMaximumLength)
	}
	return dt
}

func joinColumns(cols []string) string {
	result := ""
	for i, c := range cols {
		if i > 0 {
			result += ", "
		}
		result += c
	}
	return result
}

// RefreshSchema invalidates any cached schema data
func (i *Inspector) RefreshSchema(ctx context.Context, connID string) error {
	// Schema is queried fresh each time, so this is a no-op
	// but useful for future caching
	return nil
}

// GetFullSchema returns all schema objects for a connection
func (i *Inspector) GetFullSchema(ctx context.Context, connID string) (map[string]interface{}, error) {
	schemas, err := i.GetSchemas(ctx, connID)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	result["schemas"] = schemas

	for _, schema := range schemas {
		tables, err := i.GetTables(ctx, connID, schema.Name)
		if err != nil {
			continue
		}
		result[fmt.Sprintf("tables_%s", schema.Name)] = tables
	}

	return result, nil
}

// Helper: ensure pool exists
func ensurePool(manager *connection.Manager, connID string) (*pgxpool.Pool, error) {
	return manager.GetPool(connID)
}
