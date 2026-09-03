package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	TokenHash string    `gorm:"type:varchar(64):uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null;index"`
	RevokedAt *time.Time
	CreatedAt time.Time
}

func (t *RefreshToken) IsUsable() bool {
	return t.RevokedAt == nil && time.Now().Before(t.ExpiresAt)
}

type TokenRepository interface {
	Store(ctx context.Context, token *RefreshToken) error
	FindByHash(ctx context.Context, hash string) (*RefreshToken, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeAll(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}
