# Backups & Restoration

Mantrae provides automatic and manual backup capabilities to protect your configuration.

## Automatic Backups

Mantrae automatically backs up your database on a schedule. Configure via environment variables:

```bash
BACKUP_SCHEDULE=0 2 * * * # Daily at 2 AM (cron format)
BACKUP_RETENTION=30       # Keep for 30 days
BACKUP_AUTO_ENABLED=true  # Enable auto-backups
```

## Storage Options

### Local Storage (Default)

Backups stored in `/data/backups/` directory.

```bash
STORAGE_TYPE=local
```

### S3-Compatible Storage

Store backups in the cloud:

```bash
STORAGE_TYPE=s3
STORAGE_S3_ENDPOINT=https://s3.amazonaws.com
STORAGE_S3_REGION=us-east-1
STORAGE_S3_BUCKET=mantrae-backups
STORAGE_S3_ACCESS_KEY=your-access-key
STORAGE_S3_SECRET_KEY=your-secret-key
STORAGE_S3_PATH=backups/
```

Compatible with:
- AWS S3
- MinIO
- DigitalOcean Spaces
- Backblaze B2
- Any S3-compatible storage

## Manual Operations

Access backup controls in Settings:

### Create Backup

Generate an immediate database backup.

### Download Backup

Download existing backup files for off-site storage.

### Restore from Backup

Restore from a previous database backup.

:::warning
Restoring from a database backup replaces **all** existing data including profiles, routers, services, and middlewares.
:::

## Import Traefik Configuration

Import existing Traefik dynamic configuration files (YAML/JSON):

1. Go to Settings → Import Configuration
2. Select your Traefik configuration file
3. Review changes
4. Confirm import

**Import behavior:**
- Non-destructive - keeps existing data
- Merges new configuration with current setup
- Overwrites components with matching names
- Preserves profiles and settings

### Example Import File

```yaml
http:
  routers:
    my-router:
      rule: "Host(`example.com`)"
      service: "my-service"
      middlewares:
        - "auth"
      entryPoints:
        - "websecure"
      tls: {}

  services:
    my-service:
      loadBalancer:
        servers:
          - url: "http://backend:8080"

  middlewares:
    auth:
      basicAuth:
        users:
          - "admin:$apr1$H6uskkkW$IgXLP6ewTrSuBkTrqE8wj/"
```

## Best Practices

1. **Test restores**: Periodically verify backups can be restored
2. **Off-site storage**: Download backups or use S3 storage
3. **Before major changes**: Create manual backup
4. **Monitor retention**: Ensure old backups are cleaned up

## Backup Security

- Backups contain sensitive data (credentials, tokens)
- S3 transfers use HTTPS
- Encrypt S3 buckets at rest
- Restrict backup file access
- Use strong S3 credentials

:::tip
Set up S3 storage for automatic off-site backups and redundancy.
:::
