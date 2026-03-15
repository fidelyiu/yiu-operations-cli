# 初始化项目

## 创建项目

```sh
cargo new workspace-kit-tutorial
```

## 添加依赖

```toml
[package]
name = "workspace-kit-tutorial"
version = "0.1.0"
edition = "2024"

[dependencies]
argon2 = "0.5.3"
async-trait = "0.1.88"
chrono = { version = "0.4.41", features = ["serde"] }
dotenv = "0.15.0"
jsonwebtoken = "9.3.1"
serde = { version = "1.0.219", features = ["derive"] }
serde_json = "1.0.140"
sqlx = { version = "0.8.6", features = ["runtime-async-std-native-tls", "postgres", "chrono", "uuid"] }
uuid = { version = "1.17.0", features = ["serde", "v4"] }
validator = { version = "0.20.0", features = ["derive"] }
axum = "0.8.4"
axum-extra = { version = "0.10.1", features = ["cookie"]}
tokio = { version = "1.46.1", features = ["full"] }
tower = "0.5.2"
time = "0.3.41"
tower-http = { version = "0.6.6", features = ["cors","trace"] }
tracing-subscriber = { version = "0.3.19"}
lettre = "0.11.17"
regex = "1.11.1"
```

- [argon2](https://docs.rs/argon2/latest/argon2/): 纯 Rust 实现的 Argon2 密码哈希函数。
- [async-trait](https://docs.rs/async-trait/latest/async_trait/): 这个 crate 提供了一个属性宏，使 traits 中的 async fn 能够与 dyn traits 一起使用。
