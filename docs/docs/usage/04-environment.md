---
sidebar_position: 4
---

# Environment Variables

Mantrae provides several command-line flags and environment variables to configure the application. This guide details each option and its purpose.

## Command-Line Arguments

You can use the following flags to customize the behavior of Mantrae:

| Flag       | Type   | Default | Description                                                              |
| ---------- | ------ | ------- | ------------------------------------------------------------------------ |
| `-version` | `bool` | `false` | Prints the current version of Mantrae and exits.                         |
| `-update`  | `bool` | `false` | Updates Mantrae to the latest version. (Doesn't work inside a container) |

## Environment Variables

Environment variables can be used to set up Mantrae and configure its settings. Below is a list of the supported environment variables.

### Core Configuration

| Variable | Default | Description                                      |
| -------- | ------- | ------------------------------------------------ |
| `SECRET` |         | Secret key required for secure access. Required! |

### Server Configuration

| Variable              | Default            | Description                                                |
| --------------------- | ------------------ | ---------------------------------------------------------- |
| `SERVER_HOST`         | `0.0.0.0`          | Host address the server will bind to                       |
| `SERVER_PORT`         | `3000`             | Port which Mantrae will listen on                          |
| `SERVER_URL`          | `http://127.0.0.1` | The public URL of the Mantrae server for agent connections |
| `SERVER_BASIC_AUTH`   | `false`            | Enables basic authentication for the Mantrae server        |
| `SERVER_ENABLE_AGENT` | `true`             | Enables the Mantrae agent functionality                    |
| `SERVER_LOG_LEVEL`    | `info`             | Log verbosity level                                        |

### Admin Configuration

| Variable         | Default         | Description         |
| ---------------- | --------------- | ------------------- |
| `ADMIN_USERNAME` | `admin`         | Admin user username |
| `ADMIN_EMAIL`    | `admin@mantrae` | Admin user email    |
| `ADMIN_PASSWORD` |                 | Admin user password |

### Email Configuration

| Variable         | Default             | Description          |
| ---------------- | ------------------- | -------------------- |
| `EMAIL_HOST`     | `localhost`         | SMTP server host     |
| `EMAIL_PORT`     | `587`               | SMTP server port     |
| `EMAIL_USERNAME` |                     | SMTP server username |
| `EMAIL_PASSWORD` |                     | SMTP server password |
| `EMAIL_FROM`     | `mantrae@localhost` | Sender email address |

### Database Configuration

| Variable  | Default   | Description                                             |
| --------- | --------- | ------------------------------------------------------- |
| `DB_TYPE` | `sqlite`  | Database type. Supported options: only `sqlite` for now |
| `DB_NAME` | `mantrae` | Database/file name                                      |

### Backup Configuration

| Variable          | Default   | Description                   |
| ----------------- | --------- | ----------------------------- |
| `BACKUP_ENABLED`  | `true`    | Enable automatic backups      |
| `BACKUP_PATH`     | `backups` | Directory for storing backups |
| `BACKUP_INTERVAL` | `24h`     | Interval between backups      |
| `BACKUP_KEEP`     | `3`       | Number of backups to keep     |

### Traefik Configuration

| Variable           | Default   | Description                     |
| ------------------ | --------- | ------------------------------- |
| `TRAEFIK_PROFILE`  | `default` | Traefik profile name            |
| `TRAEFIK_URL`      |           | Traefik API URL                 |
| `TRAEFIK_USERNAME` |           | Traefik authentication username |
| `TRAEFIK_PASSWORD` |           | Traefik authentication password |
| `TRAEFIK_TLS`      | `false`   | Enable TLS for Traefik          |

### Background Jobs Configuration

| Variable                  | Default | Description                     |
| ------------------------- | ------- | ------------------------------- |
| `BACKGROUND_JOBS_TRAEFIK` | `20`    | Traefik background job interval |
| `BACKGROUND_JOBS_DNS`     | `300`   | DNS background job interval     |

### Example Usage

To run Mantrae with custom environment variables:

```bash
export SECRET="your-secret-key"
export SERVER_PORT="4000"
export ADMIN_PASSWORD="secure-password"
./mantrae
```

### Important Notes

- **SECRET** is a required environment variable and must be set; otherwise, the application will not start.
- Set **SERVER_URL** to the publicly accessible URL of Mantrae to ensure agents can connect to it.
