---
slug: yiu-ops-docs
title: Github Action push 代码
authors: [FidelYiu]
tags: [gitHubActions]
---

# Github Action push 代码

我们可以在 Github Action 上配置一个机器人的 git 信息，然后用它 push 代码 和 tag。

但是这个push的tag貌似不会触发workflow。

<!-- truncate -->

- [github 关于 github actions bot 的讨论](https://github.com/orgs/community/discussions/160496)

```
git config user.name 'github-actions[bot]'
git config user.email 'github-actions[bot]@users.noreply.github.com'
```

使用案例

```yaml
name: Changesets Tag

on:
  push:
    branches:
      - master

permissions:
  contents: write

env:
  NODE_VERSION: 24
  PNPM_VERSION: 10.28.2

jobs:
  tag:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v6
        with:
          fetch-depth: 0

      - name: Setup Node
        uses: actions/setup-node@v6
        with:
          node-version: ${{ env.NODE_VERSION }}

      - name: Setup pnpm
        uses: pnpm/action-setup@v4
        with:
          version: ${{ env.PNPM_VERSION }}

      - name: Install dependencies
        run: pnpm install --frozen-lockfile

      - name: Configure git user
        run: |
          git config user.name "github-actions[bot]"
          git config user.email 'github-actions[bot]@users.noreply.github.com'

      - name: Tag release
        run: pnpm run pkg:tag

      - name: Push tags
        run: git push origin --tags
```
