package repository

import (
	"database/sql"
	"time"

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

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id VARCHAR(36) PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			completed BOOLEAN DEFAULT FALSE,
			priority INT DEFAULT 3,
			due_date DATETIME,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (r *todoRepository) GetTodos() ([]models.Todo, error) {
	rows, err := r.db.Query(`
		SELECT 
			id, 
			title, 
			description, 
			completed, 
			priority,
			DATE_FORMAT(due_date, '%Y-%m-%dT%TZ') as due_date,
			DATE_FORMAT(created_at, '%Y-%m-%dT%TZ') as created_at,
			DATE_FORMAT(updated_at, '%Y-%m-%dT%TZ') as updated_at
		FROM todos
	`)
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

	err := r.db.QueryRow(`
		SELECT 
			id, 
			title, 
			description, 
			completed, 
			priority,
			DATE_FORMAT(due_date, '%Y-%m-%dT%TZ') as due_date,
			DATE_FORMAT(created_at, '%Y-%m-%dT%TZ') as created_at,
			DATE_FORMAT(updated_at, '%Y-%m-%dT%TZ') as updated_at
		FROM todos WHERE id = ?`, id).
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
	now := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	todo.CreatedAt = now
	todo.UpdatedAt = now

	if todo.DueDate == "" {
		todo.DueDate = now
	}

	_, err := r.db.Exec(`
		INSERT INTO todos (
			id, title, description, completed, priority, 
			due_date, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, 
			STR_TO_DATE(?, '%Y-%m-%dT%TZ'),
			STR_TO_DATE(?, '%Y-%m-%dT%TZ'),
			STR_TO_DATE(?, '%Y-%m-%dT%TZ'))`,
		todo.ID, todo.Title, todo.Description, todo.Completed, todo.Priority,
		todo.DueDate, todo.CreatedAt, todo.UpdatedAt,
	)
	if err != nil {
		return models.Todo{}, err
	}

	return todo, nil
}

func (r *todoRepository) UpdateTodo(id string, todo models.Todo) (models.Todo, error) {
	todo.UpdatedAt = time.Now().UTC().Format("2006-01-02T15:04:05Z")

	_, err := r.db.Exec(`
		UPDATE todos SET 
			title = ?, 
			description = ?, 
			completed = ?, 
			priority = ?, 
			due_date = STR_TO_DATE(?, '%Y-%m-%dT%TZ'),
			updated_at = STR_TO_DATE(?, '%Y-%m-%dT%TZ')
		WHERE id = ?`,
		todo.Title, todo.Description, todo.Completed, todo.Priority,
		todo.DueDate, todo.UpdatedAt, id,
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
	_, err := r.db.Exec("UPDATE todos SET completed = TRUE, updated_at = CURRENT_TIMESTAMP WHERE id = ?", id)
	return err
}

func (r *todoRepository) UpdatePriority(id string, priority int) error {
	_, err := r.db.Exec("UPDATE todos SET priority = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", priority, id)
	return err
}
