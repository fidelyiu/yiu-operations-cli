# 所有权

## 相关链接

- [所有权 中文文档](https://kaisery.github.io/trpl-zh-cn/ch04-00-understanding-ownership.html)

## 所有权规则

1. Rust 中的每一个值都有一个 所有者（owner）。
2. 值在任一时刻有且只有一个所有者。
3. 当所有者离开作用域，这个值将被丢弃。

## 所有权移动

### 堆栈对比

移动只会发生在堆内存变量上。用于避免深拷贝、悬垂指针，指针多引用问题。

```rust
{
    let x = 5;
    let y = x;
    // ✅ 栈内存不会移动所有权。x, y 都有效
    println!("x = {x}, y = {y}");
}
{
    let s1 = String::from("hello");
    let s2 = s1;
    // ❌ 堆内存会移动所有权。s1 不再有效, 只有 s2 有效。
    println!("{s1}, world!");
}
{
    let s1 = String::from("hello");
    let s2 = s1.clone();
    // ✅ clone 主动深拷贝，堆上两个数据。s1, s2 都有效
    println!("s1 = {s1}, s2 = {s2}");
}
```

### 赋值移动

```rust
let s1 = String::from("hello");
let s2 = s1;
// s1 不再可用, s1 值的所有权移动到 s2
```

### 函数移动

```rust
fn main() {
    let s = String::from("hello");
    takes_ownership(s);
    // s 不在可用, s 值的所有权移动到 takes_ownership 参数中

    let x = 5;
    makes_copy(x);
    // x 有效, 栈内存不移动
    println!("{}", x);
}

fn takes_ownership(some_string: String) {
    println!("{some_string}");
}

fn makes_copy(some_integer: i32) {
    println!("{some_integer}");
}
```

## 引用

使用 `&` 创建变量的引用变量没有 原值的所有权，而是指向原值的指针地址的所有权。

引用变量自动解引用，所以可以直接调用原来数据类型的方法。

而且引用是栈内存变量。

```rust
fn main() {
    let s1 = String::from("hello");
    let len = calculate_length(&s1);
    println!("The length of '{s1}' is {len}.");
}

fn calculate_length(s: &String) -> usize {
    s.len()
}
```

### 可变引用

- 在任意给定时间，要么只能有一个可变引用，要么只能有多个不可变引用。
  - 可以使用大括号创建作用域解决一些问题。
- 引用必须总是有效的。

```rust
fn main() {
    let mut s = String::from("hello");
    change(&mut s);
}

fn change(some_string: &mut String) {
    some_string.push_str(", world");
}
```

### slice

切片（slice）允许你引用集合中一段连续的元素序列，而不用引用整个集合。slice 是一种引用，所以它不拥有所有权。

```rust
fn main() {
    let mut s = String::from("hello world");
    let word = first_word(&s); // first_word 创建了一个 s 的不可变引用
    s.clear(); // ❌，s 的 clear 内部会有 s 的可变引用。变量不能同时拥有可变引用和不可变引用。可以使用大括号让word作用域结束解决。
    println!("the first word is: {word}");
}

fn first_word(s: &str) -> &str {
    let bytes = s.as_bytes();
    for (i, &item) in bytes.iter().enumerate() {
        if item == b' ' {
            return &s[0..i];
        }
    }
    &s[..]
}
```

> - `String`: 字符串类型
> - `&String`: 字符串引用
> - `&str`: 字符串 slice类型
>
> `let s = "Hello, world!";` 这里 `s` 的类型是 `&str`：它是一个指向二进制程序特定位置的 `slice`。这也就是为什么字符串字面值是不可变的；`&str` 是一个不可变引用。

## 所有权-实践方案

### 返回值与作用域

函数可以将值的所有权，再次返回给调用方。

```rust
fn main() {
    let s1 = gives_ownership();
    let s2 = String::from("hello");
    let s3 = takes_and_gives_back(s2);
}

fn gives_ownership() -> String {
    let some_string = String::from("yours");
    some_string // 返回 some_string 并将其移至调用函数
}

// 该函数将传入字符串并返回该值
fn takes_and_gives_back(a_string: String) -> String {
    a_string  // 返回 a_string 并移出给调用的函数
}
```

### 使用引用

```rust
fn main() {
    let s1 = String::from("hello");
    let len = calculate_length(&s1);
    println!("The length of '{s1}' is {len}.");
}

fn calculate_length(s: &String) -> usize {
    s.len()
}
```
