# Axum教程-12-生产就绪

优雅关机、结构化日志、健康检查、容器化。

这些不是可选项。它们决定了一个项目是业余爱好还是专业系统。

让你的 Axum 应用达到可投入生产的状态。

## 创建模块

```sh
cargo new module-12-production
```

## 添加依赖

`module-12-production/Cargo.toml`

```toml
[package]
name = "module-12-production"
version = "0.1.0"
edition = "2024"

[dependencies]
axum = { workspace = true }
tokio = { workspace = true }
serde = { workspace = true }
serde_json = { workspace = true }
tower = { workspace = true }
tower-http = { workspace = true }
tracing = { workspace = true }
tracing-subscriber = { workspace = true }
```

## Docker

Docker Compose 用于编排多个服务。

它非常适合带数据库和依赖服务的本地开发。

我们为数据库定义了 Postgres。

我们为服务器定义了 app。

`Dockerfile`

该项目包含一个使用多阶段构建的 Dockerfile。

第一阶段构建。

从 rust slim 开始。安装依赖。使用 release 标志进行编译。

第二阶段运行时。

从 debian slim 开始。

只复制二进制文件。没有 Rust 工具链、没有源代码、没有构建产物。最终镜像小于 100 MB。快速在网络中拉取，快速启动，攻击面最小。

```text
FROM rust:1.78-slim as builder
WORKDIR /app
COPY . .
RUN cargo build --release -p module-12-production

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/target/release/module-12-production /usr/local/bin/app
EXPOSE 3000
CMD ["app"]
```

应用依赖于 Postgres，docker compose 按顺序启动它们。

环境变量配置了两个服务。

- DATABASE_URL 指向 Postgres 容器。

`docker-compose.yml`

```yaml
services:
  postgres:
    image: postgres:16
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: axum_course
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

  app:
    build: .
    ports:
      - "3000:3000"
    environment:
      DATABASE_URL: postgres://postgres:postgres@postgres/axum_course
      RUST_LOG: info
    depends_on:
      - postgres

volumes:
  pgdata:
```

运行 `docker compose up`，一切就启动了。

数据库迁移由应用处理。

## 优雅关机

当服务器在没有优雅关机的情况下收到信号时，会发生什么？

连接会断开，请求会失败。用户会看到错误。

使用优雅关机时，服务器会停止接受新连接，但会完成正在进行的请求。大家都很高兴。

kubernetes 在终止 Pod 之前发送 SIGTERM。awesome 在停止实例之前发送它。你的部署脚本在更新期间发送它。

优雅地处理它，否则就面临愤怒的用户。

```rust
use axum::{extract::State, routing::get, Json, Router};
use std::{
    sync::{
        atomic::{AtomicBool, AtomicU64, Ordering},
        Arc,
    },
    time::Duration,
};
use tokio::net::TcpListener;
use tower_http::{compression::CompressionLayer, trace::TraceLayer};
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

// ============================================================================
// 应用状态
// ============================================================================

#[derive(Clone)]
struct AppState {
    ready: Arc<AtomicBool>,
    request_count: Arc<AtomicU64>,
}

impl Default for AppState {
    fn default() -> Self {
        Self {
            ready: Arc::new(AtomicBool::new(true)),
            request_count: Arc::new(AtomicU64::new(0)),
        }
    }
}

// ============================================================================
// 健康检查与就绪检查
// ============================================================================

async fn health() -> &'static str {
    "正常"
}

async fn ready(
    State(state): State<AppState>,
) -> Result<&'static str, (axum::http::StatusCode, &'static str)> {
    if state.ready.load(Ordering::SeqCst) {
        Ok("就绪")
    } else {
        Err((axum::http::StatusCode::SERVICE_UNAVAILABLE, "未就绪"))
    }
}

async fn metrics(State(state): State<AppState>) -> Json<serde_json::Value> {
    Json(serde_json::json!({
        "requests": state.request_count.load(Ordering::SeqCst),
        "ready": state.ready.load(Ordering::SeqCst)
    }))
}

async fn index(State(state): State<AppState>) -> &'static str {
    state.request_count.fetch_add(1, Ordering::SeqCst);
    "来自生产就绪版 Axum 的问候！"
}

// ============================================================================
// 主程序
// ============================================================================

#[tokio::main]
async fn main() {
    // 初始化 tracing（用于生产环境的结构化 JSON 日志）
    tracing_subscriber::registry()
        .with(
            tracing_subscriber::fmt::layer()
                .json()
                .with_target(true)
                .with_current_span(true),
        )
        .with(tracing_subscriber::EnvFilter::new(
            std::env::var("RUST_LOG").unwrap_or_else(|_| "info".into()),
        ))
        .init();

    let state = AppState::default();

    let app = Router::new()
        .route("/", get(index))
        .route("/health", get(health)) // 存活探针
        .route("/ready", get(ready)) // 就绪探针
        .route("/metrics", get(metrics))
        .with_state(state.clone())
        .layer(TraceLayer::new_for_http())
        .layer(CompressionLayer::new());

    let listener = TcpListener::bind("0.0.0.0:3000").await.unwrap();

    tracing::info!("🚀 服务启动于 http://localhost:3000");

    // 优雅关闭
    axum::serve(listener, app)
        .with_graceful_shutdown(shutdown_signal(state))
        .await
        .unwrap();

    tracing::info!("服务已优雅关闭");
}

async fn shutdown_signal(state: AppState) {
    let ctrl_c = async {
        tokio::signal::ctrl_c()
            .await
            .expect("安装 Ctrl+C 处理器失败");
    };

    #[cfg(unix)]
    let terminate = async {
        tokio::signal::unix::signal(tokio::signal::unix::SignalKind::terminate())
            .expect("安装信号处理器失败")
            .recv()
            .await;
    };

    #[cfg(not(unix))]
    let terminate = std::future::pending::<()>();

    tokio::select! {
        _ = ctrl_c => {},
        _ = terminate => {},
    }

    tracing::info!("已收到关闭信号，开始优雅关闭");

    // 标记为不再接受新连接
    state.ready.store(false, Ordering::SeqCst);

    // 给负载均衡器留出探测时间
    tokio::time::sleep(Duration::from_secs(5)).await;
}
```

## Json日志

设置 JSON 输出的跟踪（tracing）

这对于生产环境的可观测性至关重要。

JSON 日志是机器可解析的。

你的日志聚合器 CloudWatch、Datadog、Elasticsearch 可以自动索引字段。

按用户 ID 过滤。

按错误类型搜索。

按端点聚合。

可读性强的日志对开发很有帮助。

JSON 日志对生产环境是必要的。

将 Rust 的日志环境变量设置为控制日志级别。

## 健康和就绪端点

每个生产应用都需要两者。

/health 端点是存活性检查。进程在运行吗？如果是，则返回 200。如果这失败，你的编排器会重启该 Pod。

/ready 端点是就绪探针。流量应该路由到这里吗？如果就绪则返回 200。如果不是的话则返回 53。在关闭或初始化期间，应用可能仍在运行但尚未准备好接受流量。

Kubernetes、ECS、Nomad 等每个编排器都使用这些探针。

正确设置它们。

生产环境应用应当暴露指标、请求计数、延迟和错误率。

/metrics 端点返回包含计数器的 JSON。

在真实应用中，与 Prometheus 集成。

使用 prometheus crate 以标准格式暴露指标。

指标为仪表盘和告警提供数据。没有它们，你就是在盲飞。
