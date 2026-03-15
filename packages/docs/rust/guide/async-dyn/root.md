# async 和 dyn 的使用

Rust 1.75 中对 trait 中异步函数的稳定化并没有包含对将包含异步函数的 trait 用作 dyn Trait 的支持。尝试将 dyn 与异步 trait 一起使用会产生以下错误：

```rust
pub trait Trait {
    async fn f(&self);
}

pub fn make() -> Box<dyn Trait> {
    unimplemented!()
}
```

## 相关链接

- [async-trait](https://docs.rs/async-trait/latest/async_trait/)
- [为什么在 traits 中使用异步函数很难](https://smallcultfollowing.com/babysteps/blog/2019/10/26/async-fn-in-traits-are-hard/)

## 使用 async-trait

```sh
cargo add async-trait
```
