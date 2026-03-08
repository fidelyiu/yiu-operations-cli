import { DefaultTheme } from "vitepress";
import { rustGuideSidebar } from "./rust/guide";
import { rustAxumSidebar } from "./rust/axum";
import { otherSidebar } from "./other";
import { yiuOpsDocsSidebar } from "./yiu-ops/docs";
import { yiuOpsDevSidebar } from "./yiu-ops/dev";
import { feGuideSidebar } from "./fe/guide";
import { feTsSidebar } from "./fe/ts";

export const sidebar: DefaultTheme.Sidebar = {
  "/fe/guide": feGuideSidebar,
  "/fe/ts": feTsSidebar,
  "/rust/guide": rustGuideSidebar,
  "/rust/axum": rustAxumSidebar,
  "/yiu-ops/docs": yiuOpsDocsSidebar,
  "/yiu-ops/dev": yiuOpsDevSidebar,
  "/other": otherSidebar,
};
