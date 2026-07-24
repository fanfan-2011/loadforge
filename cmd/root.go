package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "loadforge",
	Short: "LoadForge - 现代化 HTTP/HTTPS 压力测试工具",
	Long: `LoadForge 是一个高性能、跨平台的 HTTP/HTTPS 压力测试工具，
是 ApacheBench (ab) 的现代化替代品。

支持高并发、异步网络模型、实时统计和 Web 可视化报告。`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(benchCmd)
}
