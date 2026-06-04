# AGENTS.md

## Project Overview

115Quick is a Go server for 115 cloud storage management. It provides a complete API for cloud file browsing, offline downloading, video playback with proxy streaming, and full OAuth token lifecycle management. Token is obtained via refreshToken input from the Chrome extension plugin.

## Build & Run

```bash
go build -o 115quick.exe .              # Windows
go build -o 115quick .                  # Linux/Mac
./115quick -f etc/quick.yaml            # Run with config
./115quick --version                    # Show version
```

## Architecture

Framework: **Gin** (HTTP server)

```
quick.go (entrypoint - server bootstrap, graceful shutdown)
  ├─ internal/config/config.go          (Config struct, YAML load/save)
  ├─ internal/handler/handler.go        (HTTP handlers, route registration)
  ├─ internal/service/
  │   ├─ auth/token.go                  (OAuth token manager: refresh, store, status)
  │   ├─ auth/client.go                 (115 API HTTP client with auto-retry)
  │   ├─ file/service.go                (File management: list, search, CRUD, download URL)
  │   ├─ video/service.go               (Video: play info, subtitles, history, transcode, proxy stream)
  │   ├─ cloudtask/service.go           (Cloud download: task list, add URLs/BT, delete, quota)
  │   └─ download/engine.go             (Local download engine with progress tracking)
  ├─ internal/store/store.go            (SQLite storage: token, task history, file info, config)
  ├─ internal/version/version.go        (Version info, GitHub release check)
  └─ etc/quick.yaml                     (Config file)
```

## Key Design Decisions

- **No go-zero dependency**: Uses Gin framework directly, no code generation required
- **Service-based architecture**: Each domain (auth, file, video, cloudtask, download) has its own service
- **Token management**: refreshToken is input via Chrome plugin POST to `/api/setToken`, then auto-refreshed
- **Video proxy streaming**: Server proxies 115 video streams to avoid CORS and auth issues
- **SQLite storage**: Pure-Go SQLite (no CGO), 4 tables: token_config, task_history, file_infos, config
- **Unified API response**: All endpoints return `{"code": 200, "data": {...}}` format
- **CORS wide open**: `Access-Control-Allow-Origin: *` for Chrome extension compatibility

## Config

`etc/quick.yaml`:
```yaml
Name: 115quick
Host: 0.0.0.0
Port: 8889
DBPath: data/115quick.db
Auth115:
  DownloadPath: data
  AccessToken: ""
  RefreshToken: ""
Log:
  Level: info
  Format: text
```

## API Routes

All under prefix `/v1/Download/api`:

### Server
- `GET /getServerInfo` - server status, version, update check

### Token Management
- `GET /getTokenStatus` - token status (configured, valid, expiresAt, message)
- `POST /setToken` - set refreshToken (body: `{"refresh_token": "..."}`)

### File Management
- `GET /files/list?cid=&offset=&limit=&type=` - list cloud files
- `GET /files/info?file_id=` - get file/folder details
- `GET /files/search?keyword=&offset=&limit=` - search files
- `POST /files/createFolder` - create folder (body: `{"pid", "file_name"}`)
- `POST /files/delete` - delete files (body: `{"file_ids"}`)
- `POST /files/rename` - rename file (body: `{"file_id", "file_name"}`)
- `POST /files/move` - move files (body: `{"file_ids", "to_cid"}`)
- `GET /files/download?pick_code=` - get download URL

### Video Playback
- `GET /video/play?pick_code=` - get video play info (URLs, definitions, tracks)
- `GET /video/subtitles?pick_code=` - get subtitle list
- `GET /video/history?pick_code=` - get play progress
- `POST /video/history` - save play progress (body: `{"pick_code", "time", "watch_end"}`)
- `POST /video/transcode` - submit transcode (body: `{"pick_code", "op"}`)
- `GET /video/proxy?pick_code=&definition=` - proxy video stream (supports Range)

### Cloud Download
- `GET /cloud/tasks?page=` - get cloud task list
- `POST /cloud/addUrls` - add download URLs (body: `{"urls", "wp_path_id"}`)
- `POST /cloud/addBt` - add BT task (body: `{"info_hash", "wanted", "save_path", ...}`)
- `POST /cloud/delete` - delete task (body: `{"info_hash", "delete_source_file"}`)
- `POST /cloud/clear` - clear tasks (body: `{"flag"}` 0=completed, 1=all, 2=failed, 3=running)
- `GET /cloud/quota` - get download quota info

### Local Download
- `GET /download/progress` - get all download progress
- `POST /download/start` - start file download (body: `{"url", "file_name", "file_size"}`)

## Token Lifecycle

1. User inputs `refresh_token` via Chrome extension -> `POST /api/setToken`
2. Server immediately calls 115 API to exchange for `access_token`
3. Token stored in SQLite (`token_config` table)
4. Auto-refresh: when token expires within 10 minutes, refresh automatically
5. Force-refresh: on API error codes 401, 40002, 40140126
6. All 115 API calls use `Authorization: Bearer {access_token}` header

## 115 API Base URLs

- File/Download/Video: `https://proapi.115.com`
- Auth/Token: `https://passportapi.115.com`

## Chrome Extension (Manifest V3)

**Features**:
- Left click opens full management panel in new tab
- Right click context menu: open panel, open settings
- Page magnet link detection with push button
- Server version display with update notification

**Pages**:
- Dashboard: server status, token status, quick add, cloud tasks
- Token Config: set 115 refreshToken
- Download List: current downloads and pending queue
- Task History: historical download records
- Settings: download mode, server info
- Cloud Files: browse and manage 115 cloud files with video playback

## macOS Installation

```bash
curl -fsSL https://raw.githubusercontent.com/ddc-111/115quick/main/scripts/install-macos.sh | bash
```

## CI/CD

GitHub Actions workflow (`.github/workflows/release.yml`):
- Triggered on version tags (`v*`)
- Builds server binaries for: windows/amd64, linux/amd64, linux/arm64, darwin/amd64, darwin/arm64
- Builds Chrome extension from `plugin/115Quick`
- Creates GitHub Release with all artifacts
- Version injected via `-ldflags` during build

## Version Management

```bash
go build -ldflags="-X '115Quick_server/internal/version.Version=x.x.x'" .
```

## Git Configuration

- Remote: `git@github.com:ddc-111/115quick.git` (SSH)
- Branch: `main`
