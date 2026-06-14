# DEV-AUDIT-TRACKER

## Technical Flow (runtime)

This document describes the runtime flow when you start the application with `go run main.go`. It explains what each major component does and where you should inspect logs or DB when debugging.

### 1) Application startup
- `main.go` configures logging (output to `logs.txt`), connects to MySQL using DSN `root:@tcp(127.0.0.1:3307)/dev_audit_db`, and runs `AutoMigrate` for `Project` and `ProjectHistory`.

### 2) Workspace initialization
- `main.go` ensures `workspacePath` exists (default configured: `/Users/bernatdsitumeang/Desktop/workspace-Bernatdev/`). It creates the directory if missing.

### 3) Background watcher goroutine
- A goroutine loops every 5 seconds calling `services.WatchWorkspace(db, workspacePath)`.
- `WatchWorkspace` enumerates immediate subdirectories under the workspace; each subdirectory name is interpreted as a `ticketID` (project identifier).
- For each `ticketID`, the watcher:
  1. Finds an existing `Project` record by `ticket_number`. If not found, it creates one.
  2. Calls `checkSubFolder(db, &project, ticketPath)` to inspect specific subfolders for trigger files.

### 4) Subfolder inspection logic (`checkSubFolder`)
- `FSD/*.pdf` → if at least 1 PDF found and `Project.HasFSD == false` then set `HasFSD=true` and call `UpdateStatus(..., "KICK_OFF", "File FSD terdeteksi")`.
- `Analysis/*.xls*` → if found and `HasFSD==true` and `HasAnalysis==false` then set `HasAnalysis=true` and call `UpdateStatus(..., "DEVELOPMENT_START", "File Analisis terdeteksi")`.
- `Document/*DONE*` → if 4+ matches and `HasAnalysis==true` and `DocCount < 4` then set `DocCount` and call `UpdateStatus(..., "DEV_DONE", "4 Dokumen DONE terdeteksi")`.
- `Revisi/revisi.txt` → if exists and `Status != "DEV_REVISION_UAT"` then call `UpdateStatus(..., "DEV_REVISION_UAT", "File revisi.txt ditemukan")`.

### 5) Persisting state and history (`UpdateStatus`)
- Updates the `projects` row with new status and timestamps using `db.Save`.
- Inserts an audit row into `project_histories` with `ProjectID`, `Status`, `Notes`, and `CreatedAt` via `db.Create`.
- Errors during `Save` or `Create` are logged to `logs.txt`. The code checks `project.ID` to avoid inserting history with `project_id = 0`.

### 6) Web server and UI
- `routes.InitRoutes(db)` registers HTTP handlers and starts the HTTP server on port `:8080`.
- `GET /` is handled by `controllers.ProjectController.Dashboard`, which queries projects and renders `views/index.html`.
- The controller populates a map with key `Projects` (slice of `models.Project`) so the template can iterate with `{{ range .Projects }}`.
- If the slice is empty the template shows: `Belum ada project terdeteksi di workspace.`

### 7) Frontend logging
- The UI includes `logFrontend(message)` script that POSTs `{ message }` to `POST /log`.
- `ProjectController.LogFrontend` writes incoming messages to `logs.txt` prefixed with `[FRONTEND]`.

### 8) End-to-end manual test (quick)
1. Start app:
```bash
cd /Users/bernatdsitumeang/Desktop/DevAuditTracker
go run main.go
```
2. Create a test ticket and FSD file:
```bash
mkdir -p /Users/bernatdsitumeang/Desktop/workspace-Bernatdev/#TEST1/FSD
touch /Users/bernatdsitumeang/Desktop/workspace-Bernatdev/#TEST1/FSD/example.pdf
```
3. Wait 5–10s and check `logs.txt` for lines like `Created new project` or `Updating status for project` and visit `http://localhost:8080`.

### 9) Troubleshooting pointers
- If UI shows "Belum ada project..." but DB has rows:
  - Verify the controller uses key `Projects` in the data map.
  - Check `logs.txt` for errors and for evidence that the watcher found and created the project.
  - Confirm the MySQL DSN port (`3307`) matches the DB you're inspecting.
- If you see duplicate-entry or foreign-key errors:
  - Inspect `projects` for duplicate `ticket_number` entries and `project_histories` referencing non-existent `project_id`s.
  - Remove inconsistent rows or reset the tables, then restart.

---

This section should help you learn how the Go app flows at runtime and where to look when testing or debugging.
