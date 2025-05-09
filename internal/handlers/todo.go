package handlers

import (
	"net/http"
	"strconv"

	"github.com/Nakib00/Todo-API-in-Go-with-SQLite/internal/models"
	"github.com/Nakib00/Todo-API-in-Go-with-SQLite/internal/repository"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	repo repository.TodoRepository
}

func NewTodoHandler(repo repository.TodoRepository) *TodoHandler {
	return &TodoHandler{repo: repo}
}

func (h *TodoHandler) GetTodos(c *gin.Context) {
	todos, err := h.repo.GetTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) GetTodo(c *gin.Context) {
	id := c.Param("id")
	todo, err := h.repo.GetTodo(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdTodo, err := h.repo.CreateTodo(todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdTodo)
}

func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	id := c.Param("id")

	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedTodo, err := h.repo.UpdateTodo(id, todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedTodo)
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id := c.Param("id")

	if err := h.repo.DeleteTodo(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TodoHandler) MarkComplete(c *gin.Context) {
	id := c.Param("id")

	if err := h.repo.MarkComplete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *TodoHandler) UpdatePriority(c *gin.Context) {
	id := c.Param("id")

	priority, err := strconv.Atoi(c.Query("priority"))
	if err != nil || priority < 1 || priority > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Priority must be between 1 and 5"})
		return
	}

	if err := h.repo.UpdatePriority(id, priority); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
