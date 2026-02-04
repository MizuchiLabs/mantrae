import { createFileRoute, Link } from "@tanstack/react-router";
import { HomeLayout } from "fumadocs-ui/layouts/home";
import { baseOptions } from "@/lib/layout.shared";
import { ArrowRight, Bot, FileCheck, Globe } from "lucide-react";
import Github from "@/styles/github.svg?react";

export const Route = createFileRoute("/")({
  component: Home,
});

function Home() {
  return (
    <HomeLayout {...baseOptions()}>
      <div className="flex flex-col items-center justify-center text-center px-6 py-24 md:py-32 flex-1 relative overflow-hidden">
        <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full max-w-4xl h-full -z-10 opacity-10">
          <div className="absolute top-0 left-0 w-64 h-64 bg-fd-primary rounded-full blur-3xl animate-pulse" />
          <div className="absolute bottom-0 right-0 w-64 h-64 bg-fd-primary rounded-full blur-3xl animate-pulse delay-700" />
        </div>

        <h1 className="font-bold text-4xl md:text-6xl mb-6 tracking-tight">
          Manage Traefik with <span className="text-fd-primary">Ease</span>
        </h1>

        <p className="max-w-2xl text-fd-muted-foreground text-lg mb-10">
          Mantræ provides a clean, intuitive interface to manage your routers, middleware, and
          services without editing YAML files manually.
        </p>

        <div className="flex flex-wrap items-center justify-center gap-4">
          <Link
            to="/docs/$/"
            className="inline-flex items-center gap-2 px-6 py-3 rounded-xl bg-fd-primary text-fd-primary-foreground font-semibold transition-all hover:scale-105 active:scale-95 shadow-lg shadow-fd-primary/20"
          >
            Get Started
            <ArrowRight className="size-4" />
          </Link>

          <a
            href="https://github.com/MizuchiLabs/mantrae"
            target="_blank"
            rel="noreferrer"
            className="inline-flex items-center gap-2 px-6 py-3 rounded-xl bg-fd-secondary text-fd-secondary-foreground font-semibold border border-fd-border transition-all hover:bg-fd-accent"
          >
            <Github className="size-4" />
            GitHub
          </a>
        </div>

        <div className="mt-20 grid grid-cols-1 md:grid-cols-3 gap-8 max-w-5xl w-full">
          <div className="p-6 rounded-2xl border border-fd-border bg-fd-card text-left">
            <div className="w-10 h-10 rounded-lg bg-fd-primary/10 flex items-center justify-center mb-4">
              <FileCheck className="size-5 text-fd-primary" />
            </div>
            <h3 className="font-bold text-lg mb-2">Dynamic Config</h3>
            <p className="text-fd-muted-foreground text-sm">
              Manage routers, services, and middlewares through a simple web UI.
            </p>
          </div>

          <div className="p-6 rounded-2xl border border-fd-border bg-fd-card text-left">
            <div className="w-10 h-10 rounded-lg bg-fd-primary/10 flex items-center justify-center mb-4">
              <Bot className="size-5 text-fd-primary" />
            </div>
            <h3 className="font-bold text-lg mb-2">Agent Discovery</h3>
            <p className="text-fd-muted-foreground text-sm">
              Automatic container discovery with the mantraed agent.
            </p>
          </div>

          <div className="p-6 rounded-2xl border border-fd-border bg-fd-card text-left">
            <div className="w-10 h-10 rounded-lg bg-fd-primary/10 flex items-center justify-center mb-4">
              <Globe className="size-5 text-fd-primary" />
            </div>
            <h3 className="font-bold text-lg mb-2">DNS Integration</h3>
            <p className="text-fd-muted-foreground text-sm">
              Automatic DNS record management for Cloudflare, PowerDNS, and more.
            </p>
          </div>
        </div>
      </div>
    </HomeLayout>
  );
}
