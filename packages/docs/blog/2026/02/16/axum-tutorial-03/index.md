---
title: Axum教程03
authors: [FidelYiu]
tags: [rust, axum]
---

# Axum教程03

Axum 是当前 Rust 中的 Web 框架。由 tokyo 团队构建，它速度快、类型安全，且使用起来非常优雅。

第二个模块：Axum路由。

<!-- truncate -->

## 创建模块

```sh
cargo new module-02-routing
```

## 添加依赖

`module-01-intro/Cargo.toml`

```toml
[package]
name = "module-02-routing"
version = "0.1.0"
edition = "2024"

[dependencies]
axum = { workspace = true }
tokio = { workspace = true }
serde = { workspace = true }
```

## 单个路径参数

```rust
async fn get_user(Path(id): Path<u64>) -> String {
    format!("Getting user with ID: {}", id)
}
```

该路由定义为 `/user/{id}`, 这就是 neoaxim 0.8 的语法。

在旧教程里，你会看到 `:id` 这样的写法。不要再使用那种写法了。它仍然可用，但已被弃用。

花括号才是未来，而且它们更易读。

- `Path` 是一个提取器, 它从请求路径中提取数据。内部的 `id` 是解构后的值。
- 尖括号中的 `u64` 告诉 `axum` 将路径段解析为无符号 64 位整数。
  - 如果有人请求 `/user/123`，`id` 变量等于 123。自动解析，自动类型检查。

### Path的语法

Path是一个元组结构体

```rust
pub struct Path<T>(pub T);
```

元组结构体的pub使用。

```rust
// 情况 1：结构体公开，字段也公开
pub struct Path<T>(pub T);
// 外部可以创建：Path(42)
// 外部可以访问：path.0
// 外部可以解构：let Path(value) = path;

// 情况 2：结构体公开，但字段私有
pub struct Path<T>(T);
// 外部不能创建：Path(42)  ❌ 编译错误
// 外部不能访问：path.0     ❌ 编译错误
// 外部不能解构：let Path(value) = path;  ❌ 编译错误
// 必须提供构造函数：
impl<T> Path<T> {
    pub fn new(value: T) -> Self {
        Path(value)
    }
}

// 情况 3：结构体私有
struct Path<T>(pub T);
// 外部完全不可见 ❌
```

## 多个路径参数

该路由定义为 `/users/{id}/posts/{post_id}`。

```rust
async fn get_user_post(Path((user_id, post_id)): Path<(u64, u64)>) -> String {
    format!("User {} - Post {}", user_id, post_id)
}
```

处理程序使用带元组的路径。

Axum 会按顺序提取两个值。

第一个段落对应 `user_id`, 第二个对应 `post_id`，类型安全、编译时检查。

## struct 参数

```rust
#[derive(Deserialize)]
struct PostPath {
    user_id: u64,
    post_id: u64,
    comment_id: u64,
}

async fn get_comment(Path(params): Path<PostPath>) -> String {
    format!(
        "User {} - Post {} - Comment {}",
        params.user_id, params.post_id, params.comment_id
    )
}
```

对于许多参数来说，这更整洁。

你会得到命名字段而不是元组索引。

该结构体必须派生 `Deserialize`，因为 Axum 在底层使用 `Serde`。

## 通配符路由

路径是 `/files/{*path}`

```rust
async fn files(Path(path): Path<String>) -> String {
    format!("Accessing file: {}", path)
}
```

`*` 表示捕获其后的一切内容。

请求`/files/doc/readme.md`。路径变量是 `doc/readme.md`。

这非常适合文件服务器、通配所有路由或提供静态内容。

只要记住通配符是贪婪的，所以要把它们放在特定路由之后。

## 查询参数

使用 `Query` 提取器的查询参数。

```rust
#[derive(Deserialize)]
struct Pagination {
    page: Option<u32>,
    limit: Option<u32>,
}

async fn list_items(Query(pagination): Query<Pagination>) -> String {
    let page = pagination.page.unwrap_or(1);
    let limit = pagination.limit.unwrap_or(10);
    format!("Listing items - Page: {}, Limit: {}", page, limit)
}
```

我们定义了一个 `Pagination` 结构体，包含 `page` 和 `limit`，类型为 `Option`。

可选字段不要求必须存在对应的查询参数。

处理程序接受分页查询, 对于像 `/items?page=2&limit=20` 这样的请求。这些值会被自动提取和解析。

默认值来自 `unwrap_or`；如果未提供 page，我们默认使用 1。

对可选参数的干净处理。

## 各种路由注册

```rust

/// 通过合并多个路由器创建 API v1 路由器
fn api_v1_routes() -> Router {
    Router::new()
        .nest("/users", user_routes())
        .nest("/posts", post_routes())
}

/// 也可以有多个 API 版本
fn api_v2_routes() -> Router {
    Router::new()
        .route("/users", get(|| async { "API v2 - Users endpoint" }))
        .route("/posts", get(|| async { "API v2 - Posts endpoint" }))
}

async fn not_found() -> (axum::http::StatusCode, &'static str) {
    (axum::http::StatusCode::NOT_FOUND, "404 - Route not found")
}

let app = Router::new()
    // 基本路由
    .route("/", get(|| async { "Welcome to the Routing Module!" }))
    // ===== HTTP 方法演示 =====
    // 每个方法都用一个独立的路由演示
    .route("/resource", get(|| async { "GET - Read resource" }))
    .route("/resource", post(|| async { "POST - Create resource" }))
    .route(
        "/resource/{id}",
        get(|Path(id): Path<u64>| async move { format!("GET - Read resource {}", id) }),
    )
    .route(
        "/resource/{id}",
        put(|Path(id): Path<u64>| async move { format!("PUT - Full update resource {}", id) }),
    )
    .route(
        "/resource/{id}",
        patch(|Path(id): Path<u64>| async move {
            format!("PATCH - Partial update resource {}", id)
        }),
    )
    .route(
        "/resource/{id}",
        delete(|Path(id): Path<u64>| async move { format!("DELETE - Remove resource {}", id) }),
    )
    // 路径参数（新语法！）
    .route("/users/{id}/posts/{post_id}", get(get_user_post))
    .route(
        "/users/{user_id}/posts/{post_id}/comments/{comment_id}",
        get(get_comment),
    )
    // 通配符路由（必须在具体路由之后）
    .route("/files/{*path}", get(files))
    // 查询参数
    .route("/items", get(list_items))
    .route("/search", get(search))
    // 嵌套路由 - 创建 /api/v1/users、/api/v1/posts 等
    .nest("/api/v1", api_v1_routes())
    .nest("/api/v2", api_v2_routes())
    // 未匹配路由的后备处理
    .fallback(not_found);
```

### 闭包move

`move` 定义在闭包前面，但这里有两层闭包：

```rust
// 没有参数的情况
|| async { "GET - Read resource" }
// ^外层闭包    ^内层 async 块

// 有参数的情况
|Path(id): Path<u64>| async move { format!("...", id) }
// ^外层闭包              ^内层 async 块
```

不需要 `move`：没有捕获外层变量

```rust
|| async { "GET - Read resource" }
// async 块内没使用外层的变量，不需要 move
```

需要 `async move`：async 块要用外层闭包的参数

```rust
|Path(id): Path<u64>| async move { format!("GET - Read resource {}", id) }
//        ^^^参数                            这里要用 id ^^^
// async 块需要把 id 移动进来，所以用 async move
```

## 总结

提取器是 Axum 的超级能力。它们从请求路径、查询字符串、头部提取数据，JSON 请求体，并将其转换为带类型的 Rust 值。
