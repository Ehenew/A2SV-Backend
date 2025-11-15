package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Task represents a task with its properties.
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"duedate"`
	Status      string    `json:"status"`
}

// Mock data for tasks
var tasks = []Task{
	{ID: "1", Title: "Task 1", Description: "First task", DueDate: time.Now(), Status: "Pending"},
	{ID: "2", Title: "Task 2", Description: "Second task", DueDate: time.Now().AddDate(0, 0, 1), Status: "In Progress"},
	{ID: "3", Title: "Task 3", Description: "Third task", DueDate: time.Now().AddDate(0, 0, 2), Status: "Completed"},
}

func main() {
	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Get all tasks
	router.GET("/tasks", getTasks)

	// Get a single task with id
	router.GET("/tasks/:id", getTask)

	// Updating existing task / PUT
	router.PUT("/tasks/:id", updateTask)

	// DELETE a task by id
	router.DELETE("/tasks/:id", deleteTask)

	// Post new task (POST)
	router.POST("/tasks", addTask)

	router.Run()
}

// Handlers
func getTasks(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, gin.H{"tasks": tasks})
}

func getTask(ctx *gin.Context) {
	id := ctx.Param("id")

	for _, task := range tasks {
		if task.ID == id {
			ctx.IndentedJSON(http.StatusOK, task)
			return
		}
	}

	ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

func updateTask(ctx *gin.Context) {
	id := ctx.Param("id")
	var updatedTask Task

	if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			// Update only the specified fields
			if updatedTask.Title != "" {
				tasks[i].Title = updatedTask.Title
			}

			if updatedTask.Description != "" {
				tasks[i].Description = updatedTask.Description
			}

			ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Task Updated"})
			return
		}
	}

	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Task not found"})
}

func deleteTask(ctx *gin.Context) {
	id := ctx.Param("id")

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Task removed"})
			return
		}
	}
	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Task not found"})
}

func addTask(ctx *gin.Context) {
	var newTask Task

	if err := ctx.ShouldBindJSON(&newTask); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tasks = append(tasks, newTask)
	ctx.IndentedJSON(http.StatusCreated, gin.H{"message": "Task Created"})
}
