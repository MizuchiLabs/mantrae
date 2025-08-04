---
sidebar_position: 5
---

# Backups & Restoration

Mantrae provides robust backup and restoration capabilities to help you manage and protect your configuration. You can perform both database backups and import configurations from Traefik dynamic config files.

## Automatic Backups

Mantrae automatically creates backups of your database according to the configuration settings detailed in the [Environment](./environment.md) documentation. These backups ensure you can recover your configuration if needed.

By default, backups are created daily and retained for 30 days. You can customize the backup schedule and retention period through environment variables.

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
      entryPoints:
        - "websecure"
      tls: {}

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

## Backup Storage

Backups can be stored in multiple locations:

### Local Storage

By default, backups are stored in the `data/backups/` directory within the Mantrae installation.

### S3-Compatible Storage

Configure S3-compatible storage for cloud backups:

```bash
STORAGE_TYPE=s3
STORAGE_S3_ENDPOINT=https://s3.amazonaws.com
STORAGE_S3_REGION=us-east-1
STORAGE_S3_BUCKET=mantrae-backups
STORAGE_S3_ACCESS_KEY=your-access-key
STORAGE_S3_SECRET_KEY=your-secret-key
STORAGE_S3_PATH=backups/
```

## Backup Configuration

Customize backup behavior through environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `BACKUP_SCHEDULE` | `0 2 * * *` | Cron expression for backup schedule (daily at 2 AM) |
| `BACKUP_RETENTION` | `30` | Number of days to keep backups |
| `BACKUP_AUTO_ENABLED` | `true` | Enable automatic backups |

## Best Practices

1. **Regular backups**: Although Mantrae handles automatic backups, consider creating manual backups before major changes
2. **Configuration versioning**: Store your Traefik configurations in version control for additional safety
3. **Test restorations**: Periodically verify your backup files can be successfully restored
4. **Off-site storage**: Download and store backups in a secure off-site location
5. **Multiple backup strategies**: Use both automatic backups and configuration imports for comprehensive protection

## Backup Security

- Backup files contain sensitive configuration data
- All backups are stored with appropriate file permissions
- When using S3 storage, data is transmitted securely
- Consider encrypting backup files when storing them off-site