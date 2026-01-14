# Authentication

Mantrae supports local authentication and OIDC (OpenID Connect) single sign-on.

## Local Authentication

Default authentication method using username and password stored in Mantrae's database.

### Initial Setup

Set admin password during first run:

```bash
ADMIN_PASSWORD=your-secure-password
```

If not set, Mantrae generates a random password (check logs).

### Password Reset

Reset passwords via CLI:

```bash
# Reset admin password
mantrae reset --password newpassword

# Reset specific user
mantrae reset --user username --password newpassword
```

## OIDC Authentication

Configure single sign-on with an OIDC provider through the Settings page.

### Supported Providers

Any OIDC-compliant provider:

- Keycloak
- Authentik
- Auth0
- Okta
- Google Workspace
- Azure AD
- Custom OIDC providers

### Configuration

1. Navigate to Settings → Authentication
2. Enable OIDC
3. Configure provider details:
   - **Issuer URL**: Your OIDC provider's issuer endpoint (e.g., `https://auth.example.com/realms/master`)
   - **Client ID**: Application client ID from your provider
   - **Client Secret**: Application client secret (leave empty for PKCE)
   - **PKCE**: Enable for public clients without client secrets
4. Save configuration

The callback URL will be automatically set to: `https://your-mantrae-domain.com/oidc/callback`

### Provider Setup

Create an OIDC application in your provider with these settings:

- **Redirect URI**: `https://mantrae.example.com/oidc/callback`
- **Scopes**: `openid`, `profile`, `email`
- **Grant Type**: Authorization Code
- **Token Endpoint Auth Method**: Client Secret Post (or PKCE for public clients)

### User Provisioning

- First-time OIDC users are automatically created in Mantrae
- Users are matched by email address
- Usernames are derived from the `preferred_username` claim or email address
- Email verification is required

### Fallback Authentication

Local authentication remains available when OIDC is enabled. Use for:

- Emergency access if OIDC provider is unavailable
- Service accounts
- Initial admin setup

## Security Considerations

1. **HTTPS Required**: OIDC authentication requires HTTPS in production
2. **Email Verification**: Users must have verified emails in the OIDC provider
3. **Secure Secrets**: Store client secrets securely, never commit to version control
4. **Exact Redirect URIs**: Ensure redirect URLs match exactly in your provider configuration
5. **PKCE for Public Clients**: Enable PKCE if your application cannot securely store client secrets

:::tip
Enable OIDC for centralized user management and leverage your provider's security features like MFA and SSO.
:::
