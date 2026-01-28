package db

import (
	"github.com/fahmihidayah/go-api-orchestrator/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(config *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  config.DatabaseURL,
		PreferSimpleProtocol: true, // REQUIRED for Neon PostgreSQL - disables implicit prepared statements
	}), &gorm.Config{
		PrepareStmt: false, // Also disable GORM's prepared statement cache
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}
