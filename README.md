# ⚡ LoadForge

> A modern ApacheBench alternative — high-performance HTTP/HTTPS load testing tool

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![Discord](https://img.shields.io/badge/Discord-Join-5865F2?style=flat&logo=discord)](https://discord.gg/zqQ6rdYT)
[![Release](https://img.shields.io/github/v/release/fanfan-2011/loadforge?style=flat)](https://github.com/fanfan-2011/loadforge/releases)

LoadForge is a high-performance, cross-platform, **single-binary** HTTP/HTTPS load testing tool — a modern replacement for ApacheBench (`ab`). It features an async concurrent engine, comprehensive latency percentiles (P50–P999), and an embedded Vue3 + ECharts web visualization dashboard.

**[🇨🇳 中文版](README.zh-CN.md)**

---

## 📦 Table of Contents

- [Features](#-features)
- [Environment Requirements](#-environment-requirements)
- [One-Click Install](#-one-click-install)
- [Quick Start](#-quick-start)
- [Usage Guide](#-usage-guide)
- [CLI Reference](#-cli-reference)
- [Web Dashboard](#-web-dashboard)
- [Sample Output](#-sample-output)
- [Real-World Scenarios](#-real-world-scenarios)
- [Data Storage](#-data-storage)
- [Project Structure](#-project-structure)
- [Tech Stack](#-tech-stack)
- [Troubleshooting](#-troubleshooting)
- [Community](#-community)
- [License](#-license)

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| **🚀 High-Performance Engine** | goroutine + channel async non-blocking I/O, HTTP/1.1 with HTTP/2 & HTTP/3 ready |
| **📊 Full Statistics** | request count, success/fail rate, QPS, throughput, latency percentiles |
| **📈 Web Dashboard** | auto-starts embedded Vue3 + ECharts on port 8899 after each test |
| **📂 Persistent Storage** | all test results saved to `~/.loadforge/tests/`, JSON exportable |
| **🔧 CLI-First** | familiar `ab`-like experience with `-n`, `-c`, `-t`, `-m`, `-H`, `-d`, `--json`, `--report` |
| **💡 Performance Tips** | automatic P99 spike detection and bottleneck analysis |
| **📦 Single Binary** | compiled with Go, embed web UI — zero runtime dependencies |
| **🌍 Cross-Platform** | runs on Linux, macOS, and Windows (build from source) |
| **♾️ Unlimited Scale** | tested with 100K+ requests and 1000+ concurrent connections |

---

## 🖥️ Environment Requirements

| Requirement | Runtime (binary) | Build from source |
|-------------|------------------|-------------------|
| **OS** | Linux, macOS, Windows | Linux, macOS, Windows |
| **CPU** | x86-64 or ARM64 | x86-64 or ARM64 |
| **RAM** | 128 MB | 512 MB |
| **Disk** | 10 MB | 100 MB |
| **Go** | Not required | 1.18+ |
| **Node.js** | Not required | 16+ |

> Pre-built binaries include the embedded web UI — **no Go or Node.js needed** at runtime.

---

## ⚡ One-Click Install

Install LoadForge anywhere on Linux/macOS with a single command:

```bash
curl -fsSL https://github.com/fanfan-2011/loadforge/raw/main/install.sh | bash
```

After installation, verify:

```bash
loadforge --help
loadforge bench -n 100 -c 10 https://example.com
```

> 🪟 **Windows users**: Install WSL + Go + Node.js, then use the build-from-source instructions below.

---

## 📥 Installation

### Option 1: Download Pre-built Binary (Recommended for Linux/macOS)

Download the latest binary from the [Releases page](https://github.com/fanfan-2011/loadforge/releases).

```bash
# Example for Linux x86-64
# Replace VERSION with the latest tag (e.g., v1.0.0)
curl -LO https://github.com/fanfan-2011/loadforge/releases/latest/download/loadforge-linux-amd64.tar.gz
tar -xzf loadforge-linux-amd64.tar.gz
sudo mv loadforge /usr/local/bin/
loadforge --help
```

> 💡 Pre-built binaries include the embedded web UI — **no Node.js or Go required** at runtime.

### Option 2: Build from Source (Universal)

This works on any platform with Go and Node.js installed.

#### Step 1: Install Prerequisites

<details>
<summary><b>Linux (Debian/Ubuntu)</b></summary>

```bash
sudo apt update
sudo apt install -y golang-go nodejs npm git

go version  # expect go 1.18+
node --version  # expect v16+
npm --version
```
</details>

<details>
<summary><b>macOS (Homebrew)</b></summary>

```bash
brew install go node git

go version  # expect go 1.18+
node --version  # expect v16+
npm --version
```
</details>

<details>
<summary><b>Windows</b></summary>

1. Install **Go**: https://go.dev/dl/ (1.18+)
2. Install **Node.js**: https://nodejs.org/ (16+)
3. Install **Git**: https://git-scm.com/
4. Open **PowerShell** or **cmd**

```powershell
go version
node --version
npm --version
git --version
```
</details>

#### Step 2: Clone and Build

```bash
# Clone the repository
git clone https://github.com/fanfan-2011/loadforge.git
cd loadforge

# Build the web frontend (Vue3 + ECharts)
cd report/ui
npm install
npm run build
cd ../..

# Build the Go binary
go build -ldflags="-s -w" -o loadforge .

# Verify
./loadforge version  # or ./loadforge --help
```

> 🐳 **Docker alternative**: You can also use the included Dockerfile (if available) or run inside a Docker container:
> ```bash
> docker run --rm -v $(pwd):/app -w /app golang:1.21 sh -c "
>   apt update && apt install -y nodejs npm &&
>   cd report/ui && npm install && npm run build && cd ../.. &&
>   go build -ldflags='-s -w' -o loadforge .
> "
> ```

#### Step 3 (Optional): Install to System PATH

```bash
# Linux / macOS
sudo mv loadforge /usr/local/bin/

# Windows (PowerShell as Admin)
# Move-Item .\loadforge.exe C:\Windows\System32\
```

### Option 3: Install via `go install`

If you already have Go 1.18+ and Go modules enabled (you can skip building the frontend, but the web UI won't work):

```bash
go install github.com/fanfan-2011/loadforge@latest

# Note: `go install` builds the CLI only — for the web dashboard,
# use the source build method above.
```

### Verify Installation

```bash
loadforge --help
```

Expected output:

```
LoadForge is a high-performance, cross-platform HTTP/HTTPS load testing tool...

Usage:
  loadforge [command]

Available Commands:
  bench       Execute an HTTP/HTTPS benchmark
  help        Help about any command
```

---

## 🚀 Quick Start

### Test a Local Server

```bash
# Start a simple test server (Python)
python3 -m http.server 8888 &

# Run LoadForge against it
loadforge bench -n 100 -c 10 http://localhost:8888/
```

### Test Any Public URL

```bash
# Lightweight test (10 requests, 2 concurrent)
loadforge bench -n 10 -c 2 https://example.com

# Standard test (1000 requests, 10 concurrent)
loadforge bench -n 1000 -c 10 https://example.com

# Heavy load test (100K requests, 500 concurrent)
loadforge bench -n 100000 -c 500 https://your-api.com

# Duration-based test (30 seconds, 100 concurrent)
loadforge bench -t 30s -c 100 https://example.com
```

---

## 📖 Usage Guide

### 1️⃣ Basic GET Test

```bash
loadforge bench -n 5000 -c 50 https://api.example.com/users
```

### 2️⃣ POST with JSON Body

```bash
loadforge bench \
  -n 10000 \
  -c 100 \
  -m POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"username":"test","email":"test@example.com"}' \
  https://api.example.com/users
```

### 3️⃣ Duration-Based Test (Recommend for Long Runs)

```bash
# Run for 5 minutes, 200 concurrent connections
loadforge bench -t 5m -c 200 https://example.com
```

### 4️⃣ JSON Output for Scripting

```bash
# Pipe to jq for processing
loadforge bench -n 10000 -c 100 --json https://example.com | jq '.qps, .latency.p99_ms'

# Save to file for later analysis
loadforge bench -n 5000 -c 50 --json https://example.com > result.json
```

### 5️⃣ Test with Web Dashboard

```bash
loadforge bench -n 50000 -c 200 --report https://example.com
# Open http://<your-LAN-IP>:8899 in browser
```

### 6️⃣ Simulate Realistic Traffic Patterns

```bash
# Mix different endpoints by running in parallel
loadforge bench -n 30000 -c 150 -m GET https://api.example.com/items &
loadforge bench -n 20000 -c 100 -m POST -d '{"action":"search"}' https://api.example.com/search &
wait
```

---

## 📋 CLI Reference

### `loadforge bench`

```
Usage:
  loadforge bench [flags] <url>

Flags:
  -n, --requests int        Total number of requests (default 1000)
  -c, --concurrency int     Number of concurrent connections (default 10)
  -t, --duration string     Test duration (e.g. "30s", "1m", "5m")
  -m, --method string       HTTP method (default "GET")
  -H, --header strings      Custom header, repeatable (e.g. -H "Key: Value")
  -d, --body string         Request body content
      --json                Output results as JSON
      --report              Auto-start web report server on port 8899
      --timeout int         Request timeout in seconds (default 30)
  -h, --help                Help for bench
```

### Flag Details

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-n` | int | `1000` | Total requests. Use with `-c` to control concurrency. |
| `-c` | int | `10` | Number of concurrent workers. Each worker runs in its own goroutine. |
| `-t` | string | — | Overrides `-n`. Runs for specified duration (`30s`, `1m`, `5m`). |
| `-m` | string | `GET` | HTTP method: `GET`, `POST`, `PUT`, `DELETE`, `PATCH`, `HEAD`, `OPTIONS`. |
| `-H` | string[] | — | Pass multiple: `-H "Accept: application/json" -H "X-API-Key: abc"` |
| `-d` | string | — | Body for POST/PUT. Can be JSON, form data, or raw text. |
| `--json` | bool | `false` | Suppresses text output, prints JSON to stdout. |
| `--report` | bool | `false` | Starts web server on `:8899`. Blocks until Ctrl+C. |
| `--timeout` | int | `30` | Per-request timeout. Increase for slow endpoints. |

### Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Test completed successfully |
| `1` | CLI error (invalid args, missing URL) |
| `2` | Runtime error (engine panic, storage failure) |

---

## 🌐 Web Dashboard

The web dashboard provides rich visual analysis of your test results.

### Starting the Dashboard

```bash
# Option A: Start automatically with --report
loadforge bench -n 100000 -c 500 --report https://example.com

# Option B: View historical results by starting a new test with --report
```

### Dashboard Pages

#### Test List
Shows all historical tests with:
- Test time, target URL, request count
- QPS (Queries Per Second)
- Status (completed)

#### Single Test Report

| Section | Content |
|---------|---------|
| **Overview Cards** | Total requests, QPS, success count, failures, avg latency, throughput |
| **Summary** | Target URL, test time, duration, request config |
| **QPS Chart** | Per-second QPS line chart during the test |
| **Latency Chart** | Per-second average latency trend |
| **Status Code Chart** | Pie chart of HTTP status code distribution |
| **Error Chart** | Breakdown of error types (timeouts, connection failures) |
| **Latency Table** | P50, P75, P90, P95, P99, P999, Min, Max |
| **Details** | HTTP method, headers, body size, response size, upload/download traffic |
| **Performance Tips** | Auto-generated suggestions for bottlenecks |

---

## 📊 Sample Output

### Standard CLI Output

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
   Max:            512.00
   P50:            14.27
   P75:            22.18
   P90:            35.66
   P95:            45.31
   P99:            128.67
   P999:           256.43

🔢 Status Codes
   200:            9982 (99.82%)
   503:            18   (0.18%)

❌ Errors
   timeout:        15   (83.33%)
   connection refused: 3 (16.67%)

💡 Performance Tips
   ⚠️  P99 latency spike detected — server may be bottlenecked
   - Possible causes: high server load, slow database, insufficient connection pool
```

### JSON Output Example

```json
{
  "total_requests": 10000,
  "success_count": 9982,
  "fail_count": 18,
  "fail_rate": 0.0018,
  "qps": 5421.33,
  "throughput_mb_s": 48.23,
  "latency": {
    "min_ms": 0.82,
    "avg_ms": 18.44,
    "max_ms": 512.0,
    "p50_ms": 14.27,
    "p95_ms": 45.31,
    "p99_ms": 128.67,
    "p999_ms": 256.43
  },
  "status_codes": { "200": 9982, "503": 18 },
  "errors": { "timeout": 15, "connection refused": 3 },
  "duration_seconds": 1.84,
  "config": {
    "URL": "https://example.com",
    "Method": "GET",
    "Concurrency": 100,
    "NumRequests": 10000
  }
}
```

---

## 🎯 Real-World Scenarios

### API Load Testing

```bash
# Test a REST API endpoint under load
loadforge bench \
  -n 50000 \
  -c 200 \
  -m POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"user_id": 12345, "action": "login"}' \
  https://api.example.com/v1/auth/login
```

### Web Server Benchmarking

```bash
# Benchmark a web server (static content)
loadforge bench -n 100000 -c 500 https://yourserver.com/

# Benchmark with keep-alive (default behavior)
loadforge bench -n 50000 -c 200 -H "Connection: keep-alive" https://yourserver.com/
```

### Continuous Monitoring (CI/CD Pipeline)

```bash
# In your CI pipeline (GitHub Actions, GitLab CI, etc.)
loadforge bench -n 10000 -c 100 --json https://staging.example.com > result.json

# Check P99 latency threshold
P99=$(cat result.json | python3 -c "import sys,json; print(json.load(sys.stdin)['latency']['p99_ms'])")
if (( $(echo "$P99 > 200" | bc -l) )); then
  echo "❌ P99 latency ${P99}ms exceeds 200ms threshold"
  exit 1
else
  echo "✅ P99 latency ${P99}ms within threshold"
fi
```

### Compare Performance Before/After

```bash
# Baseline
loadforge bench -n 50000 -c 100 --json https://old-api.example.com > baseline.json

# After optimization
loadforge bench -n 50000 -c 100 --json https://new-api.example.com > optimized.json

# Compare
python3 -c "
import json
b = json.load(open('baseline.json'))
o = json.load(open('optimized.json'))
print(f'QPS: {b[\"qps\"]:.0f} -> {o[\"qps\"]:.0f} ({(o[\"qps\"]/b[\"qps\"]-1)*100:+.1f}%)')
print(f'P99: {b[\"latency\"][\"p99_ms\"]:.1f}ms -> {o[\"latency\"][\"p99_ms\"]:.1f}ms')
"
```

---

## 🗂️ Data Storage

Test results are automatically saved to `~/.loadforge/tests/`:

```
~/.loadforge/
└── tests/
    ├── index.json                    # Test index (metadata for all tests)
    ├── 1784900665705/                # Test folder (timestamp-based ID)
    │   ├── config.json               # Test configuration (URL, method, headers, etc.)
    │   ├── result.json               # Full statistics (all metrics, latencies, errors)
    │   └── timeline.json             # Per-second timeline (QPS and latency over time)
    └── 1784900667939/                # Another test
        ├── config.json
        ├── result.json
        └── timeline.json
```

**Storage path** can be customized via environment variable:

```bash
LOADFORGE_HOME=/custom/path loadforge bench -n 1000 -c 10 https://example.com
```

---

## 🏗️ Project Structure

```
loadforge/
├── main.go                          # Entry point
├── cmd/
│   ├── root.go                      # Root CLI command (Cobra)
│   └── bench.go                     # bench subcommand
├── engine/
│   └── engine.go                    # Async HTTP benchmark engine
├── stats/
│   └── stats.go                     # Statistics computation
├── storage/
│   └── storage.go                   # File-based persistence
├── report/
│   ├── web.go                       # Web report server
│   ├── embed.go                     # Embedded frontend assets
│   └── ui/                          # Vue3 + ECharts frontend
│       ├── index.html
│       ├── vite.config.js
│       ├── package.json
│       └── src/
│           ├── main.js
│           ├── App.vue
│           └── components/
│               ├── TestList.vue     # Test history list
│               └── TestDetail.vue   # Single test detail with charts
├── README.md                        # This file (English)
├── README.zh-CN.md                  # Chinese documentation
├── .gitignore
├── go.mod
└── go.sum
```

---

## 🛠️ Tech Stack

| Layer | Technology | Purpose |
|-------|-----------|---------|
| **Backend** | Go 1.18+ (net/http) | High-concurrency HTTP client, async I/O |
| **Concurrency** | Goroutines + Channels | Worker pool pattern, non-blocking results |
| **CLI Framework** | Cobra (spf13/cobra) | Command parsing, help generation |
| **Frontend** | Vue 3 + Vite | Reactive SPA dashboard |
| **Charts** | ECharts 5 | Interactive visualization (QPS, latency, status) |
| **Embedding** | Go `embed` | Single binary — no separate frontend deployment |
| **Storage** | JSON + Filesystem | `~/.loadforge/tests/` — zero dependencies |

---

## 🔧 Troubleshooting

### Connection Issues

| Error | Likely Cause | Solution |
|-------|-------------|----------|
| `connection refused` | Server not running or wrong port | Check `curl -v http://url` first |
| `i/o timeout` | Network issue or server overload | Increase `--timeout` or reduce `-c` |
| `TLS handshake timeout` | SSL/TLS negotiation failure | Verify TLS version, try `http://` |
| `no such host` | DNS resolution failure | Check URL spelling, try IP directly |
| `403` / `401` | Authentication/authorization | Add auth headers with `-H` |

### Performance Issues

| Symptom | Check | Fix |
|---------|-------|-----|
| QPS lower than expected | `-c` too low | Increase concurrency |
| High failure rate | Server limits (rate limit, max connections) | Reduce `-c`, check server logs |
| P99 >> P95 | Occasional slow requests | Check for GC pauses, DB queries |
| Out of memory | Too many concurrent workers | Reduce `-c`, use `-t` (duration-based) |

### Build Issues

| Error | Solution |
|-------|----------|
| `go: not found` | Install Go: https://go.dev/dl/ |
| `npm: not found` | Install Node.js: https://nodejs.org/ |
| `no matching files found` (embed) | Run `npm install && npm run build` in `report/ui/` |
| `dial tcp: i/o timeout` (go get) | Set `GOPROXY=https://goproxy.cn,https://proxy.golang.org,direct` |

---

## 🤝 Community

- 💬 **Discord**: [Join the discussion](https://discord.gg/zqQ6rdYT) — get help, share ideas, report issues
- 🐛 **GitHub Issues**: [Report bugs or request features](https://github.com/fanfan-2011/loadforge/issues)
- ⭐ **Star**: If LoadForge helps you, please give it a star!
- 🔀 **Contributions**: PRs welcome! See our contributing guidelines

---

## 📄 License

MIT License © 2026 [fanfan-2011](https://github.com/fanfan-2011)

*Built with ❤️ using Go and Vue.js*
