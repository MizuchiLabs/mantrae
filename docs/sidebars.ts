import type { SidebarsConfig } from "@docusaurus/plugin-content-docs";

const sidebars: SidebarsConfig = {
   tutorialSidebar: [
      "quickstart",
      "faq",
      {
         type: "category",
         label: "Usage",
         items: [
            "usage/profiles",
            "usage/config",
            "usage/dns",
            "usage/agents",
            "usage/backups",
            "usage/oidc",
         ],
      },
   ],
};

export default sidebars;
