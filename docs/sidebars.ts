import type { SidebarsConfig } from "@docusaurus/plugin-content-docs";

const sidebars: SidebarsConfig = {
  tutorialSidebar: [
    "intro",
    "quickstart",
    "faq",
    "api",
    {
      type: "category",
      label: "Usage",
      items: [
        "usage/profiles",
        "usage/dns",
        "usage/agents",
        "usage/environment",
        "usage/backups",
      ],
    },
  ],
};

export default sidebars;
