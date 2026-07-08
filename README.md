# YxPg - Portable Postgres Studio

YxPg adalah aplikasi desktop GUI untuk PostgreSQL yang dibangun dengan Wails, Go, Vue 3, dan TypeScript. Aplikasi ini menyediakan workspace database, browser schema, editor query, query history, DDL tools, dan ekspor data dalam satu antarmuka desktop.

![Visualisasi Skema Database](frontend/public/database_schema_viz.png)

![Screenshot Aplikasi](frontend/public/sc1.png)

## Fitur

- Manajemen koneksi PostgreSQL: tambah, ubah, hapus, test, connect, dan disconnect.
- Workspace lokal untuk mengelompokkan koneksi dan objek database.
- Browser schema untuk melihat schema, table, view, function, sequence, trigger, type, index, foreign key, dan column detail.
- SQL editor dengan CodeMirror, SQL formatting, autocomplete, dan eksekusi query.
- Query history dan saved queries berbasis SQLite lokal.
- Table browser dengan pagination, sorting, dan filter.
- DDL tools untuk membuat/mengubah table, index, foreign key, rename/drop table, dan menjalankan raw DDL.
- Export hasil query ke CSV, JSON, dan SQL insert statement.
- Sinkronisasi koneksi dari central database `public._server` dan import koneksi pgAdmin 4.

## Teknologi

- Wails v2
- Go 1.22+
- Vue 3
- TypeScript
- Vite
- Tailwind CSS
- Pinia
- CodeMirror
- pgx PostgreSQL driver
- SQLite untuk query history lokal

## Prasyarat

Pastikan sudah terpasang:

- Go 1.22 atau lebih baru
- Node.js dan npm
- Wails CLI
- PostgreSQL server yang dapat diakses

Instal Wails CLI jika belum ada:

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

Verifikasi environment Wails:

```bash
wails doctor
```

## Menjalankan Development

Clone repository:

```bash
git clone <repository-url>
cd YxPg
```

Install dependency frontend:

```bash
cd frontend
npm install
cd ..
```

Jalankan aplikasi dalam mode development:

```bash
wails dev
```

## Build Aplikasi

Build aplikasi desktop:

```bash
wails build
```

Output build akan tersedia di:

```text
build/bin/
```

Untuk Windows, executable utama biasanya:

```text
build/bin/YxPg.exe
```

## Konfigurasi

YxPg menyimpan data lokal pengguna di folder home:

```text
~/.yxpg/
```

File yang digunakan:

- `connections.json` untuk daftar koneksi tersimpan.
- `workspace.json` untuk struktur workspace.
- `history.db` untuk query history dan saved queries.

File `yxpg.conf` bersifat konfigurasi lokal dan tidak perlu dipublish. Gunakan `yxpg.conf.example` jika ingin menyediakan contoh konfigurasi publik.

Contoh minimal `yxpg.conf`:

```ini
# Central database untuk sinkronisasi daftar server, opsional
central_dsn=postgres://user:password@localhost:5432/database?sslmode=disable

# Default form koneksi baru, opsional
default_host=localhost
default_port=5432
default_database=postgres
default_username=postgres
default_sslmode=disable
```

## Struktur Proyek

```text
.
|-- app.go                  # Binding utama Wails dan API backend
|-- main.go                 # Entry point aplikasi Wails
|-- backend/
|   |-- connection/         # Store, manager, workspace, helper koneksi
|   |-- ddl/                # Builder dan executor DDL
|   |-- export/             # Export CSV, JSON, SQL
|   |-- models/             # Model data backend
|   |-- query/              # Query executor, explain, history
|   `-- schema/             # Inspector metadata PostgreSQL
|-- frontend/
|   |-- src/                # Vue app
|   |-- public/             # Icon dan asset publik
|   `-- wailsjs/            # Binding hasil generate Wails
|-- build/                  # Output build Wails
|-- go.mod
`-- wails.json
```

## Script Frontend

Jalankan dari folder `frontend`:

```bash
npm run dev
npm run build
npm run preview
```

Dalam penggunaan normal, jalankan aplikasi melalui Wails:

```bash
wails dev
```

## Catatan Keamanan

- Jangan commit `yxpg.conf`, credential database, file `.env`, atau file konfigurasi lokal lain.
- Jangan publish isi folder `~/.yxpg/` karena dapat berisi koneksi dan history query.
- Pastikan binary hasil build dan dependency lokal tidak masuk repository.
- Review koneksi default, screenshot, dan contoh konfigurasi sebelum repository dibuat public.

## Checklist Sebelum Public

- Pastikan `.gitignore` sudah mengecualikan `frontend/node_modules/`, `frontend/dist/`, `build/`, `*.exe`, dan `yxpg.conf`.
- Hapus binary lokal seperti `yxpg.exe` jika sempat ter-track oleh Git.
- File `LICENSE` (MIT License) sudah ditambahkan untuk distribusi publik.
- Tambahkan `yxpg.conf.example` jika pengguna lain perlu contoh konfigurasi.
- Jalankan build terakhir:

```bash
wails build
```

## Lisensi

Proyek ini dilisensikan di bawah [MIT License](LICENSE) - bebas untuk digunakan, dimodifikasi, dan didistribusikan secara gratis.

