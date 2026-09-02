package domain

import "github.com/google/uuid"

type Task struct {
	ID uuid.UUID
	Title string
	Description string
	Status bool
	Priority string
	ReporterId uuid.UUID
}