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
          {
            text: "7. 错误处理",
            link: "/rust/axum/learning/tutorial-07/root",
          },
          {
            text: "8. 数据库",
            link: "/rust/axum/learning/tutorial-08/root",
          },
          {
            text: "9. 认证",
            link: "/rust/axum/learning/tutorial-09/root",
          },
          {
            text: "10. 高级功能",
            link: "/rust/axum/learning/tutorial-10/root",
          },
          {
            text: "11. 测试",
            link: "/rust/axum/learning/tutorial-11/root",
          },
        ],
      },
    ],
  },
];
