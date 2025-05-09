package models

// Todo represents a todo item
type Todo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	Priority    int    `json:"priority"`
	DueDate     string `json:"due_date"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// Advanced features could include:
// - Categories/Tags
// - Subtasks
// - Recurring tasks
// - Attachments
