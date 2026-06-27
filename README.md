# ⚡ Memedovski Fast Download Manager

<div align="center">
  <h3>Ultra-fast parallel download manager with GUI & CLI</h3>
  <p>İngilizce / Türkçe ve diğer diller için çoklu dil desteği</p>
</div>

---

## 🚀 Features

- **Parallel Download** — Splits files into chunks and downloads simultaneously using HTTP Range requests
- **Work Stealing** — Dynamic chunk distribution: faster workers automatically pick up remaining chunks
- **Dual Mode** — GUI (native Windows) and CLI from the same binary
- **Auto File Detection** — Automatically detects filename from `Content-Disposition` header, URL path, or `Content-Type` extension mapping (60+ MIME types)
- **Multi-Language** — English & Turkish interface (dropdown switch)
- **High Performance** — HTTP/2, connection pooling (200 idle conns), 1MB buffers, TCP keep-alive
- **Resilient** — 5 retry attempts per chunk with exponential backoff
- **Real-Time Stats** — Progress bar, speed, ETA, elapsed time per download
- **Status Bar** — Total downloaded, active download count, global speed

## 📥 Download

Download the latest `FastDownloader.exe` from [Releases](https://github.com/devmehmetemre/Memedovski-Fast-Download-Manager/releases).

**Requirements:** Windows 7+ (64-bit). No installation needed, runs standalone.

## 🖥️ Usage

### GUI Mode
Double-click `FastDownloader.exe`, paste URL, adjust connection count (1–500), click **Download**.

### CLI Mode
```cmd
FastDownloader.exe <url> [connections]
```

**Examples:**
```cmd
FastDownloader.exe https://example.com/file.zip
FastDownloader.exe https://example.com/bigfile.iso 200
```

## ⚙️ How It Works

```
┌─────────────────────────────────────────────────────┐
│  1. HEAD request → get file size + check Range      │
│  2. Split file into N×3 chunks (work stealing)      │
│  3. Spawn M parallel goroutines (connection count)  │
│  4. Each picks next available chunk from queue      │
│  5. WriteAt() directly to .part file (no merge)     │
│  6. Rename .part → filename when complete           │
└─────────────────────────────────────────────────────┘
```

### Architecture
- **Shared Transport** — HTTP/2 enabled, 200 max idle connections, TCP keep-alive
- **Direct Disk Writes** — Every chunk writes directly to its byte offset in the `.part` file using `WriteAt()` — no temporary chunk files, no merge phase
- **Work Stealing** — 3× more chunks than workers ensures fast workers don't sit idle
- **Write Worker** — Single goroutine serializes disk writes from all chunk downloads via channel

## 📊 Benchmarks

| File Size | Connections | Time | Speed |
|-----------|-----------|------|-------|
| 10 MB     | 50        | 0.7s | 14 MB/s |
| 100 MB    | 100       | 2.4s | 41 MB/s |

*(Tested on 200 Mbps connection with proof.ovh.net)*

## 🛠️ Building from Source

```cmd
go build -ldflags="-s -w -H windowsgui" -o FastDownloader.exe .
```

For CLI-only build (without GUI):
```cmd
go build -ldflags="-s -w" -o FastDownloader.exe .
```

**Dependencies:**
- Go 1.22+
- MinGW-w64 (for CGO / Walk GUI)
- `github.com/lxn/walk` (native Windows GUI toolkit)

## 🧩 Project Structure

```
FastDownloadManager/
├── main.go      # Entry point (GUI/CLI dispatch)
├── core.go      # Download engine, transport, helpers
├── cli.go       # CLI mode implementation
├── gui.go       # GUI mode (Walk Framework)
├── lang.go      # Multi-language strings (EN/TR)
└── go/          # Portable Go runtime
```

## 🌐 Language Support

Switch between English and Türkçe from the GUI dropdown. CLI output adapts to the selected language.

To add a new language:
1. Add a new `LangID` const in `lang.go`
2. Add a `LangPack` entry in the `langs` map
3. Add a `case` in the `T()` function

## 📝 License

MIT
