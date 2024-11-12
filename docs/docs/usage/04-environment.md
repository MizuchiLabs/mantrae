---
sidebar_position: 4
---

# Environment

Mantrae provides several command-line flags and environment variables to configure the application. This guide details each option and its purpose.

## Command-Line Arguments

You can use the following flags to customize the behavior of Mantrae:

| Flag        | Type     | Default | Description                                                                |
| ----------- | -------- | ------- | -------------------------------------------------------------------------- |
| `-version`  | `bool`   | `false` | Prints the current version of Mantrae and exits.                           |
| `-port`     | `string` | `3000`  | Specifies the port Mantrae will listen on.                                 |
| `-url`      | `string` |         | Specifies the URL of the Traefik instance (e.g., `http://localhost:8080`). |
| `-username` | `string` |         | Username for authenticating with the Traefik instance.                     |
| `-password` | `string` |         | Password for authenticating with the Traefik instance.                     |
| `-update`   | `bool`   | `false` | Updates Mantrae to the latest version. (Doesn't work inside a container)   |
| `-reset`    | `bool`   | `false` | Resets the default admin password and outputs a new one.                   |

### Example Usage

To start Mantrae on port 3000 and connect to a Traefik instance with authentication:

```bash
./mantrae -port 3000 -url http://localhost:8080 -username admin -password secret
```

## Environment Variables

Environment variables can be used to set up Mantrae and configure its settings. Below is a list of the supported environment variables.

| Variable            | Default                 | Description                                                           |
| ------------------- | ----------------------- | --------------------------------------------------------------------- |
| `SECRET`            |                         | Secret key required for secure access. Required!                      |
| `PORT`              | `3000`                  | Port which Mantrae will listen on.                                    |
| `AGENT_PORT`        | `8090`                  | Listen port to accept connections from agents.                        |
| `SERVER_URL`        | `http://localhost:8090` | The public URL of the Mantrae server for agent connections.           |
| `ENABLE_BASIC_AUTH` | `false`                 | Enables basic authentication for the Mantrae server if set to `true`. |
| `ENABLE_AGENT`      | `true`                  | Enables the Mantrae agent functionality.                              |
| `CONFIG_DIR`        |                         | Directory path for Mantrae's configuration files.                     |
| `BACKUP_DIR`        | `backups`               | Directory for storing backups.                                        |
| `LOG_LEVEL`         | `info`                  | Log verbosity level. Options: `debug`, `info`, `warn`, `error`.       |

### Database Configuration

| Variable  | Default   | Description                                 |
| --------- | --------- | ------------------------------------------- |
| `DB_TYPE` | `sqlite`  | Database type. Supported options: `sqlite`. |
| `DB_NAME` | `mantrae` | Database name.                              |

### Example Usage

To run Mantrae with custom environment variables:

```bash
export SECRET="your-secret-key"
export PORT="4000"
./mantrae
```

### Important Notes

- **SECRET** is a required environment variable and must be set; otherwise, the application will not start.
- Set **SERVER_URL** to the publicly accessible URL of Mantrae to ensure agents can connect to it.
