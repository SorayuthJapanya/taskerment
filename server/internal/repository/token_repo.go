package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sorayuth/task-manager-go/server/internal/domain"
	"gorm.io/gorm"
)

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) domain.TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) Store(ctx context.Context, token *domain.RefreshToken) error {
	if token.ID == uuid.Nil {
		token.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(token).Error; err != nil {
		return fmt.Errorf("Store refresh token: %w", err)
	}

	return nil
}

func (r *tokenRepository) FindByHash(ctx context.Context, hash string) (*domain.RefreshToken, error) {
	var token domain.RefreshToken
	err := r.db.WithContext(ctx).First(&token, "token_hash = ?", hash).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("refresh token: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("find refresh token: %w", err)
	}
	return &token, nil
}

func (r *tokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	now := time.Now().UTC()
	err := r.db.WithContext(ctx).Model(&domain.RefreshToken{}).
		Where("id = ? AND revoked_at IS NULL", id).
		Update("revoked_at", now).Error
	if err != nil {
		return fmt.Errorf("revoke token: %w", err)
	}
	return nil
}

func (r *tokenRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	now := time.Now().UTC()
	err := r.db.WithContext(ctx).Model(&domain.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now).Error
	if err != nil {
		return fmt.Errorf("revoke user tokens: %w", err)
	}
	return nil
}

func (r *tokenRepository) DeleteExpired(ctx context.Context) error {
	err := r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now().UTC()).
		Delete(&domain.RefreshToken{}).Error
	if err != nil {
		return fmt.Errorf("delete expired tokens: %w", err)
	}
	return nil
}
