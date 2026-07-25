package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nousresearch/loadforge/engine"
	"github.com/nousresearch/loadforge/report"
	"github.com/nousresearch/loadforge/stats"
	"github.com/nousresearch/loadforge/storage"
	"github.com/spf13/cobra"
)

var (
	flagN       int
	flagC       int
	flagT       string
	flagM       string
	flagH       []string
	flagD       string
	flagJSON    bool
	flagReport  bool
	flagNoReport bool
	flagTimeout int
)

var benchCmd = &cobra.Command{
	Use:   "bench [options] <url>",
	Short: "执行 HTTP/HTTPS 压力测试",
	Long: `对指定 URL 执行 HTTP/HTTPS 压力测试。

示例:
  loadforge bench -n 10000 -c 100 https://example.com
  loadforge bench -t 30s -c 500 -m POST -d '{"key":"value"}' https://api.example.com
  loadforge bench -n 100000 -c 1000 --json https://example.com
  loadforge bench -n 50000 -c 200 --report https://example.com`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]

		// 解析持续时间
		var duration time.Duration
		if flagT != "" {
			var err error
			duration, err = time.ParseDuration(flagT)
			if err != nil {
				fmt.Fprintf(os.Stderr, "错误: 无效的持续时间格式 '%s'\n", flagT)
				os.Exit(1)
			}
		}

		// 解析自定义 Header
		headers := make(map[string]string)
		for _, h := range flagH {
			parts := strings.SplitN(h, ":", 2)
			if len(parts) == 2 {
				headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}

		var config *engine.BenchConfig
		if !flagJSON {
			config = &engine.BenchConfig{
				URL:         url,
				Method:      flagM,
				Headers:     headers,
				Body:        flagD,
				Concurrency: flagC,
				NumRequests: flagN,
				Duration:    duration,
				Timeout:     time.Duration(flagTimeout) * time.Second,
			}

			// 执行测试
			fmt.Printf("🧪 LoadForge 压力测试\n")
			fmt.Printf("   URL:         %s\n", url)
			fmt.Printf("   方法:        %s\n", flagM)
			fmt.Printf("   并发数:      %d\n", flagC)
			if flagN > 0 {
				fmt.Printf("   请求总数:    %d\n", flagN)
			}
			if flagT != "" {
				fmt.Printf("   持续时间:    %s\n", flagT)
			}
			fmt.Println(strings.Repeat("─", 50))
		} else {
			config = &engine.BenchConfig{
				URL:         url,
				Method:      flagM,
				Headers:     headers,
				Body:        flagD,
				Concurrency: flagC,
				NumRequests: flagN,
				Duration:    duration,
				Timeout:     time.Duration(flagTimeout) * time.Second,
			}
		}

		startTime := time.Now()
		result := engine.RunBench(config)
		elapsed := time.Since(startTime)

		// 计算统计
		stat := stats.Calculate(result, elapsed)

		// 命令行输出
		if flagJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.Encode(stat)
		} else {
			displayResult(stat, config)
		}

		// 保存结果
		testID := fmt.Sprintf("%d", time.Now().UnixMilli())
		s := storage.New()
		s.SaveConfig(testID, config)
		stat.Config = config
		s.SaveResult(testID, stat)
		s.SaveTimeline(testID, result.Timeline)

		// Web 报告（默认开启，用 --no-report 关闭）
		if !flagNoReport {
			fmt.Printf("\n📊 正在启动 Web 报告服务...\n")
			ip := report.GetLocalIP()
			go func() {
				report.StartServer(s, ip)
			}()
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("   报告地址: http://%s:8899\n", ip)
			fmt.Println("   按 Ctrl+C 停止服务")
			select {}
		}
	},
}

func displayResult(stat *stats.TestStats, config *engine.BenchConfig) {
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║        LoadForge 压力测试报告             ║")
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Println()

	// 请求统计
	fmt.Println("📊 请求统计")
	fmt.Printf("   总请求数:     %d\n", stat.TotalRequests)
	fmt.Printf("   成功数:       %d\n", stat.SuccessCount)
	fmt.Printf("   失败数:       %d  (%.2f%%)\n", stat.FailCount, stat.FailRate*100)
	fmt.Println()

	// 性能
	fmt.Println("⚡ 性能")
	fmt.Printf("   QPS:          %.2f req/s\n", stat.QPS)
	fmt.Printf("   吞吐量:       %.2f MB/s\n", stat.Throughput)
	fmt.Printf("   上传流量:     %.2f MB\n", float64(stat.BytesSent)/1024/1024)
	fmt.Printf("   下载流量:     %.2f MB\n", float64(stat.BytesReceived)/1024/1024)
	fmt.Printf("   平均响应大小: %.2f KB\n", float64(stat.AvgResponseSize)/1024)
	fmt.Println()

	// 延迟
	fmt.Println("⏱️  延迟分析 (毫秒)")
	fmt.Printf("   Min:          %.2f\n", stat.Latency.Min)
	fmt.Printf("   Avg:          %.2f\n", stat.Latency.Avg)
	fmt.Printf("   Max:          %.2f\n", stat.Latency.Max)
	fmt.Printf("   P50:          %.2f\n", stat.Latency.P50)
	fmt.Printf("   P75:          %.2f\n", stat.Latency.P75)
	fmt.Printf("   P90:          %.2f\n", stat.Latency.P90)
	fmt.Printf("   P95:          %.2f\n", stat.Latency.P95)
	fmt.Printf("   P99:          %.2f\n", stat.Latency.P99)
	fmt.Printf("   P999:         %.2f\n", stat.Latency.P999)
	fmt.Println()

	// HTTP 状态码分布
	fmt.Println("🔢 状态码分布")
	for code, count := range stat.StatusCodes {
		pct := float64(count) / float64(stat.TotalRequests) * 100
		fmt.Printf("   %d:           %d (%.2f%%)\n", code, count, pct)
	}
	if len(stat.StatusCodes) == 0 {
		fmt.Println("   (无数据)")
	}
	fmt.Println()

	// 错误
	if len(stat.Errors) > 0 {
		fmt.Println("❌ 错误统计")
		errTotal := 0
		for _, c := range stat.Errors {
			errTotal += c
		}
		for errType, count := range stat.Errors {
			pct := float64(count) / float64(errTotal) * 100
			fmt.Printf("   %s:  %d (%.2f%%)\n", errType, count, pct)
		}
		fmt.Println()
	}

	// 性能建议
	if stat.Latency.P99 > stat.Latency.P95*2 {
		fmt.Println("💡 性能建议:")
		if stat.Latency.P99 > 1000 {
			fmt.Println("   - ⚠️  发现 P99 延迟突增，服务端可能存在瓶颈")
		}
		fmt.Println("   - 可能原因: 服务端压力过高 / 数据库响应慢 / 连接池不足")
		fmt.Println("   - 建议: 检查服务端资源使用情况，考虑扩容或优化")
		fmt.Println()
	}
}

func init() {
	benchCmd.Flags().IntVarP(&flagN, "requests", "n", 1000, "请求总数")
	benchCmd.Flags().IntVarP(&flagC, "concurrency", "c", 10, "并发数")
	benchCmd.Flags().StringVarP(&flagT, "duration", "t", "", "持续时间 (如 30s, 1m)")
	benchCmd.Flags().StringVarP(&flagM, "method", "m", "GET", "HTTP 方法")
	benchCmd.Flags().StringArrayVarP(&flagH, "header", "H", []string{}, "请求 Header (可重复)")
	benchCmd.Flags().StringVarP(&flagD, "body", "d", "", "请求 Body")
	benchCmd.Flags().BoolVar(&flagJSON, "json", false, "以 JSON 格式输出")
	benchCmd.Flags().BoolVar(&flagReport, "report", true, "启动 Web 报告")
	benchCmd.Flags().BoolVar(&flagNoReport, "no-report", false, "不启动 Web 报告")
	benchCmd.Flags().IntVar(&flagTimeout, "timeout", 30, "超时时间 (秒)")
}
