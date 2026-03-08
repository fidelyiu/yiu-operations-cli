import type { Config } from "@docusaurus/types";
import { i18nConfig } from "./config/i18n";
import { presetsConfig } from "./config/presets";
import { pluginsConfig } from "./config/plugins";
import { themeConfig } from "./config/themeConfig";

// 这在 Node.js 中运行 - 不要在这里使用客户端代码（浏览器 API、JSX...）

const config: Config = {
  title: "Yiu Operations CLI",
  tagline:
    "Yiu Operations CLI 是一个运维命令行工具，旨在简化和自动化各种操作任务。",
  favicon: "img/Yiu/favicon.ico",

  // Future flags，参见 https://docusaurus.io/docs/api/docusaurus-config#future
  future: {
    v4: true, // 提高与即将发布的 Docusaurus v4 的兼容性
  },

  // 在此设置站点的生产环境 URL
  url: "http://localhost:8282",
  // 设置站点提供服务的 /<baseUrl>/ 路径名
  // 对于 GitHub pages 部署，通常是 '/<projectName>/'
  baseUrl: "/",

  // GitHub pages 部署配置。
  // 如果你不使用 GitHub pages，则不需要这些。
  organizationName: "fidel-yiu", // 通常是你的 GitHub 组织/用户名。
  projectName: "yiu-operations-cli", // 通常是你的仓库名。

  onBrokenLinks: "throw",

  // 即使你不使用国际化，也可以使用此字段设置
  // 有用的元数据，如 html lang。例如，如果你的站点是中文的，
  // 你可能想将 "en" 替换为 "zh-Hans"。
  i18n: i18nConfig,
  presets: presetsConfig,
  plugins: pluginsConfig,
  themeConfig: themeConfig,
  markdown: {
    mermaid: true,
  },
  themes: [
    "@docusaurus/theme-mermaid",
    [
      require.resolve("@easyops-cn/docusaurus-search-local"),
      /** @type {import("@easyops-cn/docusaurus-search-local").PluginOptions} */
      {
        hashed: true,
        language: ["en", "zh"],
      },
    ],
  ],
};

// 2430080847&13124698770

export default config;
