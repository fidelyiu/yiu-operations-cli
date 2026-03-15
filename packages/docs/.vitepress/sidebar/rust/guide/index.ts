import { DefaultTheme } from "vitepress";

export const rustGuideSidebar: DefaultTheme.SidebarItem[] = [
  {
    text: "Rust指南",
    link: "/rust/guide/root",
    items: [
      {
        text: "所有权",
        link: "/rust/guide/ownership/root",
      },
      {
        text: "trait",
        link: "/rust/guide/trait-use/root",
      },
      {
        text: "async",
        link: "/rust/guide/async-use/root",
      },
      {
        text: "智能指针",
        link: "/rust/guide/smart-pointer/root",
      },
      {
        text: "async 和 dyn",
        link: "/rust/guide/async-dyn/root",
      },
    ],
  },
];
