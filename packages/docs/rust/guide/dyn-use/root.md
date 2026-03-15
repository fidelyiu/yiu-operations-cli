# dyn使用

`dyn Trait` 用来表示 trait object，也就是“某个实现了这个 trait 的具体类型”，但在当前上下文里不关心它到底是哪一个具体类型。

它的核心目的，是让一份代码可以在运行时处理多种实现同一 trait 的类型。

<!-- truncate -->

## dyn 是什么

Rust 默认更偏向静态分发，也就是编译期就知道具体类型，然后把调用展开到对应实现上。

`dyn Trait` 则是动态分发：

- 编译期只知道它实现了某个 trait
- 运行时通过 vtable 找到真正要调用的方法
- 因为具体类型大小未知，所以 `dyn Trait` 是 DST，必须放在指针后面使用

常见形式：

- `&dyn Trait`
- `Box<dyn Trait>`
- `Rc<dyn Trait>`
- `Arc<dyn Trait + Send + Sync>`

## 基本例子

```rust
trait Draw {
	fn draw(&self);
}

struct Button;
struct Input;

impl Draw for Button {
	fn draw(&self) {
		println!("draw button");
	}
}

impl Draw for Input {
	fn draw(&self) {
		println!("draw input");
	}
}

fn render(widget: &dyn Draw) {
	widget.draw();
}

fn main() {
	let button = Button;
	let input = Input;

	render(&button);
	render(&input);
}
```

这里的 `render` 不需要知道传入的是 `Button` 还是 `Input`，只需要它实现了 `Draw`。

## 为什么要用 dyn

最常见的场景是“同一个集合里放多种不同类型，但它们都实现了同一个 trait”。

```rust
trait Draw {
	fn draw(&self);
}

struct Button;
struct SelectBox;

impl Draw for Button {
	fn draw(&self) {
		println!("button");
	}
}

impl Draw for SelectBox {
	fn draw(&self) {
		println!("select box");
	}
}

fn main() {
	let screens: Vec<Box<dyn Draw>> = vec![
		Box::new(Button),
		Box::new(SelectBox),
	];

	for screen in screens {
		screen.draw();
	}
}
```

如果不用 `dyn`，`Vec` 要求所有元素都是同一个具体类型，这里就做不到同时放 `Button` 和 `SelectBox`。

## 常见写法

### 作为参数

```rust
fn run(task: &dyn Job) {
	task.execute();
}
```

如果只是临时借用，通常用 `&dyn Trait` 就够了。

### 作为返回值

```rust
trait Animal {
	fn speak(&self);
}

struct Dog;

impl Animal for Dog {
	fn speak(&self) {
		println!("wang");
	}
}

fn create_animal() -> Box<dyn Animal> {
	Box::new(Dog)
}
```

返回 `dyn Trait` 时，通常需要包一层 `Box`、`Rc`、`Arc` 之类的指针，因为 `dyn Trait` 本身大小未知。

### 和自动 trait 一起使用

```rust
use std::sync::Arc;

trait Job {
	fn execute(&self);
}

fn spawn_job(job: Arc<dyn Job + Send + Sync>) {
	// 可以在线程间安全共享
}
```

这里的 `Send + Sync` 是额外约束，表示这个 trait object 对并发环境也是安全的。

### 带生命周期

```rust
fn pick_printer<'a>(value: &'a str) -> Box<dyn std::fmt::Display + 'a> {
	Box::new(value)
}
```

如果 trait object 内部借用了外部数据，就需要显式标出生命周期。

## dyn 和泛型的区别

### 泛型 / `impl Trait`

- 静态分发
- 编译器知道具体类型
- 性能通常更好，容易被内联
- 但一个具体实例里只能对应一个确定类型

```rust
fn render<T: Draw>(widget: &T) {
	widget.draw();
}
```

### `dyn Trait`

- 动态分发
- 运行时通过 vtable 调用
- 能把不同具体类型统一放在一起处理
- 会有一点间接调用开销

经验上：

- 如果你关心性能，且类型在编译期已知，优先用泛型
- 如果你需要异构集合、运行时替换实现、降低 API 对具体类型的暴露，可以用 `dyn`

## object safety

不是所有 trait 都能写成 `dyn Trait`。一个 trait 想变成 trait object，通常要满足 object safety。

下面这种就不行：

```rust
trait CloneLike {
	fn clone_me(&self) -> Self;
}
```

原因是 `Self` 代表具体类型，而 `dyn CloneLike` 已经抹掉了具体类型，运行时无法知道返回值到底多大。

还有这类泛型方法也不能直接放进 trait object：

```rust
trait Mapper {
	fn map<T>(&self, value: T);
}
```

因为泛型参数要求编译期单态化，和 trait object 的动态分发模型冲突。

实际开发里常见处理方式：

- 避免在 trait object 里返回 `Self`
- 避免在 trait object 里定义泛型方法
- 如果某个方法只想给具体类型使用，可以加 `where Self: Sized`

```rust
trait Builder {
	fn build(&self);

	fn boxed(self) -> Box<Self>
	where
		Self: Sized,
	{
		Box::new(self)
	}
}
```

这样 `build` 仍然可以给 `dyn Builder` 用，而 `boxed` 只能给具体类型用。

## dyn 为什么总是跟指针一起出现

因为 `dyn Trait` 的大小在编译期不知道：

```rust
trait Draw {
	fn draw(&self);
}

fn bad(arg: dyn Draw) {}
```

上面这种写法不成立。必须改成：

```rust
fn ok(arg: &dyn Draw) {}
fn ok2(arg: Box<dyn Draw>) {}
```

本质上是因为 Rust 需要知道栈上值的大小，而 `dyn Trait` 只有在和胖指针一起出现时，才能同时携带数据地址和 vtable 信息。

## 什么时候不该用 dyn

- 只有一种具体实现类型
- 类型在编译期完全已知
- 对性能和内联比较敏感
- 可以直接用泛型表达得更清楚

很多时候 `dyn` 不是“更高级”，只是适合“运行时抽象”这个场景。

## 相关链接

- [Rust trait 使用](../trait-use/root.md)
- [async 和 dyn 的使用](../async-dyn/root.md)
- [使用特性对象抽象共享行为](https://doc.rust-lang.org/book/ch18-02-trait-objects.html)
- [traits 动态兼容性](https://doc.rust-lang.org/reference/items/traits.html#dyn-compatibility)
