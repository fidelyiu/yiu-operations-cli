---
slug: yiu-ops-docs
title: Yiu Ops 的 docs 命令开发
authors: [FidelYiu]
tags: [go, yiuOps]
---

# Yiu Ops 的 docs 命令开发

Yiu Ops 的 docs 命令开发过程中的设计和实现。

<!-- truncate -->

## 目标

```bash
# 当执行 docs 子命令之后
yiu-ops docs
# 我们就可以在浏览器中访问 yiu ops cli 的文档
# http://localhost:8281
```

## 添加子命令

```bash
cobra-cli add docs
```

## go的文件系统使用

我们在 `packages/docs/package.json` 中将web资源构建到了 `packages/cli/internal/docs/build`。

```go
//go:embed build
var docsFS embed.FS

// ...

buildFS, err := fs.Sub(docsFS, "build")

// ...

server := &http.Server{
    Addr:    fmt.Sprintf("%s:%d", host, port),
    Handler: http.FileServer(http.FS(s.buildFS)),
}
```

## 监听 Ctrl C

### root中监听

```go
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
```

### docs中使用

```go
func (s *Service) Serve(ctx context.Context, host string, port int) error {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: http.FileServer(http.FS(s.buildFS)),
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
    // 用户点击 Ctrl C 之后，会执行到这里
	case <-ctx.Done():
        // 这里会创建一个5s超时的上下文
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
        // 然后这个服务5s内关闭，不在接收新的服务
		_ = server.Shutdown(shutdownCtx)
        // 将 error 返回，让cli知道这是用户点击了 Ctrl C
		return ctx.Err()
	case err := <-errCh:
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}
}
```
