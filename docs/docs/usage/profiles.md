# Profiles

Profiles let you manage multiple Traefik instances from a single Mantrae installation. Each profile is an isolated environment with its own configuration.

## What's in a Profile

Each profile contains:
- Routers (HTTP, TCP, UDP)
- Services (HTTP, TCP, UDP)  
- Middlewares (HTTP, TCP)
- Entry points
- Server transports
- DNS providers
- Agents

Profiles are completely isolated from each other.

## Creating a Profile

1. Click the profile dropdown in the top navigation
2. Select "Create New Profile"
3. Configure:
   - **Name**: Unique identifier (e.g., `production`, `staging`)
   - **Description**: Optional purpose description
   - **Token**: Auto-generated security token for API access

The token is used by Traefik to authenticate when fetching configuration.

## Connecting Traefik

Each profile exposes a unique endpoint:
```
http://mantrae:3000/api/PROFILE_NAME?token=TOKEN
```

### Static Configuration

```yaml
providers:
  http:
    endpoint: "http://mantrae:3000/api/production?token=TOKEN"
    pollInterval: "5s"
```

### Docker Compose

```yaml
traefik:
  image: traefik:latest
  command:
    - --providers.http.endpoint=http://mantrae:3000/api/production?token=TOKEN
    - --providers.http.pollInterval=5s
```

## Switching Profiles

Use the profile dropdown in the top navigation to switch between profiles. All actions apply to the currently selected profile.

## Profile Tokens

- Each profile has a unique token
- Find the token in profile settings
- Tokens can be regenerated if compromised
- Never commit tokens to version control
