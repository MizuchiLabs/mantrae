---
sidebar_position: 3
---

# Profiles

Profiles in Mantrae allow you to manage multiple Traefik instances from a single interface. Each profile represents a distinct Traefik configuration environment, such as development, staging, or production.

## Understanding Profiles

A profile is a logical grouping of:
- Routers (HTTP, TCP, UDP)
- Services (HTTP, TCP, UDP)
- Middlewares (HTTP, TCP)
- Entry Points
- Servers Transports
- DNS Providers
- Agents

Each profile has its own dynamic configuration endpoint that Traefik instances can consume.

## Creating a Profile

To create a new profile:

1. **Access the Web UI**: Log into Mantrae at `http://localhost:3000`
2. **Navigate to Profiles**: Click on the profile dropdown in the top navigation bar
3. **Create New Profile**: Select "Create New Profile"
4. **Configure Profile Details**:
   - **Name**: A unique identifier for this profile (e.g., `production`, `staging`)
   - **Description**: Optional description of the profile's purpose
   - **Token**: A security token for accessing this profile's configuration (auto-generated)

Once saved, this profile will serve as a dedicated space for managing Traefik configurations specific to this environment.

## Using Profiles with Traefik

Each profile in Mantrae exposes a unique API endpoint that Traefik can use to fetch its dynamic configuration.

### Dynamic Configuration Endpoint

For a profile named `production`, the dynamic configuration endpoint would be:
```
http://mantrae:3000/api/production?token=GENERATED_TOKEN
```

### Configure Traefik to Use Mantrae

Configure your Traefik instance to use Mantrae as its dynamic configuration provider:

#### Using Static Configuration File

```yaml
providers:
  http:
    endpoint: "http://mantrae:3000/api/production?token=GENERATED_TOKEN"
    pollInterval: "5s" # Optional: polling interval for configuration updates
```

#### Using Command Line Arguments

In Docker Compose:

```yaml
traefik:
  image: traefik:latest
  container_name: traefik
  command:
    - --providers.http.endpoint=http://mantrae:3000/api/production?token=GENERATED_TOKEN
    - --providers.http.pollInterval=5s
    # ... other Traefik configuration
```

#### Using Environment Variables

```bash
TRAEFIK_PROVIDERS_HTTP_ENDPOINT=http://mantrae:3000/api/production?token=GENERATED_TOKEN
TRAEFIK_PROVIDERS_HTTP_POLLINTERVAL=5s
```

## Profile-Specific Configuration

Each profile maintains its own separate configuration space. This means:

- Routers created in the `production` profile are completely isolated from those in the `staging` profile
- Services, middlewares, and other components are profile-specific
- Agents can be assigned to specific profiles
- DNS providers can be configured per profile

## Managing Multiple Profiles

You can easily switch between profiles using the profile dropdown in the Mantrae web interface:

1. Click the current profile name in the top navigation bar
2. Select the profile you want to work with from the dropdown
3. All subsequent actions will apply to the selected profile

This allows you to manage multiple Traefik environments without switching between different Mantrae instances.
