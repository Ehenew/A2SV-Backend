package controllers

import (
	"a2sv-backend/task_manager/data"
	"a2sv-backend/task_manager/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetTasks
func GetTasks(c *gin.Context) {
	tasks, err := data.GetAllTasks()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving tasks"})
		return
	}
	c.IndentedJSON(http.StatusOK, tasks)
}

// GetTaskByID 
func GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	task, err := data.GetTaskByID(id)

	if err != nil {
		if err.Error() == "task not found" {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid ID format"})
		}
		return
	}
	c.IndentedJSON(http.StatusOK, task)
}

// AddTask with POST
func AddTask(c *gin.Context) {
	var newTask models.Task

	if err := c.BindJSON(&newTask); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	id, err := data.AddTask(newTask)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error creating task"})
		return
	}
	newTask.ID = id
	c.IndentedJSON(http.StatusCreated, newTask)
}

// UpdateTask 
func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var updatedTask models.Task

	if err := c.BindJSON(&updatedTask); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	err := data.UpdateTask(id, updatedTask)
	if err != nil {
		if err.Error() == "task not found" {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid ID format or update error"})
		}
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Task updated"})
}

// DeleteTask
func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	err := data.DeleteTask(id)

	if err != nil {
		if err.Error() == "task not found" {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid ID format or delete error"})
		}
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Task deleted"})
}