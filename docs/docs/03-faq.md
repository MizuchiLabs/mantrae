---
sidebar_position: 4
---

# FAQ

### Why would I use this? Traefik already has configuration discovery.

Yes, Traefik has amazing configuration discovery capabilities for various providers (Docker, Kubernetes, etc.). But for all those times you can't use these features (e.g. multiple machines not connected via Docker Swarm or Kubernetes) you have to use the file provider. Mantrae helps you with that and adds additional automation features like managing DNS records as well, similar to external-dns for Kubernetes.

### I want to use multiple DNS providers of the same type (e.g. multiple cloudflare accounts), how do I do that?

Traefik doesn't support multiple DNS Challenge providers, so you have to use CNAME records to manage multiple accounts.
E.g. if you have a domain `example.com` on account "Foo" and a domain `example.org` on account "Bar", you can add the API Key for account "Foo" normally, but to get letsencrypt certificates for `example.org` you need add a CNAME record for `example.org` with these values:

- Type: `CNAME`
- Name: `_acme-challenge.example.org`
- Target: `_acme-challenge.example.com`

Now you can request certificates for `sub.example.org` as well.
