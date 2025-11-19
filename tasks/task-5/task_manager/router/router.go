package router

import (
	"a2sv-backend/task_manager/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Task routes
	taskRoutes := router.Group("/tasks")
	{
		taskRoutes.GET("", controllers.GetTasks)
		taskRoutes.GET("/:id", controllers.GetTask)
		taskRoutes.POST("", controllers.AddTask)
		taskRoutes.PUT("/:id", controllers.UpdateTask)
		taskRoutes.DELETE("/:id", controllers.DeleteTask)
	}

	return router
}
