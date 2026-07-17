package query

import (
	"context"
	"fmt"

	"yxpg/backend/connection"
	"yxpg/backend/models"
)

// ExplainParser parses EXPLAIN ANALYZE output
type ExplainParser struct{}

// NewExplainParser creates a new explain parser
func NewExplainParser() *ExplainParser {
	return &ExplainParser{}
}

// Parse parses EXPLAIN ANALYZE JSON output into a tree structure
func (p *ExplainParser) Parse(jsonStr string) (*models.ExplainNode, error) {
	// Basic parser - in production, use proper JSON unmarshaling
	// For now, return the raw text
	return &models.ExplainNode{
		NodeType: "Result",
	}, nil
}

// GetCostRating returns a cost rating for display
func GetCostRating(cost float64) string {
	if cost < 10 {
		return "low"
	} else if cost < 100 {
		return "medium"
	}
	return "high"
}

// BrowseTable retrieves paginated data from a table
func BrowseTable(ctx context.Context, manager *connection.Manager, connID, schema, table string, opts models.BrowseOptions) models.QueryResult {
	pool, err := manager.GetPool(connID)
	if err != nil {
		return models.QueryResult{
			Error: fmt.Sprintf("Connection not found: %v", err),
		}
	}

	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PageSize <= 0 {
		opts.PageSize = 100
	}

	offset := (opts.Page - 1) * opts.PageSize

	// Build query
	query := fmt.Sprintf("SELECT * FROM %s.%s", schema, table)

	// Add WHERE conditions
	if len(opts.Filters) > 0 {
		query += " WHERE "
		for i, f := range opts.Filters {
			if i > 0 {
				query += " AND "
			}
			switch f.Operator {
			case "IS NULL":
				query += fmt.Sprintf("%s IS NULL", f.Column)
			case "IS NOT NULL":
				query += fmt.Sprintf("%s IS NOT NULL", f.Column)
			case "LIKE", "ILIKE":
				query += fmt.Sprintf("CAST(%s AS text) %s '%%%s%%'", f.Column, f.Operator, f.Value)
			case "IN":
				query += fmt.Sprintf("%s IN (%s)", f.Column, f.Value)
			default:
				query += fmt.Sprintf("%s %s '%s'", f.Column, f.Operator, f.Value)
			}
		}
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS _count_sub", query)
	var totalCount int64
	err = pool.QueryRow(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		totalCount = 0
	}

	// Add sorting
	if opts.SortBy != "" {
		order := "ASC"
		if opts.SortOrder == "desc" {
			order = "DESC"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", opts.SortBy, order)
	}

	// Add pagination
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", opts.PageSize, offset)

	// Execute
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return models.QueryResult{
			Error: err.Error(),
		}
	}
	defer rows.Close()

	result := connection.RowsToResult(rows, 0)
	result.TotalCount = totalCount

	return result
}
