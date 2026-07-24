package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nousresearch/loadforge/engine"
	"github.com/nousresearch/loadforge/stats"
)

// Storage 数据存储
type Storage struct {
	BaseDir string
	tests   []TestMeta
}

// TestMeta 测试元数据
type TestMeta struct {
	ID        string  `json:"id"`
	URL       string  `json:"url"`
	Timestamp string  `json:"timestamp"`
	QPS       float64 `json:"qps"`
	TotalReqs int     `json:"total_requests"`
	Status    string  `json:"status"`
}

// New 创建存储实例
func New() *Storage {
	home, _ := os.UserHomeDir()
	baseDir := filepath.Join(home, ".loadforge", "tests")

	s := &Storage{
		BaseDir: baseDir,
	}
	s.ensureDirs()
	s.loadTestList()
	return s
}

func (s *Storage) ensureDirs() {
	os.MkdirAll(s.BaseDir, 0755)
}

func (s *Storage) testDir(id string) string {
	return filepath.Join(s.BaseDir, id)
}

// SaveConfig 保存测试配置
func (s *Storage) SaveConfig(id string, config *engine.BenchConfig) {
	dir := s.testDir(id)
	os.MkdirAll(dir, 0755)

	data, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile(filepath.Join(dir, "config.json"), data, 0644)
}

// SaveResult 保存测试结果
func (s *Storage) SaveResult(id string, stat *stats.TestStats) {
	dir := s.testDir(id)
	os.MkdirAll(dir, 0755)

	data, _ := json.MarshalIndent(stat, "", "  ")
	os.WriteFile(filepath.Join(dir, "result.json"), data, 0644)

	// 更新测试列表
	url := ""
	if stat.Config != nil {
		url = stat.Config.URL
	}
	s.tests = append(s.tests, TestMeta{
		ID:        id,
		URL:       url,
		Timestamp: id,
		QPS:       stat.QPS,
		TotalReqs: stat.TotalRequests,
		Status:    "completed",
	})
	s.saveTestList()
}

// SaveTimeline 保存时间线数据
func (s *Storage) SaveTimeline(id string, timeline []*engine.TimelinePoint) {
	if len(timeline) == 0 {
		return
	}
	dir := s.testDir(id)
	os.MkdirAll(dir, 0755)

	data, _ := json.MarshalIndent(timeline, "", "  ")
	os.WriteFile(filepath.Join(dir, "timeline.json"), data, 0644)
}

// ListTests 获取测试列表
func (s *Storage) ListTests() []TestMeta {
	s.loadTestList()
	return s.tests
}

// GetTest 获取单个测试
func (s *Storage) GetTest(id string) (map[string]interface{}, error) {
	dir := s.testDir(id)

	result := make(map[string]interface{})

	// 读取 config.json
	configData, err := os.ReadFile(filepath.Join(dir, "config.json"))
	if err == nil {
		var config interface{}
		json.Unmarshal(configData, &config)
		result["config"] = config
	}

	// 读取 result.json
	resultData, err := os.ReadFile(filepath.Join(dir, "result.json"))
	if err == nil {
		var res interface{}
		json.Unmarshal(resultData, &res)
		result["result"] = res
	}

	// 读取 timeline.json
	timelineData, err := os.ReadFile(filepath.Join(dir, "timeline.json"))
	if err == nil {
		var timeline interface{}
		json.Unmarshal(timelineData, &timeline)
		result["timeline"] = timeline
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("测试 %s 不存在", id)
	}

	return result, nil
}

func (s *Storage) loadTestList() {
	data, err := os.ReadFile(filepath.Join(s.BaseDir, "index.json"))
	if err != nil {
		s.tests = []TestMeta{}
		return
	}
	json.Unmarshal(data, &s.tests)
}

func (s *Storage) saveTestList() {
	data, _ := json.MarshalIndent(s.tests, "", "  ")
	os.WriteFile(filepath.Join(s.BaseDir, "index.json"), data, 0644)
}
