<p align="center">
<img src="./web/src/lib/images/logo.svg" width="80">
<br><br>
<img alt="GitHub Tag" src="https://img.shields.io/github/v/tag/MizuchiLabs/mantrae?label=Version">
<img alt="GitHub License" src="https://img.shields.io/github/license/MizuchiLabs/mantrae">
<img alt="GitHub Issues or Pull Requests" src="https://img.shields.io/github/issues/MizuchiLabs/mantrae">
</p>

# Mantræ

**Mantræ** is a user-friendly web interface designed to simplify the management of Traefik's dynamic configurations. Similar to Nginx Proxy Manager (NPM), this application allows you to manage your dynamic Traefik configuration from the comfort of a simple web ui.

## Features

- **Simplified UI**: A clean and intuitive interface that keeps the complexity to a minimum.
- **Router Configuration**: Create and manage Traefik routers with custom rules, entrypoints, and middleware configurations.
- **Middleware Management**: Add middlewares to your routers, including rate limiting, authentication, and more.
- **Service Status**: Monitor the status of your services and see their health information.
- **DNS Providers**: Support for multiple DNS providers (currently Cloudflare and PowerDNS) for automatic DNS record updates.

## Getting Started

### Prerequisites

- **Traefik**: Ensure you have Traefik set up with a static configuration file that defines your entrypoints, certificates, and any other necessary static settings.
- **Docker**: Optionally, use Docker for easier deployment. (See [Docker Compose](#docker-compose) for an example)

### Installation

1. Generate a random secret with `openssl rand -hex 32`

1. Use docker `docker run --name mantrae -e SECRET=<secret> -d -p 3000:3000 ghcr.io/mizuchilabs/mantrae:latest`

1. Or use the example docker-compose.yml file to run mantrae and traefik together

1. **Access the Web UI**:
   Open your web browser and navigate to `http://localhost:3000`

1. Use the admin password, which will be printed in the logs after the first start. You won't be able to access the password afterwards!

## Usage

### CLI Options:

1. `-version`: Print the version.
1. `-config`: Specify the path where the database should be stored. (Default: `./mantrae.db`)
1. `-port`: Set the port to listen on. (Default: `3000`)
1. `-update`: Update mantrae to the latest version. (Will not work inside docker container!)
1. `-reset`: Reset the admin password. Will generate a new one and print it to the console.
1. `-auth`: Enable basic auth for the profile endpoint. (Default: `false`). This will protect the `/api/<profilename>` endpoint using BasicAuth. If you enable this you need to set the correct headers in Traefik so it can fetch the data. Example default mantrae user is `admin` and password is `supersecret`, use `echo -n "admin:supersecret" | base64` to generate the Credentials. After that you can configure the http provider in Traefik: `providers.http.headers.Authorization=Basic YWRtaW46c3VwZXJzZWNyZXQ=`

**Optionally you can set the default profile by providing the flags below:**

1. `-url`: Set the URL of the Traefik instance. (e.g.: `http://localhost:8080`)
1. `-username`: Set the username for basic auth, if your Traefik instance uses basic auth. (e.g.: `admin`)
1. `-password`: Set the password for basic auth. (e.g.: `supersecret`)

**Set the secret!**

Before starting mantrae it will look for a `SECRET` environment variable. Depending if you use container or binary you need to set it manually, e.g. `export SECRET=<somelongstring>`

### Managing Routers

1. Navigate to the "Routers" section in the web UI.
1. Click "Create Router" to define a new router.
1. Assign a name, service, and rule to the router. Optionally, set entrypoints, middlewares, and other settings.
1. Save your router to apply the changes.

### Managing Middlewares

1. Open the "Middlewares" section.
1. Create new middleware by defining its type and associated settings.
1. Save the middleware and attach it to your routers as needed.

### Managing DNS providers

1. Open the "DNS Providers" section.
1. Create new DNS provider by defining its type and associated settings. (Currently only Cloudflare and PowerDNS are supported)
1. You can choose a default DNS provider for newly created routers.

### Settings

1. Various settings will be added here in the future. Currently you can use it to backup and restore your database.

### Screenshots

![Routers](./.github/screenshots/routers.png "Routers")
![Middleware](./.github/screenshots/middlewares.png "Middleware")
![CommandCenter](./.github/screenshots/command_center.png "CommandCenter")

## FAQ

### Why would use this? Traefik already has configuration discovery.

Yes, Traefik has amazing configuration discovery capabilities for various providers (Docker, Kubernetes, etc.). But for all those times you can't use these features (e.g. multiple machines not connected via Docker Swarm or Kubernetes) you have to use the file provider. Mantrae helps you with that and adds additional automation features like managing DNS records as well, similar to external-dns for Kubernetes.

### I want to use multiple DNS providers of the same type (e.g. multiple cloudflare accounts), how do I do that?

Traefik doesn't support multiple DNS Challenge providers, so you have to use CNAME records to manage multiple accounts.
E.g. if you have a domain `example.com` on account "Foo" and a domain `example.org` on account "Bar", you can add the API Key for account "Foo" normally, but to get letsencrypt certificates for `example.org` you need add a CNAME record for `example.org` with these values:

- Type: `CNAME`
- Name: `_acme-challenge.example.org`
- Target: `_acme-challenge.example.com`

Now you can request certificates for `sub.example.org` as well.

### Example Static Configuration

Below is a simple example of a Traefik static configuration in `traefik.yml`:

```yaml
entryPoints:
  web:
    address: ":80"
  websecure:
    address: ":443"

certificatesResolvers:
  myresolver:
    acme:
      email: your-email@example.com
      storage: /acme.json
      httpChallenge:
        entryPoint: web

providers:
  http:
    endpoint: "<endpoint where mantrae is running>"
```

## Roadmap

- Add support for more DNS providers
- Backup scheduler (local or S3)
- Support multiple database providers (currently only supports SQLite)
- User management and permissions system

## Contributing

Contributions are welcome! Please feel free to submit issues, fork the repository, and create pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [**Traefik**](https://traefik.io/): For providing the powerful reverse proxy that powers this application.
- [**Nginx Proxy Manager**](https://github.com/NginxProxyManager/nginx-proxy-manager): For inspiration on building a simple and effective web UI for managing reverse proxies.
- [**External-DNS**](https://github.com/kubernetes-sigs/external-dns): For inspiration on building a simple and effective DNS management system.
