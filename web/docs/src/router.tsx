import { createMemoryHistory, createRouter as createTanStackRouter } from "@tanstack/react-router";
import { routeTree } from "./routeTree.gen";
import { NotFound } from "@/components/not-found";

const memoryHistory = createMemoryHistory({
  initialEntries: ["/"],
});
export function getRouter() {
  return createTanStackRouter({
    routeTree,
    defaultPreload: "intent",
    scrollRestoration: true,
    trailingSlash: "always",
    history: memoryHistory,
    defaultNotFoundComponent: NotFound,
  });
}
