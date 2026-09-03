package domain

import (
	"container/list"
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"gorm.io/gorm"
)

type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusInReview   TaskStatus = "in_review"
	StatusDone       TaskStatus = "done"
)

func (s TaskStatus) Valid() bool {
	switch s {
	case StatusTodo, StatusInProgress, StatusInReview, StatusDone:
		return true
	}

	return false
}

func (s TaskStatus) CanTransitionTo(next TaskStatus) bool {
	allowed := map[TaskStatus][]TaskStatus{
		StatusTodo:       {StatusInProgress},
		StatusInProgress: {StatusTodo, StatusInReview},
		StatusInReview:   {StatusInProgress, StatusDone},
		StatusDone:       {StatusTodo},
	}

	for _, candidate := range allowed[s] {
		if candidate == next {
			return true
		}
	}

	return false
}

type TaskPriority string

const (
	PriorityLow    TaskPriority = "low"
	PriorityMedium TaskPriority = "medium"
	PriorityHigh   TaskPriority = "high"
	PriorityUrgent TaskPriority = "urgent"
)

func (p TaskPriority) Valid() bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh, PriorityUrgent:
		return true
	}

	return false
}

type Task struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	Key         string         `json:"key" gorm:"type:varchar(20);uniqueIndex;not null"`
	Title       string         `json:"title" gorm:"type:varchar(200):not nul"`
	Description string         `json:"description" gorm:"type:text"`
	Status      TaskStatus     `json:"status" gorm:"type:varchar(20);not null;defualt:todo;index;"`
	Priority    TaskPriority   `json:"priotiry" gorm:"type:varchar(20);not null;default:medium;index"`
	ProjectID   uuid.UUID      `json:"project_id" gorm:"type:uuid;not null;index"`
	Project     *Project       `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
	ReporterID  uuid.UUID      `json:"reporter_id" gorm:"type:uuid;not null;index"`
	Reporter    *User          `json:"reporter,omitempty" gorm:"foreignKey:ReporterID"`
	AssigneeID  *uuid.UUID     `json:"assignee_id" gorm:"type:uuid;index"`
	Assignee    *User          `json:"assignee,omitempty" gorm:"foreignKey:AssigneeID"`
	DueDate     *time.Time     `json:"due_date"`
	EstimateHrs *float64       `json:"estimate_hours"`
	Comments    []Comment      `json:"comments,omitempty" gorm:"foreignKey:TaskID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (t *Task) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type TaskFilter struct {
	ProjectID  *uuid.UUID
	AssigneeID *uuid.UUID
	ReporterID *uuid.UUID
	Status     *TaskStatus
	Priority   *TaskPriority
	Search     string
	SortBy     string // created_at | updated_at
	SortOrder  string // asc | desc
	Page       int
	Limit      int
}

func (f *TaskFilter) Normalize() {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 {
		f.Limit = 20
	}
	if f.Limit > 100 {
		f.Limit = 100
	}

	switch f.SortBy {
	case "created_at", "updated_at", "priority", "due_date", "status":
	default:
		f.SortBy = "created_at"
	}
	if f.SortOrder != "asc" {
		f.SortOrder = "desc"
	}
}

type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	FindByID(ctx context.Context, id uuid.UUID) (*Task, error)
	List(ctx context.Context, filter TaskFilter) ([]Task, int64, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id uuid.UUID) error

	AddComment(ctx context.Context, comment *Comment) error
	ListComments(ctx context.Context, taskID uuid.UUID) ([]Comment, error)
}