package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sorayuth/task-manager-go/server/internal/domain"
	"gorm.io/gorm"
)

type projectRepository struct {
	db *gorm.DB
}

func NewProjectReposigory(db *gorm.DB) domain.ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(ctx context.Context, project *domain.Project) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(project).Error; err != nil {
			if isUniqueViolation(err) {
				return fmt.Errorf("Project key %s: %w", project.Key, domain.ErrConflict)
			}
			return fmt.Errorf("Create project: %w", err)
		}

		return tx.Exec(
			`INSERT INTO project_members (project_id, user_id) VALUES (?, ?)`,
			project.ID, project.OwnerID,
		).Error
	})
}

func (r *projectRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	var project domain.Project

	err := r.db.WithContext(ctx).
		Preload("Owner").
		First(&project, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Projet %s: %w", id, domain.ErrNotFound)
		}

		return nil, fmt.Errorf("Find project: %w", err)
	}

	return &project, nil
}

func (r *projectRepository) FindByKey(ctx context.Context, key string) (*domain.Project, error) {
	var project domain.Project

	err := r.db.WithContext(ctx).First(&project, "key = ?", key).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Project %s: %w", key, domain.ErrNotFound)
		}

		return nil, fmt.Errorf("Find project by key: %w", err)
	}

	return &project, nil
}

func (r *projectRepository) ListForUser(ctx context.Context, userID uuid.UUID) ([]domain.Project, error) {
	var projects []domain.Project

	err := r.db.WithContext(ctx).
		Joins("JOIN project_members pm ON pm.project_id = projects.id").
		Where("pm.user_id = ?", userID).
		Preload("Owner").
		Order("projects.created_at desc").
		Find(&projects).Error
	if err != nil {
		return nil, fmt.Errorf("List projects: %w", err)
	}

	return projects, nil
}

func (r *projectRepository) Update(ctx context.Context, project *domain.Project) error {
	if err := r.db.WithContext(ctx).Save(project).Error; err != nil {
		return fmt.Errorf("Failed to update project: %w", err)

	}

	return nil
}

func (r *projectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&domain.Project{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("Failed to delete project: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("Project %s: %w", id, domain.ErrNotFound)
	}

	return nil
}

func (r *projectRepository) IsMember(ctx context.Context, projectID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("project_members").
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("check membership: %w", err)
	}
	return count > 0, nil
}

func (r *projectRepository) AddMember(ctx context.Context, projectID, userID uuid.UUID) error {
	err := r.db.WithContext(ctx).Exec(
		`INSERT INTO project_members (project_id, user_id) VALUES (?, ?)
		 ON CONFLICT DO NOTHING`,
		projectID, userID,
	).Error
	if err != nil {
		return fmt.Errorf("add member: %w", err)
	}
	return nil
}

func (r *projectRepository) RemoveMember(ctx context.Context, projectID, userID uuid.UUID) error {
	err := r.db.WithContext(ctx).Exec(
		`DELETE FROM project_members WHERE project_id = ? AND user_id = ?`,
		projectID, userID,
	).Error
	if err != nil {
		return fmt.Errorf("remove member: %w", err)
	}
	return nil
}

func (r *projectRepository) NextTaskNumber(ctx context.Context, projectID uuid.UUID) (int, error) {
	var next int
	err := r.db.WithContext(ctx).Raw(
		`UPDATE projects
		    SET task_counter = task_counter + 1
		  WHERE id = ?
		RETURNING task_counter`,
		projectID,
	).Scan(&next).Error
	if err != nil {
		return 0, fmt.Errorf("next task number: %w", err)
	}
	if next == 0 {
		return 0, fmt.Errorf("project %s: %w", projectID, domain.ErrNotFound)
	}
	return next, nil
}
