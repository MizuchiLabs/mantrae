# DNS Providers

Mantrae integrates with DNS providers to automatically manage domain records for your routers and handle Let's Encrypt DNS challenges.

## Supported Providers

- Cloudflare
- PowerDNS
- Technitium

## Adding a Provider

1. Navigate to DNS section (globe icon)
2. Click "Add Provider"
3. Select provider type
4. Enter credentials:

### Cloudflare
- **API Token**: Scoped token with DNS edit permissions

### PowerDNS
- **API URL**: PowerDNS API endpoint
- **API Key**: Authentication key

### Technitium
- **API URL**: Technitium DNS API endpoint
- **API Key**: Authentication key
- **Zone Type**: `primary` or `forwarder`

## Using DNS Providers

### Automatic DNS Records

When you assign a DNS provider to a router:
- Mantrae creates DNS records for the router's domain
- Existing records are skipped to prevent overwrites
- Set a provider as "Default" to auto-assign it to new routers

### Certificate Management

DNS providers enable automatic TLS certificates:
1. Mantrae passes provider credentials to Traefik
2. Traefik performs DNS-01 ACME challenges
3. Certificates are automatically issued and renewed

No manual DNS configuration needed.

## Security

- Credentials are encrypted in the database
- Each profile has its own providers
- API tokens should use minimal required permissions
- Credentials are only accessible within their profile

:::tip
Use scoped API tokens with only DNS edit permissions for security.
:::

