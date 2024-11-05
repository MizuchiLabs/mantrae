import { themes as prismThemes } from "prism-react-renderer";
import type { Config } from "@docusaurus/types";
import type * as Preset from "@docusaurus/preset-classic";

// This runs in Node.js - Don't use client-side code here (browser APIs, JSX...)

const config: Config = {
  title: "Mantrae",
  tagline:
    "Simple yet powerful Traefik manager, enhanced with advanced features.",
  favicon: "img/favicon.ico",
  url: "https://mizuchi.dev/",
  baseUrl: "/mantrae/",
  organizationName: "mizuchilabs",
  projectName: "mantrae",
  onBrokenLinks: "throw",
  onBrokenMarkdownLinks: "warn",
  deploymentBranch: "gh-pages",

  i18n: {
    defaultLocale: "en",
    locales: ["en"],
  },

  presets: [
    [
      "classic",
      {
        docs: {
          sidebarPath: "./sidebars.ts",
        },
        blog: {
          showReadingTime: true,
          feedOptions: {
            type: ["rss", "atom"],
            xslt: true,
          },
          onInlineTags: "warn",
          onInlineAuthors: "warn",
          onUntruncatedBlogPosts: "warn",
        },
        theme: {
          customCss: "./src/css/custom.css",
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    // Replace with your project's social card
    //image: "img/mantrae-social-card.jpg",
    navbar: {
      title: "Mantrae",
      logo: {
        alt: "Mantrae Logo",
        src: "img/logo.svg",
      },
      items: [
        {
          type: "docSidebar",
          sidebarId: "tutorialSidebar",
          position: "left",
          label: "Docs",
        },
        //{ to: "/blog", label: "Blog", position: "left" },
        {
          href: "https://github.com/mizuchilabs/mantrae",
          position: "right",
          className: "header-github-link",
          "aria-label": "GitHub repository",
        },
        {
          href: "https://buymeacoffee.com/d34dscene",
          position: "right",
          className: "header-coffee-link",
          "aria-label": "Support the project",
        },
      ],
    },

    footer: {
      style: "dark",
      links: [
        {
          title: "Docs",
          items: [
            {
              label: "Introduction",
              to: "/docs/intro",
            },
            {
              label: "Usage",
              to: "/docs/category/usage",
            },
          ],
        },
        {
          title: "More",
          items: [
            {
              label: "Mizuchi Labs",
              to: "https://mizuchi.dev",
            },
            {
              label: "GitHub",
              href: "https://github.com/mizuchilabs/mantrae",
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} Mizuchi Labs, Inc. Built with Docusaurus.`,
    },

    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
    },

    announcementBar: {
      id: "github-star",
      content: `If you like Mantrae, <a href=https://github.com/mizuchilabs/mantrae rel="noopener noreferrer" target="_blank">give us a star on GitHub</a>! ⭐️`,
      backgroundColor: "var(--ifm-color-primary)",
      textColor: "var(--ifm-color-white)",
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
