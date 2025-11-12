# Agents

Agents in Mantrae extend your Traefik setup by enabling automatic configuration discovery from remote Docker hosts. Instead of manually defining routers in the web UI, you can label your containers with standard Traefik labels, and the agent will automatically synchronize this information with your Mantrae server.

## How Agents Work

An agent is a lightweight binary that runs on any machine where you want to manage Docker containers. Each agent:

- Discovers Docker containers and their Traefik labels on the local machine
- Communicates with the Mantrae server to synchronize container information
- Automatically creates routers, services, and middlewares based on container labels
- Updates the Mantrae server when containers are added, removed, or changed

## Setting Up an Agent

### Step 1: Create an Agent in Mantrae

1. Log into the Mantrae web interface
2. Navigate to the "Agents" section
3. Click "Add Agent"

### Step 2: Copy the Agent Token

After creating the agent, you'll see a configuration section with:
- Agent Token (automatically generated)
- Docker Run command
- Docker Compose configuration

Copy the agent token for the next step.

### Step 3: Deploy the Agent

You can run the agent in several ways:

#### Option 1: Direct Binary Execution

1. Download the agent binary for your platform from the [releases page](https://github.com/mizuchilabs/mantrae/releases)
2. Run the agent with the token:
   ```bash
   TOKEN=YOUR_AGENT_TOKEN HOST=https://mantrae.example.com ./mantrae-agent
   ```

#### Option 2: Docker Run

Use the pre-generated Docker run command from the agent configuration page:
```bash
docker run -d \
   --name mantrae-agent \
   -e TOKEN=YOUR_AGENT_TOKEN \
   -e HOST=https://mantrae.example.com \
   -v /var/run/docker.sock:/var/run/docker.sock:ro \
   ghcr.io/mizuchilabs/mantrae-agent:latest
```

#### Option 3: Docker Compose

Use the pre-generated Docker Compose configuration:
```yaml
services:
  mantrae-agent:
    image: ghcr.io/mizuchilabs/mantrae-agent:latest
    container_name: mantrae-agent
    network_mode: host
    environment:
      - TOKEN=YOUR_AGENT_TOKEN
      - HOST=https://mantrae.example.com
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
    restart: unless-stopped
```

## Using Traefik Labels with Agents

Once the agent is running, you can label your containers as you normally would with Traefik. The agent will automatically collect these labels and synchronize them with the Mantrae server.

### Example HTTP Service

```yaml
services:
  whoami:
    image: containous/whoami:latest
    ports:
      - "80:80"
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.whoami.rule=Host(`whoami.example.com`)"
      - "traefik.http.routers.whoami.entrypoints=websecure"
      - "traefik.http.routers.whoami.tls=true"
      - "traefik.http.services.whoami.loadbalancer.server.port=80"
```

## Agent Configuration Options

The agent supports several environment variables for configuration:

| Variable | Description | Default |
|----------|-------------|---------|
| `TOKEN` | Agent token from Mantrae server | Required |
| `HOST` | Mantrae server URL | http://localhost:3000 |

## Network Configuration

Agents can automatically detect their network configuration:

1. **Public IP Detection**: Agents will automatically detect their public IP address
2. **Private Network**: For internal networks, agents can use their private IP
3. **Manual Override**: You can manually specify the IP address that should be used for services

This information is visible in the agent details page in the Mantrae web interface.

## Monitoring Agent Status

In the Mantrae web interface, you can:

- View the status of all agents (online/offline)
- See when each agent last synchronized
- View the containers discovered by each agent
- Rotate agent tokens for security

## Security Considerations

- Each agent has a unique token that should be kept secret
- Tokens can be rotated at any time from the web interface
- Agents only need read-only access to the Docker socket
