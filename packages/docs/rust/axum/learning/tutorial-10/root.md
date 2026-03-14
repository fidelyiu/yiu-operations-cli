# Axum教程-10-高级功能

这些功能将你的 API 从请求-响应模型转变为一个实时运行的系统。

## 创建模块

```sh
cargo new module-10-advanced
```

## 添加依赖

`module-10-advanced/Cargo.toml`

```toml
[package]
name = "module-10-advanced"
version = "0.1.0"
edition = "2024"

[dependencies]
axum = { workspace = true, features = ["ws", "multipart"] }
axum-extra = { workspace = true }
tokio = { workspace = true }
tokio-stream = "0.1"
serde = { workspace = true }
serde_json = { workspace = true }
tower-http = { workspace = true }
futures = { workspace = true }
```

## WebSocket

WebSocket 是双向通信通道。

与 HTTP 的请求-响应模式不同，WebSocket 会保持连接打开。

客户端和服务器都可以随时发送消息。

http 就像寄信。你写好、寄出，然后等待回复。

websocket 就像打电话。一旦连接建立，双方可以随时讲话。

聊天应用、实时更新、多人游戏、协同编辑，凡是需要即时通信的地方都是 websocket 的领域。

### 处理器

它不同于常规处理器。

该处理器使用 websocket 升级提取器。这是特殊的。

它告诉 axum 这个端点处理 websocket 升级。

我们调用 ws upgrade 并传入一个处理该套接字的闭包。

这会将 HTTP 连接升级为 WebSocket 协议。

该闭包在一个单独的 Tokio 任务中运行，负责处理消息。升级是自动完成的。

客户端发送升级请求。axom 会进行协议切换。你的闭包在连接的整个生命周期内运行。

- 文本消息包含字符串、聊天消息、JSON 命令等。
- 二进制消息是原始字节、文件片段、图像、音频流。
  - 根据你的使用场景按需处理它们。
- ping 消息是连接健康检查。
  - axom 会自动处理 pong 响应。
- close 消息会干净地结束连接。

真正的应用会解析消息、更新状态并广播给其他连接。

对于聊天应用，你需要将消息广播给所有已连接的客户端。

这种模式将所有连接存储在共享状态中，使用 Arc 包裹的 Vec，里面是 socket 发送者。

当一条消息到达时，遍历并发送给每一个接收者。

使用 tokio broadcast 或 tokio watch 通道来实现高效的广播分发。

这些通道是为一对多通信而设计的。

我们并没有在模块中实现这一点，但这种模式很直观。一旦你理解了单个（接收者/处理）机制连接处理。

```rust
async fn ws_handler(ws: WebSocketUpgrade) -> impl IntoResponse {
    ws.on_upgrade(handle_socket)
}

async fn handle_socket(mut socket: WebSocket) {
    // 我们用循环。这会在消息到达时读取它们。
    while let Some(msg) = socket.recv().await {
        // 对于每条消息，我们都会根据类型进行匹配。
        if let Ok(Message::Text(text)) = msg {
            // 我们用一个前缀将它们回显回去。
            let response = format!("回显：{}", text);
            if socket.send(Message::Text(response.into())).await.is_err() {
                break;
            }
        }
    }
}
```

### 路由

```rust
let app = Router::new().route("/ws", get(ws_handler));
```

## 服务端发送事件

SSE 比 WebSocket 更简单，提供服务器到客户端的单向流传输。

服务器推送，客户端接收。

非常适合实时信息流、通知、进度更新、股票行情。

适用于任何服务器需要推送数据而客户端只需监听的场景。

### 处理器

```rust
// 返回一个流的 Sse。每个Item都会成为发送给客户端的事件。
async fn sse_handler() -> Sse<impl Stream<Item = Result<Event, Infallible>>> {
    // 该流使用 tokio 作为流的 repeat_until_throttle（注：原文可能指 repeat 或 throttle），用于以微秒为单位重复。
    // 带闭包的 repeat_until_throttle 会无限运行。
    let stream = stream::repeat_with(|| {
        // 每个事件都有一个事件类型、数据和可选的 id。
        // 客户端可以按事件类型筛选。
        Event::default().data(format!("服务器时间：{:?}", std::time::SystemTime::now()))
    })
    .map(Ok)
    // 我们每秒生成一个带有当前时间戳的事件。
    // throttle 在项目之间添加延迟。
    .throttle(Duration::from_secs(1));

    // 客户端保持连接打开并实时接收更新。浏览器支持非常好。
    // 事件源 API 非常简单，只需指向你的端点。
    // 保持连接存活至关重要。它能防止代理和负载均衡器导致连接超时。
    // 定期发送空消息以保持连接活跃。
    Sse::new(stream).keep_alive(KeepAlive::default())
}
```

### 路由

```rust
let app = Router::new().route("/sse", get(sse_handler));
```

## 文件上传

这就是用户将文件发送到你的服务器的方式。

该处理器使用 multipart 提取器。multipart 是表单上传的标准。

每个文件都是请求中的一个部分。

使用 tokio fs write 将其保存到磁盘。生成唯一的文件名以避免冲突。

为安全起见验证内容类型。

完整的生产环境。配置请求体限制。使用默认的请求体限制层来设置最大请求大小。

对于大多数应用来说，10 MB 是合理的。巨大的上传可能会耗尽内存或填满磁盘。根据你的使用场景设置限制。

验证文件类型。不要相信内容类型头。攻击者可以伪造它们。检查文件magic bytes或扫描文件。

将上传文件存储在你的网页路由之外。

永远不要直接提供用户上传的文件。它们可能包含脚本。

生成随机文件名。不要使用用户提供的名称。路径遍历攻击是真实存在的。

### 处理器

```rust
async fn upload(mut multipart: Multipart) -> impl IntoResponse {
    let mut files = Vec::new();

    // 我们用 while let 循环遍历字段，直到 field
    while let Some(field) = multipart.next_field().await.unwrap() {
        // 每次迭代都会给我们一个上传的文件。
        // 每个字段都有元数据。
        // 来自表单的 name、来自用户电脑的文件名、像 image/png 这样的内容类型。
        let name = field.name().unwrap_or("未知").to_string();
        // 调用 bytes 来读取文件数据。
        // 对于大文件，使用流式处理以避免将所有内容加载到内存中。
        let data = field.bytes().await.unwrap();
        files.push(format!("{}：{} 字节", name, data.len()));
    }

    if files.is_empty() {
        "没有上传任何文件".to_string()
    } else {
        format!("已上传：{}", files.join(", "))
    }
}
```

### 路由

```rust
let app = Router::new().route("/upload", post(upload))
```

## HTML

该页面提供会打开 WebSocket 的静态 HTML。

在接收文件上传的同时通过 SSE 推送更新。

Axum 干净利落地处理这一切。

每个功能只是合并到你的应用中的另一个路由或路由器。

相同的模式适用。

### 处理器

```rust
async fn demo_page() -> Html<&'static str> {
    Html(
        r#"
<!DOCTYPE html>
<html>
<head>
    <title>Axum 高级特性</title>
    <style>
        body { font-family: system-ui; max-width: 800px; margin: 50px auto; padding: 20px; }
        .demo { background: #f5f5f5; padding: 20px; margin: 20px 0; border-radius: 8px; }
        button { padding: 10px 20px; margin: 5px; cursor: pointer; }
        #ws-output, #sse-output { height: 100px; overflow-y: auto; background: #fff;
                                   border: 1px solid #ddd; padding: 10px; margin-top: 10px; }
    </style>
</head>
<body>
    <h1>🚀 Axum 高级特性</h1>

    <div class="demo">
        <h2>WebSocket 回显</h2>
        <input type="text" id="ws-input" placeholder="输入消息">
        <button onclick="sendWs()">发送</button>
        <div id="ws-output"></div>
    </div>

    <div class="demo">
        <h2>服务端发送事件</h2>
        <button onclick="startSse()">启动 SSE</button>
        <button onclick="stopSse()">停止</button>
        <div id="sse-output"></div>
    </div>

    <div class="demo">
        <h2>文件上传</h2>
        <form action="/upload" method="post" enctype="multipart/form-data">
            <input type="file" name="file" multiple>
            <button type="submit">上传</button>
        </form>
    </div>

    <script>
        let ws, sse;

        ws = new WebSocket('ws://localhost:3000/ws');
        ws.onmessage = (e) => {
            document.getElementById('ws-output').innerHTML += e.data + '<br>';
        };

        function sendWs() {
            ws.send(document.getElementById('ws-input').value);
        }

        function startSse() {
            sse = new EventSource('/sse');
            sse.onmessage = (e) => {
                document.getElementById('sse-output').innerHTML = e.data;
            };
        }

        function stopSse() { if(sse) sse.close(); }
    </script>
</body>
</html>
"#,
    )
}
```

### 路由

```rust
let app = Router::new().route("/", get(demo_page))
```

## 静态文件

axum 处理正确的内容类型、缓存头、范围请求等，所有位于像 /static 这样的路径下的内容均能被处理。

现在 `/static/styles.css` 会在生产环境中提供你的 CSS 文件。

在前面放一个 CDN，比如 Cloudflare、CloudFront 或类似服务。

它们在全局缓存文件以实现快速传送。真正的应用结合了这些功能。

```rust
#[tokio::main]
async fn main() {
    // 如果需要则创建静态目录
    std::fs::create_dir_all("static").ok();
    std::fs::write("static/hello.txt", "来自静态文件的问候！").ok();

    let app = Router::new().nest_service("/static", ServeDir::new("static"));
}
```
