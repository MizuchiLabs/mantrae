import { createRootRoute, HeadContent, Outlet, Scripts } from "@tanstack/react-router";
import * as React from "react";
import appCss from "@/styles/app.css?url";
import logoUrl from "@/styles/logo.svg?url";
import { RootProvider } from "fumadocs-ui/provider/tanstack";
import SearchDialog from "@/components/search";

export const Route = createRootRoute({
  head: () => ({
    meta: [
      {
        charSet: "utf-8",
      },
      {
        name: "viewport",
        content: "width=device-width, initial-scale=1",
      },
      {
        title: "Mantræ Docs",
      },
    ],
    links: [
      { rel: "stylesheet", href: appCss },
      { rel: "icon", href: logoUrl, type: "image/svg+xml" },
    ],
  }),
  component: RootComponent,
});

function RootComponent() {
  return (
    <RootDocument>
      <Outlet />
    </RootDocument>
  );
}

function RootDocument({ children }: { children: React.ReactNode }) {
  return (
    <html suppressHydrationWarning>
      <head>
        <HeadContent />
      </head>
      <body className="flex flex-col min-h-screen">
        <RootProvider search={{ SearchDialog }}>{children}</RootProvider>
        <Scripts />
      </body>
    </html>
  );
}
