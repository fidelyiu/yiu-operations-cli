# 网站

本网站使用 [Docusaurus](https://docusaurus.io/) 构建，这是一个现代化的静态网站生成器。

## 安装

```bash
yarn
```

## 本地开发

```bash
yarn start
```

此命令会启动本地开发服务器并打开浏览器窗口。大多数更改都会实时反映，无需重启服务器。

## 构建

```bash
yarn build
```

此命令会将静态内容生成到 `build` 目录中，可以使用任何静态内容托管服务进行部署。

## 部署

使用 SSH：

```bash
USE_SSH=true yarn deploy
```

不使用 SSH：

```bash
GIT_USER=<你的 GitHub 用户名> yarn deploy
```

如果你使用 GitHub Pages 进行托管，此命令可以方便地构建网站并推送到 `gh-pages` 分支。
