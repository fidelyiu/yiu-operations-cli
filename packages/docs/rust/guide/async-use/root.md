# async 使用

trpl 是 “The Rust Programming Language” 的缩写。它重新导出了本章需要的所有类型、trait 和函数，主要来自 [futures](https://crates.io/crates/futures) 和 [tokio](https://tokio.rs/) crate。futures crate 是 Rust 异步代码实验的官方阵地，Future trait 最初就是在那里设计出来的。Tokio 则是目前 Rust 中使用最广泛的异步运行时（async runtime），尤其常见于 Web 应用。生态中也还有其他很优秀的运行时，而且它们可能更适合你的实际用途。我们在 trpl 的底层使用 tokio，是因为它经过了充分测试，也足够常用。

## 相关链接

- [异步编程基础：Async、Await、Future 与 Stream](https://kaisery.github.io/trpl-zh-cn/ch17-00-async-await.html)

## 概念

- `future` 是一个现在也许还没准备好，但会在将来某个时刻准备好的值。
- `Future` trait 作为基础构件，让不同的异步操作可以用不同的数据结构来实现，同时又拥有统一的接口。
- `async` 关键字可以用于代码块和函数，表示它们可以被中断和恢复。
- 在 `async` 块或 `async` 函数中，你可以使用 `await` 关键字来 `await` 一个 `future`，也就是等待它变为就绪。

## 异步编译

- 当 Rust 遇到一个 `async` 关键字标记的**代码块**时，会将其编译为一个实现了 Future trait 的唯一的、匿名的数据类型。
- 当 Rust 遇到一个被标记为 `async` 的**函数**时，会将其编译成一个函数体是异步代码块的非异步函数。异步函数的返回值类型是编译器为异步代码块所创建的匿名数据类型。

```rust
use trpl::Html;

async fn page_title(url: &str) -> Option<String> {
    let response = trpl::get(url).await;
    let response_text = response.text().await;
    Html::parse(&response_text)
        .select_first("title")
        .map(|title| title.inner_html())
}
```

编译之后

```rust
use std::future::Future;
use trpl::Html;

fn page_title(url: &str) -> impl Future<Output = Option<String>> {
    async move {
        let text = trpl::get(url).await.text().await;
        Html::parse(&text)
            .select_first("title")
            .map(|title| title.inner_html())
    }
}
```

因为会改变方法签名，并且维护 await 点，所以不能在 main 函数上添加 async 关键字。
