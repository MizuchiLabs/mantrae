# Configuration

Configure Mantrae through environment variables and command-line arguments.

## Environment Variables

| Variable         | Default        | Required | Description                                 |
| ---------------- | -------------- | -------- | ------------------------------------------- |
| `SECRET`         | -              | Yes      | Encryption key (16, 24, or 32 bytes)        |
| `ADMIN_PASSWORD` | auto-generated | No       | Initial admin password                      |
| `LOG_LEVEL`      | `info`         | No       | Log level: `debug`, `info`, `warn`, `error` |
| `LOG_FORMAT`     | `text`         | No       | Log format: `text`, `json`                  |

:::warning
The `SECRET` environment variable must be 16, 24, or 32 bytes. Mantrae will not start without it.
:::

### Backup Configuration

| Variable              | Default     | Description                   |
| --------------------- | ----------- | ----------------------------- |
| `BACKUP_SCHEDULE`     | `0 2 * * *` | Cron schedule (daily at 2 AM) |
| `BACKUP_RETENTION`    | `30`        | Days to retain backups        |
| `BACKUP_AUTO_ENABLED` | `true`      | Enable automatic backups      |

### Storage Configuration

#### Local Storage (Default)

```bash
STORAGE_TYPE=local
```

Backups stored in `/data/backups/`

#### S3-Compatible Storage

```bash
STORAGE_TYPE=s3
STORAGE_S3_ENDPOINT=https://s3.amazonaws.com
STORAGE_S3_REGION=us-east-1
STORAGE_S3_BUCKET=mantrae-backups
STORAGE_S3_ACCESS_KEY=your-access-key
STORAGE_S3_SECRET_KEY=your-secret-key
STORAGE_S3_PATH=backups/
```

Works with AWS S3, MinIO, DigitalOcean Spaces, etc.

## Command-Line Interface

### Server Commands

```bash
# Start server
mantrae

# Check version
mantrae --version
mantrae -v
```

### Password Management

```bash
# Reset admin password
mantrae reset --password newpassword
mantrae reset -p newpassword

# Reset specific user
mantrae reset --user username --password newpassword
mantrae reset -u username -p newpassword
```

### Flags

| Flag         | Short | Description                           |
| ------------ | ----- | ------------------------------------- |
| `--version`  | `-v`  | Display version and exit              |
| `--password` | `-p`  | Password for reset command            |
| `--user`     | `-u`  | Username for reset (default: `admin`) |

## Docker Example

```yaml
services:
  mantrae:
    image: ghcr.io/mizuchilabs/mantrae:latest
    environment:
      - SECRET=${SECRET}
      - ADMIN_PASSWORD=${ADMIN_PASSWORD}
      - LOG_LEVEL=info
      - LOG_FORMAT=json
      - BACKUP_SCHEDULE=0 3 * * *
      - BACKUP_RETENTION=90
      - STORAGE_TYPE=s3
      - STORAGE_S3_BUCKET=my-backups
      - STORAGE_S3_ACCESS_KEY=${S3_ACCESS_KEY}
      - STORAGE_S3_SECRET_KEY=${S3_SECRET_KEY}
    ports:
      - "3000:3000"
    volumes:
      - ./mantrae:/data
    restart: unless-stopped
```
