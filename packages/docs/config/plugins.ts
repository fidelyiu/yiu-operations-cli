import type { Config } from "@docusaurus/types";
import { configPath } from "./utils.js";

export const pluginsConfig: Config["plugins"] = [
  "docusaurus-plugin-sass",
  [
    "@docusaurus/plugin-content-docs",
    {
      id: "examples",
      path: "examples",
      routeBasePath: "examples",
      sidebarPath: configPath("/sidebars/sidebarsExamples.ts"),
    },
  ],
];
