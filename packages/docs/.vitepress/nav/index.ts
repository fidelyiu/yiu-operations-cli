import { DefaultTheme } from "vitepress";

export const nav: DefaultTheme.NavItem[] = [
  { text: "首页", link: "/" },
  {
    text: "Yiu Ops",
    activeMatch: "/yiu-ops/",
    items: [
      { text: "指南", link: "/yiu-ops/docs/root" },
      { text: "开发", link: "/yiu-ops/dev/root" },
    ],
  },
  {
    text: "前端",
    activeMatch: "/fe/",
    items: [
      { text: "指南", link: "/fe/guide/root" },
      { text: "TS", link: "/fe/ts/root" },
    ],
  },
  {
    text: "Rust",
    activeMatch: "/rust/",
    items: [
      { text: "指南", link: "/rust/guide/root" },
      { text: "Axum", link: "/rust/axum/root" },
    ],
  },
  {
    text: "其他",
    activeMatch: "/other/",
    link: "/other/root",
  },
];
