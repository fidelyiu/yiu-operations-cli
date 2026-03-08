---
sidebar_position: 5
---

# 部署你的站点

Docusaurus 是一个**静态站点生成器**（也称为 **[Jamstack](https://jamstack.org/)**）。

它将你的站点构建为简单的**静态 HTML、JavaScript 和 CSS 文件**。

## 构建你的站点

为**生产环境**构建站点：

```bash
npm run build
```

静态文件会生成在 `build` 文件夹中。

## 部署你的站点

在本地测试生产构建：

```bash
npm run serve
```

`build` 文件夹现在在 [http://localhost:3000/](http://localhost:3000/) 提供服务。

现在你可以**几乎在任何地方**轻松部署 `build` 文件夹，**免费**或成本非常低（阅读**[部署指南](https://docusaurus.io/docs/deployment)**）。
