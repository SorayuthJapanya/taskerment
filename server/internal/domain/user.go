package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID      `json:"id"         gorm:"type:uuid;primaryKey"`
	Username     string         `json:"username" gorm:"type:varchar(100);uniqueIndex;not null"`
	Name         string         `json:"name" gorm:"type:varchar(100);not null"`
	PasswordHash string         `json:"-" gorm:"type:varchar(255);not null"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-"          gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

type UserRepository interface {
	Create(ctx *context.Context, user *User) error
	FindByID(ctx *context.Context, id uuid.UUID) (*User, error)
	FindByUsername(ctx *context.Context, username string) (*User, error)
	ExistsUsername(ctx *context.Context, username string) (bool, error)
	Update(ctx context.Context, user *User) error
}