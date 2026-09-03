package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	TaskID    uuid.UUID      `json:"task_id" gorm:"type:uuid;not null;index"`
	AuthorID  uuid.UUID      `json:"author_id" gorm:"type:uuid;not null;index"`
	Author    *User          `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	Body      string         `json:"body" gorm:"type:text;not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}

	return nil
}
