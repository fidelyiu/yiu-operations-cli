package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// uiCmd 代表 ui 命令
var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "启动基于 Web 的用户界面",
	Long: `yiu-operations 提供了一个基于 Web 的用户界面，允许用户通过浏览器管理和监控系统运维任务。
该界面直观且易于使用，提供了丰富的功能，帮助用户更高效地完成运维工作。
`,
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
