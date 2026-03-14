# Axum教程-08-数据库集成

使用 sqlx 将 axum 连接到 PostgreSQL。

这是你的应用变成真实数据并能在重启后保留的地方。

PostgreSQL 经受过考验、开源，并且与 Rust 配合得非常好。

sqlx 是一个用于 SQL 数据库的 Rust 库。

它是原生异步的、编译时检查的，并且类型安全。

但为什么要选择它而不是像 diesel 或 comm 这样的替代品呢？

首先，编译时的查询检查。sqlx 会将你的查询与实际的数据库模式进行验证。拼写错误的列名会在编译时被捕获。错误类型在编译时被捕获。与您的模式不匹配的查询在编译时被捕获。

其次，它从一开始就是异步的。其他框架是后期为异步改造的。sqlx 是为此而设计的。

第三，你编写真正的 SQL。无需学习任何 DSL。如果你会 SQL，那么你就会 sqlx。

第四，无需代码生成。diesel 需要 diesel doublegi。sqlx 只需使用 cargo build 即可运行。这是为了正确性和开发者体验的游戏规则改变者。

## 运行 postgresql

你需要运行 postgresql。最简单的方式是使用 docker。只需一个命令就完成了。

```sh
docker run -d --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres
```

- `-d` 表示以分离模式（后台）运行。
- `--name` 给容器命名，这样以后可以停止它。
- `-e` 将密码设置为环境变量。
- `-p` 将容器的 5432 端口映射到在本机上以 `localhost:5432` 启动的 postgresql。

## 设置环境变量

现在在 Linux 或 macOS 上设置你的环境变量。

```sh
export DATABASE_URL=postgres://postgres:postgres@localhost/axum_course
```

我们也可以在`dotenv`中配置他

.env

```sh
# Database (Module 08)
DATABASE_URL=postgres://postgres:postgres@localhost/axum_course

# JWT Secret (Module 09)
JWT_SECRET=your-super-secret-jwt-key

# Logging
RUST_LOG=info
```

## 创建模块

```sh
cargo new module-08-database
```

## 添加依赖

`module-08-database/Cargo.toml`

```toml
[package]
name = "module-08-database"
version = "0.1.0"
edition = "2024"

[dependencies]
axum = { workspace = true }
tokio = { workspace = true }
serde = { workspace = true }
serde_json = { workspace = true }
sqlx = { workspace = true }
uuid = { workspace = true }
chrono = { workspace = true }
dotenvy = { workspace = true }
thiserror = { workspace = true }
```

## 创建数据库连接池

打开数据库连接需要时间。tcp 握手、认证、协议协商。如果你为每个请求都打开一个连接，你的服务器会很慢。

连接池可以解决这个问题。

它们只打开一次连接并重复使用。连接池管理一切。当处理器需要连接时，它会借用一个。完成后，会将其归还给连接池。快速、高效，每个请求零连接开销。

有一点重要的注意事项，不要把最大连接数设置得太高。每个连接都会占用服务器内存。5 到 20 个连接可处理大多数工作负载。

只有在你证明需要更多时才增加。

```rust
#[tokio::main]
async fn main() {
    dotenvy::dotenv().ok();

    let database_url = std::env::var("DATABASE_URL")
        .unwrap_or_else(|_| "postgres://postgres:postgres@localhost/axum_course".to_string());

    let pool = PgPoolOptions::new()
        // 将最大连接数设置为五
        .max_connections(5)
        // 并使用数据库 URL 调用 connect。
        .connect(&database_url)
        .await
        .expect("连接数据库失败");
}
```

## 连接池状态

将连接池作为状态传递给你的路由器。每个处理器都可以访问它。

所有处理器共享相同的连接池。

```rust
let app = Router::new()
    .route("/users", get(list_users).post(create_user))
    .route(
        "/users/{id}",
        get(get_user).put(update_user).delete(delete_user),
    )
    .with_state(pool);
```

## 数据库模型

这是映射的魔法。

`FromRow` 让 sqlx 自动将数据库行转换为你的结构体

列名与字段名匹配 postgrql

返回列 sqlx 将它们映射到结构体字段，如果列名不同 sqlx 重命名属性，如果某列可为空，则将字段设为 Option

sqlx 全部处理好了

```rust
#[derive(Debug, Serialize, sqlx::FromRow)]
struct User {
    // UUID 作为主键，UUID 比自增整数更适合分布式系统。
    // 在服务器之间不需要协调，而且它们是不可预测的，这对安全性更好。
    id: Uuid,
    name: String,
    email: String,
    created_at: chrono::DateTime<chrono::Utc>,
}
```

## 查询处理器

- `fetch_all`
  - 返回 `Vec<User>`, 所有行都有完整类型。每一行都会变成一个 user strct。
- `fetch_one`
  - 这是用于必须返回且恰好为一行的查询。如果为零行或多行，则会报错。
  - 当存在多行时只取一行，使用 fetch_one。如果没有行则会失败。
- `fetch_optional`
  - 对于可选结果，使用 fetch_optional。
  - 如果没有行则返回 None，如果找到则返回 Some。
  - 将此用于按 id 查找。
  - 使用 fetch_optional 优雅地处理未找到的情况。

```rust

async fn list_users(State(pool): State<PgPool>) -> Result<Json<Vec<User>>, DbError> {
    let users = sqlx::query_as::<_, User>("SELECT * FROM users ORDER BY created_at DESC")
        .fetch_all(&pool)
        .await?;
    Ok(Json(users))
}
```

## sql参数

`$`号那个是占位符。bind 会填充它。这可防止 SQL 注入。这是最常见的 Web 漏洞。

如果你把用户输入拼接到 SQL 中，攻击者就能操纵你的查询。像 "drillable users" 这样的输入可能会破坏你的数据。

参数化查询可以防止这种情况。数据库知道什么是 SQL，什么是数据。数据永远不会作为 SQL 被执行。

永远不要将用户输入拼接到 SQL 字符串中。始终使用绑定参数。这是不可谈判的。

```rust
async fn get_user(State(pool): State<PgPool>, Path(id): Path<Uuid>) -> Result<Json<User>, DbError> {
    let user = sqlx::query_as::<_, User>("SELECT * FROM users WHERE id = $1")
        .bind(id)
        .fetch_optional(&pool)
        .await?
        .ok_or(DbError::NotFound)?;
    Ok(Json(user))
}
```

## 插入数据

`RETURNING *`是 PostgreSQL 特有的，但非常有用。

它会把插入的那一行返回给我们，这样我们就可以把它返回给客户端。不需要单独的查询。

```rust
#[derive(Debug, Deserialize)]
struct CreateUser {
    name: String,
    email: String,
}

async fn create_user(
    State(pool): State<PgPool>,
    Json(input): Json<CreateUser>,
) -> Result<(StatusCode, Json<User>), DbError> {
    let user = sqlx::query_as::<_, User>(
        "INSERT INTO users (id, name, email, created_at) VALUES ($1, $2, $3, NOW()) RETURNING *",
    )
    .bind(Uuid::new_v4())
    .bind(&input.name)
    .bind(&input.email)
    .fetch_one(&pool)
    .await?;
    // 状态码 201 告诉客户端资源已创建。
    Ok((StatusCode::CREATED, Json(user)))
}
```

## 修改数据

绑定新值后，`execute` 方法会返回受影响的行数。rows 检查是否有零行受影响，这意味着资源不存在返回 404 NotFound.

如果你需要更新后的行，使用 `returning *` 再次 `fetch` 而不是 `execute`。

```rust
#[derive(Debug, Deserialize)]
struct UpdateUser {
    name: Option<String>,
    email: Option<String>,
}

async fn update_user(
    State(pool): State<PgPool>,
    Path(id): Path<Uuid>,
    Json(input): Json<UpdateUser>,
) -> Result<Json<User>, DbError> {
    let user = sqlx::query_as::<_, User>(
        "UPDATE users SET name = COALESCE($2, name), email = COALESCE($3, email) WHERE id = $1 RETURNING *"
    )
    .bind(id)
    .bind(&input.name)
    .bind(&input.email)
    .fetch_optional(&pool)
    .await?
    .ok_or(DbError::NotFound)?;
    Ok(Json(user))
}
```

## 删除数据

```rust
async fn delete_user(
    State(pool): State<PgPool>,
    Path(id): Path<Uuid>,
) -> Result<StatusCode, DbError> {
    let result = sqlx::query("DELETE FROM users WHERE id = $1")
        .bind(id)
        .execute(&pool)
        .await?;

    if result.rows_affected() == 0 {
        Err(DbError::NotFound)
    } else {
        Ok(StatusCode::NO_CONTENT)
    }
}
```

## 数据库错误

sqlx返回结果类型每个数据库操作操作。

将 sql x 错误映射到你的应用错误。什么类型的错误？

- 数据库链接错误，数据库宕机。
- 超时错误。查询花费太长时间。
  - 连接问题返回 500，这不是用户的错。
- 违反约束。唯一键冲突。外键缺失。
  - 违反约束可能返回 400：错误的请求。
  - 用户提供了无效的数据。
- 当零行时，为 fetch_one 未找到。
  - 未找到返回 404。

使用 match 或 if let 来区分错误类型。抛出完整的错误以便调试。向用户返回友好消息。

## 事物

对于必须一起成功的多个操作，请使用事务。

要么全部成功，要么全部失败。

经典示例。

在账户之间转账。借记一方，贷记另一方。如果贷记失败，借记必须回滚。

使用 调用 begin 来进行事务。执行你的查询。调用 transaction。

如果在提交前发生任何错误，事务在被释放时会自动回滚。保证原子性。

在查询中使用 amperand transaction，而不是 ampersand pool。

事务会跟踪使用哪个连接。

## 数据库迁移

我应该提到迁移。我们还没有涉及如何运行它们。

但你需要创建表。创建一个 migrations 文件夹。把 SQL 文件放在那里。

运行 `sqlx migrate run SQL` 会提取哪些迁移已经运行。

它不会重新运行已执行的迁移。保留用于生产环境。

对于本课程，你可以手动运行创建表操作，但生产应用应使用迁移。

## 总结

一些用于生产数据库工作的提示。

- 只选择你需要的列。`SELECT *`方便但如果你只需要两列则很浪费。
- 在你用于过滤或排序的列上添加索引。没有索引时，PostgreSQL 会扫描整个表。使用解释性分析来了解查询性能。在大型表上查找顺序扫描。
- 如果在等待连接时，监控连接池的使用情况。增加连接池大小或优化慢查询。
