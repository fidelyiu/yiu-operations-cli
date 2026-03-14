# Axum教程-11-测试

## 创建模块

```sh
cargo new module-11-testing
```

## 添加依赖

`module-11-testing/Cargo.toml`

```toml
[package]
name = "module-11-testing"
version = "0.1.0"
edition = "2024"

[dependencies]
axum = { workspace = true, features = ["ws", "multipart"] }
axum-extra = { workspace = true }
tokio = { workspace = true }
tokio-stream = "0.1"
serde = { workspace = true }
serde_json = { workspace = true }
tower-http = { workspace = true }
futures = { workspace = true }
```

## 案例

未经测试的代码就是有缺陷的代码。你只是还不知道而已。

每一个生产环境的错误都是未被测试的行为。

axom 让测试变得愉快。

让我们看看 axom 如何让测试变得容易，因为路由器本身就是一个 tower 服务。

我们可以直接发送请求，而无需启动真实的服务器。

无需绑定端口、无需网络、无需竞争条件，只需带有完整类型安全性的直接函数调用。

oneshot 方法发送一次请求并返回响应。快速、独立、确定性。

这就是为什么 Rust 的 Web 测试比其他语言好得多的原因。

无需 HTTP 客户端模拟、无需搭建假服务器，只需调用路由器。

```rust
//! # 模块 11：测试
//!
//! 测试 Axum 应用：
//! - 对处理函数进行单元测试
//! - 使用 TestClient 进行集成测试
//! - 使用模拟状态进行测试

use axum::{
    Json, Router,
    extract::{Path, State},
    http::StatusCode,
    routing::get,
};
use serde::{Deserialize, Serialize};
use std::{
    collections::HashMap,
    sync::{Arc, RwLock},
};

// ============================================================================
// 应用代码
// ============================================================================

#[derive(Clone, Serialize, Deserialize, Debug, PartialEq)]
struct User {
    id: u64,
    name: String,
}

#[derive(Deserialize)]
struct CreateUser {
    name: String,
}

type UserStore = Arc<RwLock<HashMap<u64, User>>>;

async fn list_users(State(store): State<UserStore>) -> Json<Vec<User>> {
    let users = store.read().unwrap();
    Json(users.values().cloned().collect())
}

async fn get_user(
    State(store): State<UserStore>,
    Path(id): Path<u64>,
) -> Result<Json<User>, StatusCode> {
    let users = store.read().unwrap();
    users
        .get(&id)
        .cloned()
        .map(Json)
        .ok_or(StatusCode::NOT_FOUND)
}

async fn create_user(
    State(store): State<UserStore>,
    Json(input): Json<CreateUser>,
) -> (StatusCode, Json<User>) {
    let mut users = store.write().unwrap();
    let id = users.len() as u64 + 1;
    let user = User {
        id,
        name: input.name,
    };
    users.insert(id, user.clone());
    (StatusCode::CREATED, Json(user))
}

async fn health() -> &'static str {
    "OK"
}

fn create_app(store: UserStore) -> Router {
    Router::new()
        .route("/health", get(health))
        .route("/users", get(list_users).post(create_user))
        .route("/users/{id}", get(get_user))
        .with_state(store)
}

// ============================================================================
// 测试
// ============================================================================

#[cfg(test)]
mod tests {
    use super::*;
    use axum::{body::Body, http::Request};
    use http_body_util::BodyExt;
    use tower::ServiceExt; // 用于 oneshot

    // 使用一个辅助函数创建测试状态。
    // 将其传递给 ccreate_app，它会构建我们的路由器。
    // 现在我们有一个准备好接受请求的路由器。每个测试都会创建它自己的状态。
    // 测试是相互隔离的。
    // 没有共享的可变状态意味着不会出现不稳定的测试。
    // 全在内存中，极其快速。
    // 你的整个测试套件在几秒钟内运行完毕。
    fn test_store() -> UserStore {
        // 测试存储返回一个空的内存存储。
        Arc::new(RwLock::new(HashMap::new()))
    }

    #[tokio::test]
    async fn test_health_check() {
        let app = create_app(test_store());

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/health")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        // 断言状态为 200
        assert_eq!(response.status(), StatusCode::OK);

        // 使用 collect 提取主体字节并断言相等。这只需毫秒。
        // 你可以快速运行数千个测试。
        let body = response.into_body().collect().await.unwrap().to_bytes();
        assert_eq!(&body[..], b"OK");
    }

    #[tokio::test]
    async fn test_create_user() {
        let app = create_app(test_store());

        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/users")
                    .header("content-type", "application/json")
                    .body(Body::from(r#"{"name":"Alice"}"#))
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::CREATED);

        let body = response.into_body().collect().await.unwrap().to_bytes();
        let user: User = serde_json::from_slice(&body).unwrap();
        assert_eq!(user.name, "Alice");
    }

    #[tokio::test]
    async fn test_get_user_not_found() {
        let app = create_app(test_store());

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/users/999")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::NOT_FOUND);
    }

    #[tokio::test]
    async fn test_list_users() {
        let store = test_store();
        store.write().unwrap().insert(
            1,
            User {
                id: 1,
                name: "Bob".to_string(),
            },
        );

        let app = create_app(store);

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/users")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);

        let body = response.into_body().collect().await.unwrap().to_bytes();
        let users: Vec<User> = serde_json::from_slice(&body).unwrap();
        assert_eq!(users.len(), 1);
    }
}

// ============================================================================
// 主程序
// ============================================================================

#[tokio::main]
async fn main() {
    let store = Arc::new(RwLock::new(HashMap::new()));
    let app = create_app(store);

    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();

    println!("🚀 模块 11：测试");
    println!("   服务地址：http://localhost:3000\n");
    println!("📝 接口列表：");
    println!("   GET  /health    - 健康检查");
    println!("   GET  /users     - 获取用户列表");
    println!("   POST /users     - 创建用户");
    println!("   GET  /users/:id - 获取用户\n");
    println!("🧪 运行测试：cargo test");

    axum::serve(listener, app).await.unwrap();
}
```

## 总结

不要在测试之间共享可变状态。

如果必须共享，请使用 argar 锁，但更推荐每个测试使用独立状态。
