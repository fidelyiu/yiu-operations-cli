import { DefaultTheme } from "vitepress";

export const rustAxumSidebar: DefaultTheme.SidebarItem[] = [
  {
    text: "Axum指南",
    link: "/rust/axum/root",
    items: [
      {
        text: "Axum学习",
        link: "/rust/axum/learning/root",
        collapsed: true,
        items: [
          {
            text: "1. Hello World",
            link: "/rust/axum/learning/tutorial-01/root",
          },
          {
            text: "2. 路由",
            link: "/rust/axum/learning/tutorial-02/root",
          },
          {
            text: "3. 请求提取器",
            link: "/rust/axum/learning/tutorial-03/root",
          },
          {
            text: "4. 响应处理",
            link: "/rust/axum/learning/tutorial-04/root",
          },
          {
            text: "5. 状态共享",
            link: "/rust/axum/learning/tutorial-05/root",
          },
          {
            text: "6. 中间件",
            link: "/rust/axum/learning/tutorial-06/root",
          },
        ],
      },
    ],
  },
];
