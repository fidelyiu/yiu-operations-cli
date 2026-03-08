---
title: Axum教程04
authors: [FidelYiu]
tags: [rust, axum]
---

Axum 是当前 Rust 中的 Web 框架。由 tokyo 团队构建，它速度快、类型安全，且使用起来非常优雅。

第四个模块：涵盖响应类型字符串、JSON、HTML、自定义头部、重定向，以及实现你自己的 IntoResponse。

<!-- truncate -->

## 创建模块

```sh
cargo new module-04-responses
```

## 添加依赖

`module-04-responses/Cargo.toml`

```toml
[package]
name = "module-04-responses"
version = "0.1.0"
edition = "2024"

[dependencies]
axum = { workspace = true }
tokio = { workspace = true }
serde = { workspace = true }
```

## 简单响应类型

```rust
/// 返回一个 &'static str
async fn static_string() -> &'static str {
    "来自静态字符串的问候！"
}

```

最简单的响应是字符串。返回 `&'static` 字符串.

axom 用 200 响应,以及 text-slame `content-type`。

---

```rust
/// 返回一个拥有所有权的 String
async fn owned_string() -> String {
    format!("当前时间戳问候：{}", chrono_lite())
}
```

返回相同的分配空间的字符串。

或者你还可以返回一个空元组 204。无内容。

---

```rust
fn chrono_lite() -> u64 {
    std::time::SystemTime::now()
        .duration_since(std::time::UNIX_EPOCH)
        .unwrap()
        .as_secs()
}
```

```rust
/// 返回带状态码的元组
async fn with_status() -> (StatusCode, &'static str) {
    (StatusCode::CREATED, "资源已成功创建！")
}
```

只返回状态码，而不是状态消息。

Axum 在这方面非常灵活。

## JSON 响应

对于 JSON，请将数据包装在 json 类型中。

```rust
#[derive(Serialize)]
struct User {
    id: u64,
    name: String,
    email: String,
    active: bool,
}

async fn json_user() -> Json<User> {
    Json(User {
        id: 1,
        name: "张三".to_string(),
        email: "zhangsan@example.com".to_string(),
        active: true,
    })
}

/// 返回用户列表
#[derive(Serialize)]
struct UsersResponse {
    users: Vec<User>,
    total: usize,
    page: u32,
}

async fn json_users() -> Json<UsersResponse> {
    let users = vec![
        User {
            id: 1,
            name: "张三".to_string(),
            email: "zhangsan@example.com".to_string(),
            active: true,
        },
        User {
            id: 2,
            name: "李四".to_string(),
            email: "lisi@example.com".to_string(),
            active: true,
        },
    ];
    let total = users.len();
    Json(UsersResponse {
        users,
        total,
        page: 1,
    })
}

async fn json_with_status() -> (StatusCode, Json<User>) {
    (
        StatusCode::CREATED,
        Json(User {
            id: 3,
            name: "新用户".to_string(),
            email: "newuser@example.com".to_string(),
            active: true,
        }),
    )
}
```

当你返回用户的 `JSON` 时，`Axum` 会对其进行序列化，并将内容类型设置为 `application/json` 然后发送它。

对于 API 开发来说，这就是你的基本功。

结构化数据输入，结构化数据输出。

## HTML 响应

使用 axum 响应的 html 类型返回的 html 响应。

```rust
async fn html_page() -> Html<&'static str> {
    Html(
        r#"
        <!DOCTYPE html>
        <html>
        <head>
            <title>Axum HTML Response</title>
            <style>
                body {
                    font-family: system-ui, sans-serif;
                    max-width: 800px;
                    margin: 50px auto;
                    padding: 20px;
                    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
                    min-height: 100vh;
                }
                .card {
                    background: white;
                    border-radius: 12px;
                    padding: 30px;
                    box-shadow: 0 10px 40px rgba(0,0,0,0.2);
                }
                h1 { color: #333; }
                p { color: #666; line-height: 1.6; }
            </style>
        </head>
        <body>
            <div class="card">
                <h1>🦀 欢迎来到 Axum！</h1>
                <p>这是来自你的 Axum 服务器的 HTML 响应。</p>
                <p>你可以返回完整的 HTML 页面、模板或片段。</p>
            </div>
        </body>
        </html>
        "#,
    )
}

/// 动态 HTML
async fn dynamic_html() -> Html<String> {
    let items = vec!["Routing", "Extractors", "Responses", "Middleware"];
    let list_items: String = items
        .iter()
        .map(|item| format!("<li>{}</li>", item))
        .collect();

    Html(format!(
        r#"
        <!DOCTYPE html>
        <html>
        <head>
            <title>Axum 课程模块</title>
            <style>
                body {{ font-family: system-ui; padding: 20px; }}
                ul {{ list-style-type: none; padding: 0; }}
                li {{
                    padding: 10px 15px;
                    margin: 5px 0;
                    background: #f0f0f0;
                    border-radius: 5px;
                }}
            </style>
        </head>
        <body>
            <h1>课程主题</h1>
            <ul>{}</ul>
        </body>
        </html>
        "#,
        list_items
    ))
}
```

返回一个字符串的 `html` 并且 `axum` 将内容类型设置为 `text/html`。

你可以嵌入完整的 html 页面，使用模板引擎，或任何你需要的东西。

这对于管理仪表盘、登录页或混合应用非常合适。

## 带响应头的自定义响应

设置 x 自定义头部。

```rust
async fn with_headers() -> (HeaderMap, &'static str) {
    let mut headers = HeaderMap::new();
    headers.insert(
        header::CONTENT_TYPE,
        HeaderValue::from_static("text/plain; charset=utf-8"),
    );
    headers.insert(
        header::CACHE_CONTROL,
        HeaderValue::from_static("max-age=3600"),
    );
    headers.insert("X-Custom-Header", HeaderValue::from_static("你好！"));

    (headers, "带自定义响应头的响应")
}

/// 状态码 + 响应头 + 响应体
async fn full_response() -> (StatusCode, HeaderMap, &'static str) {
    let mut headers = HeaderMap::new();
    headers.insert(header::CONTENT_TYPE, HeaderValue::from_static("text/plain"));
    headers.insert("X-Request-Id", HeaderValue::from_static("12345"));

    (StatusCode::OK, headers, "对响应拥有完全控制！")
}
```

你可以设置任意你想要的头部。

## 重定向

带有重定向管道的重定向。

```rust
async fn redirect_permanent() -> Redirect {
    Redirect::permanent("/new-location")
}
```

返回 redirect `permanent` 到 `/new-location` 表示返回 HTTP 308 Permanent Redirect 状态码重定向。

```rust
async fn redirect_temporary() -> Redirect {
    Redirect::temporary("/temp-location")
}
```

返回 redirect `temporary` 表示 307。

```rust
async fn redirect_see_other() -> Redirect {
    // 常用于表单提交之后
    Redirect::to("/success")
}

async fn new_location() -> &'static str {
    "你已被重定向到这里！"
}
```

`redirect` to 表示 303。用于已移动资源、登录重定向或提交后重定向。

这里就是它变得强大的地方。

## IntoResponse trait

IntoResponse trait 允许你定义自定义响应类型。

```rust
/// 实现 IntoResponse 的自定义响应类型
struct CustomResponse {
    message: String,
    status: StatusCode,
}

impl IntoResponse for CustomResponse {
    fn into_response(self) -> Response {
        let body = format!(
            r#"{{"message": "{}", "status": {}}}"#,
            self.message,
            self.status.as_u16()
        );

        Response::builder()
            .status(self.status)
            .header(header::CONTENT_TYPE, "application/json")
            .body(Body::from(body))
            .unwrap()
    }
}

async fn custom_response() -> CustomResponse {
    CustomResponse {
        message: "这是一个自定义响应类型！".to_string(),
        status: StatusCode::OK,
    }
}

/// 用于统一 JSON 响应的 API 响应包装器
#[derive(Serialize)]
struct ApiResponse<T: Serialize> {
    success: bool,
    data: Option<T>,
    error: Option<String>,
}

impl<T: Serialize> IntoResponse for ApiResponse<T> {
    fn into_response(self) -> Response {
        let status = if self.success {
            StatusCode::OK
        } else {
            StatusCode::BAD_REQUEST
        };

        (status, Json(self)).into_response()
    }
}

async fn api_success() -> ApiResponse<User> {
    ApiResponse {
        success: true,
        data: Some(User {
            id: 1,
            name: "张三".to_string(),
            email: "zhangsan@example.com".to_string(),
            active: true,
        }),
        error: None,
    }
}

async fn api_error() -> ApiResponse<()> {
    ApiResponse {
        success: false,
        data: None,
        error: Some("出现了问题".to_string()),
    }
}
```

## Result 响应类型

```rust
/// 处理器可以返回 Result 以进行错误处理
async fn maybe_error() -> Result<Json<User>, (StatusCode, String)> {
    let success = true; // 切换这个值以查看不同响应

    if success {
        Ok(Json(User {
            id: 1,
            name: "成功用户".to_string(),
            email: "success@example.com".to_string(),
            active: true,
        }))
    } else {
        Err((StatusCode::NOT_FOUND, "未找到用户".to_string()))
    }
}
```

## 总结

- 对于 API，请使用 JSON。
- 对于网页，请使用 JSON 或 HTML。
- 对于文件，请使用字节并设置适当的头信息。
- 对于重定向，请使用 redirect。
- 对于流式传输，请使用服务器发送事件或 WebSocket。
- 对于复杂的应用程序，构建自定义响应类型以使你的 API 格式标准化。
