import { DefaultTheme } from "vitepress";

export const rustWorkspaceKitSidebar: DefaultTheme.SidebarItem[] = [
  {
    text: "WorkspaceKit学习项目",
    link: "/rust/workspace-kit/root",
    items: [
      {
        text: "1. 初始化项目",
        link: "/rust/workspace-kit/wk-1-init-project/root",
        collapsed: true,
        items: [
          {
            text: "postgresql 触发器",
            link: "/rust/workspace-kit/wk-1-init-project/postgresql-trigger/root",
          },
        ],
      },
      {
        text: "2. config 模块",
        link: "/rust/workspace-kit/wk-2-config-module/root",
      },
    ],
  },
];
