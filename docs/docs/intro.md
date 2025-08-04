---
sidebar_position: 1
---

# Welcome to Mantrae

**Mantrae** is a modern web interface for managing [Traefik](https://traefik.io/) configurations, built using **Go** for backend services and **SvelteKit** for a reactive frontend experience. With Mantrae, you can easily create, update, and manage routers, middlewares, services, and other Traefik components across multiple profiles.

To understand how Mantrae works, check out our [Architecture](./architecture) documentation.

## Key Features

- **Multi-Profile Management**: Create and manage multiple profiles to handle different Traefik instances (development, staging, production)
- **Agent-Based Configuration**: Deploy agents to remote machines to automatically collect Docker container information and Traefik labels
- **Unified Configuration Interface**: Manage HTTP, TCP, and UDP routers, services, and middlewares from a single dashboard
- **Dynamic Config Generation**: Generate Traefik dynamic configuration files in real-time
- **DNS Integration**: Built-in support for Cloudflare, PowerDNS, and Technitium DNS providers
- **Backup & Restore**: Automatic backup of configurations with restore capabilities
- **Real-Time Updates**: Get instant feedback with a SvelteKit-powered reactive UI

## How It Works

1. **Create Profiles**: Set up profiles for different environments (dev, staging, production)
2. **Deploy Agents** (Optional): Install agents on remote machines to automatically collect container information
3. **Configure Traefik**: Point your Traefik instances to Mantrae's dynamic configuration endpoints
4. **Manage Configuration**: Use the web UI to create and manage your routing configuration

## Getting Started

To install and set up Mantrae, refer to our [Installation Guide](./installation).