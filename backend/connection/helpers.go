package connection

import (
	"context"
	"fmt"
	"strconv"

	"yxpg/backend/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ExtractInt extracts an int value from a pgx.Row
func ExtractInt(value interface{}) int {
	switch v := value.(type) {
	case int32:
		return int(v)
	case int64:
		return int(v)
	case int:
		return v
	case float64:
		return int(v)
	case string:
		n, _ := strconv.Atoi(v)
		return n
	default:
		return 0
	}
}

// ExtractInt64 extracts an int64 value from various types
func ExtractInt64(value interface{}) int64 {
	switch v := value.(type) {
	case int32:
		return int64(v)
	case int64:
		return v
	case int:
		return int64(v)
	case float64:
		return int64(v)
	default:
		return 0
	}
}

// ExtractString extracts a string value from various types
func ExtractString(value interface{}) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%v", value)
}

// RowToMap converts a pgx.Row to a map
func RowToMap(row pgx.Row, fields []pgconn.FieldDescription) (map[string]interface{}, error) {
	values := make([]interface{}, len(fields))
	valuePtrs := make([]interface{}, len(fields))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	if err := row.Scan(valuePtrs...); err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for i, field := range fields {
		result[field.Name] = values[i]
	}

	return result, nil
}

// RowsToResult converts pgx.Rows to a QueryResult
func RowsToResult(rows pgx.Rows, duration int64) models.QueryResult {
	fields := rows.FieldDescriptions()
	columns := make([]models.ColumnMeta, len(fields))
	for i, field := range fields {
		columns[i] = models.ColumnMeta{
			Name:     field.Name,
			DataType: fmt.Sprintf("%d", field.DataTypeOID),
		}
	}

	allRows := [][]interface{}{}
	for rows.Next() {
		values := make([]interface{}, len(fields))
		valuePtrs := make([]interface{}, len(fields))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			continue
		}

		// Convert values to JSON-serializable types
		row := make([]interface{}, len(values))
		for i, v := range values {
			switch val := v.(type) {
			case []byte:
				row[i] = string(val)
			case int32:
				row[i] = int64(val)
			case int64:
				row[i] = val
			case float64:
				row[i] = val
			case bool:
				row[i] = val
			case string:
				row[i] = val
			case nil:
				row[i] = nil
			default:
				row[i] = fmt.Sprintf("%v", val)
			}
		}

		allRows = append(allRows, row)
	}

	err := rows.Err()

	result := models.QueryResult{
		Columns:  columns,
		Rows:     allRows,
		RowCount: len(allRows),
		Duration: duration,
	}

	if err != nil {
		result.Error = err.Error()
	}

	return result
}

// DetectQueryType determines the type of SQL query
func DetectQueryType(sql string) string {
	// Simple heuristic - uppercase the SQL and trim whitespace
	if len(sql) == 0 {
		return "unknown"
	}

	// Trim and uppercase for comparison
	trimmed := sql
	for len(trimmed) > 0 && (trimmed[0] == ' ' || trimmed[0] == '\t' || trimmed[0] == '\n') {
		trimmed = trimmed[1:]
	}

	if len(trimmed) < 3 {
		return "unknown"
	}

	upper := ""
	for _, c := range trimmed {
		if c >= 'A' && c <= 'Z' {
			upper += string(c)
		} else if c >= 'a' && c <= 'z' {
			upper += string(c - 32)
		} else {
			break
		}
	}

	switch {
	case len(upper) >= 6 && upper[:6] == "SELECT":
		return "select"
	case len(upper) >= 6 && upper[:6] == "INSERT":
		return "insert"
	case len(upper) >= 6 && upper[:6] == "UPDATE":
		return "update"
	case len(upper) >= 6 && upper[:6] == "DELETE":
		return "delete"
	case len(upper) >= 4 && upper[:4] == "WITH":
		return "select"
	case len(upper) >= 6 && upper[:6] == "CREATE":
		return "ddl"
	case len(upper) >= 5 && upper[:5] == "ALTER":
		return "ddl"
	case len(upper) >= 4 && upper[:4] == "DROP":
		return "ddl"
	default:
		return "other"
	}
}

// GetTotalCount gets the total count of rows matching a query for pagination
func GetTotalCount(ctx context.Context, pool *pgxpool.Pool, baseQuery string) int64 {
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS _count_subquery", baseQuery)
	var count int64
	err := pool.QueryRow(ctx, countQuery).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}
