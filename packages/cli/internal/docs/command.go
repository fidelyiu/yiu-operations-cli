package docs

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"

	"yiu-ops/internal/app"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:embed build
var docsFS embed.FS

func NewCommand(appCtx *app.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docs",
		Short: "启动基于 Web 的文档服务",
		Long: `yiu-ops 提供了一个基于 Web 的文档服务，允许用户通过浏览器访问和浏览系统运维相关的文档资源。
该服务提供了丰富的文档内容，帮助用户更好地理解和使用 yiu-ops 工具集。

使用案例:
    yiu-ops docs --port 8383
`,
		Run: func(cmd *cobra.Command, args []string) {
			host := viper.GetString("docs.host")
			port := viper.GetInt("docs.port")

			buildFS, err := fs.Sub(docsFS, "build")
			if err != nil {
				cmd.PrintErrf("无法加载文档资源: %v\n", err)
				return
			}

			if appCtx != nil && appCtx.Logger != nil {
				appCtx.Logger.Info(fmt.Sprintf("Docs 服务启动, 访问 http://%s:%d", host, port))
			}

			service := NewService(buildFS)
			err = service.Serve(cmd.Context(), host, port)
			if err != nil && !errors.Is(err, context.Canceled) {
				cmd.PrintErrf("Docs server 启动失败: %v\n", err)
			}
		},
	}

	cmd.Flags().String("host", "localhost", "运行文档服务的主机地址")
	if err := viper.BindPFlag("docs.host", cmd.Flags().Lookup("host")); err != nil {
		cmd.PrintErrf("绑定 docs.host 失败: %v\n", err)
	}
	cmd.Flags().IntP("port", "p", 8383, "运行文档服务的端口")
	if err := viper.BindPFlag("docs.port", cmd.Flags().Lookup("port")); err != nil {
		cmd.PrintErrf("绑定 docs.port 失败: %v\n", err)
	}
	return cmd
}
