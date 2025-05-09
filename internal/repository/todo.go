package repository

import (
	"database/sql"

	"github.com/Nakib00/Todo-API-in-Go-with-SQLite/internal/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type TodoRepository interface {
	GetTodos() ([]models.Todo, error)
	GetTodo(id string) (models.Todo, error)
	CreateTodo(todo models.Todo) (models.Todo, error)
	UpdateTodo(id string, todo models.Todo) (models.Todo, error)
	DeleteTodo(id string) error
	MarkComplete(id string) error
	UpdatePriority(id string, priority int) error
}

type todoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) TodoRepository {
	return &todoRepository{db: db}
}

func InitDB(databasePath string) (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/todo-go")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`DROP TABLE IF EXISTS todos`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id VARCHAR(36) PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			completed BOOLEAN DEFAULT FALSE,
			priority INT DEFAULT 3,
			due_date VARCHAR(255),
			created_at VARCHAR(255),
			updated_at VARCHAR(255)
		)
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (r *todoRepository) GetTodos() ([]models.Todo, error) {
	rows, err := r.db.Query("SELECT id, title, description, completed, priority, due_date, created_at, updated_at FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Completed,
			&todo.Priority,
			&todo.DueDate,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func (r *todoRepository) GetTodo(id string) (models.Todo, error) {
	var todo models.Todo

	err := r.db.QueryRow("SELECT id, title, description, completed, priority, due_date, created_at, updated_at FROM todos WHERE id = ?", id).
		Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Completed,
			&todo.Priority,
			&todo.DueDate,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
	if err != nil {
		return models.Todo{}, err
	}

	return todo, nil
}

func (r *todoRepository) CreateTodo(todo models.Todo) (models.Todo, error) {
	todo.ID = uuid.New().String()
	now := "2025-05-10T00:00:00Z" // Current time in ISO format
	todo.CreatedAt = now
	todo.UpdatedAt = now

	if todo.DueDate == "" {
		todo.DueDate = now
	}

	_, err := r.db.Exec(
		"INSERT INTO todos (id, title, description, completed, priority, due_date, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		todo.ID, todo.Title, todo.Description, todo.Completed, todo.Priority, todo.DueDate, todo.CreatedAt, todo.UpdatedAt,
	)
	if err != nil {
		return models.Todo{}, err
	}

	return todo, nil
}

func (r *todoRepository) UpdateTodo(id string, todo models.Todo) (models.Todo, error) {
	todo.UpdatedAt = "2025-05-10T00:00:00Z" // Current time in ISO format

	_, err := r.db.Exec(
		"UPDATE todos SET title = ?, description = ?, completed = ?, priority = ?, due_date = ?, updated_at = ? WHERE id = ?",
		todo.Title, todo.Description, todo.Completed, todo.Priority, todo.DueDate, todo.UpdatedAt, id,
	)
	if err != nil {
		return models.Todo{}, err
	}

	return r.GetTodo(id)
}

func (r *todoRepository) DeleteTodo(id string) error {
	_, err := r.db.Exec("DELETE FROM todos WHERE id = ?", id)
	return err
}

func (r *todoRepository) MarkComplete(id string) error {
	now := "2025-05-10T00:00:00Z" // Current time in ISO format
	_, err := r.db.Exec("UPDATE todos SET completed = TRUE, updated_at = ? WHERE id = ?", now, id)
	return err
}

func (r *todoRepository) UpdatePriority(id string, priority int) error {
	now := "2025-05-10T00:00:00Z" // Current time in ISO format
	_, err := r.db.Exec("UPDATE todos SET priority = ?, updated_at = ? WHERE id = ?", priority, now, id)
	return err
}
