package database

import (
	"fmt"

	"github.com/sorayuth/task-manager-go/server/internal/domain"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`).Error; err != nil {
		return fmt.Errorf("create extension: %w", err)
	}

	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Project{},
		&domain.Task{},
		&domain.Comment{},
		&domain.RefreshToken{},
	); err != nil {
		return fmt.Errorf("Auth migrate: %w", err)
	}

	return nil
}
