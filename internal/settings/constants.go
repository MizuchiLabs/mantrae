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

	// Agent settings
	KeyAgentCleanupEnabled  = "agent_cleanup_enabled"
	KeyAgentCleanupInterval = "agent_cleanup_interval"
)
