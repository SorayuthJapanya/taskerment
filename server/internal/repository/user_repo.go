package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sorayuth/task-manager-go/server/internal/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx *context.Context, user *domain.User) error {
	if err := r.db.WithContext(*ctx).Create(user).Error; err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("username %s: %w", user.Username, domain.ErrConflict)
		}
		return fmt.Errorf("Create user: %w", err)
	}
	return nil
}

func (r *userRepository) FindByID(ctx *context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(*ctx).First(&user, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user %s: %w", id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("find user: %w", err)
	}

	return &user, nil
}

func (r *userRepository) FindByUsername(ctx *context.Context, username string) (*domain.User, error) {
	var user domain.User

	err := r.db.WithContext(*ctx).First(&user, "username = ?", username).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("username %s: %w", username, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("Find by username: %w", err)
	}

	return &user, nil
}

func (r *userRepository) ExistsUsername(ctx *context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(*ctx).Model(&domain.User{}).Where("username = ?", username).Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("count users: %w", err)
	}

	return count > 0, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("Update user: %w", err)
	}

	return nil
}
