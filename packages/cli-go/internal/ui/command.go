package ui

import (
	"yiu-ops/internal/app"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCommand(appCtx *app.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ui",
		Short: "启动基于 Web 的用户界面",
		Long: `yiu-ops 提供了一个基于 Web 的用户界面，允许用户通过浏览器管理和监控系统运维任务。
该界面直观且易于使用，提供了丰富的功能，帮助用户更高效地完成运维工作。
`,
		Run: func(cmd *cobra.Command, args []string) {
			port := viper.GetInt("ui.port")
			service := NewService(appCtx)
			if err := service.Start(cmd.Context(), port); err != nil {
				cmd.PrintErrf("UI 启动失败: %v\n", err)
			}
		},
	}

	cmd.Flags().Int("ui.port", 8282, "运行服务器的端口")
	return cmd
}
