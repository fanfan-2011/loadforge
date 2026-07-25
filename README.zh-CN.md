# ⚡ LoadForge

> 现代化 ApacheBench 替代工具 — 高性能 HTTP/HTTPS 压力测试

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-GPLv3-blue)](LICENSE)
[![Discord](https://img.shields.io/badge/Discord-加入讨论-5865F2?style=flat&logo=discord)](https://discord.gg/zqQ6rdYT)
[![Release](https://img.shields.io/github/v/release/fanfan-2011/loadforge?style=flat)](https://github.com/fanfan-2011/loadforge/releases)

LoadForge 是一个高性能、跨平台、**单文件发布** 的 HTTP/HTTPS 压力测试工具，是 ApacheBench (`ab`) 的现代化替代品。支持高并发异步网络模型、实时延迟百分位统计，以及内嵌的 Vue3 + ECharts Web 可视化报告。

**[English](README.md)**

---

## 📦 目录

- [特性](#-特性)
- [环境要求](#-环境要求)
- [一键安装](#-一键安装)
- [快速开始](#-快速开始)
- [使用教程](#-使用教程)
- [CLI 参考](#-cli-参考)
- [Web 报告](#-web-报告)
- [输出示例](#-输出示例)
- [实际场景](#-实际场景)
- [数据存储](#-数据存储)
- [项目结构](#-项目结构)
- [技术栈](#-技术栈)
- [常见问题](#-常见问题)
- [社区](#-社区)
- [许可证](#-许可证)

---

## ✨ 特性

| 特性 | 说明 |
|------|------|
| **🚀 高性能引擎** | goroutine + channel 异步非阻塞 IO，HTTP/1.1（可扩展 HTTP/2、HTTP/3） |
| **📊 全面统计** | 请求数、成功率、QPS、吞吐量、延迟百分位 P50/P75/P90/P95/P99/P999 |
| **📈 Web 可视化报告** | 测试完成后自动启动内嵌 Vue3 + ECharts 图表界面，端口 8899 |
| **📂 持久化存储** | 测试结果自动保存至 `~/.loadforge/tests/`，支持 JSON 导出 |
| **🔧 CLI 优先** | 类 `ab` 命令行体验：`-n`、`-c`、`-t`、`-m`、`-H`、`-d`、`--json`、`--report` |
| **💡 性能分析建议** | 自动检测 P99 延迟突增等瓶颈信号 |
| **📦 单文件发布** | Go 编译 + 内嵌 Web UI，零运行时依赖 |
| **🌍 跨平台** | 支持 Linux、macOS、Windows |
| **♾️ 无限扩展** | 测试过 10 万+ 请求，1000+ 并发连接 |
| **🔄 自动更新** | 内置版本检测 + `loadforge update` 一键升级 |

---

## 🖥️ 环境要求

| 项目 | 运行时（二进制） | 源码编译 |
|------|----------------|---------|
| **操作系统** | Linux, macOS, Windows | Linux, macOS, Windows |
| **CPU** | x86-64 或 ARM64 | x86-64 或 ARM64 |
| **内存** | 128 MB | 512 MB |
| **磁盘** | 10 MB | 100 MB |
| **Go** | 不需要 | 1.18+ |
| **Node.js** | 不需要 | 16+ |

> 💡 预编译二进制已内嵌 Web 前端，**运行时不需要 Go 或 Node.js**。

---

## ⚡ 一键安装

Linux/macOS 用户一行命令安装：

```bash
curl -fsSL https://github.com/fanfan-2011/loadforge/raw/main/install.sh | bash
```

脚本自动多源下载，按优先级尝试：

| 优先级 | 源 | 适用 |
|--------|-----|------|
| ① | Gitee 仓库直链 | 中国大陆（快速） |
| ② | jsDelivr CDN | 中国大陆 / 亚洲（快速） |
| ③ | Raw GitHub | 备用 |
| ④ | GitHub Releases | 全球 |
| ⑤ | 源码编译 | 任意环境（需 Go + Node.js） |

所有下载源都失败后，会自动从源码编译（使用 `goproxy.cn`、`npmmirror.com` 国内加速）。

> 💡 如果 GitHub 完全不可用，也可直接从 Gitee 安装：
> ```bash
> curl -fsSL https://gitee.com/fan-haoran-01/loadforge/raw/main/install.sh | bash
> ```

安装后验证：

```bash
loadforge --help
loadforge bench -n 100 -c 10 https://example.com
```

> 🪟 **Windows 用户**：建议安装 WSL，或用源码编译方式

---

## 📥 安装指南

### 方式一：下载预编译二进制（推荐）

从 [Releases 页面](https://github.com/fanfan-2011/loadforge/releases) 下载最新版。

```bash
# Linux x86-64 示例
# 将 VERSION 替换为最新版本号（如 v1.0.0）
curl -LO https://github.com/fanfan-2011/loadforge/releases/latest/download/loadforge-linux-amd64.tar.gz
tar -xzf loadforge-linux-amd64.tar.gz
sudo mv loadforge /usr/local/bin/
loadforge --help
```

```powershell
# Windows PowerShell 示例
# 下载 loadforge-windows-amd64.zip 并解压
Move-Item .\loadforge.exe C:\Windows\System32\
loadforge --help
```

> 💡 预编译二进制已内嵌 Web 前端，**运行时不需要 Node.js 或 Go**。

### 方式二：从源码编译（通用）

适用于所有平台（需要 Go 和 Node.js）。

#### 第 1 步：安装依赖

<details>
<summary><b>Linux (Debian/Ubuntu)</b></summary>

```bash
sudo apt update
sudo apt install -y golang-go nodejs npm git

go version   # 需要 Go 1.18+
node --version  # 需要 v16+
npm --version
```
</details>

<details>
<summary><b>macOS (Homebrew)</b></summary>

```bash
brew install go node git

go version   # 需要 Go 1.18+
node --version  # 需要 v16+
npm --version
```
</details>

<details>
<summary><b>Windows</b></summary>

1. 安装 **Go**：https://go.dev/dl/（1.18+）
2. 安装 **Node.js**：https://nodejs.org/（16+）
3. 安装 **Git**：https://git-scm.com/
4. 打开 **PowerShell**

```powershell
go version
node --version
npm --version
git --version
```
</details>

#### 第 2 步：克隆并编译

```bash
# 克隆仓库
git clone https://github.com/fanfan-2011/loadforge.git
cd loadforge

# 构建 Web 前端（Vue3 + ECharts）
cd report/ui
npm install
npm run build
cd ../..

# 编译 Go 二进制
go build -ldflags="-s -w" -o loadforge .

# 验证
./loadforge --help
```

#### 第 3 步（可选）：安装到系统 PATH

```bash
# Linux / macOS
sudo mv loadforge /usr/local/bin/

# Windows（管理员 PowerShell）
# Move-Item .\loadforge.exe C:\Windows\System32\
```

### 方式三：国内加速

如果在中国大陆，编译时使用国内镜像：

```bash
# 使用 Go 国内代理
export GOPROXY=https://goproxy.cn,https://goproxy.io,direct

# 克隆仓库
git clone https://github.com/fanfan-2011/loadforge.git
cd loadforge

# 构建前端
cd report/ui && npm install && npm run build && cd ../..

# 编译
go build -ldflags="-s -w" -o loadforge .
```

### 验证安装

```bash
loadforge --help
```

期望输出：

```
LoadForge is a high-performance, cross-platform HTTP/HTTPS load testing tool...

Usage:
  loadforge [command]

Available Commands:
  bench       执行 HTTP/HTTPS 压力测试
  help        Help about any command
```

---

## 🚀 快速开始

### 测试本地服务器

```bash
# 启动一个测试服务器
python3 -m http.server 8888 &

# 运行 LoadForge
loadforge bench -n 100 -c 10 http://localhost:8888/
```

### 测试任意 URL

```bash
# 轻量测试（10 请求，2 并发）
loadforge bench -n 10 -c 2 https://example.com

# 标准测试（1000 请求，10 并发）
loadforge bench -n 1000 -c 10 https://example.com

# 压力测试（10 万请求，500 并发）
loadforge bench -n 100000 -c 500 https://your-api.com

# 持续时间模式（30 秒，100 并发）
loadforge bench -t 30s -c 100 https://example.com
```

---

## 📖 使用教程

### 1️⃣ 基本 GET 测试

```bash
loadforge bench -n 5000 -c 50 https://api.example.com/users
```

### 2️⃣ POST 带 JSON Body

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

### 3️⃣ 持续时间模式（推荐长时间运行）

```bash
# 运行 5 分钟，200 并发
loadforge bench -t 5m -c 200 https://example.com
```

### 4️⃣ JSON 输出（脚本处理）

```bash
# 管道到 jq 处理
loadforge bench -n 10000 -c 100 --json https://example.com | jq '.qps, .latency.p99_ms'

# 保存到文件
loadforge bench -n 5000 -c 50 --json https://example.com > result.json
```

### 5️⃣ 测试并查看 Web 报告

```bash
loadforge bench -n 50000 -c 200 --report https://example.com
# 浏览器打开 http://<你的局域网IP>:8899
```

### 6️⃣ 模拟真实流量

```bash
# 对不同端点并行测试
loadforge bench -n 30000 -c 150 -m GET https://api.example.com/items &
loadforge bench -n 20000 -c 100 -m POST -d '{"action":"search"}' https://api.example.com/search &
wait
```

---

## 📋 CLI 参考

### `loadforge bench`

```
用法:
  loadforge bench [flags] <url>

参数:
  -n, --requests int        请求总数（默认 1000）
  -c, --concurrency int     并发连接数（默认 10）
  -t, --duration string     测试持续时间（如 "30s"、"1m"、"5m"）
  -m, --method string       HTTP 方法（默认 "GET"）
  -H, --header strings      自定义 Header，可重复（如 -H "Key: Value"）
  -d, --body string         请求 Body
      --json                以 JSON 格式输出结果
      --no-report           不启动 Web 报告（默认启用）
      --timeout int         请求超时时间（秒，默认 30）
  -h, --help                查看帮助
```

### 退出码

| 代码 | 含义 |
|------|------|
| `0` | 测试成功完成 |
| `1` | 命令行错误（参数无效、缺少 URL） |
| `2` | 运行时错误（引擎崩溃、存储失败） |

---

## 🌐 Web 报告

Web 报告提供丰富的可视化分析。

### 启动方式

Web 报告在每次测试后**自动启动**，无需加任何参数：

```bash
# 测试完成后默认自动启动 Web 报告
loadforge bench -n 100000 -c 500 https://example.com

# 不需要 Web 报告时加 --no-report
loadforge bench -n 1000 -c 10 --no-report https://example.com
```

启动后在浏览器打开 `http://<局域网IP>:8899`。

### 报告页面

#### 测试列表
显示所有历史测试：
- 测试时间、目标 URL、请求数量
- QPS、完成状态

#### 单次报告详情

| 板块 | 内容 |
|------|------|
| **概览卡片** | 总请求数、QPS、成功数、失败数、平均延迟、吞吐量 |
| **摘要** | 目标 URL、测试时间、持续时间、请求配置 |
| **QPS 曲线图** | 每秒 QPS 折线图 |
| **延迟曲线图** | 每秒平均延迟趋势 |
| **状态码分布图** | HTTP 状态码饼图 |
| **错误统计图** | 错误类型分布（超时、连接失败） |
| **延迟分析表** | P50/P75/P90/P95/P99/P999/Min/Max |
| **详细信息** | HTTP 方法、Header、Body 大小、响应大小、流量 |
| **性能建议** | 自动生成的瓶颈分析和优化建议 |

---

## 📊 输出示例

### 标准 CLI 输出

```
╔══════════════════════════════════════════╗
║        LoadForge 压力测试报告             ║
╚══════════════════════════════════════════╝

📊 请求统计
   总请求数:     10000
   成功数:       9982
   失败数:       18  (0.18%)

⚡ 性能
   QPS:          5421.33 req/s
   吞吐量:       48.23 MB/s

⏱️  延迟分析 (毫秒)
   Min:          0.82
   Avg:          18.44
   Max:          512.00
   P50:          14.27
   P75:          22.18
   P90:          35.66
   P95:          45.31
   P99:          128.67
   P999:         256.43

🔢 状态码分布
   200:           9982 (99.82%)
   503:           18   (0.18%)

❌ 错误统计
   timeout:        15   (83.33%)
   connection refused: 3 (16.67%)

💡 性能建议
   ⚠️  发现 P99 延迟突增，服务端可能存在瓶颈
   - 可能原因: 服务端压力过高 / 数据库响应慢 / 连接池不足
```

### JSON 输出示例

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

## 🎯 实际场景

### API 压力测试

```bash
# 测试 REST API 端点
loadforge bench \
  -n 50000 \
  -c 200 \
  -m POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"user_id": 12345, "action": "login"}' \
  https://api.example.com/v1/auth/login
```

### Web 服务器基准测试

```bash
# 基准测试静态内容
loadforge bench -n 100000 -c 500 https://yourserver.com/
```

### CI/CD 管道集成

```bash
# 在 CI 中运行（GitHub Actions、GitLab CI 等）
loadforge bench -n 10000 -c 100 --json https://staging.example.com > result.json

# 检查 P99 延迟阈值
P99=$(python3 -c "import json; print(json.load(open('result.json'))['latency']['p99_ms'])")
if (( $(echo "$P99 > 200" | bc -l) )); then
  echo "❌ P99 延迟 ${P99}ms 超过 200ms 阈值"
  exit 1
else
  echo "✅ P99 延迟 ${P99}ms 在阈值内"
fi
```

### 优化前后对比

```bash
# 基准线
loadforge bench -n 50000 -c 100 --json https://old-api.example.com > baseline.json

# 优化后
loadforge bench -n 50000 -c 100 --json https://new-api.example.com > optimized.json

# 对比
python3 -c "
import json
b = json.load(open('baseline.json'))
o = json.load(open('optimized.json'))
print(f'QPS: {b[\"qps\"]:.0f} -> {o[\"qps\"]:.0f} ({(o[\"qps\"]/b[\"qps\"]-1)*100:+.1f}%)')
print(f'P99: {b[\"latency\"][\"p99_ms\"]:.1f}ms -> {o[\"latency\"][\"p99_ms\"]:.1f}ms')
"
```

---

## 🗂️ 数据存储

所有测试结果自动保存至 `~/.loadforge/tests/`：

```
~/.loadforge/
└── tests/
    ├── index.json                    # 测试索引（所有测试的元数据）
    ├── 1784900665705/                # 测试文件夹（时间戳 ID）
    │   ├── config.json               # 测试配置（URL、方法、Header 等）
    │   ├── result.json               # 完整统计（指标、延迟、错误）
    │   └── timeline.json             # 秒级时间线（QPS 和延迟变化）
    └── 1784900667939/                # 另一个测试
        ├── config.json
        ├── result.json
        └── timeline.json
```

**自定义存储路径**：

```bash
LOADFORGE_HOME=/custom/path loadforge bench -n 1000 -c 10 https://example.com
```

---

## 🏗️ 项目结构

```
loadforge/
├── main.go                          # 入口
├── cmd/
│   ├── root.go                      # CLI 根命令（Cobra）
│   └── bench.go                     # bench 子命令
├── engine/
│   └── engine.go                    # 异步 HTTP 压测引擎
├── stats/
│   └── stats.go                     # 统计计算
├── storage/
│   └── storage.go                   # 文件存储
├── report/
│   ├── web.go                       # Web 报告服务器
│   ├── embed.go                     # 内嵌前端资源
│   └── ui/                          # Vue3 + ECharts 前端
│       ├── index.html
│       ├── vite.config.js
│       ├── package.json
│       └── src/
│           ├── main.js
│           ├── App.vue
│           └── components/
│               ├── TestList.vue     # 测试历史列表
│               └── TestDetail.vue   # 测试详情（含图表）
├── README.md                        # 英文文档
├── README.zh-CN.md                  # 本文档（中文）
├── .gitignore
├── go.mod
└── go.sum
```

---

## 🛠️ 技术栈

| 层 | 技术 | 用途 |
|----|------|------|
| **后端** | Go 1.18+ (net/http) | 高并发 HTTP 客户端，异步 IO |
| **并发** | Goroutines + Channels | 工作池模式，非阻塞结果收集 |
| **CLI 框架** | Cobra (spf13/cobra) | 命令解析，帮助生成 |
| **前端** | Vue 3 + Vite | 响应式 SPA 仪表盘 |
| **图表** | ECharts 5 | 交互式可视化（QPS、延迟、状态码） |
| **嵌入** | Go `embed` | 单文件分发，无需单独部署前端 |
| **存储** | JSON + 文件系统 | `~/.loadforge/tests/`，零依赖 |

---

## 🔧 常见问题

### 连接问题

| 错误 | 可能原因 | 解决方法 |
|------|---------|---------|
| `connection refused` | 服务器未运行或端口错误 | 先用 `curl -v http://url` 确认 |
| `i/o timeout` | 网络问题或服务器过载 | 增加 `--timeout` 或减少 `-c` |
| `no such host` | DNS 解析失败 | 检查 URL 拼写，尝试直接使用 IP |
| `403` / `401` | 认证/授权失败 | 用 `-H` 添加认证 Header |

### 性能问题

| 现象 | 检查点 | 解决方法 |
|------|--------|---------|
| QPS 低于预期 | `-c` 太低 | 增加并发数 |
| 失败率高 | 服务器限流/连接限制 | 减少 `-c`，检查服务器日志 |
| P99 >> P95 | 偶发慢请求 | 检查 GC 停顿、数据库查询 |
| 内存不足 | 并发数过大 | 减少 `-c`，改用 `-t` 持续时间模式 |

### 编译问题

| 错误 | 解决方法 |
|------|---------|
| `go: not found` | 安装 Go：https://go.dev/dl/ |
| `npm: not found` | 安装 Node.js：https://nodejs.org/ |
| `no matching files found`（embed） | 在 `report/ui/` 下执行 `npm install && npm run build` |
| `dial tcp: i/o timeout`（go get） | 设置国内代理：`GOPROXY=https://goproxy.cn,https://proxy.golang.org,direct` |

### GFW 网络问题（中国大陆用户）

```bash
# SSH 推送配置
# ~/.ssh/config
Host github.com
    Hostname ssh.github.com
    Port 443
    User git
    StrictHostKeyChecking accept-new

# Go 代理
export GOPROXY=https://goproxy.cn,https://proxy.golang.org,direct

# NPM 镜像
npm config set registry https://registry.npmmirror.com
```

---

## 🤝 社区

- 💬 **Discord**: [加入讨论](https://discord.gg/zqQ6rdYT) — 获取帮助、分享想法、报告问题
- 🐛 **GitHub Issues**: [报告 Bug 或请求新功能](https://github.com/fanfan-2011/loadforge/issues)
- 🐛 **Gitee Issues**: [报告 Bug（国内）](https://gitee.com/fan-haoran-01/loadforge/issues)
- ⭐ **Star**: 如果 LoadForge 帮到了你，请点个 Star！
- 🔀 **贡献**: 欢迎 Pull Request！

---

## 📄 许可证

GNU General Public License v3.0 © 2026 [fanfan-2011](https://github.com/fanfan-2011)

详见 [LICENSE](LICENSE) 文件。

*用 ❤️ 和 Go + Vue.js 打造*
