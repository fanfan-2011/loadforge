package report

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/nousresearch/loadforge/storage"
)

// GetLocalIP 获取本机局域网 IP
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "localhost"
	}

	// 优先局域网 IP
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.String()
				if strings.HasPrefix(ip, "192.168.") ||
					strings.HasPrefix(ip, "10.") ||
					strings.HasPrefix(ip, "172.") {
					return ip
				}
			}
		}
	}

	// 其次 hostname
	hostname, err := os.Hostname()
	if err == nil {
		return hostname
	}

	return "localhost"
}

// StartServer 启动 Web 报告服务器
func StartServer(s *storage.Storage, ip string) {
	mux := http.NewServeMux()

	// API 路由
	mux.HandleFunc("/api/tests", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		tests := s.ListTests()
		json.NewEncoder(w).Encode(tests)
	})

	mux.HandleFunc("/api/tests/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// 路径: /api/tests/{id}
		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/tests/"), "/")
		if len(parts) == 0 || parts[0] == "" {
			http.Error(w, "missing test id", http.StatusBadRequest)
			return
		}

		testID := parts[0]
		test, err := s.GetTest(testID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(test)
	})

	// 静态文件 - 内嵌的 Vue3 前端
	mux.Handle("/", http.FileServer(http.FS(GetUIFS())))

	addr := fmt.Sprintf("%s:8899", ip)
	log.Printf("LoadForge 报告服务启动于 http://%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Printf("报告服务停止: %v", err)
	}
}
