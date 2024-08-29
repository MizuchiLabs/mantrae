<img src="./web/src/lib/images/logo.png" width="80">

# Mantræ

**Mantræ** is a user-friendly web interface designed to simplify the management of Traefik's dynamic configurations. Similar to Nginx Proxy Manager (NPM), this application allows you to manage your dynamic Traefik configuration from the comfort of a simple web ui.

## Features

- **Domain Management**: Easily manage your domains and assign them to specific hosts via the web interface.
- **Router Configuration**: Create and manage Traefik routers with custom rules, entrypoints, and middleware configurations.
- **Middleware Management**: Add middlewares to your routers, including rate limiting, authentication, and more.
- **Service Status**: Monitor the status of your services and see their health information.
- **Simplified UI**: A clean and intuitive interface that keeps the complexity to a minimum.

## Getting Started

### Prerequisites

- **Traefik**: Ensure you have Traefik set up with a static configuration file that defines your entrypoints, certificates, and any other necessary static settings.
- **Docker**: Optionally, use Docker for easier deployment. (See [Docker Compose](#docker-compose) for an example)

### Installation

1. Download the latest release from the [releases page](https://github.com/MizuchiLabs/mantrae/releases)

1. Extract the downloaded file

1. Run the application `./mantrae`

1. **Access the Web UI**:
   Open your web browser and navigate to `http://localhost:3000`

1. Or use docker `docker run --name mantrae -d -p 3000:3000 ghcr.io/mizuchilabs/mantrae:latest`

1. You can also use the example docker-compose.yml file to run mantrae and traefik together

1. Use the admin password, which will be printed in the logs after the first start

## Usage

### Managing Routers

1. Navigate to the "Routers" section in the web UI.
1. Click "Create Router" to define a new router.
1. Assign a name, service, and rule to the router. Optionally, set entrypoints, middlewares, and other settings.
1. Save your router to apply the changes.

### Managing Middlewares

1. Open the "Middlewares" section.
1. Create new middleware by defining its type and associated settings.
1. Save the middleware and attach it to your routers as needed.

## Static Configuration

Please note that some aspects, such as Let's Encrypt certificates, need to be configured via Traefik's static configuration. Mantrae focuses solely on managing dynamic configurations like routers, services, and middlewares.

Also Traefik doesn't support multiple DNS Challenge providers, so you have to use CNAME records to manage multiple accounts.
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

## Contributing

Contributions are welcome! Please feel free to submit issues, fork the repository, and create pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

- **Traefik**: For providing the powerful reverse proxy that powers this application.
- **Nginx Proxy Manager**: For inspiration on building a simple and effective web UI for managing reverse proxies.
