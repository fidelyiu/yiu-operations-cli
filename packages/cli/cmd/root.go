package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	EnvPrefix  = "YIU_OPERATIONS"
	ConfigName = ".yiu-operations"
)

var (
	cfgFile   string
	logLevel  string
	logFormat string
	logFile   string
	logColor  bool
)

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
	ctx, cancel := context.WithCancel(context.Background())
	// trap Ctrl+C and call cancel on the context
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c:
			// 1. 当用户按下 Ctrl+C 时，会调用 cancel() 来取消上下文。
			cancel()
		case <-ctx.Done():
		}
	}()

	// 2. 将上下文传递给 rootCmd.ExecuteContext，以便在命令执行期间可以响应取消信号。
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		/*
			3. 在执行长时间运行的任务时，命令的实现应定期检查上下文的状态。
			例如：
			func (s *StackService) RunStack() error {
			    for i := 0; i < 1000000; i++ {
			        select {
			        case <-s.ctx.Done():
			            slog.Info("收到取消信号，正在清理...")
			            s.cleanup()
			            return s.ctx.Err()  // 提前返回
			        default:
			            doWork(i)
			        }
			    }
			}
		*/
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("配置文件(默认是 $HOME/%s.yaml)", ConfigName))
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "日志级别 (debug|info|warn|error)")
	rootCmd.PersistentFlags().StringVar(&logFormat, "log-format", "text", "日志格式 (text|json)")
	rootCmd.PersistentFlags().StringVar(&logFile, "log-file", "", "日志文件路径 (默认输出到标准输出)")
	rootCmd.PersistentFlags().BoolVar(&logColor, "log-color", true, "启用彩色日志输出 (仅对 text 格式有效)")
}

func initializeConfig(cmd *cobra.Command) error {
	// 1. 设置 Viper 使用环境变量。
	viper.SetEnvPrefix(EnvPrefix)
	// 允许在环境变量中使用嵌套键（例如 `YIU_OPERATIONS_DATABASE_HOST`）
	// 配置键 database.host → 环境变量 YIU_OPERATIONS_DATABASE_HOST
	// 配置键 api-key → 环境变量 YIU_OPERATIONS_API_KEY
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
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
		viper.AddConfigPath(home + "/" + ConfigName)
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

	// 5. 初始化 slog。
	if err := initLogger(); err != nil {
		return fmt.Errorf("初始化日志失败: %w", err)
	}

	// 这是一个可选但有用的步骤,用于调试你的配置。
	slog.Info("配置已初始化", "config_file", viper.ConfigFileUsed())
	return nil
}

// initLogger 初始化全局 slog 日志记录器
func initLogger() error {
	// 获取日志级别
	level := parseLogLevel(viper.GetString("log-level"))

	// 获取日志输出目标
	var writer io.Writer = os.Stdout
	logFilePath := viper.GetString("log-file")
	isFile := logFilePath != ""

	if isFile {
		file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("无法打开日志文件 %s: %w", logFilePath, err)
		}
		writer = file
		// 注意：这里不关闭文件，因为它将在整个程序生命周期中使用
	}

	// 创建处理器选项
	opts := &slog.HandlerOptions{
		Level: level,
	}

	// 根据格式创建处理器
	var handler slog.Handler
	logFormatValue := viper.GetString("log-format")
	if logFormatValue == "json" {
		handler = slog.NewJSONHandler(writer, opts)
	} else {
		// 对于 text 格式，检查是否启用颜色
		// 输出到文件时自动禁用颜色
		useColor := viper.GetBool("log-color") && !isFile
		if useColor {
			handler = tint.NewHandler(writer, &tint.Options{
				Level:      level,
				TimeFormat: time.DateTime,
			})
		} else {
			handler = slog.NewTextHandler(writer, opts)
		}
	}

	// 设置全局默认日志记录器
	slog.SetDefault(slog.New(handler))

	return nil
}

// parseLogLevel 解析日志级别字符串
func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
