package dbexport

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"

	"yxpg/backend/models"
)

// ExportCSV exports query result as CSV
func ExportCSV(result models.QueryResult, delimiter string) (string, error) {
	if delimiter == "" {
		delimiter = ","
	}

	var sb strings.Builder
	writer := csv.NewWriter(&sb)
	writer.Comma = rune(delimiter[0])

	// Write header
	header := make([]string, len(result.Columns))
	for i, col := range result.Columns {
		header[i] = col.Name
	}
	if err := writer.Write(header); err != nil {
		return "", err
	}

	// Write rows
	for _, row := range result.Rows {
		record := make([]string, len(row))
		for i, val := range row {
			if val == nil {
				record[i] = "NULL"
			} else {
				record[i] = fmt.Sprintf("%v", val)
			}
		}
		if err := writer.Write(record); err != nil {
			return "", err
		}
	}

	writer.Flush()
	return sb.String(), writer.Error()
}

// ExportJSON exports query result as JSON
func ExportJSON(result models.QueryResult) (string, error) {
	// Convert rows to array of objects
	objects := make([]map[string]interface{}, 0, len(result.Rows))
	for _, row := range result.Rows {
		obj := make(map[string]interface{})
		for i, col := range result.Columns {
			if i < len(row) {
				obj[col.Name] = row[i]
			}
		}
		objects = append(objects, obj)
	}

	data, err := json.MarshalIndent(objects, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// ExportSQL exports query result as INSERT statements
func ExportSQL(result models.QueryResult, schema, table string) (string, error) {
	var sb strings.Builder

	for _, row := range result.Rows {
		columns := make([]string, 0, len(result.Columns))
		values := make([]string, 0, len(row))

		for i, col := range result.Columns {
			if i < len(row) {
				columns = append(columns, col.Name)
				values = append(values, formatSQLValue(row[i]))
			}
		}

		sb.WriteString(fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES (%s);\n",
			schema, table,
			strings.Join(columns, ", "),
			strings.Join(values, ", ")))
	}

	return sb.String(), nil
}

func formatSQLValue(val interface{}) string {
	if val == nil {
		return "NULL"
	}

	switch v := val.(type) {
	case string:
		escaped := strings.ReplaceAll(v, "'", "''")
		return fmt.Sprintf("'%s'", escaped)
	case bool:
		if v {
			return "TRUE"
		}
		return "FALSE"
	case int, int32, int64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	default:
		escaped := strings.ReplaceAll(fmt.Sprintf("%v", v), "'", "''")
		return fmt.Sprintf("'%s'", escaped)
	}
}
