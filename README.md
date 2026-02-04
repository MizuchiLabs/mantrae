<p align="center">
<img src="./.github/logo.svg" width="80">
<br><br>
<img alt="GitHub Tag" src="https://img.shields.io/github/v/tag/MizuchiLabs/mantrae?label=Version">
<img alt="GitHub License" src="https://img.shields.io/github/license/MizuchiLabs/mantrae">
<img alt="GitHub Issues or Pull Requests" src="https://img.shields.io/github/issues/MizuchiLabs/mantrae">
</p>

# Mantræ

**Mantræ** is a web-based configuration manager for Traefik's dynamic configuration file. It provides a clean, intuitive interface to manage your routers, middleware, and services without editing YAML or TOML files manually.

> **Important**: Mantræ is **not** a dashboard for Traefik. It operates independently and does not monitor Traefik's status. Instead, Traefik connects to Mantræ to fetch its dynamic configuration.

> **Note**: The Mantræ agent has moved to a separate repository: [mantraed](https://github.com/MizuchiLabs/mantraed). The agent is now called `mantraed` (Mantræ daemon) and has its own container image at `ghcr.io/mizuchilabs/mantraed`. The old `mantrae-agent` image will remain available for compatibility but is deprecated.

## Features

- **Clean Interface**: Manage your Traefik configuration through a simple web UI
- **Router Management**: Create and configure routers with custom rules, entrypoints, and middleware
- **Middleware Support**: Add rate limiting, authentication, headers, and other middleware
- **Agent Mode**: Label your containers with standard Traefik labels and let the agent automatically sync them
- **DNS Integration**: Automatic DNS record management for Cloudflare, PowerDNS, Technitium and PiHole DNS

## 🚧 Development Status

This project is in active development and not yet production-ready. Expect breaking changes before the first stable release.

## Quick Start

### Installation

**Using the install script:**

```bash
curl -fsSL https://raw.githubusercontent.com/mizuchilabs/mantrae/main/install.sh | sh
```

**Manual installation:**
Download the latest release from [releases](https://github.com/mizuchilabs/mantrae/releases) and extract to `~/.local/bin`.

**Docker (recommended for production):**
See the [documentation](https://mantrae.pages.dev) for Docker setup instructions.

### Usage

```bash
# Start the server
mantrae

# Display version
mantrae --version

# Reset admin password
mantrae reset --password newpassword

# Reset password for a specific user
mantrae reset --user username --password newpassword
```

## Command Reference

| Command             | Description                            |
| ------------------- | -------------------------------------- |
| `mantrae`           | Start the Mantræ server                |
| `mantrae reset`     | Reset user password (admin by default) |
| `mantrae --version` | Display version information            |

### Flags

| Flag         | Aliases | Default | Description                    |
| ------------ | ------- | ------- | ------------------------------ |
| `--version`  | `-v`    |         | Display version and exit       |
| `--password` | `-p`    |         | New password (used with reset) |
| `--user`     | `-u`    | `admin` | Username for password reset    |

## Documentation

Full documentation is available [here](https://mizuchilabs.github.io/mantrae/)

## Screenshot

![Dashboard](./.github/screenshots/dash.png "Dashboard")

## Contributing

Contributions are welcome! Feel free to submit issues, fork the repository, and create pull requests.

## License

MIT License - See [LICENSE](LICENSE)

## Acknowledgements

- [**Traefik**](https://traefik.io/) - The powerful reverse proxy that this project manages
- [**Nginx Proxy Manager**](https://github.com/NginxProxyManager/nginx-proxy-manager) - Inspiration for the UI approach
- [**External-DNS**](https://github.com/kubernetes-sigs/external-dns) - Inspiration for DNS management
