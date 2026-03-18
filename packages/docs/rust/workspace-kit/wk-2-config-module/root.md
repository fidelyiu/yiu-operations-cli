# config 模块

## 全局 config

`src/config/config.rs`

```rust
#[derive(Debug, Clone)]
pub struct Config {
    pub database_url: String,
    pub jwt_secret: String,
    pub jwt_maxage: i64,
    pub port: u16,
    pub backend_base_url: String,
    pub frontend_base_url: String,
}

impl Config {
    pub fn init() -> Self {
        let database_url = std::env::var("DATABASE_URL").expect("DATABASE_URL must be set");
        let jwt_secret = std::env::var("JWT_SECRET").expect("JWT_SECRET must be set");
        let jwt_maxage = std::env::var("JWT_MAXAGE")
            .expect("JWT_MAXAGE must be set")
            .parse::<i64>()
            .expect("JWT_MAXAGE must be a valid integer");
        let port = std::env::var("PORT")
            .unwrap_or_else(|_| "8000".to_string())
            .parse::<u16>()
            .expect("PORT must be a valid integer");
        let backend_base_url =
            std::env::var("BACKEND_BASE_URL").expect("BACKEND_BASE_URL must be set");
        let frontend_base_url =
            std::env::var("FRONTEND_BASE_URL").expect("FRONTEND_BASE_URL must be set");
        Config {
            database_url,
            jwt_secret,
            jwt_maxage,
            port,
            backend_base_url,
            frontend_base_url,
        }
    }
}
```

## 邮件config

`src/config/mail_config.rs`

```rust
#[derive(Debug, Clone)]
pub struct MailConfig {
    pub smtp_server: String,
    pub smtp_port: u16,
    pub smtp_username: String,
    pub smtp_password: String,
    pub mail_template_path: String,
}

impl MailConfig {
    pub fn init() -> Self {
        let smtp_server = std::env::var("SMTP_SERVER").expect("SMTP_SERVER must be set");
        let smtp_port = std::env::var("SMTP_PORT")
            .unwrap_or_else(|_| "587".to_string())
            .parse::<u16>()
            .expect("SMTP_PORT must be a valid integer");
        let smtp_username = std::env::var("SMTP_USERNAME").expect("SMTP_USERNAME must be set");
        let smtp_password = std::env::var("SMTP_PASSWORD").expect("SMTP_PASSWORD must be set");
        let mail_template_path =
            std::env::var("MAIL_TEMPLATE_PATH").expect("MAIL_TEMPLATE_PATH must be set");
        MailConfig {
            smtp_server,
            smtp_port,
            smtp_username,
            smtp_password,
            mail_template_path,
        }
    }
}
```

## config 模块

`src/config.rs`

```rust
mod config;
mod mail_config;
```

## main 入口

`src/main.rs`

```rust
mod constants;
mod models;
mod config;  // [!code ++]

fn main() {
    println!("Hello, world!");
}
```