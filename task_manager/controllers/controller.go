package controllers

import (
	"a2sv-backend/task_manager/data"
	"a2sv-backend/task_manager/models"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtSecret []byte

func init() {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		s = "your_jwt_secret_key"
	}
	jwtSecret = []byte(s)
}

// RegisterUser
func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := data.RegisterUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// LoginUser
func LoginUser(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := data.LoginUser(credentials.Username, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID.Hex(),
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "user_id": user.ID.Hex()})
}

// GetMe returns current user info
func GetMe(c *gin.Context) {
	role, _ := c.Get("role")
	username, _ := c.Get("username")
	c.JSON(http.StatusOK, gin.H{"username": username, "role": role})
}

// PromoteUser handles promoting a user to admin
func PromoteUser(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := data.PromoteUser(req.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User promoted to admin"})
}

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
