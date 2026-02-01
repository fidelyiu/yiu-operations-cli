package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// uiCmd 代表 ui 命令
var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "命令的简要描述",
	Long: `更长的描述，可以跨越多行，通常包含命令的
示例和用法。例如：

Cobra 是一个为 Go 应用程序赋能的 CLI 库。
此应用程序是一个用于生成所需文件的工具，
可以快速创建 Cobra 应用程序。`,
	Run: func(cmd *cobra.Command, args []string) {
		port := viper.GetInt("ui.port")
		fmt.Printf("在端口 %d 上启动服务器\\n", port)
		fmt.Println("ui 已调用")
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
	uiCmd.Flags().Int("port", 8080, "运行服务器的端口")

	// 在这里定义你的标志和配置设置。

	// Cobra 支持持久化标志，它将对此命令及其所有子命令生效，例如：
	// uiCmd.PersistentFlags().String("foo", "", "foo 的帮助信息")

	// Cobra 支持本地标志，它仅在直接调用此命令时运行，例如：
	// uiCmd.Flags().BoolP("toggle", "t", false, "toggle 的帮助信息")
}
