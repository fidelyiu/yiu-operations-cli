# 初始化项目

## 创建项目

```sh
cargo new workspace-kit-tutorial
```

## 修改ignore

`.gitignore`

```text
/target
temp_data
```

## 添加依赖

`Cargo.toml`

::: details 文件内容

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

:::

- [argon2](https://docs.rs/argon2/latest/argon2/): 纯 Rust 实现的 Argon2 密码哈希函数。
- [async-trait](https://docs.rs/async-trait/latest/async_trait/): 这个 crate 提供了一个属性宏，使 traits 中的 async fn 能够与 dyn traits 一起使用。
- [chrono](https://docs.rs/chrono/latest/chrono/): Rust 的日期和时间
- [dotenv](https://docs.rs/dotenv/latest/dotenv/): 环境变量配置加载器
- [jsonwebtoken](https://docs.rs/jsonwebtoken/latest/jsonwebtoken/): JWT实现
- [serde](https://serde.rs/): 序列化
- [serde_json](https://docs.rs/serde_json/latest/serde_json/): json 序列化
- [sqlx](https://docs.rs/sqlx_wasi/latest/sqlx/): 数据库连接
- [uuid](https://docs.rs/uuid/latest/uuid/): UUID实现
- [validator](https://docs.rs/validator/latest/validator/): 结构体验证
- [axum](https://docs.rs/axum/latest/axum/): web框架
- [axum-extra](https://crates.io/crates/axum-extra): axum 的额外实用程序。
- [tokio](https://docs.rs/tokio/latest/tokio/): rust异步运行时
- [tower](https://docs.rs/tower/latest/tower/): Tower 是一个模块化和可重用组件库，用于构建强大的网络客户端和服务器。
- [time](https://docs.rs/time/latest/time/): time 工具
- [tower-http](https://crates.io/crates/tower-http): Tower 中间件和用于 HTTP 客户端和服务器的实用程序。
- [tracing-subscriber](https://docs.rs/tracing-subscriber/latest/tracing_subscriber/): 追踪订阅者
- [lettre](https://docs.rs/lettre/latest/lettre/): 邮件库
- [regex](https://docs.rs/regex/latest/regex/): 正则

## 创建环境变量文档

`.env`

::: details 文件内容

```sh
# -----------------------------------------------------------------------------
# 数据库（PostgreSQL）
# -----------------------------------------------------------------------------
DATABASE_URL=postgresql://postgres:password@localhost:5432/workspace-kit

# -----------------------------------------------------------------------------
# JSON Web Token 凭证
# -----------------------------------------------------------------------------
JWT_SECRET_KEY=my_ultra_secure_jwt_secret_key
JWT_MAXAGE=86400 # 1 天（以秒为单位）

# -----------------------------------------------------------------------------
# 端口
# -----------------------------------------------------------------------------
PORT=8000

# -----------------------------------------------------------------------------
# 后端 URL
# -----------------------------------------------------------------------------
BACKEND_BASE_URL=http://localhost:8000/api

# -----------------------------------------------------------------------------
# 前端 URL
# -----------------------------------------------------------------------------
FRONTEND_BASE_URL=http://localhost:3000

# -----------------------------------------------------------------------------
# SMTP 服务器设置
# -----------------------------------------------------------------------------
SMTP_SERVER=smtp.your-email-provider.com
SMTP_PORT=587 # 常用端口：587（TLS）、465（SSL）、25（非安全）
SMTP_USERNAME=your_email@example.com
SMTP_PASSWORD=your_email_password
SMTP_FROM_ADDRESS=no-reply@yourdomain.com
MAIL_TEMPLATE_PATH=src/mail/templates
```

:::

## 启动数据库

```sh
docker run --name wk-postgres -p 5432:5432 -e POSTGRES_PASSWORD=password -e POSTGRES_USER=postgres -e POSTGRES_DB=workspace-kit -d postgres
```

这里的 `-p 5432:5432` 表示：

- 前面的 `5432` 是宿主机端口
- 后面的 `5432` 是容器内 PostgreSQL 的端口

这样本机就可以通过 `localhost:5432` 访问容器里的 PostgreSQL。

`docker-compose.yaml`

::: details 文件内容

```yaml
services:
  postgres:
    image: postgres
    container_name: workspace-kit-postgres
    # 容器异常退出或 Docker 重启后会自动重启；如果是手动停止，则不会自动拉起。
    restart: unless-stopped
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: workspace-kit
    ports:
      - "5432:5432"
    volumes:
      - ./temp_data/postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres -d workspace-kit"]
      interval: 10s
      timeout: 5s
      retries: 5
```

:::

命令

```sh
docker compose -f docker-compose.yaml up -d
docker compose -f docker-compose.yaml down
```

如果你的文件名不是默认的 `compose.yaml`、`compose.yml`、`docker-compose.yaml`、`docker-compose.yml`，就需要显式加上 `-f`。

## 迁移数据库

- [sqlx crate](https://docs.rs/sqlx_wasi/latest/sqlx/)
- [sqlx github](https://github.com/launchbadge/sqlx)
- [sqlx cli](https://github.com/launchbadge/sqlx/tree/main/sqlx-cli)

```sh
# 支持 SQLx 支持的所有数据库
cargo install sqlx-cli
```

迁移数据库

```sh
sqlx migrate add init-tables
```

它创建了带有日期标签时间戳和表 sqlx 的迁移。

`migrations/xxxxxxxxxxxxxx_init-tables.sql`

::: details 文件内容

```sql
-- 在此添加迁移脚本
-- UUID 扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 用户表
CREATE TABLE "users" (
    -- 插入用户时如果没传 id，会自动生成一个 UUID
    id UUID NOT NULL PRIMARY KEY DEFAULT (uuid_generate_v4()),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    email_verified BOOLEAN DEFAULT FALSE,
    pending_email VARCHAR(255),
    pending_email_token UUID,
    pending_email_expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_users_email ON "users"(email);

-- 工作区表
CREATE TABLE "workspaces" (
    id UUID NOT NULL PRIMARY KEY DEFAULT (uuid_generate_v4()),
    name TEXT NOT NULL,
    -- 用户删除时，关联工作区也会级联删除
    owner_user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    -- 用随机数和当前时间拼接后取 MD5，再截取前 25 个字符作为邀请码
    invite_code VARCHAR(25) UNIQUE NOT NULL DEFAULT (
      substr(md5(random()::text || clock_timestamp()::text), 0, 26)
      ),
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    -- 限制同一个用户不能创建同名工作区
    UNIQUE (owner_user_id, name)
);

CREATE INDEX idx_workspaces_owner_user_id ON "workspaces"(owner_user_id);

-- 角色表
CREATE TABLE "roles" (
    id UUID NOT NULL PRIMARY KEY DEFAULT (uuid_generate_v4()),
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    UNIQUE(workspace_id, name)
);

CREATE INDEX idx_roles_workspace_id ON "roles"(workspace_id);

-- 权限表（全局）
CREATE TABLE "permissions" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL,
    description TEXT
);

-- 角色权限关联表
CREATE TABLE "role_permissions" (
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID REFERENCES permissions(id) ON DELETE CASCADE,
    -- 联合主键，避免同一角色重复绑定同一权限
    PRIMARY KEY (role_id, permission_id)
);

CREATE INDEX idx_role_permissions_role_id ON "role_permissions"(role_id);
CREATE INDEX idx_role_permissions_permission_id ON "role_permissions"(permission_id);

-- 工作区用户关联表
CREATE TABLE "workspace_users" (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    -- 如果角色被删，成员记录保留，只是 role_id 变成 NULL
    role_id UUID REFERENCES roles(id) ON DELETE SET NULL,
    status TEXT DEFAULT 'active',
    -- 一个用户在一个工作区里只能出现一次
    PRIMARY KEY (user_id, workspace_id)
);

CREATE INDEX idx_workspace_users_user_id ON "workspace_users"(user_id);
CREATE INDEX idx_workspace_users_workspace_id ON "workspace_users"(workspace_id);


-- 邮箱验证表
CREATE TABLE email_verifications (
     user_id UUID REFERENCES users(id) ON DELETE CASCADE,
     token UUID UNIQUE NOT NULL,
     expires_at TIMESTAMPTZ NOT NULL,
     PRIMARY KEY (user_id)
);

CREATE INDEX idx_email_verifications_token ON email_verifications(token);

-- 密码重置表
CREATE TABLE password_resets (
     user_id UUID REFERENCES users(id) ON DELETE CASCADE,
     token UUID UNIQUE NOT NULL,
     expires_at TIMESTAMPTZ NOT NULL,
     PRIMARY KEY (user_id)
);

CREATE INDEX idx_password_resets_token ON password_resets(token);

INSERT INTO permissions (id, name, description) VALUES
    (gen_random_uuid(), 'update_workspace', 'Update workspace details'),
    (gen_random_uuid(), 'delete_workspace', 'Delete the workspace'),
    (gen_random_uuid(), 'manage_roles', 'Manage workspace roles'),
    (gen_random_uuid(), 'manage_permissions', 'Assign permissions to roles'),
    (gen_random_uuid(), 'invite_members', 'Invite new members to workspace'),
    (gen_random_uuid(), 'view_members', 'View workspace members'),
    (gen_random_uuid(), 'view_roles', 'View available roles'),
    (gen_random_uuid(), 'view_permissions', 'View available permissions'),
    (gen_random_uuid(), 'remove_members', 'Ability to remove members from the workspace'),
    (gen_random_uuid(), 'assign_roles_to_members', 'Ability to assign roles to workspace members');

CREATE OR REPLACE FUNCTION create_default_roles_for_workspace()
RETURNS TRIGGER AS $$
DECLARE
    admin_role_id UUID;
    manager_role_id UUID;
BEGIN
    INSERT INTO roles (workspace_id, name, description)
    VALUES
        (NEW.id, 'Admin', 'Workspace administrator with full permissions'),
        (NEW.id, 'Manager', 'Workspace manager with limited permissions');
    SELECT id INTO admin_role_id FROM roles WHERE workspace_id = NEW.id AND name = 'Admin';
    SELECT id INTO manager_role_id FROM roles WHERE workspace_id = NEW.id AND name = 'Manager';
    INSERT INTO role_permissions (role_id, permission_id)
    SELECT admin_role_id, id FROM permissions;
    INSERT INTO role_permissions (role_id, permission_id)
    SELECT manager_role_id, id FROM permissions
    WHERE name IN ('view_members', 'view_roles', 'view_permissions', 'invite_members');
    INSERT INTO workspace_users (user_id, workspace_id, role_id, status)
    VALUES (NEW.owner_user_id, NEW.id, admin_role_id, 'active');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER create_default_roles_for_workspace
AFTER INSERT ON workspaces
FOR EACH ROW
EXECUTE FUNCTION create_default_roles_for_workspace();

CREATE OR REPLACE FUNCTION ensure_single_default_workspace()
RETURNS TRIGGER AS $$
DECLARE
    default_count INT;
BEGIN
    SELECT COUNT(*) INTO default_count
    FROM workspaces
    WHERE owner_user_id = NEW.owner_user_id
     AND is_default = TRUE;
    IF default_count > 0 THEN
        NEW.is_default = FALSE;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER ensure_single_default_workspace_trigger
BEFORE INSERT ON workspaces
FOR EACH ROW
EXECUTE FUNCTION ensure_single_default_workspace();
```

:::

## 执行数据库迁移

所有命令都需要提供数据库 URL。这可以通过 `--database-url` 命令行选项或设置 `DATABASE_URL` 来实现，后者可以在环境变量中设置，也可以在当前工作目录下的 `.env` 文件中设置。

```sh
# 读取 DATABASE_URL 然后去尝试创建数据库
sqlx database create
# 执行迁移脚本
sqlx migrate run
```

> 当你运行 `sqlx migrate run`：
>
> 如果数据库里还没有记录 `20260315061043`，这个文件就会执行。
> 如果数据库里已经有这条记录，这个文件就会跳过，不会重复跑。
>
> 本地 `migrations` 目录 = “应该有哪些变更”
> 数据库里的 `_sqlx_migrations` = “已经执行到哪里了”

## 创建 constants 模块

`constants.rs`

::: details 文件内容

```rust
pub mod permissions {
    pub const UPDATE_WORKSPACE: &str = "update_workspace";
    pub const DELETE_WORKSPACE: &str = "delete_workspace";
    pub const MANAGE_ROLES: &str = "manage_roles";
    pub const MANAGE_PERMISSIONS: &str = "manage_permissions";
    pub const INVITE_MEMBERS: &str = "invite_members";
    pub const VIEW_MEMBERS: &str = "view_members";
    pub const VIEW_ROLES: &str = "view_roles";
    pub const VIEW_PERMISSIONS: &str = "view_permissions";
    pub const REMOVE_MEMBERS: &str = "remove_members";
    pub const ASSIGN_ROLES_TO_MEMBERS: &str = "assign_roles_to_members";

    pub const ALL: [&str; 10] = [
        UPDATE_WORKSPACE,
        DELETE_WORKSPACE,
        MANAGE_ROLES,
        MANAGE_PERMISSIONS,
        INVITE_MEMBERS,
        VIEW_MEMBERS,
        VIEW_ROLES,
        VIEW_PERMISSIONS,
        REMOVE_MEMBERS,
        ASSIGN_ROLES_TO_MEMBERS,
    ];
}
```

:::

注册模块

```rust
mod constants;

fn main() {
    println!("Hello, world!");
}
```

## 创建 models 模块

`models.rs`

::: details 文件内容

```rust
use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sqlx::types::Uuid;

#[derive(Debug, Serialize, sqlx::FromRow, Deserialize, Clone)]
pub struct User {
    pub id: Uuid,
    pub name: String,
    pub email: String,
    pub password: String,
    pub email_verified: Option<bool>,
    pub pending_email: Option<String>,
    pub pending_email_token: Option<Uuid>,
    pub pending_email_expires_at: Option<DateTime<Utc>>,
    pub created_at: Option<DateTime<Utc>>,
    pub updated_at: Option<DateTime<Utc>>,
}

#[derive(Debug, Serialize, sqlx::FromRow)]
pub struct Workspace {
    pub id: Uuid,
    pub name: String,
    pub owner_user_id: Option<Uuid>,
    pub invite_code: String,
    pub is_default: Option<bool>,
    pub created_at: Option<DateTime<Utc>>,
    pub updated_at: Option<DateTime<Utc>>,
}

#[derive(Debug, Serialize, sqlx::FromRow)]
pub struct WorkspaceUser {
    pub user_id: Uuid,
    pub workspace_id: Uuid,
    pub role_id: Option<Uuid>,
    pub status: String,
}

#[derive(Debug, Serialize, sqlx::FromRow)]
pub struct Role {
    pub id: Uuid,
    pub workspace_id: Uuid,
    pub name: String,
    pub description: Option<String>,
}

#[derive(Debug, Serialize, sqlx::FromRow)]
pub struct Permission {
    pub id: Uuid,
    pub name: String,
    pub description: String,
}

#[derive(Debug, Serialize, sqlx::FromRow)]
pub struct RolePermission {
    pub role_id: Uuid,
    pub permission_id: Uuid,
}

#[derive(Debug, Serialize, sqlx::FromRow)]
pub struct EmailVerification {
    pub user_id: Uuid,
    pub token: Uuid,
    pub expires_at: DateTime<Utc>,
}

#[derive(Debug, Serialize, sqlx::FromRow)]
pub struct PasswordReset {
    pub user_id: Uuid,
    pub token: Uuid,
    pub expires_at: DateTime<Utc>,
}
```

:::

注册模块

`main.rs`

```rust
// ...
mod models;

fn main() {
    println!("Hello, world!");
}
```

## mail 模板

`src/mail/templates/email-change-verification.html`

::: details

```html
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Verify Email Change - workspace-kit</title>
    <style>
      * {
        margin: 0;
        padding: 0;
        box-sizing: border-box;
      }
      body {
        font-family:
          "Inter",
          -apple-system,
          BlinkMacSystemFont,
          "Segoe UI",
          Roboto,
          sans-serif;
        background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
        min-height: 100vh;
        padding: 20px;
      }
      .email-wrapper {
        max-width: 600px;
        margin: 0 auto;
        background: #ffffff;
        border-radius: 24px;
        overflow: hidden;
        box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
      }
      .header-section {
        background: #ffffff;
        padding: 48px 40px 32px;
        text-align: center;
        position: relative;
      }
      .change-icon {
        width: 80px;
        height: 80px;
        background: linear-gradient(135deg, #f59e0b, #d97706);
        border-radius: 50%;
        margin: 0 auto 24px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 36px;
        animation: rotate 3s linear infinite;
      }
      @keyframes rotate {
        from {
          transform: rotate(0deg);
        }
        to {
          transform: rotate(360deg);
        }
      }
      .change-badge {
        display: inline-flex;
        align-items: center;
        gap: 8px;
        background: linear-gradient(135deg, #f59e0b, #d97706);
        color: white;
        padding: 8px 16px;
        border-radius: 50px;
        font-size: 14px;
        font-weight: 600;
        margin-bottom: 24px;
      }
      .main-title {
        font-size: 28px;
        font-weight: 800;
        color: #1a1a1a;
        margin-bottom: 12px;
        line-height: 1.2;
      }
      .subtitle {
        font-size: 16px;
        color: #6b7280;
        font-weight: 400;
        line-height: 1.5;
      }
      .content-section {
        padding: 0 40px 48px;
      }
      .personal-greeting {
        background: linear-gradient(135deg, #fefbf3 0%, #fef3c7 100%);
        border: 2px solid #fed7aa;
        border-radius: 20px;
        padding: 32px;
        margin-bottom: 32px;
        text-align: center;
      }
      .greeting-text {
        font-size: 20px;
        font-weight: 700;
        color: #92400e;
        margin-bottom: 12px;
      }
      .greeting-message {
        font-size: 16px;
        color: #d97706;
        line-height: 1.6;
      }
      .email-comparison {
        background: #ffffff;
        border: 2px solid #e5e7eb;
        border-radius: 16px;
        padding: 24px;
        margin: 32px 0;
      }
      .comparison-title {
        font-size: 18px;
        font-weight: 700;
        color: #374151;
        margin-bottom: 20px;
        text-align: center;
      }
      .email-row {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 16px;
        margin-bottom: 12px;
        background: #f9fafb;
        border-radius: 12px;
        border: 1px solid #e5e7eb;
      }
      .email-row:last-child {
        margin-bottom: 0;
        background: #fefbf3;
        border-color: #fed7aa;
      }
      .email-label {
        font-weight: 600;
        color: #374151;
        font-size: 14px;
      }
      .email-value {
        font-family: "Monaco", "Menlo", monospace;
        font-size: 14px;
        color: #6b7280;
        background: #f3f4f6;
        padding: 6px 12px;
        border-radius: 6px;
      }
      .new-email {
        color: #d97706 !important;
        background: #fef3c7 !important;
        font-weight: 600;
      }
      .verification-card {
        background: #ffffff;
        border: 3px solid #f59e0b;
        border-radius: 20px;
        padding: 40px;
        text-align: center;
        margin: 32px 0;
        position: relative;
        overflow: hidden;
      }
      .verification-card::before {
        content: "";
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        height: 4px;
        background: linear-gradient(90deg, #fbbf24, #f59e0b, #d97706);
      }
      .verification-title {
        font-size: 20px;
        font-weight: 700;
        color: #92400e;
        margin-bottom: 16px;
      }
      .confirm-button {
        display: inline-flex;
        align-items: center;
        gap: 12px;
        background: linear-gradient(135deg, #f59e0b, #d97706);
        color: white;
        text-decoration: none;
        padding: 20px 40px;
        border-radius: 16px;
        font-weight: 700;
        font-size: 18px;
        transition: all 0.3s ease;
        box-shadow: 0 8px 32px rgba(245, 158, 11, 0.3);
      }
      .confirm-button:hover {
        transform: translateY(-2px);
        box-shadow: 0 12px 40px rgba(245, 158, 11, 0.4);
      }
      .security-info {
        background: #fffbeb;
        border: 2px solid #fcd34d;
        border-radius: 16px;
        padding: 24px;
        margin: 32px 0;
      }
      .security-title {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 16px;
        font-weight: 700;
        color: #92400e;
        margin-bottom: 12px;
      }
      .security-text {
        font-size: 14px;
        color: #92400e;
        line-height: 1.6;
      }
      .danger-alert {
        background: #fef2f2;
        border: 2px solid #fca5a5;
        border-radius: 16px;
        padding: 24px;
        margin: 32px 0;
      }
      .danger-title {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 16px;
        font-weight: 700;
        color: #dc2626;
        margin-bottom: 12px;
      }
      .danger-text {
        font-size: 14px;
        color: #dc2626;
        line-height: 1.6;
      }
      .next-steps {
        background: #f0f9ff;
        border: 2px solid #93c5fd;
        border-radius: 16px;
        padding: 24px;
        margin: 32px 0;
      }
      .steps-title {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 16px;
        font-weight: 700;
        color: #1e40af;
        margin-bottom: 16px;
      }
      .steps-list {
        list-style: none;
        padding: 0;
      }
      .steps-list li {
        font-size: 14px;
        color: #1e40af;
        line-height: 1.6;
        margin-bottom: 8px;
        padding-left: 24px;
        position: relative;
      }
      .steps-list li::before {
        content: counter(step-counter);
        counter-increment: step-counter;
        position: absolute;
        left: 0;
        top: 0;
        background: #3b82f6;
        color: white;
        width: 18px;
        height: 18px;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 12px;
        font-weight: bold;
      }
      .steps-list {
        counter-reset: step-counter;
      }
      .link-section {
        background: #f8fafc;
        border-radius: 12px;
        padding: 20px;
        margin: 24px 0;
      }
      .link-label {
        font-size: 14px;
        color: #64748b;
        margin-bottom: 8px;
        font-weight: 600;
      }
      .link-text {
        font-family: "Monaco", "Menlo", monospace;
        font-size: 12px;
        color: #f59e0b;
        word-break: break-all;
        background: #fefbf3;
        padding: 12px;
        border-radius: 8px;
        border: 1px solid #fed7aa;
      }
      .footer-section {
        background: #f8fafc;
        padding: 32px 40px;
        text-align: center;
        border-top: 1px solid #e2e8f0;
      }
      .footer-text {
        color: #64748b;
        font-size: 14px;
        margin-bottom: 20px;
      }
      .footer-links {
        display: flex;
        justify-content: center;
        gap: 32px;
        flex-wrap: wrap;
      }
      .footer-link {
        color: #f59e0b;
        text-decoration: none;
        font-weight: 600;
        font-size: 14px;
        transition: color 0.2s ease;
      }
      .footer-link:hover {
        color: #d97706;
      }
      @media (max-width: 640px) {
        .email-wrapper {
          margin: 0;
          border-radius: 0;
        }
        .header-section,
        .content-section {
          padding-left: 24px;
          padding-right: 24px;
        }
        .email-row {
          flex-direction: column;
          align-items: flex-start;
          gap: 8px;
        }
        .footer-links {
          flex-direction: column;
          gap: 16px;
        }
      }
    </style>
  </head>
  <body>
    <div class="email-wrapper">
      <div class="header-section">
        <div class="change-icon">🔄</div>
        <div class="change-badge">
          <span>📧</span>
          Email Change
        </div>
        <h1 class="main-title">Hi {{ .Name }}, confirm your email change</h1>
        <p class="subtitle">Verify your new email address</p>
      </div>

      <div class="content-section">
        <div class="personal-greeting">
          <div class="greeting-text">Hello {{ .Name }}! 👋</div>
          <div class="greeting-message">
            We received a request to change the email address for your
            workspace-kit account. Please verify this change to complete the
            process.
          </div>
        </div>

        <div class="email-comparison">
          <div class="comparison-title">📋 Email Address Change Summary</div>
          <div class="email-row">
            <span class="email-label">Current Email:</span>
            <span class="email-value">{{ .Email }}</span>
          </div>
          <div class="email-row">
            <span class="email-label">New Email:</span>
            <span class="email-value new-email">{{ .NewEmail }}</span>
          </div>
        </div>

        <div class="verification-card">
          <div class="verification-title">✅ Confirm Email Change</div>
          <p style="color: #6b7280; margin-bottom: 24px; font-size: 16px">
            Click the button below to confirm {{ .NewEmail }} as your new email
            address
          </p>
          <a href="{{ .ConfirmationURL }}" class="confirm-button">
            <span>🔐</span>
            Confirm Change
          </a>
        </div>

        <div class="security-info">
          <div class="security-title">
            <span>⏰</span>
            Time Sensitive
          </div>
          <div class="security-text">
            This verification link expires in 24 hours. After confirmation,
            you'll need to use {{ .NewEmail }} to sign in to workspace-kit.
          </div>
        </div>

        <div class="danger-alert">
          <div class="danger-title">
            <span>🚨</span>
            Security Alert
          </div>
          <div class="danger-text">
            If you didn't request this email change, please contact our support
            team immediately at support@workspace-kit.com. Someone may be trying
            to access your account.
          </div>
        </div>

        <div class="link-section">
          <div class="link-label">
            Having trouble with the button? Copy this link:
          </div>
          <div class="link-text">{{ .ConfirmationURL }}</div>
        </div>

        <div class="next-steps">
          <div class="steps-title">
            <span>📝</span>
            What Happens Next
          </div>
          <ul class="steps-list">
            <li>
              Use <strong>{{ .NewEmail }}</strong> to sign in to workspace-kit
            </li>
            <li>Update your email in any connected third-party services</li>
            <li>Check your new inbox for future workspace-kit notifications</li>
            <li>Update your notification preferences in account settings</li>
          </ul>
        </div>

        <div
          style="
            text-align: center;
            margin-top: 32px;
            padding: 24px;
            background: #fefbf3;
            border-radius: 16px;
            border: 1px solid #fed7aa;
          "
        >
          <p style="color: #d97706; font-size: 16px; margin-bottom: 8px">
            Questions about this change?
          </p>
          <p style="color: #92400e; font-weight: 600">
            The workspace-kit Account Security Team
          </p>
        </div>
      </div>

      <!-- <div class="footer-section">
        <p class="footer-text">&copy; 2024 workspace-kit. All rights reserved.</p>
        <div class="footer-links">
            <a href="{{ .SiteURL }}/privacy" class="footer-link">Privacy</a>
            <a href="{{ .SiteURL }}/terms" class="footer-link">Terms</a>
            <a href="{{ .SiteURL }}/contact" class="footer-link">Support</a>
        </div>
    </div> -->
    </div>
  </body>
</html>
```

:::

`src/mail/templates/reset-password-email.html`

::: details

```html
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Reset Your Password - workspace-kit</title>
    <style>
      * {
        margin: 0;
        padding: 0;
        box-sizing: border-box;
      }
      body {
        font-family:
          "Inter",
          -apple-system,
          BlinkMacSystemFont,
          "Segoe UI",
          Roboto,
          sans-serif;
        background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
        min-height: 100vh;
        padding: 20px;
      }
      .email-wrapper {
        max-width: 600px;
        margin: 0 auto;
        background: #ffffff;
        border-radius: 24px;
        overflow: hidden;
        box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
      }
      .header-section {
        background: #ffffff;
        padding: 48px 40px 32px;
        text-align: center;
        position: relative;
      }
      .security-icon {
        width: 80px;
        height: 80px;
        background: linear-gradient(135deg, #ef4444, #dc2626);
        border-radius: 50%;
        margin: 0 auto 24px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 36px;
        animation: shake 0.5s ease-in-out;
      }
      @keyframes shake {
        0%,
        100% {
          transform: translateX(0);
        }
        25% {
          transform: translateX(-5px);
        }
        75% {
          transform: translateX(5px);
        }
      }
      .alert-badge {
        display: inline-flex;
        align-items: center;
        gap: 8px;
        background: linear-gradient(135deg, #ef4444, #dc2626);
        color: white;
        padding: 8px 16px;
        border-radius: 50px;
        font-size: 14px;
        font-weight: 600;
        margin-bottom: 24px;
      }
      .main-title {
        font-size: 28px;
        font-weight: 800;
        color: #1a1a1a;
        margin-bottom: 12px;
        line-height: 1.2;
      }
      .subtitle {
        font-size: 16px;
        color: #6b7280;
        font-weight: 400;
        line-height: 1.5;
      }
      .content-section {
        padding: 0 40px 48px;
      }
      .personal-greeting {
        background: linear-gradient(135deg, #fef2f2 0%, #fee2e2 100%);
        border: 2px solid #fca5a5;
        border-radius: 20px;
        padding: 32px;
        margin-bottom: 32px;
        text-align: center;
      }
      .greeting-text {
        font-size: 20px;
        font-weight: 700;
        color: #991b1b;
        margin-bottom: 12px;
      }
      .greeting-message {
        font-size: 16px;
        color: #dc2626;
        line-height: 1.6;
      }
      .reset-card {
        background: #ffffff;
        border: 3px solid #ef4444;
        border-radius: 20px;
        padding: 40px;
        text-align: center;
        margin: 32px 0;
        position: relative;
        overflow: hidden;
      }
      .reset-card::before {
        content: "";
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        height: 4px;
        background: linear-gradient(90deg, #ef4444, #dc2626, #b91c1c);
      }
      .reset-title {
        font-size: 20px;
        font-weight: 700;
        color: #991b1b;
        margin-bottom: 16px;
      }
      .reset-button {
        display: inline-flex;
        align-items: center;
        gap: 12px;
        background: linear-gradient(135deg, #ef4444, #dc2626);
        color: white;
        text-decoration: none;
        padding: 20px 40px;
        border-radius: 16px;
        font-weight: 700;
        font-size: 18px;
        transition: all 0.3s ease;
        box-shadow: 0 8px 32px rgba(239, 68, 68, 0.3);
      }
      .reset-button:hover {
        transform: translateY(-2px);
        box-shadow: 0 12px 40px rgba(239, 68, 68, 0.4);
      }
      .warning-card {
        background: #fffbeb;
        border: 2px solid #fcd34d;
        border-radius: 16px;
        padding: 24px;
        margin: 32px 0;
      }
      .warning-title {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 16px;
        font-weight: 700;
        color: #92400e;
        margin-bottom: 12px;
      }
      .warning-text {
        font-size: 14px;
        color: #92400e;
        line-height: 1.6;
      }
      .danger-alert {
        background: #fef2f2;
        border: 2px solid #fca5a5;
        border-radius: 16px;
        padding: 24px;
        margin: 32px 0;
      }
      .danger-title {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 16px;
        font-weight: 700;
        color: #dc2626;
        margin-bottom: 12px;
      }
      .danger-text {
        font-size: 14px;
        color: #dc2626;
        line-height: 1.6;
      }
      .security-tips {
        background: #f0f9ff;
        border: 2px solid #93c5fd;
        border-radius: 16px;
        padding: 24px;
        margin: 32px 0;
      }
      .tips-title {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 16px;
        font-weight: 700;
        color: #1e40af;
        margin-bottom: 16px;
      }
      .tips-list {
        list-style: none;
        padding: 0;
      }
      .tips-list li {
        font-size: 14px;
        color: #1e40af;
        line-height: 1.6;
        margin-bottom: 8px;
        padding-left: 20px;
        position: relative;
      }
      .tips-list li::before {
        content: "✓";
        position: absolute;
        left: 0;
        color: #10b981;
        font-weight: bold;
      }
      .link-section {
        background: #f8fafc;
        border-radius: 12px;
        padding: 20px;
        margin: 24px 0;
      }
      .link-label {
        font-size: 14px;
        color: #64748b;
        margin-bottom: 8px;
        font-weight: 600;
      }
      .link-text {
        font-family: "Monaco", "Menlo", monospace;
        font-size: 12px;
        color: #ef4444;
        word-break: break-all;
        background: #fef2f2;
        padding: 12px;
        border-radius: 8px;
        border: 1px solid #fca5a5;
      }
      .footer-section {
        background: #f8fafc;
        padding: 32px 40px;
        text-align: center;
        border-top: 1px solid #e2e8f0;
      }
      .footer-text {
        color: #64748b;
        font-size: 14px;
        margin-bottom: 20px;
      }
      .footer-links {
        display: flex;
        justify-content: center;
        gap: 32px;
        flex-wrap: wrap;
      }
      .footer-link {
        color: #ef4444;
        text-decoration: none;
        font-weight: 600;
        font-size: 14px;
        transition: color 0.2s ease;
      }
      .footer-link:hover {
        color: #dc2626;
      }
      @media (max-width: 640px) {
        .email-wrapper {
          margin: 0;
          border-radius: 0;
        }
        .header-section,
        .content-section {
          padding-left: 24px;
          padding-right: 24px;
        }
        .footer-links {
          flex-direction: column;
          gap: 16px;
        }
      }
    </style>
  </head>
  <body>
    <div class="email-wrapper">
      <div class="header-section">
        <div class="security-icon">🔐</div>
        <div class="alert-badge">
          <span>🚨</span>
          Password Reset
        </div>
        <h1 class="main-title">Hi {{ .Name }}, let's reset your password</h1>
        <p class="subtitle">Secure your workspace-kit account</p>
      </div>

      <div class="content-section">
        <div class="personal-greeting">
          <div class="greeting-text">Hello {{ .Name }}! 👋</div>
          <div class="greeting-message">
            We received a request to reset the password for your workspace-kit
            account ({{ .Email }}). If this was you, click the button below to
            create a new password.
          </div>
        </div>

        <div class="reset-card">
          <div class="reset-title">🔑 Create New Password</div>
          <p style="color: #6b7280; margin-bottom: 24px; font-size: 16px">
            This secure link will take you to a page where you can set a new
            password
          </p>
          <a href="{{ .ConfirmationURL }}" class="reset-button">
            <span>🔒</span>
            Reset My Password
          </a>
        </div>

        <div class="warning-card">
          <div class="warning-title">
            <span>⏰</span>
            Time Sensitive
          </div>
          <div class="warning-text">
            This password reset link expires in 1 hour for your security. You
            can only use this link once.
          </div>
        </div>

        <div class="danger-alert">
          <div class="danger-title">
            <span>🚨</span>
            Didn't Request This?
          </div>
          <div class="danger-text">
            If you didn't request a password reset, please ignore this email.
            Your password remains secure and unchanged. Consider enabling
            two-factor authentication for extra security.
          </div>
        </div>

        <div class="link-section">
          <div class="link-label">
            Having trouble with the button? Copy this link:
          </div>
          <div class="link-text">{{ .ConfirmationURL }}</div>
        </div>

        <div class="security-tips">
          <div class="tips-title">
            <span>🛡️</span>
            Password Security Tips
          </div>
          <ul class="tips-list">
            <li>
              Use at least 12 characters with a mix of letters, numbers, and
              symbols
            </li>
            <li>Avoid using personal information or common words</li>
            <li>Don't reuse passwords from other accounts</li>
            <li>Consider using a password manager for better security</li>
          </ul>
        </div>

        <div
          style="
            text-align: center;
            margin-top: 32px;
            padding: 24px;
            background: #fef2f2;
            border-radius: 16px;
            border: 1px solid #fca5a5;
          "
        >
          <p style="color: #dc2626; font-size: 16px; margin-bottom: 8px">
            Need immediate help?
          </p>
          <p style="color: #991b1b; font-weight: 600">
            The workspace-kit Security Team
          </p>
        </div>
      </div>

      <!-- <div class="footer-section">
        <p class="footer-text">&copy; 2024 workspace-kit. All rights reserved.</p>
        <div class="footer-links">
            <a href="{{ .SiteURL }}/privacy" class="footer-link">Privacy</a>
            <a href="{{ .SiteURL }}/terms" class="footer-link">Terms</a>
            <a href="{{ .SiteURL }}/contact" class="footer-link">Support</a>
        </div>
    </div> -->
    </div>
  </body>
</html>
```

:::

`src/mail/templates/verification-email.html`

::: details

```html
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Verify Your Email - workspace-kit</title>
    <style>
      * {
        margin: 0;
        padding: 0;
        box-sizing: border-box;
      }
      body {
        font-family:
          "Inter",
          -apple-system,
          BlinkMacSystemFont,
          "Segoe UI",
          Roboto,
          sans-serif;
        background: linear-gradient(135deg, #10b981 0%, #059669 100%);
        min-height: 100vh;
        padding: 20px;
      }
      .email-wrapper {
        max-width: 600px;
        margin: 0 auto;
        background: #ffffff;
        border-radius: 24px;
        overflow: hidden;
        box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
      }
      .header-section {
        background: #ffffff;
        padding: 48px 40px 32px;
        text-align: center;
        position: relative;
      }
      .verification-icon {
        width: 80px;
        height: 80px;
        background: linear-gradient(135deg, #10b981, #059669);
        border-radius: 50%;
        margin: 0 auto 24px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 36px;
        animation: pulse 2s infinite;
      }
      @keyframes pulse {
        0%,
        100% {
          transform: scale(1);
        }
        50% {
          transform: scale(1.05);
        }
      }
      .status-badge {
        display: inline-flex;
        align-items: center;
        gap: 8px;
        background: linear-gradient(135deg, #10b981, #059669);
        color: white;
        padding: 8px 16px;
        border-radius: 50px;
        font-size: 14px;
        font-weight: 600;
        margin-bottom: 24px;
      }
      .main-title {
        font-size: 28px;
        font-weight: 800;
        color: #1a1a1a;
        margin-bottom: 12px;
        line-height: 1.2;
      }
      .subtitle {
        font-size: 16px;
        color: #6b7280;
        font-weight: 400;
        line-height: 1.5;
      }
      .content-section {
        padding: 0 40px 48px;
      }
      .personal-greeting {
        background: linear-gradient(135deg, #ecfdf5 0%, #d1fae5 100%);
        border: 2px solid #a7f3d0;
        border-radius: 20px;
        padding: 32px;
        margin-bottom: 32px;
        text-align: center;
      }
      .greeting-text {
        font-size: 20px;
        font-weight: 700;
        color: #065f46;
        margin-bottom: 12px;
      }
      .greeting-message {
        font-size: 16px;
        color: #047857;
        line-height: 1.6;
      }
      .verification-card {
        background: #ffffff;
        border: 3px solid #10b981;
        border-radius: 20px;
        padding: 40px;
        text-align: center;
        margin: 32px 0;
        position: relative;
        overflow: hidden;
      }
      .verification-card::before {
        content: "";
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        height: 4px;
        background: linear-gradient(90deg, #10b981, #059669, #047857);
      }
      .verification-title {
        font-size: 20px;
        font-weight: 700;
        color: #065f46;
        margin-bottom: 16px;
      }
      .verify-button {
        display: inline-flex;
        align-items: center;
        gap: 12px;
        background: linear-gradient(135deg, #10b981, #059669);
        color: white;
        text-decoration: none;
        padding: 20px 40px;
        border-radius: 16px;
        font-weight: 700;
        font-size: 18px;
        transition: all 0.3s ease;
        box-shadow: 0 8px 32px rgba(16, 185, 129, 0.3);
      }
      .verify-button:hover {
        transform: translateY(-2px);
        box-shadow: 0 12px 40px rgba(16, 185, 129, 0.4);
      }
      .security-info {
        background: #fffbeb;
        border: 2px solid #fcd34d;
        border-radius: 16px;
        padding: 24px;
        margin: 32px 0;
      }
      .security-title {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 16px;
        font-weight: 700;
        color: #92400e;
        margin-bottom: 12px;
      }
      .security-text {
        font-size: 14px;
        color: #92400e;
        line-height: 1.6;
      }
      .link-section {
        background: #f8fafc;
        border-radius: 12px;
        padding: 20px;
        margin: 24px 0;
      }
      .link-label {
        font-size: 14px;
        color: #64748b;
        margin-bottom: 8px;
        font-weight: 600;
      }
      .link-text {
        font-family: "Monaco", "Menlo", monospace;
        font-size: 12px;
        color: #10b981;
        word-break: break-all;
        background: #ecfdf5;
        padding: 12px;
        border-radius: 8px;
        border: 1px solid #a7f3d0;
      }
      .footer-section {
        background: #f8fafc;
        padding: 32px 40px;
        text-align: center;
        border-top: 1px solid #e2e8f0;
      }
      .footer-text {
        color: #64748b;
        font-size: 14px;
        margin-bottom: 20px;
      }
      .footer-links {
        display: flex;
        justify-content: center;
        gap: 32px;
        flex-wrap: wrap;
      }
      .footer-link {
        color: #10b981;
        text-decoration: none;
        font-weight: 600;
        font-size: 14px;
        transition: color 0.2s ease;
      }
      .footer-link:hover {
        color: #059669;
      }
      @media (max-width: 640px) {
        .email-wrapper {
          margin: 0;
          border-radius: 0;
        }
        .header-section,
        .content-section {
          padding-left: 24px;
          padding-right: 24px;
        }
        .footer-links {
          flex-direction: column;
          gap: 16px;
        }
      }
    </style>
  </head>
  <body>
    <div class="email-wrapper">
      <div class="header-section">
        <div class="verification-icon">✉️</div>
        <div class="status-badge">
          <span>🔐</span>
          Email Verification
        </div>
        <h1 class="main-title">Almost there, {{ .Name }}!</h1>
        <p class="subtitle">Just one click to secure your account</p>
      </div>

      <div class="content-section">
        <div class="personal-greeting">
          <div class="greeting-text">Hi {{ .Name }}! 👋</div>
          <div class="greeting-message">
            Thanks for joining workspace-kit! We just need to verify your email
            address to complete your account setup and ensure your security.
          </div>
        </div>

        <div class="verification-card">
          <div class="verification-title">🎯 Verify Your Email Address</div>
          <p style="color: #6b7280; margin-bottom: 24px; font-size: 16px">
            Click the button below to confirm {{ .Email }} and activate your
            account
          </p>
          <a href="{{ .ConfirmationURL }}" class="verify-button">
            <span>✅</span>
            Verify My Email
          </a>
        </div>

        <div class="security-info">
          <div class="security-title">
            <span>🛡️</span>
            Security Notice
          </div>
          <div class="security-text">
            This verification link expires in 24 hours for your protection. If
            you didn't create this account, you can safely ignore this email.
          </div>
        </div>

        <div class="link-section">
          <div class="link-label">
            Having trouble with the button? Copy this link:
          </div>
          <div class="link-text">{{ .ConfirmationURL }}</div>
        </div>

        <div
          style="
            text-align: center;
            margin-top: 32px;
            padding: 24px;
            background: #f0fdf4;
            border-radius: 16px;
            border: 1px solid #bbf7d0;
          "
        >
          <p style="color: #047857; font-size: 16px; margin-bottom: 8px">
            Need help? We're here for you!
          </p>
          <p style="color: #065f46; font-weight: 600">The workspace-kit Team</p>
        </div>
      </div>

      <div class="footer-section">
        <p class="footer-text">
          &copy; 2024 workspace-kit. All rights reserved.
        </p>
        <!--        <div class="footer-links">-->
        <!--            <a href="{{ .SiteURL }}/privacy" class="footer-link">Privacy</a>-->
        <!--            <a href="{{ .SiteURL }}/terms" class="footer-link">Terms</a>-->
        <!--            <a href="{{ .SiteURL }}/contact" class="footer-link">Support</a>-->
        <!--        </div>-->
      </div>
    </div>
  </body>
</html>
```

:::

`src/mail/templates/welcome-email.html`

::: details

```html
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Welcome to workspace-kit</title>
    <style>
      * {
        margin: 0;
        padding: 0;
        box-sizing: border-box;
      }
      body {
        font-family:
          "Inter",
          -apple-system,
          BlinkMacSystemFont,
          "Segoe UI",
          Roboto,
          sans-serif;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        min-height: 100vh;
        padding: 20px;
      }
      .email-wrapper {
        max-width: 640px;
        margin: 0 auto;
        background: #ffffff;
        border-radius: 24px;
        overflow: hidden;
        box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
      }
      .header-section {
        background: #ffffff;
        padding: 48px 40px 32px;
        text-align: center;
        position: relative;
      }
      .brand-logo {
        width: 56px;
        height: 56px;
        background: linear-gradient(135deg, #667eea, #764ba2);
        border-radius: 16px;
        margin: 0 auto 24px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-weight: 800;
        color: white;
        font-size: 24px;
        letter-spacing: -1px;
      }
      .welcome-badge {
        display: inline-flex;
        align-items: center;
        gap: 8px;
        background: linear-gradient(135deg, #667eea, #764ba2);
        color: white;
        padding: 8px 16px;
        border-radius: 50px;
        font-size: 14px;
        font-weight: 600;
        margin-bottom: 24px;
      }
      .main-title {
        font-size: 32px;
        font-weight: 800;
        color: #1a1a1a;
        margin-bottom: 12px;
        line-height: 1.2;
      }
      .subtitle {
        font-size: 18px;
        color: #6b7280;
        font-weight: 400;
        line-height: 1.5;
      }
      .content-section {
        padding: 0 40px 48px;
      }
      .greeting-card {
        background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
        border: 1px solid #e2e8f0;
        border-radius: 20px;
        padding: 32px;
        margin-bottom: 32px;
        text-align: center;
      }
      .greeting-text {
        font-size: 20px;
        font-weight: 700;
        color: #1e293b;
        margin-bottom: 12px;
      }
      .greeting-subtext {
        font-size: 16px;
        color: #64748b;
        line-height: 1.6;
      }
      .cta-container {
        text-align: center;
        margin: 40px 0;
      }
      .primary-button {
        display: inline-flex;
        align-items: center;
        gap: 12px;
        background: linear-gradient(135deg, #667eea, #764ba2);
        color: white;
        text-decoration: none;
        padding: 18px 36px;
        border-radius: 16px;
        font-weight: 700;
        font-size: 16px;
        transition: all 0.3s ease;
        box-shadow: 0 8px 32px rgba(102, 126, 234, 0.3);
      }
      .primary-button:hover {
        transform: translateY(-2px);
        box-shadow: 0 12px 40px rgba(102, 126, 234, 0.4);
      }
      .features-grid {
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        gap: 20px;
        margin: 40px 0;
      }
      .feature-item {
        background: #ffffff;
        border: 2px solid #f1f5f9;
        border-radius: 16px;
        padding: 24px;
        text-align: center;
        transition: all 0.3s ease;
      }
      .feature-item:hover {
        border-color: #667eea;
        transform: translateY(-4px);
      }
      .feature-emoji {
        font-size: 32px;
        margin-bottom: 16px;
        display: block;
      }
      .feature-title {
        font-size: 16px;
        font-weight: 700;
        color: #1e293b;
        margin-bottom: 8px;
      }
      .feature-desc {
        font-size: 14px;
        color: #64748b;
        line-height: 1.5;
      }
      .footer-section {
        background: #f8fafc;
        padding: 32px 40px;
        text-align: center;
        border-top: 1px solid #e2e8f0;
      }
      .footer-text {
        color: #64748b;
        font-size: 14px;
        margin-bottom: 20px;
      }
      .footer-links {
        display: flex;
        justify-content: center;
        gap: 32px;
        flex-wrap: wrap;
      }
      .footer-link {
        color: #667eea;
        text-decoration: none;
        font-weight: 600;
        font-size: 14px;
        transition: color 0.2s ease;
      }
      .footer-link:hover {
        color: #764ba2;
      }
      @media (max-width: 640px) {
        .email-wrapper {
          margin: 0;
          border-radius: 0;
        }
        .header-section,
        .content-section {
          padding-left: 24px;
          padding-right: 24px;
        }
        .features-grid {
          grid-template-columns: 1fr;
          gap: 16px;
        }
        .footer-links {
          flex-direction: column;
          gap: 16px;
        }
      }
    </style>
  </head>
  <body>
    <div class="email-wrapper">
      <div class="header-section">
        <div class="brand-logo">WK</div>
        <div class="welcome-badge">
          <span>🎉</span>
          Welcome aboard!
        </div>
        <h1 class="main-title">You're all set, {{ .Name }}!</h1>
        <p class="subtitle">Your workspace-kit journey begins now</p>
      </div>

      <div class="content-section">
        <div class="greeting-card">
          <div class="greeting-text">
            Hey {{ .Name }}, welcome to the team! 👋
          </div>
          <div class="greeting-subtext">
            We're excited to have you join thousands of productive teams using
            workspace-kit to streamline their workflow and boost collaboration.
          </div>
        </div>

        <div class="cta-container">
          <a href="{{ .SiteURL }}" class="primary-button">
            <span>🚀</span>
            Start Building
          </a>
        </div>

        <div class="features-grid">
          <div class="feature-item">
            <span class="feature-emoji">⚡</span>
            <div class="feature-title">Lightning Fast</div>
            <div class="feature-desc">
              Get up and running in minutes with our intuitive setup process
            </div>
          </div>
          <div class="feature-item">
            <span class="feature-emoji">🤝</span>
            <div class="feature-title">Team Sync</div>
            <div class="feature-desc">
              Real-time collaboration tools that keep everyone aligned
            </div>
          </div>
          <div class="feature-item">
            <span class="feature-emoji">📊</span>
            <div class="feature-title">Smart Analytics</div>
            <div class="feature-desc">
              Insights and metrics to optimize your team's performance
            </div>
          </div>
          <div class="feature-item">
            <span class="feature-emoji">🔒</span>
            <div class="feature-title">Enterprise Security</div>
            <div class="feature-desc">
              Bank-level security with end-to-end encryption
            </div>
          </div>
        </div>

        <div
          style="
            text-align: center;
            margin-top: 40px;
            padding: 24px;
            background: #f8fafc;
            border-radius: 16px;
          "
        >
          <p style="color: #64748b; font-size: 16px; margin-bottom: 8px">
            Questions? We're here to help!
          </p>
          <p style="color: #1e293b; font-weight: 600">The workspace-kit Team</p>
        </div>
      </div>

      <!--    <div class="footer-section">-->
      <!--        <p class="footer-text">&copy; 2024 workspace-kit. All rights reserved.</p>-->
      <!--        <div class="footer-links">-->
      <!--            <a href="{{ .SiteURL }}/privacy" class="footer-link">Privacy</a>-->
      <!--            <a href="{{ .SiteURL }}/terms" class="footer-link">Terms</a>-->
      <!--            <a href="{{ .SiteURL }}/contact" class="footer-link">Support</a>-->
      <!--        </div>-->
      <!--    </div>-->
    </div>
  </body>
</html>
```

:::
