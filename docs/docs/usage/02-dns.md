---
sidebar_position: 2
---

# DNS

In Mantrae, you can set up DNS providers and configure them to automatically manage domain entries for your routers. This section explains how to add DNS providers, link them to routers, and manage DNS records for seamless integration with Traefik.

## Supported DNS Providers

Mantrae currently supports the following DNS providers:

- **Cloudflare**
- **PowerDNS**
- **Technitium**

### Adding a DNS Provider

To add a DNS provider:

1. Navigate to the **DNS** section (globe icon).
2. Select your provider from the available options and enter the necessary credentials.
3. For PowerDNS and Technitium you also need to set the endpoint where they are running.
4. Save the provider. It will now be available for selection when configuring routers.

### Setting a DNS Provider in Routers

Once a DNS provider is configured, you can assign it to specific routers. When you assign a DNS provider to a router:

- Mantrae will automatically attempt to add the router's domain name to the configured DNS provider.
- **Duplicate Check**: If the domain already exists, Mantrae will skip it to avoid overwriting any existing records.
- **Default**: Setting a provider as "Default" will automatically use it on newly created routers, so if no DNS provider is assigned to the router, Mantrae will use the default DNS provider.

> **Note**: This DNS automation only applies if no entry for the domain exists. Ensure your domain records are unique to prevent conflicts.
