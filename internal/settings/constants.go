package settings

// Constants for setting keys
const (
	// General settings
	KeyServerURL = "server_url"
	KeyStorage   = "storage_select"

	// Backup settings
	KeyBackupEnabled  = "backup_enabled"
	KeyBackupInterval = "backup_interval"
	KeyBackupKeep     = "backup_keep"

	// S3 settings
	KeyS3Endpoint     = "s3_endpoint"
	KeyS3Bucket       = "s3_bucket"
	KeyS3Region       = "s3_region"
	KeyS3AccessKey    = "s3_access_key"
	KeyS3SecretKey    = "s3_secret_key" // #nosec G101
	KeyS3UsePathStyle = "s3_use_path_style"

	// Email settings
	KeyEmailHost     = "email_host"
	KeyEmailPort     = "email_port"
	KeyEmailUser     = "email_user"
	KeyEmailPassword = "email_password"
	KeyEmailFrom     = "email_from"

	// OIDC settings
	KeyOIDCEnabled          = "oidc_enabled"
	KeyOIDCClientID         = "oidc_client_id"
	KeyOIDCClientSecret     = "oidc_client_secret" // #nosec G101
	KeyOIDCProviderName     = "oidc_provider_name"
	KeyOIDCIssuerURL        = "oidc_issuer_url"
	KeyOIDCScopes           = "oidc_scopes"
	KeyOIDCPKCE             = "oidc_pkce"
	KeyPasswordLoginEnabled = "password_login_enabled"

	// Agent settings
	KeyAgentCleanupEnabled  = "agent_cleanup_enabled"
	KeyAgentCleanupInterval = "agent_cleanup_interval"

	// Background jobs settings
	KeyTraefikSyncInterval = "traefik_sync_interval"
	KeyDNSSyncInterval     = "dns_sync_interval"
	KeyAgentSyncInterval   = "agent_sync_interval"
)
