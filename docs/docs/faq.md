# FAQ

## General Questions

### Why would I use Mantrae instead of Traefik's built-in configuration discovery?

Traefik's automatic discovery works great for Docker Swarm and Kubernetes, but many users run Traefik in simpler setups where these features aren't available. Mantrae fills this gap by:

- Providing a visual interface for managing Traefik's file provider configuration
- Supporting multi-host Docker setups without requiring Swarm or Kubernetes
- Offering centralized management of multiple Traefik instances
- Automatically managing DNS records alongside routing configuration
- Enabling team collaboration with user management and audit logs

### Is Mantrae a replacement for Traefik's dashboard?

No. Mantrae is a **configuration manager**, not a monitoring dashboard. It generates the dynamic configuration that Traefik consumes via its HTTP provider. Traefik's dashboard shows real-time routing status and metrics, while Mantrae helps you create and manage those routes.

### Can I use Mantrae alongside Traefik's other providers (Docker, Kubernetes, etc.)?

Yes! Mantrae uses Traefik's HTTP provider, which works alongside other providers. You can have Traefik auto-discover some services while managing others through Mantrae. Just be careful to avoid naming conflicts between configurations.

### What happens if Mantrae goes down?

Traefik caches the configuration it receives from Mantrae. If Mantrae becomes unavailable:
- Existing routes continue to work with the last cached configuration
- You cannot make configuration changes until Mantrae is back online
- DNS updates (if configured) won't happen until Mantrae recovers

## Profiles

### What are profiles and when should I use them?

Profiles are isolated configuration environments within Mantrae. Each profile has its own:
- Routers, services, and middleware
- DNS provider settings
- Agent connections
- API token for Traefik

Common use cases:
- Separate configurations for development, staging, and production
- Different Traefik instances for different teams or projects
- Isolating configurations by customer in managed hosting scenarios

## Agents

### What do agents do and when do I need them?

Agents watch Docker containers on remote hosts and automatically create Mantrae configurations based on container labels. Use agents when:

- Your Docker hosts can't be reached directly by Mantrae (behind NAT/firewall)
- You want automatic discovery similar to Traefik's Docker provider
- Managing multiple Docker hosts from a single Mantrae instance
- You prefer labeling containers over manual configuration

If you're managing static configurations or non-Docker services, you don't need agents.

### How do agents communicate with Mantrae?

Agents establish **outbound** connections to the Mantrae server using gRPC. This means:
- Agents work behind firewalls and NAT (only outbound access needed)
- No ports need to be opened on the agent's host
- Communication can be secured with HTTPS/TLS
- Mantrae server must be reachable from agent hosts

### What happens when an agent disconnects?

- The agent continues monitoring containers locally
- When reconnected, it syncs any changes that occurred offline
- Existing Traefik configurations remain active
- Mantrae shows the agent as offline in the interface

## DNS Integration

### Which DNS providers are supported?

Currently supported:
- Cloudflare
- PowerDNS
- Technitium DNS

DNS integration automatically creates/updates DNS records when you configure routers, eliminating manual DNS management.

### How do I handle multiple accounts for the same DNS provider?

Traefik's ACME challenge only supports one provider instance. Use CNAME delegation:

1. Configure the primary account (e.g., Cloudflare account A with `example.com`)
2. For domains on other accounts (e.g., `example.org` on account B), add a CNAME:
   - Type: `CNAME`
   - Name: `_acme-challenge.example.org`
   - Target: `_acme-challenge.example.com`

This allows Traefik to validate certificates for domains across multiple accounts.

## Configuration Management

### How do I version control my Mantrae configurations?

Use the backup feature to export configurations as JSON. You can:
- Store exports in Git for version control
- Include them in your infrastructure-as-code repository
- Automate exports using Mantrae's API
- Restore from backups when needed

### Can I manage Mantrae configurations via API or CLI?

Mantrae exposes a gRPC API using the Connect protocol (definitions in `proto/mantrae/v1/`). You can:
- Programmatically create/update configurations
- Integrate with CI/CD pipelines
- Build custom tooling on top of Mantrae
- Automate configuration management

A CLI tool is planned for future releases.

### How do I migrate from manual file provider configuration to Mantrae?

1. Deploy Mantrae and create a profile
2. Import your existing dynamic yaml config
3. Verify routes work correctly
4. Remove or comment out the old file provider configuration

Start with non-critical routes first to minimize risk.

## Security

### How are credentials and API tokens protected?

- User passwords: hashed with bcrypt
- API tokens: cryptographically secure random generation
- DNS provider credentials: encrypted at rest
- Database queries: prepared statements prevent SQL injection
- Agent communication: supports HTTPS/TLS encryption

### Can multiple users access Mantrae?

Yes! Mantrae supports:
- Multiple user accounts with individual credentials
- OIDC integration
- Audit logging of all configuration changes

### How do I rotate profile tokens?

Profile tokens are used by Traefik to fetch configurations. To rotate:
1. Generate a new token in the profile settings
2. Update your Traefik configuration with the new token
3. Wait for Traefik to successfully pull configuration with the new token

Traefik will continue using the old token until you update its configuration.

## Troubleshooting

### Traefik isn't applying my Mantrae configuration

Check these common issues:

1. **Verify the HTTP provider endpoint**: Ensure the URL and token are correct in Traefik's configuration
2. **Check Traefik logs**: Look for errors fetching from the HTTP provider
3. **Test the endpoint manually**: Access `http://mantrae:3000/api/PROFILE?token=TOKEN` in a browser
4. **Verify network connectivity**: Ensure Traefik can reach Mantrae (especially in Docker networks)
5. **Check poll interval**: Default is 5s, but you can adjust if needed

### DNS records aren't being created

Verify:
1. DNS provider credentials are correct in Mantrae
2. API token has sufficient permissions for DNS management
3. The domain is managed by the configured DNS provider
4. Router has DNS provider assigned in configuration
5. Check Mantrae logs for DNS provider errors

### Agent shows as offline but Docker containers are running

Check:
1. Agent can reach Mantrae server (test with curl/ping)
2. Profile token is correct in agent configuration
3. Firewall isn't blocking outbound connections
4. Agent logs for connection errors
5. Mantrae server is running and accessible

## Performance

### How often does Traefik poll Mantrae?

The default poll interval is 5 seconds. You can adjust this in Traefik's configuration:
- Lower values: faster updates, more load on Mantrae
- Higher values: reduced load, slower configuration propagation

For most use cases, 5-10 seconds is appropriate.

