---
sidebar_position: 2
---

# 翻译你的站点

让我们将 `docs/intro.md` 翻译成法语。

## 配置 i18n

修改 `docusaurus.config.js` 以添加对 `fr` 语言环境的支持：

```js title="docusaurus.config.js"
export default {
  i18n: {
    defaultLocale: "en",
    locales: ["en", "fr"],
  },
};
```

## 翻译文档

将 `docs/intro.md` 文件复制到 `i18n/fr` 文件夹：

```bash
mkdir -p i18n/fr/docusaurus-plugin-content-docs/current/

cp docs/intro.md i18n/fr/docusaurus-plugin-content-docs/current/intro.md
```

将 `i18n/fr/docusaurus-plugin-content-docs/current/intro.md` 翻译成法语。

## 启动本地化站点

在法语语言环境下启动站点：

```bash
npm run start -- --locale fr
```

你的本地化站点可以在 [http://localhost:3000/fr/](http://localhost:3000/fr/) 访问，`Getting Started` 页面已被翻译。

:::caution

在开发环境中，一次只能使用一个语言环境。

:::

## 添加语言下拉菜单

为了在语言之间无缝导航，添加语言下拉菜单。

修改 `docusaurus.config.js` 文件：

```js title="docusaurus.config.js"
export default {
  themeConfig: {
    navbar: {
      items: [
        // highlight-start
        {
          type: "localeDropdown",
        },
        // highlight-end
      ],
    },
  },
};
```

语言下拉菜单现在会出现在导航栏中：

![Locale Dropdown](./img/localeDropdown.png)

## 构建本地化站点

为特定语言环境构建站点：

```bash
npm run build -- --locale fr
```

或一次性构建包含所有语言环境的站点：

```bash
npm run build
```
