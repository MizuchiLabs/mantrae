import type { BaseLayoutProps } from "fumadocs-ui/layouts/shared";
import Logo from "@/styles/logo.svg?react";

export function baseOptions(): BaseLayoutProps {
  return {
    nav: {
      title: (
        <>
          <Logo className="w-6 h-6" />
          <span className="font-bold">Mantræ Docs</span>
        </>
      ),
    },
    githubUrl: "https://github.com/MizuchiLabs/mantrae",
    links: [
      {
        text: "Documentation",
        url: "/docs",
        on: "nav",
      },
    ],
  };
}
