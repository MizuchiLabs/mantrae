---
sidebar_position: 4
---

# Environment Variables

Mantrae provides several command-line flags and environment variables to
configure the application. This guide details each option and its purpose.

## Command-Line Arguments

You can use the following flags to customize the behavior of Mantrae:

| Flag       | Type   | Default | Description                                                              |
| ---------- | ------ | ------- | ------------------------------------------------------------------------ |
| `-version` | `bool` | `false` | Prints the current version of Mantrae and exits.                         |
| `-update`  | `bool` | `false` | Updates Mantrae to the latest version. (Doesn't work inside a container) |

## Environment Variables

Environment variables can be used to set up Mantrae and configure its settings.
Below is a list of the supported environment variables.

### Core Configuration

| Variable | Default | Description                                      |
| -------- | ------- | ------------------------------------------------ |
| `SECRET` |         | Secret key required for secure access. Required! |

### Server Configuration

| Variable            | Default            | Description                                                |
| ------------------- | ------------------ | ---------------------------------------------------------- |
| `SERVER_HOST`       | `0.0.0.0`          | Host address the server will bind to                       |
| `SERVER_PORT`       | `3000`             | Port which Mantrae will listen on                          |
| `SERVER_URL`        | `http://127.0.0.1` | The public URL of the Mantrae server for agent connections |
| `SERVER_BASIC_AUTH` | `false`            | Enables basic authentication for the Mantrae server        |

### Traefik Configuration

| Variable           | Default   | Description                     |
| ------------------ | --------- | ------------------------------- |
| `TRAEFIK_PROFILE`  | `default` | Traefik profile name            |
| `TRAEFIK_URL`      |           | Traefik API URL                 |
| `TRAEFIK_USERNAME` |           | Traefik authentication username |
| `TRAEFIK_PASSWORD` |           | Traefik authentication password |
| `TRAEFIK_TLS`      | `false`   | Enable TLS for Traefik          |

### Example Usage

To run Mantrae with custom environment variables:

```bash
export SECRET="your-secret-key"
export SERVER_PORT="4000"
export ADMIN_PASSWORD="secure-password"
./mantrae
```

### Important Notes

- **SECRET** is a required environment variable and must be set; otherwise, the
  application will not start.
- Set **SERVER_URL** to the publicly accessible URL of Mantrae to ensure agents
  can connect to it.
