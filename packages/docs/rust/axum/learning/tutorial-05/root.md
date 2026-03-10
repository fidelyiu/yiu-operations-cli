# Axum教程-05-状态共享

Axum 是当前 Rust 中的 Web 框架。由 tokyo 团队构建，它速度快、类型安全，且使用起来非常优雅。

第五个模块：真实的应用需要共享状态：数据库连接池、配置、缓存、计数器。

本模块教你如何在 axum 中[管理状态](https://docs.rs/axum/latest/axum/routing/struct.Router.html#method.with_state)。

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
#[derive(Clone)] // 为该结构实现 Clone
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

将组合后的状态传递给 with_state。

```rust
// 为复杂应用准备的组合状态
let combined_state = CombinedState {
    // 每个 Arc 在克隆时仅通过增加计数器来实现廉价克隆。
    config: config.clone(),
    todos: todo_store.clone(),
    metrics: Arc::new(RwLock::new(Metrics::default())),
};

// 构建主应用
let app = Router::new()
    .route("/metrics", get(get_metrics))
    .route("/track", get(increment_request_count))
    .with_state(combined_state);
```

## 数据库连接池

在真实应用中，这通常会是 sql xpg poolool 或类似的东西。

连接池本质上是可共享的。那就是它们存在的目的。

把连接池作为状态传递，处理器可以根据需要进行查询。

我们接下来会在第8模块中看到真实的数据库。

目前，理解为连接池只是另一种状态类型即可。

### 定义状态

```rust
/// 模拟数据库连接池
/// 在真实应用中，这里会是 sqlx::PgPool 或类似类型
#[derive(Clone)]
#[allow(dead_code)] // 仅用于演示展示字段
struct DbPool {
    connection_string: String,
    max_connections: u32,
}

impl DbPool {
    fn new(connection_string: &str) -> Self {
        Self {
            connection_string: connection_string.to_string(),
            max_connections: 10,
        }
    }

    // 模拟查询
    async fn query(&self, _sql: &str) -> Result<Vec<String>, String> {
        // 在真实应用中：sqlx::query!(...).fetch_all(&self.pool).await
        Ok(vec!["result1".to_string(), "result2".to_string()])
    }
}
```

### 定义处理器

```rust
async fn db_query(State(pool): State<DbPool>) -> Json<Vec<String>> {
    match pool.query("SELECT * FROM users").await {
        Ok(results) => Json(results),
        Err(_) => Json(vec![]),
    }
}
```

### 注册路由

```rust
// 模拟数据库连接池
let db_pool = DbPool::new("postgres://localhost/myapp");
// 构建主应用
let app = Router::new()
    // 数据库端点
    .route("/db/users", get(db_query))
    .with_state(db_pool);
```

## 扩展模式

扩展是每个请求不同的请求雕塑状态，由中间件设置

常用的认证中间件会验证令牌并将当前用户插入到 extensions 中。

处理器在提取它时知道用户已经通过了认证。

对应用范围的数据使用 state。

对请求特定的数据使用 extension。

### 定义扩展

```rust
/// 有时你希望动态添加状态（例如从中间件注入）
use axum::Extension;

#[derive(Clone)]
struct CurrentUser {
    id: String,
    name: String,
}
```

### 定义处理器

```rust
async fn get_current_user(Extension(user): Extension<CurrentUser>) -> Json<serde_json::Value> {
    Json(serde_json::json!({
        "id": user.id,
        "name": user.name
    }))
}
```

### 注册路由

```rust
// 当前用户（通常由鉴权中间件设置）
let current_user = CurrentUser {
    id: "user-123".to_string(),
    name: "Demo User".to_string(),
};

// 构建主应用
let app = Router::new()
    // 基于 Extension 的状态
    .route("/me", get(get_current_user))
    .layer(Extension(current_user));
```

## 总结

- 保持状态精简。
  - 不要把所有东西塞进一个结构体。
  - 将相关数据分组。
- 使用 Arc 来共享。
  - 它既廉价又安全。
- 在 Web 应用中优先使用 RwLock 而不是互斥锁（mutex）。
- 在构建路由器之前先初始化状态。
  - 启动时的错误比请求期间的错误更容易调试。
- 考虑为不同的状态需求使用独立的路由器。
  - 最后将它们合并。
- Arc 用于不可变共享的状态。
- RwLock 用于可变共享的状态。
- 组合状态结构以支持多种类型。
- 扩展用于请求的作用域的状态。
- 将数据库连接池作为状态。
- 状态管理是真正应用的支柱。
