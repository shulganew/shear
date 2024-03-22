package entities

// DTO for configuration from ENV parameters.
type ConfENV struct {
	Address    *string
	Response   *string
	BackupPath *string
	DSN        *string
	Pass       *string
	IsSequre   *bool
	JSONPath   *string
}
