package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "启动基于 Web 的文档服务",
	Long: `yiu-operations 提供了一个基于 Web 的文档服务，允许用户通过浏览器访问和浏览系统运维相关的文档资源。
该服务提供了丰富的文档内容，帮助用户更好地理解和使用 yiu-operations 工具集。
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("docs called")
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
	docsCmd.Flags().Int("docs.port", 8383, "运行文档服务的端口")
}
