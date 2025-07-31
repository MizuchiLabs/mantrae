---
sidebar_position: 4
---

# FAQ

## General Questions

### Why would I use this? Traefik already has configuration discovery.

Yes, Traefik has amazing configuration discovery capabilities for various providers (Docker, Kubernetes, etc.). But for all those times you can't use these features (e.g. multiple machines not connected via Docker Swarm or Kubernetes) you have to use the file provider. Mantrae helps you with that and adds additional automation features like managing DNS records as well, similar to external-dns for Kubernetes.

Mantrae is particularly useful when you need to:
- Manage multiple separate Traefik instances from a single interface
- Automatically synchronize configurations from remote Docker hosts using agents
- Have a visual interface for managing complex routing configurations
- Centralize DNS management across multiple domains and providers

### What's the difference between using Mantrae and Traefik's file provider directly?

Using Mantrae provides several advantages over direct file provider management:
- Web-based interface for easier configuration management
- Multi-profile support for different environments
- Agent-based automatic discovery of container configurations
- Built-in DNS provider integration for certificate management
- Backup and restore capabilities
- User management and audit logging
- Real-time configuration updates without manual file editing

## DNS Providers

### I want to use multiple DNS providers of the same type (e.g. multiple Cloudflare accounts), how do I do that?

Traefik doesn't support multiple DNS Challenge providers, so you have to use CNAME records to manage multiple accounts.

E.g. if you have a domain `example.com` on account "Foo" and a domain `example.org` on account "Bar", you can add the API Key for account "Foo" normally, but to get letsencrypt certificates for `example.org` you need add a CNAME record for `example.org` with these values:

- Type: `CNAME`
- Name: `_acme-challenge.example.org`
- Target: `_acme-challenge.example.com`

Now you can request certificates for `sub.example.org` as well.

### Can I use different DNS providers for different domains within the same profile?

Yes, you can configure multiple DNS providers within a single profile and assign different providers to different routers. This allows you to manage domains across multiple DNS providers from a single Mantrae profile.

## Agents

### How do agents handle network connectivity when behind firewalls?

Agents establish outbound connections to the Mantrae server, so they work well in environments where the Mantrae server is in a DMZ or cloud environment, and agents are behind firewalls. The Mantrae server does not need to initiate connections to the agents.

### What happens if an agent loses connectivity to the Mantrae server?

When an agent loses connectivity:
1. The agent will continue to monitor Docker containers locally
2. Upon reconnection, the agent will synchronize any changes that occurred while offline
3. The Mantrae server will mark the agent as offline after a timeout period
4. Existing configurations will continue to function in Traefik

### Can I run multiple agents on the same host?

Yes, you can run multiple agents on the same host, each configured for different profiles or with different tokens. This is useful in environments where you need to separate configurations for different purposes.

## Profiles & Multi-Environment Management

### Can I share configurations between profiles?

Currently, configurations are profile-specific for security and isolation. However, you can export configurations from one profile and import them into another through the backup/import functionality.

### How do I migrate from a single-profile setup to multi-profile?

1. Create new profiles for each environment
2. Export your existing configuration
3. Import the configuration into each new profile as needed
4. Update your Traefik instances to point to the new profile endpoints

## Security

### How are credentials and secrets protected?

- Passwords are hashed using bcrypt with a cost factor of 12
- API tokens are generated using cryptographically secure random generators
- Database connections use prepared statements to prevent SQL injection
- Communication between agents and the server can be secured with HTTPS
- DNS provider credentials are encrypted at rest in the database

### Can I enforce specific security policies?

Mantrae supports:
- OIDC/SAML integration for enterprise authentication
- Audit logging of all configuration changes
- Regular backup capabilities for disaster recovery

## Troubleshooting

### My routers aren't working, how can I debug this?

1. Check the Traefik logs for any error messages related to configuration loading
2. Verify that the profile token in your Traefik configuration is correct
3. Ensure the Mantrae server is accessible from your Traefik instance
4. Check that the router rules are syntactically correct
5. Confirm that services referenced by routers exist and are correctly configured

## Development & Customization

### Is there an API I can integrate with?

Yes, Mantrae exposes a gRPC API using the Connect protocol. The API definitions are in the `proto/mantrae/v1/` directory of the source code.

### How do I request new features?

Feature requests can be submitted as GitHub issues in the Mantrae repository. Community contributions are also welcome through pull requests.
