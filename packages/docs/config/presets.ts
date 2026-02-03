import type { Config } from "@docusaurus/types";
import type * as Preset from "@docusaurus/preset-classic";
import { configPath, srcPath } from "./utils.js";

export const presetsConfig: Config["presets"] = [
  [
    "@docusaurus/preset-classic",
    {
      theme: {
        customCss: [srcPath("/css/custom.scss")],
      },
      docs: {
        sidebarPath: configPath("/sidebars/sidebarsDocs.ts"),
      },
      blog: {
        showReadingTime: true,
        feedOptions: {
          type: ["rss", "atom"],
          xslt: true,
        },
        // Please change this to your repo.
        // Remove this to remove the "edit this page" links.
        // editUrl:
        //   "https://github.com/facebook/docusaurus/tree/main/packages/create-docusaurus/templates/shared/",
        // Useful options to enforce blogging best practices
        onInlineTags: "warn",
        onInlineAuthors: "warn",
        onUntruncatedBlogPosts: "warn",
      },
    } satisfies Preset.Options,
  ],
];
