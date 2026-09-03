package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	Key         string         `json:"key" gorm:"type:varchar(10);uniqueIndex;not null"`
	Name        string         `json:"name" gorm:"type:varchar(150);not null"`
	Description string         `json:"description" gorm:"type:text"`
	OwnerID     uuid.UUID      `json:"owner_id" gorm:"type:uuid;not null"`
	Owner       *User          `json:"owner,omitempty"`
	Members     []User         `json:"members,omitempty" gorm:"many2many:project_members"`
	TaskCounter int            `json:"-" gorm:"not null;default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-"           gorm:"index"`
}

func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}

	return  nil
}

type ProjectRepository interface {
	Create(ctx context.Context, project *Project) error
	FindByID(ctx context.Context, id uuid.UUIDs) (*Project, error)
	FindByKey(ctx context.Context, key string) (*Project, error)
	ListForUser(ctx context.Context, userID uuid.UUIDs) ([]Project, error)
	Update(ctx context.Context, project *Project) error
	Delete(ctx context.Context, id uuid.UUID) error

	IsMember(ctx context.Context, projectID, userID uuid.UUID) (bool, error)
	AddMember(ctx context.Context, projectID, userID uuid.UUID) error
	RemoveMember(ctx context.Context, projectID, userID uuid.UUID) error

	NextTaskNumber(ctx context.Context, projectID uuid.UUID) (int, error)
}