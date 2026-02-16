---
title: Axum教程02
authors: [FidelYiu]
tags: [rust, axum]
---

# Axum教程02

Axum 是当前 Rust 中的 Web 框架。由 tokyo 团队构建，它速度快、类型安全，且使用起来非常优雅。

第一个模块：Hello World。

<!-- truncate -->

## 创建模块

```sh
cargo new module-01-intro
```

## 添加依赖

`module-01-intro/Cargo.toml`

```toml
[package]
name = "module-01-intro"
version = "0.1.0"
edition = "2024"

[dependencies]
axum = { workspace = true }
tokio = { workspace = true }
```

## 引入 Axum

在顶部，我们从 Axum 导入所需内容，用于定义路由的 Router，以及带有 get 和 post 函数的 routing 模块，以及用于 HTTP 的状态码.

```rust
use axum::{
    Router,
    routing::{get, post},
};
```

我们还导入了 tokio 的 tcp 监听器，因为 axum 需要一个网络监听器来接受连接。

## main 函数

main 函数带有 tokio 的 main 属性。这是至关重要的。

axom 是异步的。它使用 Rust 的 async 运行时来进行非阻塞 I/O。

但是异步函数不能自行运行。它们需要一个运行时，用来调度和执行这些函数。tokio 就是这样的运行时。

tokio 的 main 属性很神奇。

它会把你的 main 函数转换成异步函数，并将其封装在 tokio 运行时中。

没有它的话，你得手动创建运行时，这会带来更多样板代码。

```rust
#[tokio::main]
async fn main() {
    // ...
}
```

在 main 内部，我们构建了一个路由器。

## 创建路由

路由器是你的中央枢纽。它将 URL 映射到处理函数。

我们调用 `Router::new`，然后链式调用 `route`。

把路由器想象成电话交换机。当有一个呼叫——我的意思是请求——进来时，路由器会查看号码，也就是 URL 并把它连接到正确的处理器。

对于根路径 / 我们使用 GET 并传入 hello 处理器。

```rust
let app = Router::new()
    // 基本 GET 路由
    .route("/", get(hello_world))
    .route("/hello", get(hello_axum))
    .route("/health", get(health_check))
    // 带不同状态码的路由
    .route("/created", get(with_status))
    .route("/status", get(conditional_response))
    // POST 路由（后面会更深入）
    .route("/echo", post(echo));
```

每个`route`方法接受两样东西。路径模式和处理函数。

## hello 处理器

让我们看看这些处理函数。`hello` 函数非常简洁。

异步函数 `hello` 箭头 与 `&`符号 `static` `str`。它只是返回一个字符串字面量。

```rust
async fn hello_world() -> &'static str {
    "Hello, World! 🦀"
}
```

那是一个有效的 axum 处理器。没有宏，没有注解，仅仅是一个函数。

这是 axum 的优势之一。处理程序只是函数。

你可以对它们进行单独测试。你可以组合它们。你可以重用它们。没有框架魔法隐藏行为。

将这与其他一些框架进行比较，那些框架中处理函数是带有装饰器和继承的类的方法。

axom 保持简单，函数接收输入并返回输出。

axom 在返回类型上很灵活。你可以返回一个静态字符串切片，一个拥有的字符串，一个包含状态码的元组，JSON 数据，或者甚至实现你自己的 `IntoResponse` 的 trait。

## 返回状态码

health 处理器返回一个状态码和字符串的元组。这会显式设置 HTTP 状态码。

```rust
async fn health_check() -> &'static str {
    "OK"
}

use axum::http::StatusCode;

async fn with_status() -> (StatusCode, &'static str) {
    (StatusCode::CREATED, "Resource created!")
}
```

如果没有这个，Axum 默认为 200。但有时你会想要 201 Created 或 204 No Content。元组让你可以控制这一点。

因为 axum 有一个叫做 `IntoResponse` 的 trait。

任何实现了 `IntoResponse` 的类型都可以作为处理器的返回类型。

字符串实现了它。元组 实现了它。json 包装器实现了它。你也可以为你自己的类型实现它。

这是 Rust 交易系统的体现：可扩展性而无需继承。

## 创建监听器

在构建路由器之后，我们创建一个绑定到 3000 端口的 TCP 监听器。

```rust
let listener = tokio::net::TcpListener::bind("0.0.0.0:3000")
    .await
    .expect("Failed to bind to port 3000");
```

这里的 await 关键字，我们在等待操作系统为我们分配端口。

这可能会失败。如果 3000 端口已经被占用怎么办？

这就是我们使用 `expect` 调用的原因。如果绑定失败，它会提供有用的错误信息。

## 启动服务

然后我们用 `listener` 和 `router` 调用 `axum::serve` 并 `await` 它。

```rust
axum::serve(listener, app)
    .await
    .expect("Server failed to start");
```

这就是整个服务器的设置。

三行，构建路由器，绑定监听器，启动服务。干净且简洁。

## 运行项目

```sh
cargo run -p module-01-intro
```

## 总结

当一个请求到来时，处理器会运行。

如果该处理器执行 I/O、数据库查询、文件读取、API 调用等操作，它会等待。

在传统的每请求一线程模型中，该线程会被阻塞，就那样坐在那里一动不动。使用内存。

在异步中，当我们等待一个 I/O 操作时，Tokio 可以调度其他任务。

在等待数据库时，我们可以处理其他请求。

一个线程可以管理数千个连接。
