package ddl

import (
	"fmt"
	"strings"

	"yxpg/backend/models"
)

// Builder helps construct DDL statements
type Builder struct{}

// NewBuilder creates a new DDL builder
func NewBuilder() *Builder {
	return &Builder{}
}

// BuildCreateTable builds a CREATE TABLE statement
func (b *Builder) BuildCreateTable(def models.TableDefinition) string {
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

		if col.Comment != "" {
			colStr += fmt.Sprintf(" -- %s", col.Comment)
		}

		parts = append(parts, colStr)
	}

	// Primary key
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

// BuildAlterTable builds ALTER TABLE statements
func (b *Builder) BuildAlterTable(def models.TableDefinition, changes []ColumnChange) []string {
	var statements []string

	for _, change := range changes {
		switch change.Type {
		case "add":
			stmt := fmt.Sprintf("ALTER TABLE %s.%s ADD COLUMN %s %s",
				def.Schema, def.TableName, change.Column.Name, change.Column.DataType)
			if !change.Column.IsNullable {
				stmt += " NOT NULL"
			}
			if change.Column.DefaultValue != "" {
				stmt += fmt.Sprintf(" DEFAULT %s", change.Column.DefaultValue)
			}
			statements = append(statements, stmt)

		case "drop":
			statements = append(statements, fmt.Sprintf(
				"ALTER TABLE %s.%s DROP COLUMN %s",
				def.Schema, def.TableName, change.Column.Name))

		case "alter":
			if change.OldColumn.DataType != change.Column.DataType {
				statements = append(statements, fmt.Sprintf(
					"ALTER TABLE %s.%s ALTER COLUMN %s TYPE %s",
					def.Schema, def.TableName, change.Column.Name, change.Column.DataType))
			}
			if change.OldColumn.IsNullable != change.Column.IsNullable {
				if change.Column.IsNullable {
					statements = append(statements, fmt.Sprintf(
						"ALTER TABLE %s.%s ALTER COLUMN %s DROP NOT NULL",
						def.Schema, def.TableName, change.Column.Name))
				} else {
					statements = append(statements, fmt.Sprintf(
						"ALTER TABLE %s.%s ALTER COLUMN %s SET NOT NULL",
						def.Schema, def.TableName, change.Column.Name))
				}
			}
		}
	}

	return statements
}

// BuildCreateIndex builds a CREATE INDEX statement
func (b *Builder) BuildCreateIndex(def models.IndexDefinition) string {
	unique := ""
	if def.IsUnique {
		unique = "UNIQUE "
	}

	indexType := ""
	if def.IndexType != "" && def.IndexType != "btree" {
		indexType = fmt.Sprintf(" USING %s", strings.ToUpper(def.IndexType))
	}

	return fmt.Sprintf("CREATE %sINDEX %s ON %s%s (%s)",
		unique, def.Name, def.TableName, indexType, strings.Join(def.Columns, ", "))
}

// BuildDropIndex builds a DROP INDEX statement
func (b *Builder) BuildDropIndex(name string) string {
	return fmt.Sprintf("DROP INDEX IF EXISTS %s", name)
}

// BuildAddForeignKey builds an ALTER TABLE ADD CONSTRAINT for FK
func (b *Builder) BuildAddForeignKey(tableName string, def models.FKDefinition) string {
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
		tableName,
		def.Name,
		strings.Join(def.SourceColumns, ", "),
		def.TargetSchema,
		def.TargetTable,
		strings.Join(def.TargetColumns, ", "),
		onUpdate,
		onDelete,
	)
}

// BuildDropConstraint builds a DROP CONSTRAINT statement
func (b *Builder) BuildDropConstraint(tableName, constraintName string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s", tableName, constraintName)
}

// ColumnChange represents a change to a column
type ColumnChange struct {
	Type      string                  // "add", "drop", "alter"
	Column    models.ColumnDefinition // new column definition
	OldColumn models.ColumnDefinition // old column definition (for alter)
}

// PreviewAlter generates a preview of ALTER statements for given changes
func (b *Builder) PreviewAlter(def models.TableDefinition, changes []ColumnChange) string {
	statements := b.BuildAlterTable(def, changes)
	if len(statements) == 0 {
		return "-- No changes to apply"
	}
	return strings.Join(statements, ";\n") + ";"
}
