import type { Config } from "@docusaurus/types";
import { configPath } from "./utils.js";

export const pluginsConfig: Config["plugins"] = [
  [
    "@docusaurus/plugin-content-docs",
    {
      id: "community",
      path: "community",
      routeBasePath: "community",
      sidebarPath: configPath("/sidebars/sidebarsCommunity.ts"),
    },
  ],
];
