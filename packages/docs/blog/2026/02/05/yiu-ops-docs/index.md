---
slug: yiu-ops-docs
title: Yiu Ops 的 docs 命令开发
authors: [FidelYiu]
tags: [go, yiuOps]
---

# Yiu Ops 的 docs 命令开发

Yiu Ops 的 docs 命令开发过程中的设计和实现。

<!-- truncate -->

## 目标

```bash
# 当执行 docs 子命令之后
yiu-ops docs
# 我们就可以在浏览器中访问 yiu ops cli 的文档
# http://localhost:8281
```

## 添加子命令

```bash
cobra-cli add docs
```
