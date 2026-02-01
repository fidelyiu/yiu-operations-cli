package cmd

import (
	"fmt"
	"log/slog"

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
		slog.Info(fmt.Sprintf("在端口 %d 上启动服务器", port))
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
	uiCmd.Flags().Int("ui.port", 8282, "运行服务器的端口")
}
