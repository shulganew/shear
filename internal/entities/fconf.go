package entities

// DTO for configuration from Flag parameters.
type ConfFlag struct {
	Address    *string
	Response   *string
	BackupPath *string
	DSN        *string
	Pass       *string
	IsBackup   *bool
	Pprof      *bool
	IsSequre   *bool
	JSONPath   *string
}
