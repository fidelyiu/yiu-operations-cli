# Axum教程-09-认证

第 9 模块增加了安全性：JWT 令牌、密码哈希、受保护路由和基于角色的访问控制。

这是使你的 API 从开放到安全的关键一步。

认证回答了一个问题。谁在发出这个请求。

授权回答了另一个问题。他们被允许做什么？

## JWT

JSON Web Tokens 是一种用于无状态认证的标准。

---

让我解释为什么这很重要。传统的会话认证是这样工作的。用户登录后，服务器创建会话，并将其存储在内存或数据库中，然后给用户一个会话 ID 的 cookie。

每次请求时，服务器查询该会话以确定请求者的身份。问题在于会话是有状态的。服务器需要记住它们。如果你有多台服务器，你需要共享会话存储。扩展（横向扩展）是复杂的。

---

JWT 不同。服务器会签发一个包含用户信息的签名令牌。该令牌本身包含验证用户所需的所有信息。服务器不存储任何内容。没有服务器端会话。

可以在无需协调的情况下添加更多服务器。横向扩展非常简单。非常适合 API 和微服务。

JWT 有三部分，用点分隔。头部、载荷和签名。

- 头部说明哪个算法对令牌进行签名，通常是 HS256 或 RS256。
- 载荷包含声明、用户 ID、以及过期时间。
- 签名用于验证令牌未被篡改。

任何人都可以解码该负载。它只是 Base64，但只有服务器能创建有效的签名，因为只有服务器知道密钥。

如果有人修改了负载，签名就不会匹配。服务器会拒绝它。设计上防篡改。

## 创建模块

```sh
cargo new module-09-auth
```

## 添加依赖

`module-09-auth/Cargo.toml`

```toml
[package]
name = "module-09-auth"
version = "0.1.0"
edition = "2024"

[dependencies]
axum = { workspace = true }
axum-extra = { workspace = true }
tokio = { workspace = true }
serde = { workspace = true }
serde_json = { workspace = true }
jsonwebtoken = { workspace = true }
argon2 = { workspace = true }
chrono = { workspace = true }
uuid = { workspace = true }
thiserror = { workspace = true }
rand = "0.8"
```

## 密码加密和验证

永远不要存储明文密码。我不能再强调这一点了。

如果你的数据库泄露，攻击者就会得到所有人的密码。

用户在不同网站重复使用密码，这样攻击者就能访问电子邮件、银行和所有内容。

使用 argon2 进行哈希。argon2 在 2015 年赢得了密码哈希竞赛。

它是“内存硬”的，意味着 GPU 攻击代价高昂。

它是目前的黄金标准。

hash_password 函数为每个密码生成唯一的随机盐。盐值可以防止彩虹表攻击。

然后它用盐值对密码进行哈希，并存储生成的哈希字符串。

哈希是单向的。你无法从哈希中恢复出原始密码。这就是关键。

相反，你用相同的盐对提供的密码进行哈希并比较结果。

用提供的盐对密码进行哈希。如果匹配，认证成功。如果不匹配，则失败。

绝不要告诉用户是哪一部分失败。无效的电子邮件或密码。

无效的密码会告诉攻击者该电子邮件存在。

```rust
// ============================================================================
// 密码哈希
// ============================================================================

fn hash_password(password: &str) -> String {
    use argon2::{password_hash::SaltString, Argon2, PasswordHasher};
    let salt = SaltString::generate(&mut rand::rngs::OsRng);
    Argon2::default()
        .hash_password(password.as_bytes(), &salt)
        .unwrap()
        .to_string()
}

/// 根据哈希值验证密码（演示函数）
#[allow(dead_code)]
fn verify_password(password: &str, hash: &str) -> bool {
    use argon2::{Argon2, PasswordHash, PasswordVerifier};
    let parsed_hash = PasswordHash::new(hash).unwrap(); // 解析存储的哈希，它包含了盐。
    Argon2::default()
        .verify_password(password.as_bytes(), &parsed_hash)
        .is_ok()
}
```

## Claims

claims 是关于用户的声明，嵌入在令牌中。

- `sub`声明是主题，通常是用户 ID。
- `exp`声明是过期时间，是一个表示令牌失效时间的 Unix 时间戳。
- `role`声明存储授权级别，例如 user、admin 等。

你可以添加自定义声明，但要保持令牌体积小。它们随每个请求发送。将用户详细信息存储在数据库中，而不是存储在令牌中。

```rust
#[derive(Debug, Serialize, Deserialize)]
struct Claims {
    sub: String, // 用户 ID
    exp: usize,  // 过期时间戳
    role: String,
}
```

## JWT 创建 token

我们构建了包含用户 ID 的声明，过期时间为从现在起 24 小时，并进行签发。

jwt 秘密必须保密。使用环境变量。生成至少 32 个字符的长随机字符串。切勿将密钥写死在代码中。切勿将它们提交到 git。

在生产环境中，考虑使用令牌刷新流程。

短期有效的访问令牌（例如 15 分钟），配合长期有效的刷新令牌。

如果令牌泄露，影响范围会被限制在较小范围内。

```rust
#[derive(Clone)]
struct AuthConfig {
    jwt_secret: String,
    jwt_expiry_hours: i64,
}

fn create_token(config: &AuthConfig, user_id: &str, role: &str) -> Result<String, StatusCode> {
    // chrono crate 用于计算时间戳。
    let expiry = Utc::now() + Duration::hours(config.jwt_expiry_hours);
    let claims = Claims {
        sub: user_id.to_string(),
        exp: expiry.timestamp() as usize,
        role: role.to_string(),
    };
    // 我们在代码中使用默认头、声明和密钥进行调用。
    encode(
        &Header::default(),
        &claims,
        &EncodingKey::from_secret(config.jwt_secret.as_bytes()),
    )
    .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)
}
```

## JWT 验证token

每个受保护的请求都会验证该令牌。

使用令牌的密钥和验证选项调用 decode。

验证会自动检查过期时间。

如果验证成功，我们就会得到声明（claims）。

如果令牌过期，返回 401；如果被篡改，返回 401；如果缺失，返回 401。

`Claims`告诉我们用户是谁以及他们可以做什么。将此插入请求`extensions`中以供处理器使用。

```rust
fn verify_token(config: &AuthConfig, token: &str) -> Result<Claims, StatusCode> {
    decode::<Claims>(
        token,
        &DecodingKey::from_secret(config.jwt_secret.as_bytes()),
        &Validation::default(),
    )
    .map(|data| data.claims)
    .map_err(|_| StatusCode::UNAUTHORIZED)
}
```

## 认证中间键

```rust

// ============================================================================
// 认证中间件
// ============================================================================

async fn auth_middleware(
    State(config): State<Arc<AuthConfig>>,
    mut request: Request,
    next: Next,
) -> Result<Response, StatusCode> {
    let auth_header = request
        .headers()
        .get("Authorization")
        .and_then(|v| v.to_str().ok())
        .and_then(|v| v.strip_prefix("Bearer "));

    let token = auth_header.ok_or(StatusCode::UNAUTHORIZED)?;
    let claims = verify_token(&config, token)?;

    let user = CurrentUser {
        id: claims.sub,
        role: claims.role,
    };
    // 我们把它插入到请求的 extensions 中。
    // extensions 是请求范围的存储，中间件可以写入。处理器可以读取。
    request.extensions_mut().insert(user);

    // 如果任何步骤失败——没有头、格式错误、令牌无效——立即返回 401 未授权。
    // 处理程序不会运行。边缘安全。
    Ok(next.run(request).await)
}
```

## 路由

现在 post /login 返回 token。get /protected/me 需要有效的 token，get /protected/admin 需要带有 admin 角色的 token。

```rust
#[derive(Deserialize)]
struct LoginRequest {
    email: String,
    password: String,
}

#[derive(Serialize)]
struct LoginResponse {
    token: String,
    expires_in: i64,
}

#[derive(Deserialize)]
#[allow(dead_code)] // Fields shown for demonstration
struct RegisterRequest {
    name: String,
    email: String,
    password: String,
}

#[derive(Debug, Clone)]
struct CurrentUser {
    id: String,
    role: String,
}

async fn register(Json(input): Json<RegisterRequest>) -> impl IntoResponse {
    let _hashed = hash_password(&input.password);
    Json(serde_json::json!({
        "message": "用户已注册",
        "email": input.email
    }))
}

async fn login(
    State(config): State<Arc<AuthConfig>>,
    Json(input): Json<LoginRequest>,
) -> Result<Json<LoginResponse>, StatusCode> {
    // 模拟用户校验
    // 在真实应用中，我们会在数据库中查询该用户。
    if input.email == "test@example.com" && input.password == "password123" {
        // 如果凭证有效，则用它们的 id 和 role 创建一个令牌。
        let token = create_token(&config, "user-1", "user")?;
        // 以包含过期信息的 JSON 返回。
        Ok(Json(LoginResponse {
            token,
            expires_in: config.jwt_expiry_hours * 3600,
        }))
    } else {
        // 如果无效，返回 401。不要说明失败原因，只说凭证无效。
        Err(StatusCode::UNAUTHORIZED)
    }
}

async fn protected(axum::Extension(user): axum::Extension<CurrentUser>) -> impl IntoResponse {
    Json(serde_json::json!({
        "message": "访问已授权！",
        "user_id": user.id,
        "role": user.role
    }))
}

async fn admin_only(axum::Extension(user): axum::Extension<CurrentUser>) -> impl IntoResponse {
    // 对于更复杂的场景，构建 permission extractors。检查多个角色、特定权限、资源所有权。
    // 但模式是相同的。
    // 在中间件或处理器的早期进行检查。
    if user.role != "admin" {
        return (StatusCode::FORBIDDEN, "需要管理员权限").into_response();
    }
    Json(serde_json::json!({"message": "管理员区域", "user": user.id})).into_response()
}

// ============================================================================
// 主程序
// ============================================================================

#[tokio::main]
async fn main() {
    let config = Arc::new(AuthConfig {
        jwt_secret: "super-secret-key-change-in-production".to_string(),
        jwt_expiry_hours: 24,
    });

    let protected_routes = Router::new()
        .route("/me", get(protected))
        .route("/admin", get(admin_only))
        .route_layer(middleware::from_fn_with_state(
            config.clone(),
            auth_middleware,
        ));

    let app = Router::new()
        .route("/register", post(register))
        .route("/login", post(login))
        .nest("/protected", protected_routes)
        .with_state(config);

    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();

    println!("🚀 模块 09：认证");
    println!("   服务地址：http://localhost:3000\n");
    println!("📝 接口列表：");
    println!("   POST /register    - 注册用户");
    println!("   POST /login       - 登录（test@example.com / password123）");
    println!("   GET  /protected/me - 受保护路由");
    println!("\n💡 使用方式：");
    println!("   1. 使用凭据调用 POST /login");
    println!("   2. 携带令牌：curl -H 'Authorization: Bearer <token>' /protected/me");

    axum::serve(listener, app).await.unwrap();
}
```

## 总结

- 将 JWT 存储在本地存储中是危险的。
  - XSS 攻击可能会窃取它们。
  - HttpOnly Cookie 更安全。
- 未检查令牌过期。总是使用弱秘密验证 exp（过期）时存在问题。
- 使用加密安全的随机生成器生成。
- 传输中的令牌可能会被拦截。
- 登录没有速率限制。
  - 暴力破解攻击。尝试数千个密码。
- 将密钥存储在环境变量中。绝不要对其进行编码。
- 在所有地方使用 HTTPS。Let's Encrypt 是免费的。
- 为访问令牌设置合理的过期时间：24 小时。
- 刷新令牌有效期为 7 天。
  - 实现令牌刷新以在不使用长期有效令牌的情况下提升用户体验。
- 对每个 IP 的登录尝试进行速率限制，每分钟五次。
- 记录认证事件。成功登录。失败的尝试。令牌刷新。
- 对高价值应用考虑多因素认证。
