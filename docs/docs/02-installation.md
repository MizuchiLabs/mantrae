---
sidebar_position: 2
---

# Installation

Follow these instructions to install and start using Mantrae with Traefik.

## Prerequisites

1. **Install Traefik**: Ensure you have a running instance of [Traefik](https://traefik.io/).
2. **Generate a Secret**: Create a secure, random secret key to use with Mantrae:
   ```bash
   openssl rand -hex 32
   ```
   Copy the generated secret as you’ll need it in the next steps.

## Installation Options

### Option 1: Download the Binary

1. Visit the [releases page](https://github.com/mizuchilabs/mantrae/releases) to download the latest Mantrae binary for your platform.
2. Run the binary with:
   ```bash
   export SECRET=<your_secret>
   ./mantrae
   ```

### Option 2: Use Docker

Run Mantrae in a Docker container:

```bash
docker run --name mantrae -e SECRET=<your_secret> -d -p 3000:3000 ghcr.io/mizuchilabs/mantrae:latest
```

### Option 3: Docker Compose (Mantrae + Traefik)

Use the provided `docker-compose.yml` file to deploy Mantrae and Traefik together.

1. Download the example [docker-compose.yml](https://github.com/mizuchilabs/mantrae/blob/main/docker-compose.yml).
2. Run the following command to start both services:
   ```bash
   docker-compose up -d
   ```

## Access the Web UI

Once Mantrae is running, open your browser and navigate to:

[http://localhost:3000](http://localhost:3000)

Log in using the user `admin` and the generated password, which will display on the first run.
