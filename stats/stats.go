package stats

import (
	"sort"
	"time"

	"github.com/nousresearch/loadforge/engine"
)

// TestStats 测试统计结果
type TestStats struct {
	TotalRequests   int                `json:"total_requests"`
	SuccessCount    int                `json:"success_count"`
	FailCount       int                `json:"fail_count"`
	FailRate        float64            `json:"fail_rate"`
	QPS             float64            `json:"qps"`
	Throughput      float64            `json:"throughput_mb_s"`
	BytesSent       int64              `json:"bytes_sent"`
	BytesReceived   int64              `json:"bytes_received"`
	AvgResponseSize float64            `json:"avg_response_size_bytes"`
	Latency         LatencyStats       `json:"latency"`
	StatusCodes     map[int]int        `json:"status_codes"`
	Errors          map[string]int     `json:"errors"`
	Duration        float64            `json:"duration_seconds"`
	Config          *engine.BenchConfig `json:"config"`
	PerformanceTips []string           `json:"performance_tips"`
}

// LatencyStats 延迟统计
type LatencyStats struct {
	Min  float64 `json:"min_ms"`
	Avg  float64 `json:"avg_ms"`
	Max  float64 `json:"max_ms"`
	P50  float64 `json:"p50_ms"`
	P75  float64 `json:"p75_ms"`
	P90  float64 `json:"p90_ms"`
	P95  float64 `json:"p95_ms"`
	P99  float64 `json:"p99_ms"`
	P999 float64 `json:"p999_ms"`
}

// Calculate 计算测试统计
func Calculate(result *engine.BenchResult, elapsed time.Duration) *TestStats {
	stat := &TestStats{
		StatusCodes: make(map[int]int),
		Errors:      make(map[string]int),
	}

	stat.TotalRequests = len(result.Results)
	stat.Duration = elapsed.Seconds()

	latencies := make([]float64, 0, len(result.Results))

	for _, r := range result.Results {
		if r.Error != "" {
			stat.FailCount++
			stat.Errors[r.Error]++
		} else {
			stat.SuccessCount++
		}

		latencies = append(latencies, float64(r.Latency/time.Millisecond))

		stat.BytesSent += r.BytesSent
		stat.BytesReceived += r.BytesReceived

		if r.StatusCode > 0 {
			stat.StatusCodes[r.StatusCode]++
		}
	}

	// 延迟统计
	if len(latencies) > 0 {
		sort.Float64s(latencies)

		stat.Latency.Min = latencies[0]
		stat.Latency.Max = latencies[len(latencies)-1]

		var sum float64
		for _, v := range latencies {
			sum += v
		}
		stat.Latency.Avg = sum / float64(len(latencies))

		stat.Latency.P50 = percentile(latencies, 50)
		stat.Latency.P75 = percentile(latencies, 75)
		stat.Latency.P90 = percentile(latencies, 90)
		stat.Latency.P95 = percentile(latencies, 95)
		stat.Latency.P99 = percentile(latencies, 99)
		stat.Latency.P999 = percentile(latencies, 99.9)
	}

	// 失败率
	if stat.TotalRequests > 0 {
		stat.FailRate = float64(stat.FailCount) / float64(stat.TotalRequests)
	}

	// QPS 和吞吐量
	if elapsed.Seconds() > 0 {
		stat.QPS = float64(stat.TotalRequests) / elapsed.Seconds()
		stat.Throughput = float64(stat.BytesReceived) / elapsed.Seconds() / 1024 / 1024
	}

	// 平均响应大小
	if stat.SuccessCount > 0 {
		stat.AvgResponseSize = float64(stat.BytesReceived) / float64(stat.TotalRequests)
	}

	// 性能建议
	if stat.Latency.P99 > stat.Latency.P95*1.5 {
		stat.PerformanceTips = append(stat.PerformanceTips,
			"发现 P99 延迟突增，服务端可能存在瓶颈")
	}
	if stat.FailRate > 0.05 {
		stat.PerformanceTips = append(stat.PerformanceTips,
			"失败率超过 5%，请检查服务端状态")
	}

	return stat
}

// percentile 计算百分位值
func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	if p <= 0 {
		return sorted[0]
	}
	if p >= 100 {
		return sorted[len(sorted)-1]
	}

	idx := float64(len(sorted)-1) * p / 100.0
	lower := int(idx)
	upper := lower + 1
	if upper >= len(sorted) {
		return sorted[lower]
	}
	fraction := idx - float64(lower)

	return sorted[lower]*(1-fraction) + sorted[upper]*fraction
}
