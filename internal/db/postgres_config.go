package db

import "github.com/albertojnk/stonks/internal/common"

// PostgresConfig contains the configurations for database.
type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
	SSL      string
}

// NewPostgresConfig load the configuration for database.
func NewPostgresConfig() *PostgresConfig {
	config := &PostgresConfig{
		Host:     common.GetEnv("DATABASE_HOST", "localhost"),
		Port:     common.GetEnv("DATABASE_PORT", "5444"),
		User:     common.GetEnv("DATABASE_USER", "postgres"),
		Password: common.GetEnv("DATABASE_PASSWORD", "postgres"),
		DbName:   common.GetEnv("DATABASE_DBNAME", "postgres"),
		SSL:      common.GetEnv("SSL", "disable"),
	}

	return config
}
