# Axum教程-06-中间件

中间件在处理器之前或之后运行，用于记录、认证、压缩。

第六模块教你使用内置中间件并构建自定义中间件。

tower 是 axom 使用的服务层。

layer 包裹服务以添加行为。

把层想象成洋葱。

请求进入，依次通过每一层(layer)，直到到达处理器。

响应再通过这些层返回。

每一层都可以修改请求或响应。

## 创建模块

```sh
cargo new module-06-middleware
```

## 添加依赖

`module-06-middleware/Cargo.toml`

```toml
[package]
name = "module-06-middleware"
version = "0.1.0"
edition = "2024"

[dependencies]
axum = { workspace = true }
tokio = { workspace = true }
serde = { workspace = true }
serde_json = { workspace = true }
uuid = { workspace = true }
```

## 使用案例

```rust
//! # 模块 06：中间件与分层
//!
//! Axum 中对 Tower 中间件的集成：
//! - 内置中间件（CORS、压缩、超时）
//! - 使用 from_fn 自定义中间件
//! - 针对特定路由的层

use axum::{
    extract::Request,
    http::{header, HeaderValue, Method, StatusCode},
    middleware::{self, Next},
    response::{IntoResponse, Response},
    routing::get,
    Router,
};
use std::time::{Duration, Instant};
use tower::ServiceBuilder;
use tower_http::{
    compression::CompressionLayer,
    cors::{Any, CorsLayer},
    trace::TraceLayer,
};
use tracing::Level;

// ============================================================================
// 课程 1：使用 from_fn 自定义中间件
// ============================================================================

/// 日志中间件：记录每一个请求
async fn logging_middleware(request: Request, next: Next) -> Response {
    let method: Method = request.method().clone();
    let uri = request.uri().clone();
    let start = Instant::now();

    let response = next.run(request).await;

    tracing::info!(
        method = %method,
        uri = %uri,
        status = %response.status().as_u16(),
        duration_ms = %start.elapsed().as_millis(),
        "Request completed"
    );
    response
}

/// 计时中间件：添加 X-Response-Time 响应头
async fn timing_middleware(request: Request, next: Next) -> Response {
    let start = Instant::now();
    let mut response = next.run(request).await;

    response.headers_mut().insert(
        "X-Response-Time",
        HeaderValue::from_str(&format!("{}ms", start.elapsed().as_millis())).unwrap(),
    );
    response
}

/// 认证中间件
async fn auth_middleware(request: Request, next: Next) -> Result<Response, StatusCode> {
    let auth_header = request
        .headers()
        .get("X-API-Key")
        .and_then(|v| v.to_str().ok());

    match auth_header {
        Some("secret-key") => Ok(next.run(request).await),
        _ => Err(StatusCode::UNAUTHORIZED),
    }
}

// ============================================================================
// 课程 2：内置的 Tower-HTTP 中间件
// ============================================================================

fn cors_layer() -> CorsLayer {
    CorsLayer::new()
        .allow_origin(Any)
        .allow_methods([Method::GET, Method::POST, Method::PUT, Method::DELETE])
        .allow_headers([header::CONTENT_TYPE, header::AUTHORIZATION])
}

// ============================================================================
// 处理函数
// ============================================================================

async fn index() -> &'static str {
    "Welcome to Axum Middleware Module!"
}

async fn public_data() -> impl IntoResponse {
    axum::Json(serde_json::json!({"message": "Public data", "accessible": true}))
}

async fn protected_data() -> impl IntoResponse {
    axum::Json(serde_json::json!({"message": "Secret data", "authorized": true}))
}

async fn slow_endpoint() -> &'static str {
    tokio::time::sleep(Duration::from_secs(1)).await;
    "Slow operation done!"
}

// ============================================================================
// 主函数
// ============================================================================

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt().with_max_level(Level::INFO).init();

    // 受保护的路由（需要认证）
    let protected = Router::new()
        .route("/data", get(protected_data))
        // 加在指定路由上
        .route_layer(middleware::from_fn(auth_middleware));

    // 带有分层中间件的主应用
    let app = Router::new()
        .route("/", get(index))
        .route("/public", get(public_data))
        .route("/slow", get(slow_endpoint))
        .nest("/protected", protected)
        // 加在全局路由上
        .layer(middleware::from_fn(timing_middleware))
        .layer(middleware::from_fn(logging_middleware))
        .layer(
            ServiceBuilder::new()
                // TraceLayer 会记录每个请求的方法、路径、状态码和持续时间。
                // 在生产环境中这是不可或缺的。
                .layer(TraceLayer::new_for_http())
                .layer(cors_layer())
                // CompressionLayer 会自动压缩响应，
                // 浏览器会发送 accept-encoding: gzip，你的服务器会以压缩数据响应。从而节省带宽。
                .layer(CompressionLayer::new()),
        );

    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();

    println!("🚀 Module 06: Middleware & Layers");
    println!("   Server: http://localhost:3000");
    println!("\n📝 Endpoints:");
    println!("   GET /              - Welcome");
    println!("   GET /public        - Public data");
    println!("   GET /slow          - Slow endpoint");
    println!("   GET /protected/data - Auth required (X-API-Key: secret-key)");

    axum::serve(listener, app).await.unwrap();
}
```

## 顺序

layer 的顺序很重要。

layer 以相反的顺序应用。最后添加的 layer 最先运行。

如果你先添加了日志记录，然后添加了认证。请求会先经过认证，然后再经过日志记录。

## ServiceBuilder

整洁地应用多个 layer。

ServiceBuilder 中新链层的调用会传递给路由器的 layer 方法。

它从上到下读取，但应用时从下到上。

## 总结

下面是一个典型的中间件堆栈。

- 最外层的跟踪层用于日志记录。
- 下一个压缩层。
- 下一个超时层（如果有的话）。
- 下一个速率限制层。
- 最内层的认证，用于受保护的路由。
