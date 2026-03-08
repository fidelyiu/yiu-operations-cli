import { DefaultTheme } from "vitepress";

export const nav: DefaultTheme.NavItem[] = [
  { text: "首页", link: "/" },
  {
    text: "Rust",
    activeMatch: "/rust/",
    items: [
      { text: "指南", link: "/rust/guide/root" },
      { text: "Axum", link: "/rust/axum/root" },
    ],
  },
];
