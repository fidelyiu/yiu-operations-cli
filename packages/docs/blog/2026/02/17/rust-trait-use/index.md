---
title: Rust的trait使用
authors: [FidelYiu]
tags: [rust]
---

# Rust的trait使用

- 普通的 tarit
- tarit 默认实现
- trait 作为返回值
- trait 限制结构体实现
- 扩展现有类型的方法
- trait 范型 和 关联类型

<!-- truncate -->

## 普通的 tarit

```rust
pub trait Summary {
    fn summarize(&self) -> String;
}
```

类型可以直接去实现它。

```rust
pub struct NewsArticle {
    pub headline: String,
    pub location: String,
    pub author: String,
    pub content: String,
}

impl Summary for NewsArticle {
    fn summarize(&self) -> String {
        format!("{}, by {} ({})", self.headline, self.author, self.location)
    }
}
```

## tarit 默认实现

```rust
pub trait Summary {
    fn summarize_author(&self) -> String;

    fn summarize(&self) -> String {
        format!("(Read more from {}...)", self.summarize_author())
    }
}
```

## trait 作为参数

```rust
pub fn notify(item: &impl Summary) {
    println!("Breaking news! {}", item.summarize());
}
pub fn notifyMultiple(item: &(impl Summary + Display)) {}
```

实际语法

```rust
pub fn notify<T: Summary>(item: &T) {
    println!("Breaking news! {}", item.summarize());
}
```

where语法

```rust
fn some_function<T, U>(t: &T, u: &U) -> i32
where
    T: Display + Clone,
    U: Clone + Debug,
{}
```

## trait 作为返回值

因为rust的原理是通过编译生成不同的方法。所以不能在方法体中返回多种类型的不同trait。

```rust
fn returns_summarizable() -> impl Summary {
    SocialPost {
        username: String::from("horse_ebooks"),
        content: String::from(
            "of course, as you probably already know, people",
        ),
        reply: false,
        repost: false,
    }
}
```

## trait 限制结构体实现

只有当 `Pair` 的范型 `T` 传入的类型实现了 `Display + PartialOrd`，这时 `Pair` 的值才会有 `cmp_display` 方法

```rust
use std::fmt::Display;

struct Pair<T> {
    x: T,
    y: T,
}

impl<T> Pair<T> {
    fn new(x: T, y: T) -> Self {
        Self { x, y }
    }
}

impl<T: Display + PartialOrd> Pair<T> {
    fn cmp_display(&self) {
        if self.x >= self.y {
            println!("The largest member is x = {}", self.x);
        } else {
            println!("The largest member is y = {}", self.y);
        }
    }
}
```

## 扩展现有类型的方法

给所有实现了 `Display` 的类型默认实现上 `ToString` trait

```rust
impl<T: Display> ToString for T {
    // --snip--
}
```

## trait 范型 和 关联类型

### trait 范型

```rust
trait Convert<T> {
    fn convert(&self) -> T;
}

struct Wrapper(String);

impl Convert<usize> for Wrapper {
    fn convert(&self) -> usize {
        self.0.len()
    }
}

impl Convert<bool> for Wrapper {
    fn convert(&self) -> bool {
        !self.0.is_empty()
    }
}
```

当 `Wrapper` 同时实现了多个 `Convert<T>` 时，调用 `convert()` 必须有类型上下文，否则会歧义。

```rust
let w = Wrapper("hi".to_string());

let len: usize = w.convert(); // 使用 Convert<usize>
let ok: bool = w.convert();    // 使用 Convert<bool>

let len = <Wrapper as Convert<usize>>::convert(&w);
let ok = <Wrapper as Convert<bool>>::convert(&w);
```

#### Axum 中的 FromRequestParts

`FromRequestParts` 使用关联类型 `Rejection` 表示提取失败的错误类型。

- `S`：应用的共享状态类型（state），由你在 Router 上设置，`FromRequestParts<S>` 的实现可以读取该状态。
- `Rejection`：提取失败时返回的错误类型，用于把失败转换成响应（如 `StatusCode`、自定义错误）。

```rust
use axum::{
    async_trait,
    extract::FromRequestParts,
    extract::State,
    http::{request::Parts, StatusCode},
    response::IntoResponse,
    routing::get,
    Router,
};
use std::sync::Arc;

struct AuthUser {
    user_id: String,
}

#[derive(Clone)]
struct AppState {
    app_name: Arc<String>,
}

#[async_trait]
impl<S> FromRequestParts<S> for AuthUser
where
    S: Send + Sync,
{
    type Rejection = StatusCode;

    async fn from_request_parts(parts: &mut Parts, _state: &S) -> Result<Self, Self::Rejection> {
        let user_id = parts
            .headers
            .get("x-user-id")
            .and_then(|value| value.to_str().ok())
            .ok_or(StatusCode::UNAUTHORIZED)?;

        Ok(Self {
            user_id: user_id.to_string(),
        })
    }
}

async fn profile(State(state): State<AppState>, user: AuthUser) -> impl IntoResponse {
    format!("{}: {}", state.app_name, user.user_id)
}

fn app() -> Router {
    let state = AppState {
        app_name: Arc::new("demo".to_string()),
    };

    Router::new()
        .route("/profile", get(profile))
        .with_state(state)
}
```

### 关联类型

```rust
trait Container {
    type Item;
    fn add(&mut self, item: Self::Item);
    fn len(&self) -> usize;
}

struct Bag {
    items: Vec<String>,
}

impl Container for Bag {
    type Item = String;

    fn add(&mut self, item: Self::Item) {
        self.items.push(item);
    }

    fn len(&self) -> usize {
        self.items.len()
    }
}
```

### 实现方法上的范型

```rust
struct Cache;

impl Cache {
    fn save<T: ToString>(&self, value: T) -> String {
        value.to_string()
    }
}
```

### 总结

- trait 范型：`trait Convert<T>` 会为每个 `T` 生成一组实现，同一个类型可以实现多次（如 `Wrapper: Convert<usize>` 和 `Wrapper: Convert<bool>`）。
- 关联类型：`trait Container` 里 `type Item` 由实现者确定，一个实现只能绑定一个 `Item`，调用方更简洁（`Self::Item`）。
- 方法范型：把范型放在具体方法上，适合单个方法需要泛化而不是整个 trait 或类型都泛化。

---

- trait 范型: 更多的是给 这个 trait 添加一个未知类型，从而给实现方法中提供统一的类型。
- 关联类型: 更多的是给 实现者的 对应未知类型。一个类型的一个关联类型只能绑定一个类型。
- 方法范型: 这个只是给方法添加了范型，如果 trait 中有多个需要统一的未知类型就不好用了。
