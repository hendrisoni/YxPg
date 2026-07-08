package query

import (
	"context"
	"fmt"
	"time"

	"yxpg/backend/connection"
	"yxpg/backend/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Executor handles query execution
type Executor struct {
	manager *connection.Manager
	history *History
}

// NewExecutor creates a new query executor
func NewExecutor(manager *connection.Manager, history *History) *Executor {
	return &Executor{
		manager: manager,
		history: history,
	}
}

// Execute runs a SQL query and returns the result
func (e *Executor) Execute(ctx context.Context, connID, sql string, timeout int) models.QueryResult {
	if timeout <= 0 {
		timeout = 30
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	e.manager.StartQuery(connID, cancel)
	defer e.manager.FinishQuery(connID)

	return e.ExecuteWithConnection(ctx, connID, sql, timeout)
}

// ExecuteMultiple splits by semicolon and runs multiple statements, returning array of results
func (e *Executor) ExecuteMultiple(ctx context.Context, connID, sql string, timeout int) []models.QueryResult {
	if timeout <= 0 {
		timeout = 30
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	e.manager.StartQuery(connID, cancel)
	defer e.manager.FinishQuery(connID)

	// Very basic split by semicolon, ignoring those inside quotes (simplified)
	// A proper parser is better, but this handles 90% of basic scripts
	statements := splitSQLStatements(sql)
	var results []models.QueryResult

	for _, stmt := range statements {
		if stmt == "" {
			continue
		}
		if ctx.Err() != nil {
			break
		}
		res := e.ExecuteWithConnection(ctx, connID, stmt, timeout)
		res.RawSQL = stmt // pony: add RawSQL to result for UI tab labels
		results = append(results, res)
	}

	if len(results) == 0 {
		if ctx.Err() != nil {
			return []models.QueryResult{{
				Error:    ctx.Err().Error(),
				Duration: 0,
			}}
		}
		// Fallback if split fails or empty
		return []models.QueryResult{e.ExecuteWithConnection(ctx, connID, sql, timeout)}
	}

	return results
}

// splitSQLStatements splits a script by ';' but tries to ignore semicolons inside single quotes
func splitSQLStatements(sql string) []string {
	var stmts []string
	var current []rune
	inQuote := false

	for _, char := range sql {
		if char == '\'' {
			inQuote = !inQuote
		}
		if char == ';' && !inQuote {
			stmts = append(stmts, string(current))
			current = []rune{}
		} else {
			current = append(current, char)
		}
	}
	if len(current) > 0 {
		// add remaining
		stmts = append(stmts, string(current))
	}
	return stmts
}

// ExecuteWithConnection executes a SQL query via pgxpool
func (e *Executor) ExecuteWithConnection(ctx context.Context, connID, sql string, timeout int) models.QueryResult {
	if timeout <= 0 {
		timeout = 30
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	start := time.Now()

	pool, err := e.manager.GetPool(connID)
	if err != nil {
		return models.QueryResult{
			Error:    fmt.Sprintf("Connection not found: %v", err),
			Duration: time.Since(start).Milliseconds(),
		}
	}

	queryType := connection.DetectQueryType(sql)

	switch queryType {
	case "select":
		return e.executeSelect(ctx, pool, connID, sql, start)
	case "insert", "update", "delete":
		return e.executeDML(ctx, pool, connID, sql, start)
	case "ddl":
		return e.executeDDL(ctx, pool, connID, sql, start)
	default:
		return e.executeOther(ctx, pool, connID, sql, start)
	}
}

func (e *Executor) executeSelect(ctx context.Context, pool *pgxpool.Pool, connID, sql string, start time.Time) models.QueryResult {
	rows, err := pool.Query(ctx, sql)
	if err != nil {
		duration := time.Since(start).Milliseconds()
		e.saveHistory(connID, "", sql, duration, 0, err.Error())
		return models.QueryResult{
			Error:    err.Error(),
			Duration: duration,
		}
	}
	defer rows.Close()

	result := connection.RowsToResult(rows, time.Since(start).Milliseconds())
	e.saveHistory(connID, "", sql, result.Duration, result.RowCount, "")
	return result
}

func (e *Executor) executeDML(ctx context.Context, pool *pgxpool.Pool, connID, sql string, start time.Time) models.QueryResult {
	tag, err := pool.Exec(ctx, sql)
	duration := time.Since(start).Milliseconds()

	if err != nil {
		e.saveHistory(connID, "", sql, duration, 0, err.Error())
		return models.QueryResult{
			Error:    err.Error(),
			Duration: duration,
		}
	}

	rowsAffected := tag.RowsAffected()
	e.saveHistory(connID, "", sql, duration, 0, "")
	return models.QueryResult{
		QueryType:    connection.DetectQueryType(sql),
		Duration:     duration,
		RowsAffected: rowsAffected,
	}
}

func (e *Executor) executeDDL(ctx context.Context, pool *pgxpool.Pool, connID, sql string, start time.Time) models.QueryResult {
	_, err := pool.Exec(ctx, sql)
	duration := time.Since(start).Milliseconds()

	if err != nil {
		e.saveHistory(connID, "", sql, duration, 0, err.Error())
		return models.QueryResult{
			Error:     err.Error(),
			Duration:  duration,
			QueryType: "ddl",
		}
	}

	e.saveHistory(connID, "", sql, duration, 0, "")
	return models.QueryResult{
		QueryType: "ddl",
		Duration:  duration,
	}
}

func (e *Executor) executeOther(ctx context.Context, pool *pgxpool.Pool, connID, sql string, start time.Time) models.QueryResult {
	_, err := pool.Exec(ctx, sql)
	duration := time.Since(start).Milliseconds()

	if err != nil {
		e.saveHistory(connID, "", sql, duration, 0, err.Error())
		return models.QueryResult{
			Error:     err.Error(),
			Duration:  duration,
			QueryType: "other",
		}
	}

	e.saveHistory(connID, "", sql, duration, 0, "")
	return models.QueryResult{
		QueryType: "other",
		Duration:  duration,
	}
}

func (e *Executor) saveHistory(connID, database, sql string, durationMs int64, rowsReturned int, errMsg string) {
	if e.history != nil {
		entry := models.QueryHistoryEntry{
			ConnectionID: connID,
			Database:     database,
			SQL:          sql,
			DurationMs:   durationMs,
			RowsReturned: rowsReturned,
			ExecutedAt:   time.Now(),
			Error:        errMsg,
		}
		e.history.Save(entry)
	}
}

// GetHistory returns query history for a connection
func (e *Executor) GetHistory(connID string, limit int) ([]models.QueryHistoryEntry, error) {
	if e.history == nil {
		return nil, nil
	}
	return e.history.ListByConnection(connID, limit)
}

// SaveQuery saves a named query
func (e *Executor) SaveQuery(name, sql string) error {
	if e.history == nil {
		return fmt.Errorf("history not initialized")
	}
	return e.history.SaveQuery(name, sql)
}

// GetSavedQueries returns all saved queries
func (e *Executor) GetSavedQueries() ([]models.SavedQuery, error) {
	if e.history == nil {
		return nil, nil
	}
	return e.history.ListSavedQueries()
}

// CancelQuery cancels an active query
func (e *Executor) CancelQuery(connID string) error {
	return e.manager.CancelQuery(connID)
}

// ExplainQuery runs EXPLAIN ANALYZE and returns structured result
func (e *Executor) ExplainQuery(ctx context.Context, connID, sql string) models.ExplainResult {
	explainSQL := fmt.Sprintf("EXPLAIN (FORMAT JSON, ANALYZE, BUFFERS) %s", sql)

	pool, err := e.manager.GetPool(connID)
	if err != nil {
		return models.ExplainResult{
			RawText: fmt.Sprintf("Error: %v", err),
		}
	}

	var resultJSON []byte
	err = pool.QueryRow(ctx, explainSQL).Scan(&resultJSON)
	if err != nil {
		return models.ExplainResult{
			RawText: fmt.Sprintf("Error executing EXPLAIN: %v", err),
		}
	}

	return models.ExplainResult{
		RawText: string(resultJSON),
	}
}
