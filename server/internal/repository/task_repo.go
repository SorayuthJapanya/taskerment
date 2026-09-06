package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sorayuth/task-manager-go/server/internal/domain"
	"gorm.io/gorm"
)

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) domain.TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(ctx context.Context, task *domain.Task) error {
	if err := r.db.WithContext(ctx).Create(task).Error; err != nil {
		return fmt.Errorf("Create task failed: %w", err)
	}

	return nil
}

func (r *taskRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	var task domain.Task

	err := r.db.WithContext(ctx).
		Preload("Repoter").
		Preload("Assignee").
		Preload("Project").
		First(&task, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Task %s: %w", id, domain.ErrNotFound)
		}

		return nil, fmt.Errorf("Task not found: %w", err)
	}

	return &task, nil
}

func (r *taskRepository) List(ctx context.Context, f domain.TaskFilter) ([]domain.Task, int64, error) {

	// Normalize the filter
	f.Normalize()

	query := r.db.WithContext(ctx).Model(&domain.Task{})

	if f.ProjectID != nil {
		query = query.Where("project_id = ?", *f.ProjectID)
	}
	if f.AssigneeID != nil {
		query = query.Where("assignee_id = ?", *f.AssigneeID)
	}
	if f.ReporterID != nil {
		query = query.Where("reporter_id = ?", *f.ReporterID)
	}
	if f.Status != nil {
		query = query.Where("status = ?", *f.Status)
	}
	if f.Priority != nil {
		query = query.Where("priority = ?", *f.Priority)
	}
	if f.Search != "" {
		like := "%" + f.Search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ? OR key ILIKE ?", like, like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("Count task: %w", err)
	}

	order := fmt.Sprintf("%s %s", f.SortBy, f.SortOrder)

	var tasks []domain.Task
	err := query.
		Preload("Assignee").
		Preload("Reporter").
		Order(order).
		Limit(f.Limit).
		Offset((f.Page - 1) * f.Limit).
		Find((&tasks)).Error

	if err != nil {
		return nil, 0, fmt.Errorf("List tasks: %w", err)
	}

	return tasks, total, nil

}

func (r *taskRepository) Update(ctx context.Context, task *domain.Task) error {
	if err := r.db.WithContext(ctx).Save(task).Error; err != nil {
		return fmt.Errorf("Update task: %w", err)
	}

	return nil
}

func (r *taskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&domain.Task{}, "id = ?", id)

	if res.Error != nil {
		return fmt.Errorf("Delete task: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("tasl %s: %w", id, domain.ErrNotFound)
	}

	return nil
}

func (r *taskRepository) AddComment(ctx context.Context, comment *domain.Comment) error {
	if err := r.db.WithContext(ctx).Create(comment).Error; err != nil {
		return fmt.Errorf("Create comment: %w", err)
	}

	return nil
}

func (r *taskRepository) ListComments(ctx context.Context, taskID uuid.UUID) ([]domain.Comment, error) {
	var comments []domain.Comment

	err := r.db.WithContext(ctx).Preload("Author").Where("task_id = ?", taskID).Order("created_at asc").Find(&comments).Error
	if err != nil {
		return nil, fmt.Errorf("List comment: %w", err)
	}

	return comments, nil
}
