package engine

import (
	"crypto/tls"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// BenchConfig 压测配置
type BenchConfig struct {
	URL         string
	Method      string
	Headers     map[string]string
	Body        string
	Concurrency int
	NumRequests int
	Duration    time.Duration
	Timeout     time.Duration
}

// RequestResult 单次请求结果
type RequestResult struct {
	StatusCode   int
	Latency      time.Duration
	BytesSent    int64
	BytesReceived int64
	Error        string
	Timestamp    time.Time
}

// BenchResult 压测结果
type BenchResult struct {
	Results      []*RequestResult
	Timeline     []*TimelinePoint
	TotalTime    time.Duration
}

// TimelinePoint 时间线点
type TimelinePoint struct {
	Timestamp time.Time `json:"timestamp"`
	QPS       float64   `json:"qps"`
	Latency   float64   `json:"latency"`
}

// RunBench 执行压测
func RunBench(config *BenchConfig) *BenchResult {
	// 创建 HTTP 客户端
	transport := &http.Transport{
		MaxIdleConns:        config.Concurrency * 2,
		MaxIdleConnsPerHost: config.Concurrency * 2,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  false,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   config.Timeout,
	}

	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
		client.Timeout = 30 * time.Second
	}

	// 使用持续时间模式或请求数模式
	useDuration := config.Duration > 0

	var totalRequests int64
	var completedRequests int64

	results := make(chan *RequestResult, 10000)
	timeline := make(chan *TimelinePoint, 1000)
	done := make(chan struct{})

	// 统计协程 - 收集结果
	var allResults []*RequestResult
	var allTimeline []*TimelinePoint
	var mu sync.Mutex

	go func() {
		for r := range results {
			mu.Lock()
			allResults = append(allResults, r)
			mu.Unlock()
			atomic.AddInt64(&completedRequests, 1)
		}
		close(done)
	}()

	// 时间线采样协程
	timelineDone := make(chan struct{})
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		lastCount := int64(0)

		for {
			select {
			case <-timelineDone:
				return
			case <-ticker.C:
				current := atomic.LoadInt64(&completedRequests)
				qps := float64(current-lastCount) / 1.0
				lastCount = current

				mu.Lock()
				var avgLatency float64
				if len(allResults) > 0 {
					var totalLat time.Duration
					start := len(allResults) - 100
					if start < 0 {
						start = 0
					}
					count := 0
					for i := start; i < len(allResults); i++ {
						totalLat += allResults[i].Latency
						count++
					}
					if count > 0 {
						avgLatency = float64(totalLat/time.Millisecond) / float64(count)
					}
				}
				mu.Unlock()

				select {
				case timeline <- &TimelinePoint{
					Timestamp: time.Now(),
					QPS:       qps,
					Latency:   avgLatency,
				}:
				default:
				}
			}
		}
	}()

	// 收集时间线
	go func() {
		for tp := range timeline {
			allTimeline = append(allTimeline, tp)
		}
	}()

	var wg sync.WaitGroup

	// 控制起始时间
	startTime := time.Now()

	if useDuration {
		// 持续时间模式
		endTime := startTime.Add(config.Duration)
		var stop int32

		for i := 0; i < config.Concurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					if atomic.LoadInt32(&stop) != 0 || time.Now().After(endTime) {
						return
					}
					atomic.AddInt64(&totalRequests, 1)
					r := sendRequest(client, config)
					results <- r
				}
			}()
		}

		// 等待时间结束
		time.Sleep(config.Duration)
		atomic.StoreInt32(&stop, 1)
	} else {
		// 请求数模式
		total := int64(config.NumRequests)

		for i := 0; i < config.Concurrency; i++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()
				for {
					reqNum := atomic.AddInt64(&totalRequests, 1)
					if reqNum > total {
						atomic.AddInt64(&totalRequests, -1)
						return
					}
					r := sendRequest(client, config)
					results <- r
				}
			}(i)
		}
	}

	wg.Wait()
	close(timelineDone)
	close(results)
	<-done
	close(timeline)

	elapsed := time.Since(startTime)

	return &BenchResult{
		Results:   allResults,
		Timeline:  allTimeline,
		TotalTime: elapsed,
	}
}

// sendRequest 发送单个 HTTP 请求
func sendRequest(client *http.Client, config *BenchConfig) *RequestResult {
	start := time.Now()

	var bodyReader io.Reader
	if config.Body != "" {
		bodyReader = io.NopCloser(stringsReader{config.Body})
	}

	req, err := http.NewRequest(config.Method, config.URL, bodyReader)
	if err != nil {
		return &RequestResult{
			StatusCode: 0,
			Latency:    time.Since(start),
			Error:      "请求创建失败: " + err.Error(),
			Timestamp:  start,
		}
	}

	for k, v := range config.Headers {
		req.Header.Set(k, v)
	}

	// 计算发送字节数
	bytesSent := int64(len(config.Method) + len(config.URL) + 20)
	for k, v := range config.Headers {
		bytesSent += int64(len(k) + len(v) + 4)
	}
	bytesSent += int64(len(config.Body))

	resp, err := client.Do(req)
	if err != nil {
		errMsg := err.Error()
		errType := "网络错误"
		if isTimeout(err) {
			errType = "超时"
		}
		return &RequestResult{
			StatusCode: 0,
			Latency:    time.Since(start),
			Error:      errType + ": " + errMsg,
			BytesSent: bytesSent,
			Timestamp:  start,
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RequestResult{
			StatusCode: resp.StatusCode,
			Latency:    time.Since(start),
			Error:      "读取响应失败: " + err.Error(),
			BytesSent:  bytesSent,
			Timestamp:  start,
		}
	}

	latency := time.Since(start)
	return &RequestResult{
		StatusCode:    resp.StatusCode,
		Latency:       latency,
		BytesSent:     bytesSent,
		BytesReceived: int64(len(body)),
		Timestamp:     start,
	}
}

func isTimeout(err error) bool {
	if err == nil {
		return false
	}
	if netErr, ok := err.(interface{ Timeout() bool }); ok {
		return netErr.Timeout()
	}
	return false
}

// stringsReader 用于 io.NopCloser
type stringsReader struct {
	s string
}

func (r stringsReader) Read(b []byte) (int, error) {
	return copy(b, r.s), io.EOF
}
