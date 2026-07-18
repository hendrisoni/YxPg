package dbexport

import (
	"context"
	"fmt"
	"time"

	"yxpg/backend/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// MaintenanceOptions defines the configuration for a database maintenance task
type MaintenanceOptions struct {
	ConnectionID string `json:"connection_id"`
	Database     string `json:"database"`
	Task         string `json:"task"`        // "vacuum_analyze", "vacuum_full", "reindex"
	Scope        string `json:"scope"`       // "database", "schema", "table"
	SchemaName   string `json:"schema_name"` // optional schema target
	TableName    string `json:"table_name"`  // optional table target
	Verbose      bool   `json:"verbose"`
}

// RunMaintenance executes a maintenance task over a dedicated connection and streams notices
func RunMaintenance(ctx context.Context, conn models.Connection, opts MaintenanceOptions) error {
	// 1. Build statement
	var sql string
	switch opts.Task {
	case "vacuum_analyze":
		if opts.Scope == "table" && opts.SchemaName != "" && opts.TableName != "" {
			sql = fmt.Sprintf(`VACUUM (ANALYZE%s) "%s"."%s"`, mapVerbose(opts.Verbose), opts.SchemaName, opts.TableName)
		} else {
			sql = fmt.Sprintf("VACUUM (ANALYZE%s)", mapVerbose(opts.Verbose))
		}
	case "vacuum_full":
		if opts.Scope == "table" && opts.SchemaName != "" && opts.TableName != "" {
			sql = fmt.Sprintf(`VACUUM (FULL, ANALYZE%s) "%s"."%s"`, mapVerbose(opts.Verbose), opts.SchemaName, opts.TableName)
		} else {
			sql = fmt.Sprintf("VACUUM (FULL, ANALYZE%s)", mapVerbose(opts.Verbose))
		}
	case "reindex":
		switch opts.Scope {
		case "schema":
			sql = fmt.Sprintf(`REINDEX %s SCHEMA "%s"`, mapReindexVerbose(opts.Verbose), opts.SchemaName)
		case "table":
			sql = fmt.Sprintf(`REINDEX %s TABLE "%s"."%s"`, mapReindexVerbose(opts.Verbose), opts.SchemaName, opts.TableName)
		default:
			// database
			sql = fmt.Sprintf(`REINDEX %s DATABASE "%s"`, mapReindexVerbose(opts.Verbose), opts.Database)
		}
	default:
		return fmt.Errorf("unknown maintenance task: %s", opts.Task)
	}

	wailsruntime.EventsEmit(ctx, "maintenance:log", fmt.Sprintf("[%s] Starting database maintenance on database '%s'...", time.Now().Format("15:04:05"), opts.Database))
	wailsruntime.EventsEmit(ctx, "maintenance:log", fmt.Sprintf("Executing query: %s", sql))

	// 2. Parse config and connect
	sslMode := conn.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		conn.Host, conn.Port, conn.Username, conn.Password, opts.Database, sslMode,
	)

	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		errStr := fmt.Sprintf("Failed to parse connection config: %v", err)
		wailsruntime.EventsEmit(ctx, "maintenance:log", errStr)
		wailsruntime.EventsEmit(ctx, "maintenance:status", map[string]interface{}{
			"status":  "error",
			"message": errStr,
		})
		return err
	}

	// Capture notice/warnings logs and stream them to the frontend
	config.OnNotice = func(c *pgconn.PgConn, notice *pgconn.Notice) {
		wailsruntime.EventsEmit(ctx, "maintenance:log", fmt.Sprintf("[%s] %s", notice.Severity, notice.Message))
	}

	dbConn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		errStr := fmt.Sprintf("Failed to connect: %v", err)
		wailsruntime.EventsEmit(ctx, "maintenance:log", errStr)
		wailsruntime.EventsEmit(ctx, "maintenance:status", map[string]interface{}{
			"status":  "error",
			"message": errStr,
		})
		return err
	}
	defer dbConn.Close(ctx)

	// 3. Execute
	startTime := time.Now()
	_, err = dbConn.Exec(ctx, sql)
	duration := time.Since(startTime)

	if err != nil {
		errStr := fmt.Sprintf("Maintenance failed: %v", err)
		wailsruntime.EventsEmit(ctx, "maintenance:log", errStr)
		wailsruntime.EventsEmit(ctx, "maintenance:status", map[string]interface{}{
			"status":  "error",
			"message": errStr,
		})
		return err
	}

	successMsg := fmt.Sprintf("Completed maintenance successfully in %v!", duration.Round(time.Millisecond))
	wailsruntime.EventsEmit(ctx, "maintenance:log", successMsg)
	wailsruntime.EventsEmit(ctx, "maintenance:status", map[string]interface{}{
		"status":  "success",
		"message": successMsg,
	})

	return nil
}

func mapVerbose(verbose bool) string {
	if verbose {
		return ", VERBOSE"
	}
	return ""
}

func mapReindexVerbose(verbose bool) string {
	if verbose {
		return "(VERBOSE)"
	}
	return ""
}
