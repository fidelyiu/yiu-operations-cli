# 智能指针

智能指针带来了很多更高级而复杂的规则。

- `Box<T>`，用于在堆上分配值
  - 因为rust要在编译时知道一个结构体的栈内存大小。对于树结构的嵌套类型。就需要手动在堆内存上分配大小。
- `Rc<T>`，一个引用计数类型，其数据可以有多个所有者
- `Ref<T>` 和 `RefMut<T>`，它们是通过 `RefCell<T>` 访问的。而 `RefCell<T>` 是一个在运行时而非编译时执行借用规则的类型。

## 相关链接

- [智能指针 中文文档](https://kaisery.github.io/trpl-zh-cn/ch15-00-smart-pointers.html)

## Box

使用场景

- 嵌套结构体。
  - 当有一个在编译时未知大小的类型，而又想要在需要确切大小的上下文中使用这个类型值的时候。
- 大量栈数据。
  - 当有大量数据并希望在确保数据不被拷贝的情况下转移所有权的时候。
  - 转移大量数据的所有权可能会花费很长时间，因为数据会在栈上被复制。
- 当希望拥有一个值并只关心它的类型是否实现了特定 trait 而不是其具体类型的时候

```rust
enum List {
    Cons(i32, Box<List>),
    Nil,
}

use crate::List::{Cons, Nil};

fn main() {
    let list = Cons(1, Box::new(Cons(2, Box::new(Cons(3, Box::new(Nil))))));
}
```

## Rc

Box 还是符合所有权规则，值只能有一个所有者。

```rust
enum List {
    Cons(i32, Box<List>),
    Nil,
}

use crate::List::{Cons, Nil};

fn main() {
    let a = Cons(5, Box::new(Cons(10, Box::new(Nil))));
    let b = Cons(3, Box::new(a));
    let c = Cons(4, Box::new(a)); // ❌ a 不能有两个所有者
}
```

Rc运行值有多个所有者。

```rust
enum List {
    Cons(i32, Rc<List>),
    Nil,
}

use crate::List::{Cons, Nil};
use std::rc::Rc;

fn main() {
    let a = Rc::new(Cons(5, Rc::new(Cons(10, Rc::new(Nil)))));
    println!("创建 a 后, 计数 = {}", Rc::strong_count(&a)); // 1
    let b = Cons(3, Rc::clone(&a));
    println!("创建 b 后, 计数 = {}", Rc::strong_count(&a)); // 2
    {
        let c = Cons(4, Rc::clone(&a));
        println!("创建 c 后, 计数 = {}", Rc::strong_count(&a)); // 3
    }
    println!("释放 c 后, 计数 = {}", Rc::strong_count(&a)); // 2
}
```
