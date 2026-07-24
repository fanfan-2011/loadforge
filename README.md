# ⚡ LoadForge

> A modern ApacheBench alternative — high-performance HTTP/HTTPS load testing tool

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![Discord](https://img.shields.io/badge/Discord-Join-5865F2?style=flat&logo=discord)](https://discord.gg/zqQ6rdYT)

LoadForge is a high-performance, cross-platform, **single-binary** HTTP/HTTPS load testing tool — a modern replacement for ApacheBench (`ab`). It features an async concurrent engine, comprehensive latency percentiles (P50–P999), and an embedded Vue3 + ECharts web visualization dashboard.

[🇨🇳 中文版](README.zh-CN.md)

---

## ✨ Features

- **🚀 High-Performance Engine** — goroutine + channel async non-blocking I/O, HTTP/1.1 (extensible to HTTP/2, HTTP/3)
- **📊 Full Statistics** — request count, success/fail rate, QPS, throughput, latency percentiles (P50/P75/P90/P95/P99/P999)
- **📈 Web Dashboard** — auto-starts embedded Vue3 + ECharts on port 8899 after each test
- **📂 Persistent Storage** — test results saved to `~/.loadforge/tests/`, JSON exportable
- **🔧 CLI-First** — familiar `ab`-like CLI with flags: `-n`, `-c`, `-t`, `-m`, `-H`, `-d`, `--json`, `--report`
- **💡 Performance Tips** — auto-detects P99 latency spikes and other bottlenecks
- **📦 Single Binary** — compiled with Go, zero runtime dependencies

---

## 🚀 Quick Start

### Build from source

```bash
# 1. Clone the repo
git clone https://github.com/fanfan-2011/loadforge.git
cd loadforge

# 2. Build the web frontend
cd report/ui
npm install
npm run build
cd ../..

# 3. Compile the Go binary
go build -ldflags="-s -w" -o loadforge .

# 4. Run!
./loadforge bench -n 10000 -c 100 https://example.com
```

> 💡 Or download a pre-built binary from the [Releases](https://github.com/fanfan-2011/loadforge/releases) page.

### Basic usage

```bash
# 1000 requests with 10 concurrent connections
loadforge bench -n 1000 -c 10 https://example.com

# 30-second duration test
loadforge bench -t 30s -c 200 https://example.com

# POST request with custom headers
loadforge bench -n 5000 -c 50 -m POST -H "Content-Type: application/json" -d '{"key":"value"}' https://api.example.com

# JSON output (pipe-friendly)
loadforge bench -n 10000 -c 100 --json https://example.com

# Auto-start web dashboard after test
loadforge bench -n 50000 -c 200 --report https://example.com
```

---

## 📋 CLI Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-n`  | Total number of requests | `1000` |
| `-c`  | Number of concurrent connections | `10` |
| `-t`  | Test duration (e.g. `30s`, `1m`) | (disabled) |
| `-m`  | HTTP method | `GET` |
| `-H`  | Custom request header (repeatable) | — |
| `-d`  | Request body | — |
| `--json` | Output as JSON | `false` |
| `--report` | Start web report server | `false` |
| `--timeout` | Request timeout in seconds | `30` |

---

## 📊 Sample Output

```
╔══════════════════════════════════════════╗
║        LoadForge Benchmark Report         ║
╚══════════════════════════════════════════╝

📊 Requests
   Total:          10000
   Success:        9982
   Failed:         18  (0.18%)

⚡ Performance
   QPS:            5421.33 req/s
   Throughput:     48.23 MB/s

⏱️  Latency (ms)
   Min:            0.82
   Avg:            18.44
   P50:            14.27
   P95:            45.31
   P99:            128.67
   P999:           256.43

🔢 Status Codes
   200:            9982 (99.82%)
   503:            18 (0.18%)
```

---

## 🌐 Web Dashboard

Use the `--report` flag to launch the web dashboard:

```bash
loadforge bench -n 100000 -c 500 --report https://example.com
```

After the test completes, open `http://<LAN-IP>:8899` in your browser. The dashboard includes:

- **Test List** — all historical tests at a glance
- **Overview Cards** — QPS, success rate, throughput
- **Charts** — QPS curve, latency curve, status code distribution, error breakdown
- **Latency Analysis** — P50 / P75 / P90 / P95 / P99 / P999
- **Details** — HTTP config, headers, body size, traffic stats
- **Performance Tips** — automatic bottleneck detection and suggestions

---

## 🗂️ Data Storage

Test results are saved to `~/.loadforge/tests/<timestamp>/`:

```
~/.loadforge/tests/
├── 1784900665705/
│   ├── config.json      # test configuration
│   ├── result.json      # full statistics
│   └── timeline.json    # per-second timeline data
├── 1784900667939/
│   └── ...
└── index.json           # test index
```

---

## 🏗️ Project Structure

```
loadforge/
├── main.go                    # entry point
├── cmd/
│   ├── root.go                # root CLI command
│   └── bench.go               # bench subcommand
├── engine/
│   └── engine.go              # async HTTP benchmark engine
├── stats/
│   └── stats.go               # statistics (percentiles, etc.)
├── storage/
│   └── storage.go             # file-based storage
├── report/
│   ├── web.go                 # web report server
│   ├── embed.go               # embedded frontend assets
│   └── ui/                    # Vue3 + ECharts frontend
│       └── src/
│           ├── App.vue
│           └── components/
│               ├── TestList.vue
│               └── TestDetail.vue
├── README.md                  # this file (English)
├── README.zh-CN.md            # Chinese version
└── go.mod
```

---

## 🛠️ Tech Stack

| Layer | Technology |
|-------|-----------|
| **Backend** | Go 1.18+ (net/http, goroutines, channels) |
| **CLI** | Cobra (spf13/cobra) |
| **Frontend** | Vue 3 + Vite |
| **Charts** | ECharts 5 |
| **Embedding** | Go `embed` package (single binary) |
| **Storage** | JSON filesystem |

---

## 🤝 Community

- 💬 **Discord**: [Join the discussion](https://discord.gg/zqQ6rdYT)
- 🐛 **Issues**: [GitHub Issues](https://github.com/fanfan-2011/loadforge/issues)
- ⭐ **Star**: If you find this useful, please give it a star!

---

## 📄 License

MIT License © 2026 [fanfan-2011](https://github.com/fanfan-2011)
