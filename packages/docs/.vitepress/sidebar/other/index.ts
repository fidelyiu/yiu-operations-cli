import { DefaultTheme } from "vitepress";

export const otherSidebar: DefaultTheme.SidebarItem[] = [
  {
    text: "其他",
    link: "/other/root",
    items: [
      {
        text: "网站收藏",
        link: "/other/website/root",
      },
      {
        text: "GitHub Actions",
        link: "/other/github-actions/root",
        items: [
          {
            text: "机器人 Push 代码",
            link: "/other/github-actions/push-code/root",
          },
        ],
      },
      {
        text: "postgresql 触发器",
        link: "/other/postgresql-trigger/root",
      },
    ],
  },
];
