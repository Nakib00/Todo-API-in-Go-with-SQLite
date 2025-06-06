package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Nakib00/Todo-API-in-Go-with-SQLite/internal/config"
	"github.com/Nakib00/Todo-API-in-Go-with-SQLite/internal/handlers"
	"github.com/Nakib00/Todo-API-in-Go-with-SQLite/internal/repository"
)

func main() {
	// Initialize configuration
	cfg := config.LoadConfig()
	// Initialize database
	db, err := repository.InitDB("") // Empty string as we're using the config inside InitDB
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create repository and handlers
	todoRepo := repository.NewTodoRepository(db)
	todoHandler := handlers.NewTodoHandler(todoRepo)
	// Set up router
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Add middleware (e.g., authentication, logging)
	// router.Use(middleware.AuthMiddleware())

	// Routes
	api := router.Group("/api/v1")
	{
		api.GET("/todos", todoHandler.GetTodos)
		api.GET("/todos/:id", todoHandler.GetTodo)
		api.POST("/todos", todoHandler.CreateTodo)
		api.PUT("/todos/:id", todoHandler.UpdateTodo)
		api.DELETE("/todos/:id", todoHandler.DeleteTodo)
		api.PATCH("/todos/:id/complete", todoHandler.MarkComplete)
		api.PATCH("/todos/:id/priority", todoHandler.UpdatePriority)
	}

	// Start server
	log.Printf("Server starting on port %s...", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
