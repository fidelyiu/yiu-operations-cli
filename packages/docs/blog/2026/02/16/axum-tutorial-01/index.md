---
title: Axum教程01
authors: [FidelYiu]
tags: [rust, axum]
---

# Axum教程01

Axum 是当前 Rust 中的 Web 框架。由 tokyo 团队构建，它速度快、类型安全，且使用起来非常优雅。

初始化项目。

<!-- truncate -->

## 相关链接

- [Youtubeo](https://www.youtube.com/watch?v=Ka7mRKsTCyE)
- [Github](https://github.com/aarambh-darshan/axum-full-course)
- [Cargo 手册](https://doc.rust-lang.org/cargo/index.html)
- [rust中文指南](https://kaisery.github.io/trpl-zh-cn/)
- [通过例子学 Rust](https://rustwiki.org/zh-CN/rust-by-example/)

## 初始 workspace

```sh
mkdir axum-tutorial
cd axum-tutorial
```

这是一个 cargo 工作区，由 12 个独立的 Rust crate 组成。

创建 `Cargo.toml` 文件。

> workspace 支持的[字段](https://doc.rust-lang.org/cargo/reference/workspaces.html)

```toml
[workspace]
resolver = "3"

# 定义所有模块的公共依赖版本
[workspace.dependencies]
# Axum 核心框架
axum = "0.8.8"
axum-extra = { version = "0.12", features = [
    "typed-header",
    "query",
    "multipart",
] }

# 异步运行时
tokio = { version = "1.49", features = ["full"] }

# 序列化
serde = { version = "1.0", features = ["derive"] }
serde_json = "1.0"

# Tower 生态系统
tower = { version = "0.5", features = ["full"] }
tower-http = { version = "0.6", features = [
    "cors",
    "compression-gzip",
    "timeout",
    "trace",
    "fs",
    "limit",
] }
tower-service = "0.3"

# 数据库
sqlx = { version = "0.8", features = [
    "runtime-tokio",
    "postgres",
    "uuid",
    "chrono",
] }

# 认证和安全
jsonwebtoken = "10.3"
argon2 = "0.5"

# 错误处理
thiserror = "2.0"
anyhow = "1.0"

# 日志 & 追踪
tracing = "0.1"
tracing-subscriber = { version = "0.3", features = ["env-filter", "json"] }

# 其他实用工具
uuid = { version = "1.21", features = ["v4", "serde"] }
chrono = { version = "0.4", features = ["serde"] }
dotenvy = "0.15"
futures = "0.3"

# 测试
http-body-util = "0.1"
```

## 添加 gitignore

`.gitignore`

```
/target
Cargo.lock
.env
*.log
```
