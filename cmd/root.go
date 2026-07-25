package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/nousresearch/loadforge/version"
	"github.com/spf13/cobra"
)

var (
	versionCheckOnce sync.Once
	rootCmd          = &cobra.Command{
		Use:   "loadforge",
		Short: "LoadForge - 现代化 HTTP/HTTPS 压力测试工具",
		Long: `LoadForge 是一个高性能、跨平台的 HTTP/HTTPS 压力测试工具，
是 ApacheBench (ab) 的现代化替代品。

支持高并发、异步网络模型、实时统计和 Web 可视化报告。`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// 每次运行自动检查更新（仅一次）
			versionCheckOnce.Do(checkVersion)
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(benchCmd)
}

func checkVersion() {
	// 开发版不检查
	if version.Version == "dev" {
		return
	}

	go func() {
		info := version.CheckUpdate()
		if info != nil {
			fmt.Fprintf(os.Stderr, "\n📢 A new release %s is now available.\n", info.Version)
			fmt.Fprintf(os.Stderr, "   Install it by running: \033[1mloadforge update\033[0m\n\n")
		}
	}()
}
