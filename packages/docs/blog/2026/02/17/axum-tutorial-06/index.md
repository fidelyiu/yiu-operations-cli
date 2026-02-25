---
title: Axum教程06
authors: [FidelYiu]
tags: [rust, axum]
---

# Axum教程06

Axum 是当前 Rust 中的 Web 框架。由 tokyo 团队构建，它速度快、类型安全，且使用起来非常优雅。

第五个模块：真实的应用需要共享状态：数据库连接池、配置、缓存、计数器。

本模块教你如何在 axum 中[管理状态](https://docs.rs/axum/latest/axum/routing/struct.Router.html#method.with_state)。

<!-- truncate -->

## 创建模块

```sh
cargo new module-05-state
```

## 添加依赖

`module-05-state/Cargo.toml`

```toml
[package]
name = "module-05-state"
version = "0.1.0"
edition = "2024"

[dependencies]
axum = { workspace = true }
tokio = { workspace = true }
serde = { workspace = true }
serde_json = { workspace = true }
uuid = { workspace = true }
```

共享状态的主要方式是使用 state 提取器。

## 不可变共享状态

### 定义状态

```rust
#[derive(Clone)]
struct AppConfig {
    app_name: String,
    version: String,
    max_items_per_page: usize,
}
```

### 定义处理器

```rust
async fn get_config(State(config): State<Arc<AppConfig>>) -> Json<serde_json::Value> {
    Json(serde_json::json!({
        "app_name": config.app_name,
        "version": config.version,
        "max_items": config.max_items_per_page
    }))
}
```

### 定义路由

```rust
#[tokio::main]
async fn main() {
    let config = Arc::new(AppConfig {
        app_name: "Axum Todo API".to_string(),
        version: "1.0.0".to_string(),
        max_items_per_page: 100,
    });
    let app = Router::new()
        // Config endpoint
        .route("/config", get(get_config))
        .with_state(config);
    axum::serve(listener, app).await.expect("Server failed");
}
```

### 总计

这个配置在运行时不会改变。我们把它包装在 `Arc`（原子引用计数）中。

`Arc` 允许多个处理程序共享相同的数据，而无需克隆它。

它是线程安全且几乎没有开销。

当你用 `Arc` 包装并调用时，所有处理程序都可以访问相同的配置。

内存高效且安全。

## 可变共享状态

那会改变的状态怎么办？

### 定义内存状态

```rust
/// 一个简单的内存“数据库”，用于存储待办事项
/// 使用 RwLock 获得更好的读取性能（多个读者，单个写者）
#[derive(Debug, Clone, Serialize, Deserialize)]
struct Todo {
    id: String,
    title: String,
    completed: bool,
}

/// 我们的可变状态 - 线程安全的 HashMap
type TodoStore = Arc<RwLock<HashMap<String, Todo>>>;
```

RwLock 是一个读写 锁。多个读取者可以同时访问。写入者获得独占访问权。

这种模式非常适合内存缓存、计数器或任何由处理程序修改的数据。为什么使用行锁而不是互斥锁来提高性能？

互斥锁为所有操作提供排他访问，在以读为主的工作负载中（大多数 Web 应用都是读多写少），读操作会不必要地相互阻塞。

`RwLock`帐户读取权限不需要排他性。

对于 Web 应用程序，读取远多于写入。

因此 `RwLock` 通常更好。

### 定义处理器-list

```rust
// 列出所有待办事项
async fn list_todos(State(store): State<TodoStore>) -> Json<Vec<Todo>> {
    let todos = store.read().unwrap();
    let todos_vec: Vec<Todo> = todos.values().cloned().collect();
    Json(todos_vec)
}
```

### 定义处理器-create

```rust
#[derive(Debug, Deserialize)]
struct CreateTodo {
    title: String,
}

// 创建新的待办事项
async fn create_todo(
    State(store): State<TodoStore>,
    Json(input): Json<CreateTodo>,
) -> (StatusCode, Json<Todo>) {
    let todo = Todo {
        id: Uuid::new_v4().to_string(),
        title: input.title,
        completed: false,
    };

    store.write().unwrap().insert(todo.id.clone(), todo.clone());

    (StatusCode::CREATED, Json(todo))
}
```

### 定义处理器-getItem

```rust
// 获取单个待办事项
async fn get_todo(
    State(store): State<TodoStore>,
    axum::extract::Path(id): axum::extract::Path<String>,
) -> Result<Json<Todo>, StatusCode> {
    let todos = store.read().unwrap();
    todos
        .get(&id)
        .cloned()
        .map(Json)
        .ok_or(StatusCode::NOT_FOUND)
}
```

### 定义处理器-update

```rust
#[derive(Debug, Deserialize)]
struct UpdateTodo {
    title: Option<String>,
    completed: Option<bool>,
}

// 更新待办事项
async fn update_todo(
    State(store): State<TodoStore>,
    axum::extract::Path(id): axum::extract::Path<String>,
    Json(input): Json<UpdateTodo>,
) -> Result<Json<Todo>, StatusCode> {
    let mut todos = store.write().unwrap();

    if let Some(todo) = todos.get_mut(&id) {
        if let Some(title) = input.title {
            todo.title = title;
        }
        if let Some(completed) = input.completed {
            todo.completed = completed;
        }
        Ok(Json(todo.clone()))
    } else {
        Err(StatusCode::NOT_FOUND)
    }
}
```

### 定义处理器-delete

```rust
// 删除待办事项
async fn delete_todo(
    State(store): State<TodoStore>,
    axum::extract::Path(id): axum::extract::Path<String>,
) -> StatusCode {
    let mut todos = store.write().unwrap();
    if todos.remove(&id).is_some() {
        StatusCode::NO_CONTENT
    } else {
        StatusCode::NOT_FOUND
    }
}
```

### 注册路由

```rust

#[tokio::main]
async fn main() {
    // 初始化可变的待办事项存储
    let todo_store: TodoStore = Arc::new(RwLock::new(HashMap::new()));

    // 预先添加一些待办事项
    {
        let mut store = todo_store.write().unwrap();
        let todo = Todo {
            id: Uuid::new_v4().to_string(),
            title: "Learn Axum".to_string(),
            completed: false,
        };
        store.insert(todo.id.clone(), todo);
    }

    // 构建待办事项 CRUD 路由
    let todo_routes = Router::new()
        .route("/", get(list_todos).post(create_todo))
        .route("/{id}", get(get_todo).put(update_todo).delete(delete_todo))
        .with_state(todo_store);

    // 构建主应用
    let app = Router::new()
        // 合并待办事项路由
        .merge(Router::new().nest("/todos", todo_routes));

    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000")
        .await
        .expect("Failed to bind");
    axum::serve(listener, app).await.expect("Server failed");
}

```

## 多种状态类型

大型应用需要多种状态类型。

配置、待办事项和指标合并在一个 struct 中的组合状态。

### 定义状态

```rust
/// 当你需要多个相互独立的状态类型时，将它们组合起来
#[derive(Clone)]
#[allow(dead_code)] // 仅用于演示展示字段
struct CombinedState {
    config: Arc<AppConfig>,
    todos: TodoStore,
    metrics: Arc<RwLock<Metrics>>,
}

#[derive(Debug, Default)]
struct Metrics {
    request_count: u64,
    error_count: u64,
}
```

为该 struct 实现 Clone。每个 Arc 在克隆时仅通过增加计数器来实现廉价克隆。

### 定义处理器

处理程序可以访问它们需要的任何字段。

```rust
// 你可以提取整个状态，或使用 From trait 以获得更方便的提取方式
async fn get_metrics(State(state): State<CombinedState>) -> Json<serde_json::Value> {
    let metrics = state.metrics.read().unwrap();
    Json(serde_json::json!({
        "requests": metrics.request_count,
        "errors": metrics.error_count,
        "app_version": state.config.version
    }))
}

async fn increment_request_count(State(state): State<CombinedState>) -> &'static str {
    let mut metrics = state.metrics.write().unwrap();
    metrics.request_count += 1;
    "Request counted!"
}
```

### 定义路由

将组合后的状态传递给 with_data 或 withdate（保持原文拼写）。

在真实应用中，这通常会是 sql xpg poolool 或类似的东西。

```rust

```
