---
title: Axum教程04
authors: [FidelYiu]
tags: [rust, axum]
---

# Axum教程04

Axum 是当前 Rust 中的 Web 框架。由 tokyo 团队构建，它速度快、类型安全，且使用起来非常优雅。

第二个模块：Axum路由。

<!-- truncate -->

## 创建模块

```sh
cargo new module-03-extractors
```

## 添加依赖

`module-03-extractors/Cargo.toml`

```toml
[package]
name = "module-03-extractors"
version = "0.1.0"
edition = "2024"

[dependencies]
axum = { workspace = true }
tokio = { workspace = true }
serde = { workspace = true }
```

## 内置提取器

### Path

提取路径参数

```rust
/// 路由：GET /users/{id}
async fn get_user(Path(id): Path<u64>) -> String {
    format!("User ID: {}", id)
}
```

### Query

提取查询字符串参数

```rust
#[derive(Debug, Deserialize)]
struct ListParams {
    page: Option<u32>,
    limit: Option<u32>,
    sort: Option<String>,
}

async fn list_users(Query(params): Query<ListParams>) -> String {
    format!(
        "Page: {}, Limit: {}, Sort: {}",
        params.page.unwrap_or(1),
        params.limit.unwrap_or(10),
        params.sort.unwrap_or_else(|| "id".to_string())
    )
}
```

### Json

提取并反序列化 JSON 请求体

```rust
#[derive(Debug, Deserialize)]
struct CreateUserRequest {
    name: String,
    email: String,
}

#[derive(Debug, Serialize)]
struct CreateUserResponse {
    id: u64,
    name: String,
    email: String,
}

async fn create_user(Json(payload): Json<CreateUserRequest>) -> Json<CreateUserResponse> {
    Json(CreateUserResponse {
        id: 1,
        name: payload.name,
        email: payload.email,
    })
}
```

`CreateUserRequest` 结构体有 `name` 和 `email` 字段，都是字符串。

处理器接受 `CreateUserRequest` 的 json。

当收到一个 POST 请求时, 以 JSON 格式的 body 传入。

axum 会自动将其解析到你的结构体中。

如果解析失败，比如 JSON 格式错误、类型错误或缺少必需字段，

axum 会返回 422 Unprocessable Entity（无法处理的实体）错误。

你不需要验证代码。类型系统会处理它。

### Headers

访问请求头

```rust
async fn show_headers(headers: HeaderMap) -> String {
    let user_agent = headers
        .get("user-agent")
        .and_then(|v| v.to_str().ok())
        .unwrap_or("Unknown");

    let content_type = headers
        .get("content-type")
        .and_then(|v| v.to_str().ok())
        .unwrap_or("Not specified");

    format!("User-Agent: {}\nContent-Type: {}", user_agent, content_type)
}
```

处理器接受名为 headers 的 `HeaderMap`，并且可以访问任何头部。

使用 `headers.get("content-type")` 获取内容类型。

使用 `headers.get` 获取授权（authorization）以进行身份验证。

完全访问客户端发送的所有内容。

对于特定的头，你可以只提取所需的部分。

### 原始请求体

```rust
async fn raw_body(body: Bytes) -> String {
    format!("Received {} bytes", body.len())
}
```

有时你需要原始字节。

bytes 提取器会以字节缓冲区的形式给你请求体。

这对于文件上传、二进制协议等很有用，

或者当你需要对请求体计算哈希时。

不做解析，只提供原始数据。

### 多个提取器

```rust
async fn combined_extractors(
    Path(id): Path<u64>,
    Query(params): Query<ListParams>,
    headers: HeaderMap,
    Json(body): Json<CreateUserRequest>, // Must be last!
) -> String {
    format!(
        "ID: {}\nPage: {:?}\nUser-Agent: {:?}\nName: {}",
        id,
        params.page,
        headers.get("user-agent"),
        body.name
    )
}
```

处理器可以接受多个提取器。axum 按顺序提取每一项。

先是 state，然后是 path，再是 query，最后是 json。

全都强类型、全都经过验证、在你的代码运行之前全部完成。

这就是最精妙的组合。每个提取器只做一件事。将它们组合以实现复杂的行为。

这里有一件关键的事。提取器的顺序很重要。

像 `json` 和 `bytes` 这样的会消耗请求体的提取器会读取整个请求体。

一旦被消耗，就不复存在。

你不能在同一个处理器中使用两个会消耗请求体的提取器。

像 `Path`、`Query` 和 `HeaderMap` 这样的非消耗性提取器可以按任意顺序放置。

把body提取器放在最后。

如果你搞错了，会得到关于body已被消耗的晦涩错误。

### 可选提取器

```rust
async fn optional_query(Query(params): Query<Option<ListParams>>) -> String {
    match params {
        Some(p) => format!("Got params: page={:?}", p.page),
        None => "No query params provided".to_string(),
    }
}
```

如果提取器可能没有数据怎么办？使用 query 的 Option 或 header 的 Option。

如果数据不存在，你会得到 None，而不是错误。

这对可选的查询参数或可能存在也可能不存在的头非常适用。

在 Axum 0.8 中，`Option<T>` 作为提取器需要实现 `OptionalFromRequestParts` 或 `OptionalFromRequest`

这允许更好的错误处理 - 拒绝可以转换为错误响应而不是被静默忽略。

像 `Query` 和 `HeaderMap` 这样的内置类型已经实现了这个。

## 自定义提取器

### 获取特定请求头

```rust
#[derive(Debug)]
struct ApiKeyError;

impl IntoResponse for ApiKeyError {
    fn into_response(self) -> Response {
        (
            StatusCode::UNAUTHORIZED,
            "Missing or invalid API key. Provide X-API-Key header.",
        )
            .into_response()
    }
}

// AXUM 0.8 新增：不需要 #[async_trait]！
impl<S> FromRequestParts<S> for ApiKey
where
    S: Send + Sync,
{
    type Rejection = ApiKeyError;

    // 注意：我们返回 `impl Future` 而不是使用 #[async_trait]
    fn from_request_parts(
        parts: &mut Parts,
        _state: &S,
    ) -> impl std::future::Future<Output = Result<Self, Self::Rejection>> + Send {
        // 我们使用 async 块来创建 future
        let api_key = parts
            .headers
            .get("x-api-key")
            .and_then(|v| v.to_str().ok())
            .map(|s| s.to_string());

        async move {
            match api_key {
                Some(key) if !key.is_empty() => Ok(ApiKey(key)),
                _ => Err(ApiKeyError),
            }
        }
    }
}

async fn protected_endpoint(ApiKey(key): ApiKey) -> String {
    format!("Access granted! Your API key: {}", key)
}
```

我们定义了一个 `ApiKey` 结构体，用来保存一个字符串。

我们为它实现了 `FromRequestParts`。

这个 `trait` 让 `axom` 能把它当作一个抽取器使用。

该实现读取 `x-api-key` 头。

如果找到，返回应用密钥头。

如果缺失，返回 401 未经授权。

现在任何处理函数都可以将应用密钥作为参数。

如果该头缺失，axom 会在处理器执行之前就运行你的拒绝逻辑。

在提取器层面的安全性。

避免每个处理器中的样板代码。

这种模式非常强大。

你可以为认证、速率限制、请求验证以及任何你反复检查的东西构建提取器。

状态提取器值得拥有自己的模块。不过我们先预览一下。

状态保存了共享的应用状态：数据库连接、配置、缓存。现在先知道，状态和其他提取器一样。

#### 解析语法

Axum 为许多常见的元组组合都实现了 `IntoResponse` trait。

```rust
// 这些都实现了 IntoResponse
(StatusCode, String)
(StatusCode, &str)
(StatusCode, Body)
(StatusCode, HeaderMap, Body)
(StatusCode, Headers, Body)
// 等等...
```

---

```rust
impl<S> FromRequestParts<S> for ApiKey
```

为 `ApiKey` 类型实现 `FromRequestParts<S> trait`
`<S>` 是泛型参数，代表任意状态类型

---

```rust
where
    S: Send + Sync,
```

- `S` 必须实现 `Send` 和 `Sync` trait
- `Send`：可以在线程间安全传递所有权
- `Sync`：可以在线程间安全共享引用
- 这保证状态对象在异步环境中是安全的
- 目前暂且认为他们是Axum必须要的trait吧，太复杂了。

---

```rust
type Rejection = ApiKeyError;
```

关联类型

定义提取失败时返回的错误类型
如果没有 API key 或为空，就返回 `ApiKeyError`

---

```rust
fn from_request_parts(
    parts: &mut Parts,
    _state: &S,
) -> impl std::future::Future<Output = Result<Self, Self::Rejection>> + Send
```

| 部分                                     | 说明                                                         |
| :--------------------------------------- | :----------------------------------------------------------- |
| `parts: &mut Parts`                      | 请求头部分（包含请求头信息）                                 |
| `_state: &S`                             | 应用状态（下划线表示不使用）                                 |
| `-> impl Future<...>`                    | 不需要 #[async_trait]，直接返回 impl Future                  |
| `Output = Result<Self, Self::Rejection>` | Future 的输出类型：成功返回 `ApiKey`，失败返回 `ApiKeyError` |
| `+ Send`                                 | 这个 Future 可以在线程间移动                                 |

---

```rust
let api_key = parts
    .headers
    .get("x-api-key")      // from 请求头中获取 "x-api-key"
    .and_then(|v| v.to_str().ok())  // 转换为字符串
    .map(|s| s.to_string());         // 转换为 String（拥有所有权）
```

此时 `api_key` 的类型是 `Option<String>`：

如果存在且有效：`Some(String)`
否则：`None`

---

```rust
async move {
    match api_key {
        Some(key) if !key.is_empty() => Ok(ApiKey(key)),  // API key 存在且非空
        _ => Err(ApiKeyError),                            // 不存在或为空字符串
    }
}
```

- `async move` 闭包捕获 `api_key` 的所有权
- `if !key.is_empty()` 是守卫条件（guard），额外检查是否为空字符串
- 成功时返回 `Ok(ApiKey(key))`
- 失败时返回 `Err(ApiKeyError)`

### 验证 JSON 请求体

```rust
#[derive(Debug, Deserialize)]
struct ValidatedUser {
    name: String,
    email: String,
}

struct ValidatedJson<T>(T);

#[derive(Debug)]
enum ValidationError {
    InvalidJson(String),
    InvalidEmail,
    NameTooShort,
}

impl IntoResponse for ValidationError {
    fn into_response(self) -> Response {
        let (status, message) = match self {
            ValidationError::InvalidJson(e) => {
                (StatusCode::BAD_REQUEST, format!("无效的 JSON：{}", e))
            }
            ValidationError::InvalidEmail => {
                (StatusCode::BAD_REQUEST, "无效的邮箱格式".to_string())
            }
            ValidationError::NameTooShort => {
                (StatusCode::BAD_REQUEST, "名字至少需要 2 个字符".to_string())
            }
        };
        (status, message).into_response()
    }
}

// 验证请求体的自定义提取器
// 注意：对于请求体提取器，我们实现 FromRequest 而不是 FromRequestParts
impl<S> FromRequest<S> for ValidatedJson<ValidatedUser>
where
    S: Send + Sync,
{
    type Rejection = ValidationError;

    fn from_request(
        req: Request,
        state: &S,
    ) -> impl std::future::Future<Output = Result<Self, Self::Rejection>> + Send {
        async move {
            // 首先提取 JSON
            let Json(user): Json<ValidatedUser> = Json::from_request(req, state)
                .await
                .map_err(|e| ValidationError::InvalidJson(e.to_string()))?;

            // 验证名字长度
            if user.name.len() < 2 {
                return Err(ValidationError::NameTooShort);
            }

            // 验证邮箱（简单检查）
            if !user.email.contains('@') {
                return Err(ValidationError::InvalidEmail);
            }

            Ok(ValidatedJson(user))
        }
    }
}

async fn create_validated_user(ValidatedJson(user): ValidatedJson<ValidatedUser>) -> String {
    format!("Created user: {} <{}>", user.name, user.email)
}
```

## 状态作为提取器

```rust
#[derive(Clone)]
struct AppState {
    db_pool: String, // 在真实应用中，这将是数据库连接池
    api_version: String,
}

async fn with_state(State(state): State<Arc<AppState>>) -> String {
    format!("API Version: {}, DB: {}", state.api_version, state.db_pool)
}
```
