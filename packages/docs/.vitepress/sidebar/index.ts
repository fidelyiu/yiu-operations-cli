import { DefaultTheme } from "vitepress";
import { rustGuideSidebar } from "./rust/guide";
import { rustAxumSidebar } from "./rust/axum";

export const sidebar: DefaultTheme.Sidebar = {
  "/rust/guide": rustGuideSidebar,
  "/rust/axum": rustAxumSidebar,
};
