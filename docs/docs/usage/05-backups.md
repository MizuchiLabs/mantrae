---
sidebar_position: 5
---

# Backups & Restoration

Mantrae provides robust backup and restoration capabilities to help you manage and protect your configuration. You can perform both database backups and import configurations from Traefik dynamic config files.

## Automatic Backups

Mantrae automatically creates backups of your database according to the configuration settings detailed in the [Environment](./04-environment.md) documentation. These backups ensure you can recover your configuration if needed.

## Manual Backups & Restoration

### Database Backups

From the Settings page, you can:

- Create an immediate backup of your entire database
- Download existing backup files
- Restore from a previous backup file

> **Note**: Restoring from a database backup will completely reset your Mantrae instance, replacing all existing data including profiles, routers, services, and middlewares.

### Traefik Configuration Import

For more flexible configuration management, Mantrae supports importing from Traefik dynamic configuration files:

1. Navigate to the Settings page
2. Select "Import Configuration"
3. Choose your Traefik YAML/JSON configuration file
4. Review the changes and confirm

Key benefits of configuration import:

- Non-destructive operation - existing data remains intact
- Merges new configuration with existing setup
- Overwrites only components with matching names
- Preserves your current profiles and settings

Example of an importable Traefik configuration:

```yaml
http:
  routers:
    my-router:
      rule: "Host(`example.com`)"
      service: "my-service"
      middlewares:
        - "auth-middleware"

  services:
    my-service:
      loadBalancer:
        servers:
          - url: "http://localhost:8080"

  middlewares:
    auth-middleware:
      basicAuth:
        users:
          - "test:$apr1$H6uskkkW$IgXLP6ewTrSuBkTrqE8wj/"
```

## Best Practices

1. Regular backups: Although Mantrae handles automatic backups, consider creating manual backups before major changes
2. Configuration versioning: Store your Traefik configurations in version control for additional safety
3. Test restorations: Periodically verify your backup files can be successfully restored
