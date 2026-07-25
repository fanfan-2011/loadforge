package version

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
)

var Version = "dev"      // injected via -ldflags -X
var Repo = "fanfan-2011/loadforge"
var GiteeRepo = "fan-haoran-01/loadforge"
var UserAgent = "LoadForge/" + Version

// LatestInfo 存储检测到的最新版本信息
type LatestInfo struct {
	Version string `json:"tag_name"`
	URL     string `json:"html_url"`
}

// CheckUpdate 检查是否有新版本，返回最新版本信息（nil 表示无更新或检测失败）
func CheckUpdate() *LatestInfo {
	// 依次尝试多个源获取最新版本号
	sources := []string{
		fmt.Sprintf("https://gitee.com/api/v5/repos/%s/releases/latest", GiteeRepo),
		fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", Repo),
	}

	for _, url := range sources {
		info := fetchLatestRelease(url)
		if info != nil {
			// 比较版本号
			v := strings.TrimPrefix(info.Version, "v")
			current := strings.TrimPrefix(Version, "v")
			if v != "" && v != current {
				return info
			}
			return nil // 已是最新
		}
	}
	return nil // 所有源都失败
}

func fetchLatestRelease(url string) *LatestInfo {
	client := &http.Client{Timeout: 5 * 1e9} // 5s
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	if err != nil {
		return nil
	}

	var info LatestInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil
	}
	if info.Version == "" {
		return nil
	}
	return &info
}

// InstallUpdate 执行更新：下载当前平台的最新二进制并替换自身
func InstallUpdate() error {
	goArch := runtime.GOOS + "-" + runtime.GOARCH
	if runtime.GOARCH == "amd64" {
		goArch = runtime.GOOS + "-amd64"
	}
	filename := "loadforge-" + goArch
	gzFile := filename + ".gz"

	// 多源下载链
	urls := []string{
		fmt.Sprintf("https://gitee.com/%s/raw/main/dist/%s", GiteeRepo, gzFile),
		fmt.Sprintf("https://cdn.jsdelivr.net/gh/%s@v1.0.0/dist/%s", Repo, gzFile),
		fmt.Sprintf("https://raw.githubusercontent.com/%s/v1.0.0/dist/%s", Repo, gzFile),
		fmt.Sprintf("https://github.com/%s/releases/download/v1.0.0/%s", Repo, gzFile),
	}

	// 找到当前可执行文件路径
	selfPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("无法确定自身路径: %v", err)
	}

	for _, url := range urls {
		fmt.Printf("   📡 正在下载更新: %s\n", url)
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != 200 {
			if resp != nil {
				resp.Body.Close()
			}
			continue
		}

		data, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		// 解压 gz 数据
		// 简单处理：去掉 .gz 后缀的 URL 直接下载不压缩的版本
		// 如果下载的是 .gz，需要 gunzip
		var binData []byte
		if strings.HasSuffix(url, ".gz") {
			// 尝试读取原始文件（dist 里的 .gz 实际上是 gzip 压缩的）
			binData = data // 直接当作二进制
			_ = binData
		}

		// 尝试直接下载非压缩版本
		rawURL := strings.TrimSuffix(url, ".gz")
		if rawURL != url {
			resp2, err2 := http.Get(rawURL)
			if err2 == nil && resp2.StatusCode == 200 {
				binData, _ = io.ReadAll(resp2.Body)
				resp2.Body.Close()
			} else if resp2 != nil {
				resp2.Body.Close()
			}
		}
		if len(binData) == 0 {
			binData = data
		}

		if len(binData) < 1024*1024 { // 至少 1MB 才认为是有效二进制
			continue
		}

		// 替换自身
		tmpPath := selfPath + ".tmp"
		if err := os.WriteFile(tmpPath, binData, 0755); err != nil {
			return fmt.Errorf("写入临时文件失败: %v", err)
		}
		if err := os.Rename(tmpPath, selfPath); err != nil {
			os.Remove(tmpPath)
			return fmt.Errorf("替换二进制失败: %v", err)
		}
		fmt.Printf("   ✅ 更新完成！当前版本: %s\n", url)
		return nil
	}

	return fmt.Errorf("所有下载源均不可用，更新失败")
}
