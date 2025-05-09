package models

import (
	"time"
)

type Todo struct {
	ID          string    `json:"id"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	Priority    int       `json:"priority"` // 1-5 where 5 is highest
	DueDate     time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Advanced features could include:
// - Categories/Tags
// - Subtasks
// - Recurring tasks
// - Attachments