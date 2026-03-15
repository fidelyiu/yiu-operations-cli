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

## 启动数据库

```sh
docker run --name wk-postgres -p 5432:5432 -e POSTGRES_PASSWORD=password -e POSTGRES_USER=postgres -e POSTGRES_DB=workspace-kit -d postgres
```

这里的 `-p 5432:5432` 表示：

- 前面的 `5432` 是宿主机端口
- 后面的 `5432` 是容器内 PostgreSQL 的端口

这样本机就可以通过 `localhost:5432` 访问容器里的 PostgreSQL。

`docker-compose.yaml`

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
