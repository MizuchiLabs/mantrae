import { createRouter as createTanStackRouter } from "@tanstack/react-router";
import { routeTree } from "./routeTree.gen";
import { NotFound } from "@/components/not-found";

export function getRouter() {
  return createTanStackRouter({
    routeTree,
    basepath: "/mantrae",
    defaultPreload: "intent",
    scrollRestoration: true,
    trailingSlash: "always",
    defaultNotFoundComponent: NotFound,
  });
}
