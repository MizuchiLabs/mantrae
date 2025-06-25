export const settingGroups = {
	general: {
		title: "General Settings",
		description: "Core configuration for connecting to backend services.",
		keys: [
			{
				key: "server_url",
				label: "Server URL",
				type: "text",
				description:
					"The base URL of your backend server, including protocol (e.g., https://example.com).",
			},
			{
				key: "storage_select",
				label: "Storage Type",
				type: "select",
				description:
					"Select the storage backend for backups (e.g., local, S3).",
			},
		],
	},
	backup: {
		title: "Backup Settings",
		description: "Control how and where the system creates database backups.",
		keys: [
			{
				key: "backup_enabled",
				label: "Enable Backups",
				type: "boolean",
				description: "Toggle automatic backups of the database.",
			},
			{
				key: "backup_interval",
				label: "Backup Interval",
				type: "duration",
				description: "How often backups should run (e.g., 24h, 1h30m).",
			},
			{
				key: "backup_keep",
				label: "Backups to Keep",
				type: "number",
				description: "The number of recent backups to retain.",
			},
		],
	},
	s3: {
		title: "S3 Storage Settings",
		description:
			"Configure connection to Amazon S3 or any compatible object storage.",
		keys: [
			{
				key: "s3_endpoint",
				label: "S3 Endpoint",
				type: "text",
				description:
					"Custom endpoint for your S3-compatible storage (e.g., https://s3.amazonaws.com).",
			},
			{
				key: "s3_bucket",
				label: "Bucket Name",
				type: "text",
				description: "The name of the S3 bucket used for storing backups.",
			},
			{
				key: "s3_region",
				label: "Region",
				type: "text",
				description: "AWS region of your bucket (e.g., us-east-1).",
			},
			{
				key: "s3_access_key",
				label: "Access Key",
				type: "text",
				description: "Access key ID for your S3 storage credentials.",
			},
			{
				key: "s3_secret_key",
				label: "Secret Key",
				type: "text",
				description: "Secret access key for S3 storage. Keep it safe.",
			},
			{
				key: "s3_use_path_style",
				label: "Use Path Style",
				type: "boolean",
				description:
					"Enable if your S3 provider requires path-style URLs (instead of virtual-host style).",
			},
		],
	},
	email: {
		title: "Email Settings",
		description: "Set up SMTP for sending system and user notification emails.",
		keys: [
			{
				key: "email_host",
				label: "SMTP Host",
				type: "text",
				description: "Hostname of the SMTP server (e.g., smtp.mailgun.org).",
			},
			{
				key: "email_port",
				label: "SMTP Port",
				type: "number",
				description:
					"Port used to connect to the SMTP server (e.g., 587 or 465).",
			},
			{
				key: "email_user",
				label: "Username",
				type: "text",
				description: "Login username for the SMTP server.",
			},
			{
				key: "email_password",
				label: "Password",
				type: "password",
				description: "SMTP password or app-specific token.",
			},
			{
				key: "email_from",
				label: "From Address",
				type: "text",
				description: "Email address to use as the sender in outgoing emails.",
			},
		],
	},
	oauth: {
		title: "OIDC Authentication",
		description:
			"Configure OpenID Connect for secure single sign-on (SSO) and identity management.",
		keys: [
			{
				key: "oidc_enabled",
				label: "Enable OIDC",
				type: "boolean",
				description: "Turn on OpenID Connect authentication.",
			},
			{
				key: "password_login_enabled",
				label: "Password Login",
				type: "boolean",
				description:
					"Force users to log in only via OIDC when disabled (no local passwords).",
			},
			{
				key: "oidc_client_id",
				label: "Client ID",
				type: "text",
				description: "OIDC client ID issued by your identity provider.",
			},
			{
				key: "oidc_client_secret",
				label: "Client Secret",
				type: "password",
				description: "OIDC client secret issued by your identity provider.",
			},
			{
				key: "oidc_issuer_url",
				label: "Issuer URL",
				type: "text",
				description:
					"URL of the OIDC provider (e.g., https://auth.example.com).",
			},
			{
				key: "oidc_provider_name",
				label: "Provider Name",
				type: "text",
				description:
					"Friendly name shown in the login UI for the OIDC provider.",
			},
			{
				key: "oidc_scopes",
				label: "OIDC Scopes",
				type: "text",
				description:
					"Requested scopes (space-separated) for the OIDC flow (e.g., openid email profile).",
			},
			{
				key: "oidc_pkce",
				label: "Use PKCE",
				type: "boolean",
				description:
					"Enable PKCE (Proof Key for Code Exchange) for better security.",
			},
		],
	},
	agents: {
		title: "Agent Configuration",
		description: "Manage automated cleanup tasks for connected agents.",
		keys: [
			{
				key: "agent_cleanup_enabled",
				label: "Enable Cleanup",
				type: "boolean",
				description:
					"Automatically remove stale or offline agents on a schedule.",
			},
			{
				key: "agent_cleanup_interval",
				label: "Cleanup Interval",
				type: "duration",
				description: "Frequency of cleanup jobs (e.g., 1h, 24h).",
			},
		],
	},
};

// Storage types for the select dropdown
export const storageTypes = [
	{ value: "local", label: "Local Storage" },
	{ value: "s3", label: "S3 Storage" },
];
