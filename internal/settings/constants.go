package settings

// Constants for setting keys
const (
	// General settings
	KeyServerURL = "server_url"

	// Backup settings
	KeyBackupEnabled  = "backup_enabled"
	KeyBackupInterval = "backup_interval"
	KeyBackupKeep     = "backup_keep"
	KeyBackupStorage  = "backup_storage_select"
	KeyBackupPath     = "backup_path"

	// S3 settings
	KeyS3Endpoint     = "s3_endpoint"
	KeyS3Bucket       = "s3_bucket"
	KeyS3Region       = "s3_region"
	KeyS3AccessKey    = "s3_access_key"
	KeyS3SecretKey    = "s3_secret_key"
	KeyS3UsePathStyle = "s3_use_path_style"

	// Email settings
	KeyEmailHost     = "email_host"
	KeyEmailPort     = "email_port"
	KeyEmailUser     = "email_user"
	KeyEmailPassword = "email_password"
	KeyEmailFrom     = "email_from"

	// OAuth settings
	KeyOIDCEnabled      = "oidc_enabled"
	KeyOIDCClientID     = "oidc_client_id"
	KeyOIDCClientSecret = "oidc_client_secret"
	KeyOIDCProviderName = "oidc_provider_name"
	KeyOIDCIssuerURL    = "oidc_issuer_url"
	KeyOIDCScopes       = "oidc_scopes"
	KeyOIDCPKCE         = "oidc_pkce"

	// Agent settings
	KeyAgentCleanupEnabled  = "agent_cleanup_enabled"
	KeyAgentCleanupInterval = "agent_cleanup_interval"

	// Background jobs settings
	KeyTraefikSyncInterval = "traefik_sync_interval"
	KeyDNSSyncInterval     = "dns_sync_interval"
	KeyAgentSyncInterval   = "agent_sync_interval"
)
