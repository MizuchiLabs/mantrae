---
sidebar_position: 2
---

# Quick Start

Get Mantrae up and running quickly with this guide. This will walk you through installing Mantrae, creating your first profile, and configuring Traefik to use it.

## Prerequisites

- Docker (recommended) or ability to run Go binaries
- A running Traefik instance
- OpenSSL or similar tool to generate a secret

## Step 1: Generate a Secret

First, generate a secure secret for Mantrae:

```bash
openssl rand -hex 32
```

Save this secret for the next step.

## Step 2: Run Mantrae

### Using Docker (Recommended)

```bash
docker run --name mantrae \
    -e SECRET=your-generated-secret \
    -e ADMIN_PASSWORD=your-admin-password \
    -p 3000:3000 \
    -v mantrae-data:/app/data \
    ghcr.io/mizuchilabs/mantrae:latest
```

### Using Docker Compose

Create a `docker-compose.yml` file:

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
      - mantrae-data:/app/data
    restart: unless-stopped

volumes:
  mantrae-data:
```

Then run:

```bash
docker compose up -d
```

## Step 3: Access the Web Interface

Open your browser and navigate to:

[http://localhost:3000](http://localhost:3000)

Log in with the username `admin` and your admin password.
- Username: `admin`
- Password: `your-admin-password`

**Important**: Change the default password immediately after first login!

## Step 4: Use the default profile or create your own

1. Click on the profile dropdown in the top navigation
2. Select "Create New Profile"
3. Enter profile details:
   - Name: `another-profile`
   - Description: `New profile for another site`
4. Click "Create"

## Step 5: Configure Traefik

Update your Traefik configuration to use Mantrae as a dynamic configuration provider.

### Using Traefik Configuration File

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

**Note**: Replace `PROFILE_TOKEN` with the actual token from your profile (visible in the profile settings).

## Step 6: Create Your First Router

1. In the Mantrae web interface, navigate to "Routers"
2. Click "Add Router"
3. Configure your router:
   - Name: `whoami`
   - Rule: `Host(\`whoami.local\`)`
   - Service: Create a new service pointing to your backend
   - Optional: Add middlewares or entry points
4. Save the router

## Step 7: Test Your Configuration

If you've set up everything correctly, Traefik should now be routing requests based on your Mantrae configuration.

You can verify the dynamic configuration by accessing:
```
http://localhost:3000/api/default?token=PROFILE_TOKEN
```

This should return the JSON configuration that Traefik is using.

## Next Steps

- [Learn about Profiles](./usage/profiles) to manage multiple environments
- [Set up Agents](./usage/agents) for automatic container discovery
- [Configure DNS Providers](./usage/dns) for automatic certificate management
- [Explore Backups](./usage/backups) to protect your configurations

## Troubleshooting

### "Invalid token" errors

Ensure that the token in your Traefik configuration matches the profile token in Mantrae.

### Traefik not picking up changes

- Check that Traefik can reach the Mantrae server
- Verify the poll interval is reasonable (5-30 seconds)
- Check Traefik logs for any error messages

### Can't access the web interface

- Ensure the port mapping is correct in your Docker configuration
- Check that no other service is using port 3000
- Verify the container is running with `docker ps`
