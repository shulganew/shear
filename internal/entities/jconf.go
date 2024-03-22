package entities

// DTO for configuration from JSON parameters.
type ConfJSON struct {
	Address    *string `json:"server_address,omitempty"`
	Response   *string `json:"response_address,omitempty"`
	BackupPath *string `json:"file_storage_path,omitempty"`
	DSN        *string `json:"database_dsn,omitempty"`
	Pass       *string `json:"pass,omitempty"`
	IsBackup   *bool   `json:"enable_backup,omitempty"`
	Pprof      *bool   `json:"enable_pprof,omitempty"`
	IsSequre   *bool   `json:"enable_https,omitempty"`
}
