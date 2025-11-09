---
sidebar_position: 4
---

# Environment Variables

Mantrae provides several command-line flags and environment variables to configure the application. This guide details each option and its purpose.

## Command-Line Arguments

You can use the following flags to customize the behavior of Mantrae:

| Flag              | Type     | Default | Description                                                              |
| ----------------- | -------- | ------- | ------------------------------------------------------------------------ |
| `-version`        | `bool`   | `false` | Prints the current version of Mantrae and exits.                         |
| `-update`         | `bool`   | `false` | Updates Mantrae to the latest version. (Doesn't work inside a container) |
| `-reset-password` | `string` |         | Resets the admin password to the specified value.                        |
| `-reset-user`     | `string` | `admin` | Choose the username to reset the password for.                           |

## Environment Variables

Environment variables can be used to set up Mantrae and configure its settings. Below is a list of the supported environment variables.

### Core Configuration

| Variable         | Default | Description                                                              |
| ---------------- | ------- | ------------------------------------------------------------------------ |
| `SECRET`         |         | Secret key required for secure access (required)                         |
| `ADMIN_PASSWORD` |         | Admin password for the web interface (will be auto-generated if not set) |
| `LOG_LEVEL`      | `info`  | Logging level (debug, info, warn, error)                                 |
| `LOG_FORMAT`     | `text`  | Log format (text, json)                                                  |

### Important Notes

- **SECRET** is a required environment variable and must be set; otherwise, the application will not start.
- Database migrations are automatically applied on startup.
- In production, always use a strong secret and HTTPS.
