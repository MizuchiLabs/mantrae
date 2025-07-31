---
sidebar_position: 2
---

# DNS Providers

In Mantrae, you can set up DNS providers and configure them to automatically manage domain entries for your routers. This section explains how to add DNS providers, link them to routers, and manage DNS records for seamless integration with Traefik.

## Supported DNS Providers

Mantrae currently supports the following DNS providers:

- **Cloudflare**
- **PowerDNS**
- **Technitium**

Each provider can be configured within a profile, allowing you to use different DNS providers for different environments.

## Adding a DNS Provider

To add a DNS provider:

1. **Select Profile**: Ensure you're working with the correct profile
2. Navigate to the **DNS** section (globe icon)
3. Click "Add Provider"
4. Select your provider from the available options and enter the necessary credentials
5. For PowerDNS and Technitium you also need to set the endpoint where they are running
6. Save the provider. It will now be available for selection when configuring routers

## Provider-Specific Configuration

### Cloudflare

For Cloudflare, you'll need:

- **API Token**: A scoped API token with DNS permissions

### PowerDNS

For PowerDNS, you'll need:

- **API URL**: The URL to your PowerDNS API
- **API Key**: The API key for authentication

### Technitium

For Technitium DNS, you'll need:

- **API URL**: The URL to your Technitium DNS API
- **API Key**: The API key for authentication
- **Zone Type**: The zone type can be either `primary` or `forwarder`

## Setting a DNS Provider in Routers

Once a DNS provider is configured, you can assign it to specific routers. When you assign a DNS provider to a router:

- Mantrae will automatically attempt to add the router's domain name to the configured DNS provider
- **Duplicate Check**: If the domain already exists, Mantrae will skip it to avoid overwriting any existing records
- **Default**: Setting a provider as "Default" will automatically use it on newly created routers, so if no DNS provider is assigned to the router, Mantrae will use the default DNS provider

## Automatic Certificate Management

When using DNS providers with Traefik:

1. Traefik will automatically request certificates for configured domains
2. Mantrae will provide the DNS provider credentials to Traefik through the dynamic configuration
3. Traefik will create the necessary DNS challenge records
4. Certificates will be automatically renewed as needed

## Security Considerations

- DNS provider credentials are stored encrypted in the database
- Each profile can have its own set of DNS providers
- Credentials are only accessible to the profile they belong to
- API tokens should follow the principle of least privilege

> **Note**: This DNS automation only applies if no entry for the domain exists. Ensure your domain records are unique to prevent conflicts.
