# ⚡ LoadForge

> 现代化 ApacheBench 替代工具 — 高性能 HTTP/HTTPS 压力测试

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![Discord](https://img.shields.io/badge/Discord-加入讨论-5865F2?style=flat&logo=discord)](https://discord.gg/zqQ6rdYT)

LoadForge 是一个高性能、跨平台、**单文件发布** 的 HTTP/HTTPS 压力测试工具，是 ApacheBench (`ab`) 的现代化替代品。支持高并发异步网络模型、实时延迟百分位统计，以及内嵌的 Vue3 + ECharts Web 可视化报告。

---

## ✨ 特性

- **🚀 高性能压测引擎** — goroutine + channel 异步非阻塞 IO，支持 HTTP/1.1（可扩展 HTTP/2、HTTP/3）
- **📊 全面统计** — 请求数、成功率、QPS、吞吐量、延迟百分位 (P50/P75/P90/P95/P99/P999)
- **📈 Web 可视化报告** — 测试完成后自动启动内嵌的 Vue3 + ECharts 图表界面（端口 8899）
- **📂 持久化存储** — 测试结果自动保存至 `~/.loadforge/tests/`，支持 JSON 导出
- **🔧 CLI 优先** — 类 `ab` 命令行体验，`-n`、`-c`、`-t`、`-m`、`-H`、`-d`、`--json`、`--report`
- **💡 性能分析建议** — 自动检测 P99 延迟突增等瓶颈信号
- **📦 单文件发布** — Go 编译，零运行时依赖

---

## 🚀 快速开始

### 从源码编译

```bash
# 1. 克隆仓库
git clone https://github.com/fanfan-2011/loadforge.git
cd loadforge

# 2. 构建 Web 前端
cd report/ui
npm install
npm run build
cd ../..

# 3. 编译 Go 二进制
go build -ldflags="-s -w" -o loadforge .

# 4. 运行！
./loadforge bench -n 10000 -c 100 https://example.com
```

> 💡 或者直接下载 [Releases](https://github.com/fanfan-2011/loadforge/releases) 页面的预编译二进制文件。

### 基本用法

```bash
# 10 个并发发送 1000 个请求
loadforge bench -n 1000 -c 10 https://example.com

# 持续 30 秒压力测试
loadforge bench -t 30s -c 200 https://example.com

# POST 请求 + 自定义 Header
loadforge bench -n 5000 -c 50 -m POST -H "Content-Type: application/json" -d '{"key":"value"}' https://api.example.com

# JSON 格式输出（适合管道处理）
loadforge bench -n 10000 -c 100 --json https://example.com

# 测试完成后自动启动 Web 报告
loadforge bench -n 50000 -c 200 --report https://example.com
```

---

## 📋 命令行参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-n`  | 请求总数 | `1000` |
| `-c`  | 并发数 | `10` |
| `-t`  | 持续时间（如 `30s`、`1m`） | (禁用) |
| `-m`  | HTTP 方法 | `GET` |
| `-H`  | 请求 Header（可重复） | — |
| `-d`  | 请求 Body | — |
| `--json` | JSON 格式输出 | `false` |
| `--report` | 启动 Web 报告服务 | `false` |
| `--timeout` | 超时时间（秒） | `30` |

---

## 📊 输出示例

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
   P50:          14.27
   P95:          45.31
   P99:          128.67
   P999:         256.43

🔢 状态码分布
   200:           9982 (99.82%)
   503:           18 (0.18%)
```

---

## 🌐 Web 报告

使用 `--report` 标志启动 Web 可视化报告：

```bash
loadforge bench -n 100000 -c 500 --report https://example.com
```

测试完成后自动在 `http://<局域网IP>:8899` 启动 Web 界面，包含：

- **测试列表** — 所有历史测试一览
- **总览卡片** — QPS、成功率、吞吐量
- **图表** — QPS 曲线、延迟曲线、状态码分布、错误统计
- **延迟分析** — P50/P75/P90/P95/P99/P999
- **详细信息** — HTTP 配置、Header、Body 大小、流量统计
- **性能建议** — 自动检测瓶颈并给出优化建议

---

## 🗂️ 数据存储

所有测试结果自动保存至 `~/.loadforge/tests/<时间戳>/`：

```
~/.loadforge/tests/
├── 1784900665705/
│   ├── config.json      # 测试配置
│   ├── result.json      # 统计结果
│   └── timeline.json    # 时间线数据
├── 1784900667939/
│   └── ...
└── index.json           # 测试索引
```

---

## 🏗️ 项目结构

```
loadforge/
├── main.go                    # 入口
├── cmd/
│   ├── root.go                # CLI 根命令
│   └── bench.go               # bench 子命令
├── engine/
│   └── engine.go              # 异步 HTTP 压测引擎
├── stats/
│   └── stats.go               # 统计计算（百分位等）
├── storage/
│   └── storage.go             # 文件存储
├── report/
│   ├── web.go                 # Web 报告服务器
│   ├── embed.go               # 内嵌前端资源
│   └── ui/                    # Vue3 + ECharts 前端源码
│       └── src/
│           ├── App.vue
│           └── components/
│               ├── TestList.vue
│               └── TestDetail.vue
└── go.mod
```

---

## 🛠️ 技术栈

| 层 | 技术 |
|----|------|
| **后端** | Go 1.18+ (net/http, goroutine, channel) |
| **CLI** | Cobra (spf13/cobra) |
| **前端** | Vue 3 + Vite |
| **图表** | ECharts 5 |
| **嵌入** | Go `embed` 包（单文件分发） |
| **存储** | JSON 文件系统 |

---

## 🤝 社区

- 💬 **Discord**: [加入讨论](https://discord.gg/zqQ6rdYT)
- 🐛 **Issues**: [GitHub Issues](https://github.com/fanfan-2011/loadforge/issues)
- ⭐ **Star**: 如果觉得有用，请点个 Star！

---

## 📄 许可证

MIT License © 2026 [fanfan-2011](https://github.com/fanfan-2011)
