package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "yiu-operations",
	Short: "yiu-operations 系统运维工具",
	Long: `yiu-operations 是一个用于系统运维的强大工具集。
	
你可以使用它来简化和自动化各种运维任务，提高效率和可靠性。
	`,
	// 如果你的基础应用程序有与之关联的操作，
	// 请取消注释以下行：
	// Run: func(cmd *cobra.Command, args []string) { },
	// PersistentPreRunE 在解析标志之后、命令的 RunE 函数调用之前被调用。
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initializeConfig(cmd)
	},
}

// Execute 将所有子命令添加到根命令并适当地设置标志。
// 这由 main.main() 调用。它只需要对 rootCmd 执行一次。
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件(默认是 $HOME/.yiu-operations.yaml)")
}

func initializeConfig(cmd *cobra.Command) error {
	// 1. 设置 Viper 使用环境变量。
	viper.SetEnvPrefix("MYAPP")
	// 允许在环境变量中使用嵌套键（例如 `MYAPP_DATABASE_HOST`）
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "*", "-", "*"))
	viper.AutomaticEnv()
	// 2. 处理配置文件。
	if cfgFile != "" {
		// 使用标志指定的配置文件。
		viper.SetConfigFile(cfgFile)
	} else {
		// 在默认位置搜索配置文件。
		home, err := os.UserHomeDir()
		// 只有在无法获取主目录时才 panic。
		cobra.CheckErr(err)

		// 搜索名为 "config" 的配置文件（不带扩展名）。
		viper.AddConfigPath(".")
		viper.AddConfigPath(home + "/.yiu-operations")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	// 3. 读取配置文件。
	// 如果找到配置文件，则读取它。我们使用健壮的错误检查
	// 来忽略 "文件未找到" 错误，但对任何其他错误进行 panic。
	if err := viper.ReadInConfig(); err != nil {
		// 如果配置文件不存在，这是可以接受的。
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return err
		}
	}

	// 4. 将 Cobra 标志绑定到 Viper。
	// 这是使标志值通过 Viper 可用的魔法。
	// 它绑定传入命令的完整标志集。
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return err
	}

	// 这是一个可选但有用的步骤，用于调试你的配置。
	fmt.Println("配置已初始化。使用的配置文件：", viper.ConfigFileUsed())
	return nil
}
