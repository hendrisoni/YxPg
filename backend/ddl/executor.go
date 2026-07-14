package ddl

import (
	"context"
	"fmt"
	"strings"

	"yxpg/backend/connection"
	"yxpg/backend/models"
)

// Executor handles DDL operations
type Executor struct {
	manager *connection.Manager
}

// NewExecutor creates a new DDL executor
func NewExecutor(manager *connection.Manager) *Executor {
	return &Executor{manager: manager}
}

// CreateTable creates a new table
func (e *Executor) CreateTable(ctx context.Context, connID string, def models.TableDefinition) (string, error) {
	ddl := e.buildCreateTableDDL(def)

	pool, err := e.manager.GetPool(connID)
	if err != nil {
		return "", err
	}

	_, err = pool.Exec(ctx, ddl)
	if err != nil {
		return "", fmt.Errorf("failed to create table: %w", err)
	}

	return ddl, nil
}

// AlterTable alters an existing table
func (e *Executor) AlterTable(ctx context.Context, connID string, original, modified models.TableDefinition) ([]string, error) {
	statements := e.generateAlterStatements(original, modified)

	if len(statements) == 0 {
		return nil, nil
	}

	pool, err := e.manager.GetPool(connID)
	if err != nil {
		return nil, err
	}

	for _, stmt := range statements {
		_, err = pool.Exec(ctx, stmt)
		if err != nil {
			return statements, fmt.Errorf("failed to execute: %s: %w", stmt, err)
		}
	}

	return statements, nil
}

// DropTable drops a table
func (e *Executor) DropTable(ctx context.Context, connID, schema, table string, cascade bool) error {
	pool, err := e.manager.GetPool(connID)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("DROP TABLE %s.%s", schema, table)
	if cascade {
		query += " CASCADE"
	}

	_, err = pool.Exec(ctx, query)
	return err
}

// RenameTable renames a table
func (e *Executor) RenameTable(ctx context.Context, connID, schema, oldName, newName string) error {
	pool, err := e.manager.GetPool(connID)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("ALTER TABLE %s.%s RENAME TO %s", schema, oldName, newName)
	_, err = pool.Exec(ctx, query)
	return err
}

// CreateIndex creates a new index
func (e *Executor) CreateIndex(ctx context.Context, connID string, def models.IndexDefinition) error {
	pool, err := e.manager.GetPool(connID)
	if err != nil {
		return err
	}

	ddl := e.buildCreateIndexDDL(def)
	_, err = pool.Exec(ctx, ddl)
	return err
}

// DropIndex drops an index
func (e *Executor) DropIndex(ctx context.Context, connID, indexName string) error {
	pool, err := e.manager.GetPool(connID)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("DROP INDEX %s", indexName)
	_, err = pool.Exec(ctx, query)
	return err
}

// AddForeignKey adds a foreign key constraint
func (e *Executor) AddForeignKey(ctx context.Context, connID string, def models.FKDefinition) error {
	pool, err := e.manager.GetPool(connID)
	if err != nil {
		return err
	}

	query := e.buildAddFKDDL(def)
	_, err = pool.Exec(ctx, query)
	return err
}

// DropConstraint drops a constraint from a table
func (e *Executor) DropConstraint(ctx context.Context, connID, schema, table, constraint string) error {
	pool, err := e.manager.GetPool(connID)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("ALTER TABLE %s.%s DROP CONSTRAINT %s", schema, table, constraint)
	_, err = pool.Exec(ctx, query)
	return err
}

// ExecuteRaw executes a raw DDL statement
func (e *Executor) ExecuteRaw(ctx context.Context, connID, sql string) error {
	pool, err := e.manager.GetPool(connID)
	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx, sql)
	return err
}

// buildCreateTableDDL generates CREATE TABLE DDL from a definition
func (e *Executor) buildCreateTableDDL(def models.TableDefinition) string {
	var parts []string

	for _, col := range def.Columns {
		colStr := fmt.Sprintf("    %s %s", col.Name, col.DataType)

		if col.Length != nil {
			colStr += fmt.Sprintf("(%d)", *col.Length)
		}

		if !col.IsNullable {
			colStr += " NOT NULL"
		}

		if col.DefaultValue != "" {
			colStr += fmt.Sprintf(" DEFAULT %s", col.DefaultValue)
		}

		if col.IsUnique {
			colStr += " UNIQUE"
		}

		parts = append(parts, colStr)
	}

	// Add primary key constraint
	pkCols := make([]string, 0)
	for _, col := range def.Columns {
		if col.IsPrimaryKey {
			pkCols = append(pkCols, col.Name)
		}
	}
	if len(pkCols) > 0 {
		parts = append(parts, fmt.Sprintf("    CONSTRAINT %s_pkey PRIMARY KEY (%s)",
			def.TableName, strings.Join(pkCols, ", ")))
	}

	return fmt.Sprintf("CREATE TABLE %s.%s (\n%s\n);", def.Schema, def.TableName, strings.Join(parts, ",\n"))
}

// buildCreateIndexDDL generates CREATE INDEX DDL
func (e *Executor) buildCreateIndexDDL(def models.IndexDefinition) string {
	unique := ""
	if def.IsUnique {
		unique = "UNIQUE "
	}

	indexType := ""
	if def.IndexType != "" && def.IndexType != "btree" {
		indexType = fmt.Sprintf(" USING %s", strings.ToUpper(def.IndexType))
	}

	where := ""
	if def.Where != "" {
		where = fmt.Sprintf(" WHERE %s", def.Where)
	}

	return fmt.Sprintf("CREATE %sINDEX %s ON %s%s (%s)%s",
		unique, def.Name, def.TableName, indexType, strings.Join(def.Columns, ", "), where)
}

// buildAddFKDDL generates ALTER TABLE ADD CONSTRAINT for foreign key
func (e *Executor) buildAddFKDDL(def models.FKDefinition) string {
	onUpdate := "NO ACTION"
	if def.OnUpdate != "" {
		onUpdate = def.OnUpdate
	}
	onDelete := "NO ACTION"
	if def.OnDelete != "" {
		onDelete = def.OnDelete
	}

	return fmt.Sprintf(
		"ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s.%s (%s) ON UPDATE %s ON DELETE %s",
		def.SourceColumns[0], // Simplified - in production use table name
		def.Name,
		strings.Join(def.SourceColumns, ", "),
		def.TargetSchema,
		def.TargetTable,
		strings.Join(def.TargetColumns, ", "),
		onUpdate,
		onDelete,
	)
}

// generateAlterStatements compares two table definitions and generates ALTER statements
func (e *Executor) generateAlterStatements(original, modified models.TableDefinition) []string {
	var statements []string

	// Create maps for comparison
	origCols := make(map[string]models.ColumnDefinition)
	for _, col := range original.Columns {
		origCols[col.Name] = col
	}

	modCols := make(map[string]models.ColumnDefinition)
	for _, col := range modified.Columns {
		modCols[col.Name] = col
	}

	// Find added columns
	for name, col := range modCols {
		if _, exists := origCols[name]; !exists {
			stmt := fmt.Sprintf("ALTER TABLE %s.%s ADD COLUMN %s %s",
				modified.Schema, modified.TableName, col.Name, col.DataType)
			if !col.IsNullable {
				stmt += " NOT NULL"
			}
			if col.DefaultValue != "" {
				stmt += fmt.Sprintf(" DEFAULT %s", col.DefaultValue)
			}
			statements = append(statements, stmt)
		}
	}

	// Find dropped columns
	for name := range origCols {
		if _, exists := modCols[name]; !exists {
			statements = append(statements, fmt.Sprintf(
				"ALTER TABLE %s.%s DROP COLUMN %s",
				original.Schema, original.TableName, name))
		}
	}

	// Find modified columns
	for name, modCol := range modCols {
		origCol, exists := origCols[name]
		if !exists {
			continue
		}

		if origCol.DataType != modCol.DataType {
			statements = append(statements, fmt.Sprintf(
				"ALTER TABLE %s.%s ALTER COLUMN %s TYPE %s",
				modified.Schema, modified.TableName, name, modCol.DataType))
		}

		if origCol.IsNullable != modCol.IsNullable {
			if modCol.IsNullable {
				statements = append(statements, fmt.Sprintf(
					"ALTER TABLE %s.%s ALTER COLUMN %s DROP NOT NULL",
					modified.Schema, modified.TableName, name))
			} else {
				statements = append(statements, fmt.Sprintf(
					"ALTER TABLE %s.%s ALTER COLUMN %s SET NOT NULL",
					modified.Schema, modified.TableName, name))
			}
		}
	}

	return statements
}
