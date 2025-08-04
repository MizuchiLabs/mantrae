---
sidebar_position: 9
---

# API Reference

Mantrae exposes a comprehensive API for programmatic access to all configuration management features. The API is built using gRPC with the Connect protocol, providing both gRPC and HTTP/JSON interfaces.

## API Endpoints

The Mantrae API documentation is available at:
```
http://localhost:3000/docs
```

Each service has its own endpoint under this base path.

## Authentication

All API requests require authentication. You can authenticate using:

1. **Session cookies**: Automatically set when logging in through the web interface
2. **Authorization header**: Bearer token for programmatic access

For programmatic access, obtain an API token through the web interface or create one using the UserService API.

## Dynamic Configuration API

In addition to the management API, Mantrae provides endpoints for Traefik to consume dynamic configuration:

### Profile Configuration Endpoint

```
GET /api/{profile-name}
```

Query Parameters:
- `token`: Profile token for authentication
- `format`: Response format (json or yaml)

Headers:
- `Accept`: Alternative way to specify format (application/json or application/x-yaml)

This endpoint returns the complete Traefik dynamic configuration for the specified profile.

Example:
```bash
curl "http://localhost:3000/api/production?token=PROFILE_TOKEN&format=yaml"
```

