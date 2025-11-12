# Getting Started

Mantrae is a web interface for managing Traefik's dynamic configuration. It generates configuration files that Traefik consumes via its HTTP provider.

## What Mantrae Does

- Generates Traefik dynamic configuration (routers, services, middlewares)
- Manages multiple configuration profiles for different environments
- Provides optional agents for automatic Docker container discovery
- Integrates with DNS providers (Cloudflare, PowerDNS, Technitium)
- Backs up and restores configurations

:::note
Mantrae is **not** a Traefik dashboard. It doesn't monitor Traefik's status. Instead, Traefik pulls configuration from Mantrae.
:::

## Prerequisites

- Docker (recommended) or Go runtime
- A running Traefik instance
- OpenSSL or similar tool

## Installation

### 1. Generate a Secret

Generate a 16, 24, or 32 byte secret:

```bash
openssl rand -hex 16
```

### 2. Deploy Mantrae

Choose your deployment method:

#### Docker Compose (Recommended)

```yaml
services:
  mantrae:
    image: ghcr.io/mizuchilabs/mantrae:latest
    container_name: mantrae
    environment:
      - SECRET=your-generated-secret
      - ADMIN_PASSWORD=your-admin-password
    ports:
      - "3000:3000"
    volumes:
      - ./mantrae:/data
    restart: unless-stopped
```

Run with:
```bash
docker compose up -d
```

#### Docker

```bash
docker run --name mantrae \
   -e SECRET=your-generated-secret \
   -e ADMIN_PASSWORD=your-admin-password \
   -p 3000:3000 \
   -v ./mantrae:/data \
   ghcr.io/mizuchilabs/mantrae:latest
```

#### Binary

Download from the [releases page](https://github.com/mizuchilabs/mantrae/releases), then:

```bash
export SECRET=your-generated-secret
export ADMIN_PASSWORD=your-admin-password
./mantrae
```

### 3. Access the Interface

Open [http://localhost:3000](http://localhost:3000) and log in:
- Username: `admin`
- Password: your admin password

## Configure Traefik

Point Traefik to Mantrae's HTTP provider endpoint. You'll need a profile token (found in profile settings).

### Using traefik.yml

```yaml
providers:
  http:
    endpoint: "http://mantrae:3000/api/default?token=PROFILE_TOKEN"
    pollInterval: "5s"
```

### Using Docker Compose

```yaml
services:
  traefik:
    image: traefik:latest
    container_name: traefik
    command:
      - --providers.http.endpoint=http://mantrae:3000/api/default?token=PROFILE_TOKEN
      - --providers.http.pollInterval=5s
      - --entrypoints.web.address=:80
      - --entrypoints.websecure.address=:443
    ports:
      - "80:80"
      - "443:443"
    restart: unless-stopped
```

Replace `PROFILE_TOKEN` with your actual token from the profile settings.

## Create Your First Router

1. Navigate to "Routers" in the Mantrae interface
2. Click "Add Router"
3. Configure the router:
   - Name: `whoami`
   - Rule: Host(\`whoami.local\`)
   - Service: Create or select a service pointing to your backend
4. Save

Traefik will poll Mantrae and apply the configuration automatically.

## Verify Configuration

Check the dynamic configuration at:
```
http://localhost:3000/api/default?token=PROFILE_TOKEN
```

This returns the JSON configuration Traefik is consuming.

## Command Reference

```bash
# Start the server
mantrae

# Check version
mantrae --version

# Check for updates
mantrae update

# Update to latest version (binary only, not Docker)
mantrae update --install

# Reset admin password
mantrae reset --password newpassword

# Reset a specific user's password
mantrae reset --user username --password newpassword
```

## Next Steps

- [Profiles](./usage/profiles) - Manage multiple environments
- [Agents](./usage/agents) - Automatic container discovery
- [DNS Providers](./usage/dns) - DNS integration
- [Backups](./usage/backups) - Configuration backups
