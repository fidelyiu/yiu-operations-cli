import { defineConfig } from "vitepress";
import { nav } from "./nav";
import { sidebar } from "./sidebar";
import { socialLinks } from "./socialLinks";
import { footer } from "./footer";
import { docFooter } from "./docFooter";
import { search } from "./search";

// 2430080847&13124698770
// https://vitepress.dev/reference/site-config
export default defineConfig({
  lang: "zh-CN",
  title: "Yiu Operations CLI",
  description:
    "Yiu Operations CLI 是一个运维命令行工具，旨在简化和自动化各种操作任务。",
  head: [["link", { rel: "icon", href: "/img/Yiu/favicon.ico" }]],
  markdown: {
    config(md) {
      const defaultFence = md.renderer.rules.fence;

      md.renderer.rules.fence = (tokens, idx, options, env, self) => {
        const token = tokens[idx];

        if (token.info.trim() === "mermaid") {
          return `<pre class="mermaid">${md.utils.escapeHtml(token.content)}</pre>`;
        }

        if (defaultFence) {
          return defaultFence(tokens, idx, options, env, self);
        }

        return self.renderToken(tokens, idx, options);
      };
    },
  },
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    logo: "/img/Yiu/icononly_transparent_nobuffer.png",
    siteTitle: "Yiu Operations CLI",
    outline: { level: "deep", label: "页面大纲" },
    nav,
    sidebar,
    socialLinks,
    footer,
    docFooter,
    darkModeSwitchLabel: "外观",
    lightModeSwitchTitle: "切换亮色模式",
    darkModeSwitchTitle: "切换黑色模式",
    sidebarMenuLabel: "菜单",
    returnToTopLabel: "返回顶部",
    langMenuLabel: "语言",
    search,
  },
});
