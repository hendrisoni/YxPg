package dbexport

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"yxpg/backend/models"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// BackupOptions defines the parameters for pg_dump
type BackupOptions struct {
	ConnectionID    string `json:"connection_id"`
	Database        string `json:"database"`
	OutputPath      string `json:"output_path"`
	Format          string `json:"format"` // p (plain), c (custom), d (directory), t (tar)
	SchemaOnly      bool   `json:"schema_only"`
	DataOnly        bool   `json:"data_only"`
	Clean           bool   `json:"clean"`
	Create          bool   `json:"create"`
	Inserts         bool   `json:"inserts"`
	ColumnInserts   bool   `json:"column_inserts"`
	DisableTriggers bool   `json:"disable_triggers"`
	Verbose         bool   `json:"verbose"`
	PgBinPath       string `json:"pg_bin_path"`
}

// RunPgDump executes pg_dump with the specified options and streams logs to the frontend
func RunPgDump(ctx context.Context, conn models.Connection, opts BackupOptions) error {
	// 1. Determine pg_dump path
	pgDumpCmd := "pg_dump"
	if opts.PgBinPath != "" {
		binDir := strings.TrimSpace(opts.PgBinPath)
		pgDumpCmd = filepath.Join(binDir, "pg_dump")
	}

	// For Windows, append .exe if it doesn't have it
	if runtime.GOOS == "windows" && !strings.HasSuffix(strings.ToLower(pgDumpCmd), ".exe") {
		if strings.ContainsAny(pgDumpCmd, `/\`) {
			pgDumpCmd += ".exe"
		}
	}

	// 2. Build arguments
	var args []string
	args = append(args, "-h", conn.Host)
	args = append(args, "-p", fmt.Sprintf("%d", conn.Port))
	args = append(args, "-U", conn.Username)
	args = append(args, "-d", opts.Database)
	args = append(args, "-f", opts.OutputPath)

	if opts.Format != "" {
		args = append(args, "-F", opts.Format)
	}
	if opts.SchemaOnly {
		args = append(args, "-s")
	}
	if opts.DataOnly {
		args = append(args, "-a")
	}
	if opts.Clean {
		args = append(args, "-c")
	}
	if opts.Create {
		args = append(args, "-C")
	}
	if opts.Inserts {
		args = append(args, "--inserts")
	}
	if opts.ColumnInserts {
		args = append(args, "--column-inserts")
	}
	if opts.DisableTriggers {
		args = append(args, "--disable-triggers")
	}
	if opts.Verbose {
		args = append(args, "-v")
	}

	// Emit initial status
	wailsruntime.EventsEmit(ctx, "backup:log", fmt.Sprintf("Starting backup of database '%s' to '%s'...", opts.Database, opts.OutputPath))
	wailsruntime.EventsEmit(ctx, "backup:log", fmt.Sprintf("Running command: %s %s", pgDumpCmd, strings.Join(args, " ")))

	// 3. Prepare exec.Command
	cmd := exec.CommandContext(ctx, pgDumpCmd, args...)

	// Pass PGPASSWORD via environment
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", conn.Password))

	// Redirect stderr to read log lines from pg_dump
	stderr, err := cmd.StderrPipe()
	if err != nil {
		wailsruntime.EventsEmit(ctx, "backup:log", fmt.Sprintf("Error creating stderr pipe: %v", err))
		return err
	}

	// Redirect stdout to read (usually empty if we use -f, but useful)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		wailsruntime.EventsEmit(ctx, "backup:log", fmt.Sprintf("Error creating stdout pipe: %v", err))
		return err
	}

	// Start command
	if err := cmd.Start(); err != nil {
		errStr := fmt.Sprintf("Failed to start pg_dump process: %v\nMake sure PostgreSQL bin directory is correct in Settings.", err)
		wailsruntime.EventsEmit(ctx, "backup:log", errStr)
		wailsruntime.EventsEmit(ctx, "backup:status", map[string]interface{}{
			"status":  "error",
			"message": errStr,
		})
		return err
	}

	// Stream stderr logs
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			wailsruntime.EventsEmit(ctx, "backup:log", line)
		}
	}()

	// Stream stdout logs
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			wailsruntime.EventsEmit(ctx, "backup:log", line)
		}
	}()

	// Wait for completion in a goroutine
	go func() {
		err := cmd.Wait()
		if err != nil {
			errStr := fmt.Sprintf("pg_dump finished with error: %v", err)
			wailsruntime.EventsEmit(ctx, "backup:log", errStr)
			wailsruntime.EventsEmit(ctx, "backup:status", map[string]interface{}{
				"status":  "error",
				"message": errStr,
			})
		} else {
			successStr := "Backup completed successfully!"
			wailsruntime.EventsEmit(ctx, "backup:log", successStr)
			wailsruntime.EventsEmit(ctx, "backup:status", map[string]interface{}{
				"status":  "success",
				"message": successStr,
			})
		}
	}()

	return nil
}
