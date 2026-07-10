# YxPg — Project Context for AI Assistants

## TL;DR
Desktop PostgreSQL client (seperti DBeaver/TablePlus) berbasis **Wails v2** (Go backend + Vue 3 frontend). Tidak ada HTTP API — semua komunikasi via Wails IPC bridge langsung ke method Go.

---

## Stack

| Layer | Tech |
|---|---|
| Desktop framework | [Wails v2](https://wails.io) — WebView2 (Windows) |
| Backend | Go, `pgx/v5` (PostgreSQL), `go-sqlite3` (history) |
| Frontend | Vue 3 (Composition API + `<script setup>`), TypeScript, Vite 6 |
| State | Pinia stores |
| UI | TailwindCSS 3 (dark-only theme, navy + teal) |
| SQL Editor | CodeMirror 6 + `@codemirror/lang-sql` |
| Data Grid | Tabulator 6 |

---

## Struktur Folder

```
YxPg/
├── main.go              # Wails entry point (window config)
├── app.go               # App struct + semua method yang di-expose ke frontend
├── go.mod / go.sum
├── wails.json
├── backend/
│   ├── connection/      # Pool manager (pgxpool), koneksi CRUD, workspace persistence
│   ├── models/          # Shared structs (Connection, QueryResult, SchemaInfo, dll)
│   ├── query/           # Query executor, EXPLAIN, query history (SQLite)
│   ├── schema/          # Schema introspection via information_schema + pg_catalog
│   ├── ddl/             # DDL builder + executor (CREATE/ALTER/DROP)
│   └── export/          # CSV, JSON, SQL INSERT export
└── frontend/
    └── src/
        ├── stores/      # Pinia: connections, tabs, schema, ui, workspace
        ├── views/       # WorkspaceView, QueryView, TableView, DDLView, BuilderView, QueryLogView
        ├── components/  # layout/, connection/, query/, schema/, ddl/, builder/, shared/
        └── utils/       # shortcuts.ts, sql-formatter.ts, type-mapper.ts
```

---

## Alur Data

```
Vue component
  → wailsjs/go/main/App.js (auto-generated bindings)
  → Wails IPC (window['go']['main']['App']['Method'](args))
  → app.go method
  → backend/* packages
  → pgxpool → PostgreSQL
```

Semua binding Go method tersedia di `frontend/wailsjs/go/main/App.js` dan `App.d.ts`. **Jangan edit file ini manual** — di-generate otomatis oleh Wails dari `app.go`.

---

## Database yang Digunakan

| DB | Tujuan | Lokasi |
|---|---|---|
| PostgreSQL (multi) | Target — database yang di-manage | Remote servers |
| SQLite | Query history + saved queries | `~/.yxpg/history.db` |
| JSON file | Connection configs | `~/.yxpg/connections.json` |
| JSON file | Workspace tree | `~/.yxpg/workspace.json` |
| PostgreSQL `yxz` | Sumber sync koneksi dari server internal | `localhost:7733` |
| SQLite pgAdmin | Import server dari pgAdmin 4 | `%APPDATA%\pgAdmin\pgadmin4.db` |

---

## Method Penting di app.go

### Connection
- `GetConnections()` — list semua saved connections
- `AddConnection(conn)` / `UpdateConnection(conn)` / `DeleteConnection(id)`
- `TestConnection(conn)` → latency ms
- `Connect(id)` / `Disconnect(id)` / `GetConnectionStatus(id)`
- `SyncServerConnections()` — sync dari `localhost:7733/yxz/public._server`
- `SyncPgAdminConnections()` — import dari pgAdmin 4 SQLite

### Query
- `ExecuteQuery(connId, sql, timeoutSec)` → `QueryResult`
- `ExplainQuery(connId, sql)` → JSON string
- `BrowseTable(connId, schema, table, opts)` → paginated `QueryResult`
- `GetQueryHistory(connId, limit)` / `BookmarkQuery(id)` / `DeleteHistoryEntry(id)`
- `GetSavedQueries()` / `SaveQuery(name, sql, folder)` / `DeleteSavedQuery(id)`

### Schema
- `GetSchemas(connId)` / `GetTables(connId, schema)` / `GetColumns(connId, schema, table)`
- `GetIndexes(connId, schema, table)` / `GetForeignKeys(connId, schema, table)`
- `GetViews(connId, schema)` / `GetFunctions(connId, schema)`
- `GetSequences(connId, schema)` / `GetTriggers(connId, schema, table)`
- `GetTypes(connId, schema)` / `GetTableDDL(connId, schema, table)`
- `GetFullSchema(connId)` / `GetSearchCatalog()` (semua connections)

### DDL
- `CreateTable(connId, def)` / `AlterTable(connId, schema, table, changes)` / `DropTable(connId, schema, table)`
- `RenameTable(connId, schema, oldName, newName)`
- `CreateIndex(connId, schema, table, idx)` / `DropIndex(connId, schema, idx)`
- `AddForeignKey(connId, schema, table, fk)` / `DropConstraint(connId, schema, table, name)`
- `BuildAlterTableSQL(schema, table, changes)` — preview SQL tanpa eksekusi

### Export & Workspace
- `ExportCSV(connId, sql, filePath, delimiter)` / `ExportJSON(...)` / `ExportSQL(...)`
- `LoadWorkspace()` / `SaveWorkspace(json)` (file `~/.yxpg/workspace.json`)

---

## Tab Types

| Type | View | Buka Via |
|---|---|---|
| `query` | QueryView.vue | Ctrl+Q/T, drag tabel, klik "New Query" |
| `table` | TableView.vue | Context menu "Open Table" di sidebar |
| `builder` | BuilderView.vue | F1 |
| `ddl` | DDLView.vue | F2, context menu "Table Designer" |
| `log` | QueryLogView.vue | Dari status bar / menu |
| `home` | HomeView.vue | Default saat app buka |

Tab di-persist ke `localStorage` — reopen setelah restart.

---

## Fitur Spesifik / Non-obvious

### Date Template
Query dengan `:date_from` atau `:date_to` akan trigger modal `QueryWinDate.vue` sebelum eksekusi. Placeholder diganti jadi `'YYYY-MM-DD'::date`.

### Drag-and-Drop ke Editor
Drag tabel dari sidebar ke CodeMirror → auto-generate `SELECT col1, col2, ... FROM schema.table LIMIT 100`.

### Auto-detect CREATE TABLE
Setelah `CREATE TABLE` berhasil dijalankan, QueryView otomatis reload schema dan tambah node ke workspace tree.

### SyncServerConnections
Pada startup (500ms delay), app connect ke `localhost:7733` DB `yxz` dan baca `public._server` untuk auto-populate connections. Ini khusus untuk environment YuServer. Fail secara silent jika tidak ada.

### Connection Pool Config
MaxConns=10, MinConns=1, MaxConnLifetime=30m, MaxConnIdleTime=5m per connection.

---

## Keyboard Shortcuts

| Shortcut | Aksi |
|---|---|
| `Ctrl+Enter` / `F5` | Run Query |
| `Ctrl+Q` / `Ctrl+T` | New Query Tab |
| `Ctrl+W` | Close Tab |
| `Ctrl+Shift+F` | Format SQL |
| `Ctrl+H` | Toggle History |
| `Ctrl+K` | Search Tables/Views (palette) |
| `Ctrl+R` | Refresh Schema |
| `Ctrl+N` | New Connection |
| `F1` | Builder |
| `F2` | Table Designer |

---

## Dev Setup

```bash
# Install frontend deps
cd frontend && npm install

# Run dev (hot-reload)
wails dev

# Build production
wails build
```

Window: 700×700px, mulai maximized, background `#0A0F1E`.
