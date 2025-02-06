---
sidebar_position: 3
---

# Agents

Agents extend the functionality of your Traefik setup by allowing you to collect Docker container information from machines that don’t have Traefik installed. This way, you can apply Traefik labels on containers across multiple hosts, and the agent will sync this data with your main Traefik instance.

## How Agents Work

An agent is a standalone binary that runs on any machine where you want to collect container information. Each agent:

- Collects Docker container metadata, including Traefik labels.
- Communicates with the Mantrae server, sending back container info for unified management.
- Regularly renews its access token to ensure a secure, persistent connection.

## Setting Up an Agent

### Step 1: Set the Server Address

In the settings, specify the **server address** for your Mantrae server. This address must be accessible by the agent to ensure successful communication.

- **Example**: If Mantrae is hosted on a public IP or domain (e.g., `https://mantrae.example.com`), configure this as the server address so agents can connect reliably.

> **Note**: The agent will automatically renew its token at regular intervals, so you don’t need to worry about re-authenticating it manually.

### Step 2: Generate and Copy the Agent Token

1. In the Mantrae UI, navigate to the **Agents** tab.
2. Add a new agent by clicking the **Add Agent** button.
3. Copy the generated agent token by using the **Copy Token** button.

### Step 3: Run the Mantrae Agent

1. Download the agent binary for your platform. (Or start the agent in a container.)
2. Run the agent with the token using the following command:
   ```bash
   TOKEN=YOUR_TOKEN ./mantrae-agent
   ```
3. Ensure that the machine running the agent has Docker installed, as it will gather container details from the Docker daemon.

## Use Case: Traefik Labels on Remote Hosts

Once the agent is running, you can set Traefik labels on containers located on the agent’s host machine as you normally would for Traefik. The agent collects these labels and sends them to the Mantrae server, where they are applied to the main Traefik instance.

This setup allows you to centralize routing configurations across multiple hosts without installing Traefik on each machine.

---

Using agents, you can scale your Traefik configuration effortlessly, managing containers across multiple machines from a single, centralized Mantrae server.
