---
sidebar_position: 3
---

# 创建博客文章

Docusaurus 为**每篇博客文章创建一个页面**，同时还提供**博客索引页面**、**标签系统**、**RSS** 订阅...

## 创建你的第一篇文章

在 `blog/2021-02-28-greetings.md` 创建文件：

```md title="blog/2021-02-28-greetings.md"
---
slug: greetings
title: Greetings!
authors:
  - name: Joel Marcey
    title: Co-creator of Docusaurus 1
    url: https://github.com/JoelMarcey
    image_url: https://github.com/JoelMarcey.png
  - name: Sébastien Lorber
    title: Docusaurus maintainer
    url: https://sebastienlorber.com
    image_url: https://github.com/slorber.png
tags: [greetings]
---

Congratulations, you have made your first post!

Feel free to play around and edit this post as much as you like.
```

现在可以在 [http://localhost:3000/blog/greetings](http://localhost:3000/blog/greetings) 访问新的博客文章了。
