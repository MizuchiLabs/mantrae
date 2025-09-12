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
openssl rand -hex 16
```

Save this secret for the next step. It has to be either of size 16, 24, or 32 bytes.

## Step 2: Run Mantrae

### Using the Binary

1. Visit the [releases page](https://github.com/mizuchilabs/mantrae/releases) to download the latest Mantrae binary for your platform.
2. Run the binary with:
   ```bash
   export SECRET=<your_secret>
   export ADMIN_PASSWORD=<your_admin_password>
   ./mantrae
   ```
### Using Docker

```bash
docker run --name mantrae \
   -e SECRET=your-generated-secret \
   -e ADMIN_PASSWORD=your-admin-password \
   -p 3000:3000 \
   -v mantrae-data:/app/data \
   ghcr.io/mizuchilabs/mantrae:latest
```

### Using Docker Compose (Recommended)

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
      - ./mantrae:/app/data
    restart: unless-stopped
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
   - Rule: Host(\`whoami.local\`)
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

