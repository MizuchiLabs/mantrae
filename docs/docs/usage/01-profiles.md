---
sidebar_position: 1
---

# Profiles

Mantrae allows you to create and manage multiple profiles, making it easy to configure multiple different Traefik instances. Profiles can be created and selected from the dropdown menu at the top of the Mantrae interface.

## Creating a Profile

To create a new profile, follow these steps:

1. **Open the Profile Dropdown**: In the top navigation bar, open the profile dropdown and select **Create New Profile**.
2. **Set Profile Details**:
   - **Profile Name**: Enter a unique name for this profile (e.g., `default`, `staging`, `production`).
   - **Traefik Instance URL**: Provide the URL for the Traefik instance you want this profile to connect to.
     - For example, if you’re running Traefik on Docker on the same host, use `http://traefik:8080`.
   - **Basic Auth (Optional)**: If your Traefik instance requires basic authentication, provide the credentials here.

Once saved, this profile will serve as a space where you can manage routers and middlewares specifically for this Traefik instance.

## Using Profiles with Traefik

Each profile in Mantrae has its own API endpoint, allowing Traefik to fetch the correct configuration based on the active profile.

- **Example**: If your profile name is `default`, the corresponding API endpoint in Mantrae will be:
  `/api/default`

Configure Traefik to use this endpoint to pull configuration details specific to this profile. E.g. by using the static config:

```yaml
providers:
  http:
    endpoint: "http://mantrae:3000/api/default"
    # Optional if you enabled basic auth on mantrae itself
    headers:
      Authorization: Basic <base64 encoded username:password>
```

---

Profiles make it easy to work with multiple Traefik setups, ensuring configurations remain organized and accessible.
