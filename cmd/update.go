package cmd

import (
	"fmt"
	"os"

	"github.com/nousresearch/loadforge/version"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "检查并安装 LoadForge 最新版本",
	Long: `从多源自动下载并安装 LoadForge 最新版本。
依次尝试 Gitee、jsDelivr CDN、GitHub 等源。`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔍 正在检查更新...")

		info := version.CheckUpdate()
		if info == nil {
			fmt.Printf("   ✅ 当前版本 (%s) 已是最新。\n", version.Version)
			return
		}

		fmt.Printf("   📢 发现新版本: %s\n", info.Version)
		fmt.Println("   ⬇️  正在下载更新...")

		if err := version.InstallUpdate(); err != nil {
			fmt.Fprintf(os.Stderr, "   ❌ 更新失败: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("   🎉 更新完成！请重新运行 loadforge 使用新版本。")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
