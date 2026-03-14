# Axum教程-07-错误处理

错误会发生，网络会出错，数据会无效。

资源不存在——你如何处理错误定义了你的 API 的。

## 创建模块

```sh
cargo new module-07-errors
```

## 添加依赖

`module-07-errors/Cargo.toml`

```toml
[package]
name = "module-07-errors"
version = "0.1.0"
edition = "2024"

[dependencies]
axum = { workspace = true }
tokio = { workspace = true }
serde = { workspace = true }
serde_json = { workspace = true }
uuid = { workspace = true }
```

## 错误处理

```rust
//! # 模块 07：错误处理
//!
//! Axum 中正确的错误处理方式：
//! - 使用 thiserror 自定义错误类型
//! - 为错误实现 IntoResponse
//! - 基于 Result 的处理函数
//! - 错误恢复模式

use axum::{
    extract::Path,
    http::StatusCode,
    response::{IntoResponse, Response},
    routing::get,
    Json, Router,
};
use serde::Serialize;
use thiserror::Error;

// ============================================================================
// 课程 1：使用 thiserror 自定义错误类型
// ============================================================================

// 带有此错误派生的应用错误枚举。
#[derive(Error, Debug)]
#[allow(dead_code)] // 这些变体仅用于演示
enum AppError {
    #[error("用户未找到：{0}")]
    UserNotFound(u64),

    #[error("输入无效：{0}")]
    InvalidInput(String),

    #[error("数据库错误：{0}")]
    DatabaseError(String),

    #[error("未授权")]
    Unauthorized,

    #[error("服务器内部错误")]
    Internal,
}

// ============================================================================
// 课程 2：为自定义错误实现 IntoResponse
// ============================================================================

#[derive(Serialize)]
struct ErrorResponse {
    error: String,
    code: u16,
}

// AppError 实现了 IntoResponse
impl IntoResponse for AppError {
    fn into_response(self) -> Response {
        // 我们对错误变体进行匹配。
        let (status, message) = match &self {
            // 用户未找到返回 404。
            AppError::UserNotFound(_) => (StatusCode::NOT_FOUND, self.to_string()),
            // 无效输入返回 404。
            AppError::InvalidInput(_) => (StatusCode::BAD_REQUEST, self.to_string()),
            // 数据库错误返回 500。
            AppError::DatabaseError(_) => (StatusCode::INTERNAL_SERVER_ERROR, self.to_string()),
            // 未认证返回401
            AppError::Unauthorized => (StatusCode::UNAUTHORIZED, self.to_string()),
            // 内部错误返回 500。
            AppError::Internal => (StatusCode::INTERNAL_SERVER_ERROR, self.to_string()),
        };

        // 我们创建了一个包含错误信息和状态码的错误响应结构体。
        let body = ErrorResponse {
            error: message,
            code: status.as_u16(),
        };

        // 将其作为 JSON 返回。
        // 现在，任何返回带有 app error 的 Result 的处理器都会自动得到适当的错误响应。
        (status, Json(body)).into_response()
    }
}

// ============================================================================
// 课程 3：基于 Result 的处理函数
// ============================================================================

#[derive(Serialize)]
struct User {
    id: u64,
    name: String,
}

// 返回 Json<User> 或 AppError 的处理器。
async fn get_user(Path(id): Path<u64>) -> Result<Json<User>, AppError> {
    // 模拟查询用户
    match id {
        // 如果找到，就返回包含用户的 JSON 响应
        1 => Ok(Json(User {
            id: 1,
            name: "Alice".to_string(),
        })),
        2 => Ok(Json(User {
            id: 2,
            name: "Bob".to_string(),
        })),
        // 如果没找到，就返回应用错误：用户未找到。
        _ => Err(AppError::UserNotFound(id)),
    }
}

async fn validate_input(Path(value): Path<String>) -> Result<String, AppError> {
    if value.len() < 3 {
        return Err(AppError::InvalidInput(
            "值至少需要 3 个字符".to_string(),
        ));
    }
    Ok(format!("输入有效：{}", value))
}

async fn protected_resource() -> Result<&'static str, AppError> {
    // 模拟认证检查
    let is_authenticated = false;
    if !is_authenticated {
        return Err(AppError::Unauthorized);
    }
    Ok("机密数据！")
}

async fn database_operation() -> Result<&'static str, AppError> {
    // 模拟数据库错误
    Err(AppError::DatabaseError("连接超时".to_string()))
}

// ============================================================================
// 课程 4：使用 ? 的可失败操作
// ============================================================================

// 包含链式可能失败操作的复杂处理器。
// 问号操作符在这里工作得非常好。
async fn complex_operation(Path(id): Path<u64>) -> Result<Json<User>, AppError> {
    // 使用 ? 运算符进行提前返回
    // 如果一个函数返回错误，问号会传播该错误。
    // 无需冗长的 match 语句。
    // 这是 Rust 的强大之处——显式的错误处理且无样板代码。
    let user = find_user(id)?;
    validate_user(&user)?;
    Ok(Json(user))
}

fn find_user(id: u64) -> Result<User, AppError> {
    if id == 0 {
        Err(AppError::InvalidInput("ID 不能为 0".to_string()))
    } else if id > 100 {
        Err(AppError::UserNotFound(id))
    } else {
        Ok(User {
            id,
            name: format!("User{}", id),
        })
    }
}

fn validate_user(user: &User) -> Result<(), AppError> {
    if user.name.is_empty() {
        Err(AppError::InvalidInput("名称不能为空".to_string()))
    } else {
        Ok(())
    }
}

// ============================================================================
// 主程序
// ============================================================================

#[tokio::main]
async fn main() {
    let app = Router::new()
        .route("/users/{id}", get(get_user))
        .route("/validate/{value}", get(validate_input))
        .route("/protected", get(protected_resource))
        .route("/database", get(database_operation))
        .route("/complex/{id}", get(complex_operation));

    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();

    println!("🚀 模块 07：错误处理");
    println!("   服务地址：http://localhost:3000\n");
    println!("📝 可以尝试以下端点：");
    println!("   GET /users/1      - 成功（用户存在）");
    println!("   GET /users/999    - 404（用户未找到）");
    println!("   GET /validate/ab  - 400（长度过短）");
    println!("   GET /protected    - 401（未授权）");
    println!("   GET /database     - 500（数据库错误）");

    axum::serve(listener, app).await.unwrap();
}
```

## 总结

- 第一条规则：处理器中绝不使用 panic。
  - panic 会导致线程崩溃并带来糟糕的用户体验。
  - 使用 Result 返回错误。不要使用 panic。
