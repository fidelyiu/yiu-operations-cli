import { defineConfig } from "vitepress";
import { nav } from "./nav";
import { sidebar } from "./sidebar";
import { socialLinks } from "./socialLinks";
import { footer } from "./footer";
import { docFooter } from "./docFooter";
import { search } from "./search";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  lang: "zh-CN",
  title: "Yiu Operations CLI",
  description:
    "Yiu Operations CLI 是一个运维命令行工具，旨在简化和自动化各种操作任务。",
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
