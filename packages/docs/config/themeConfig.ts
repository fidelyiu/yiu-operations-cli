import { themes as prismThemes } from "prism-react-renderer";
import type { Config } from "@docusaurus/types";
import type * as Preset from "@docusaurus/preset-classic";

export const themeConfig: Config["themeConfig"] = {
  // Replace with your project's social card
  image: "img/docusaurus-social-card.jpg",
  colorMode: {
    respectPrefersColorScheme: true,
  },
  docs: {
    sidebar: {
      hideable: true,
      autoCollapseCategories: true,
    },
  },
  navbar: {
    title: "Yiu Operations Docs",
    logo: {
      alt: "Yiu Operations Docs Logo",
      src: "img/logo.svg",
    },
    items: [
      {
        type: "docSidebar",
        sidebarId: "tutorialSidebar",
        position: "left",
        label: "Tutorial",
      },
      {
        type: "docSidebar",
        sidebarId: "communitySidebar",
        position: "left",
        label: "Community",
        docsPluginId: "community",
      },
      { to: "/blog", label: "Blog", position: "left" },
      // {
      //   href: "https://github.com/fidelyiu/yiu-operations-cli",
      //   label: "GitHub",
      //   position: "right",
      // },
    ],
  },
  footer: {
    style: "dark",
    // links: [
    //   {
    //     title: "Docs",
    //     items: [
    //       {
    //         label: "Tutorial",
    //         to: "/docs/intro",
    //       },
    //     ],
    //   },
    //   {
    //     title: "Community",
    //     items: [
    //       {
    //         label: "Stack Overflow",
    //         href: "https://stackoverflow.com/questions/tagged/docusaurus",
    //       },
    //       {
    //         label: "Discord",
    //         href: "https://discordapp.com/invite/docusaurus",
    //       },
    //       {
    //         label: "X",
    //         href: "https://x.com/docusaurus",
    //       },
    //     ],
    //   },
    //   {
    //     title: "More",
    //     items: [
    //       {
    //         label: "Blog",
    //         to: "/blog",
    //       },
    //       {
    //         label: "GitHub",
    //         href: "https://github.com/facebook/docusaurus",
    //       },
    //     ],
    //   },
    // ],
    copyright: `Copyright © ${new Date().getFullYear()} Yiu Operations Docs, Inc. Built with Fidel Yiu.`,
  },
  prism: {
    theme: prismThemes.github,
    darkTheme: prismThemes.dracula,
  },
  mermaid: {
    theme: { light: "neutral", dark: "forest" },
    options: {
      maxTextSize: 50,
    },
  },
} satisfies Preset.ThemeConfig;
