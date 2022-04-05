package db

import "github.com/albertojnk/stonks/internal/common"

// NewPostgresConfigMock load the configuration for database - USED FOR TEST.
func NewPostgresConfigMock() *PostgresConfig {
	config := &PostgresConfig{
		Host:     common.GetEnv("database_host_mock", "localhost"),
		Port:     common.GetEnv("database_port_mock", "5430"),
		User:     common.GetEnv("database_user", "postgres"),
		Password: common.GetEnv("database_password", ""),
		DbName:   common.GetEnv("database_dbname", "postgres"),
	}

	return config
}
