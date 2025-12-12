package routers

import (
	"a2sv-backend/task_manager_v3/Delivery/controllers"
	"a2sv-backend/task_manager_v3/Infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRouter(taskController *controllers.TaskController, userController *controllers.UserController, jwtService Infrastructure.JWTService) *gin.Engine {
	r := gin.Default()

	// Public routes
	r.POST("/register", userController.Register)
	r.POST("/login", userController.Login)

	// Protected routes
	protected := r.Group("/")
	protected.Use(Infrastructure.AuthMiddleware(jwtService))
	{
		protected.GET("/me", userController.GetProfile)
		protected.GET("/tasks", taskController.GetAllTasks)
		protected.GET("/tasks/:id", taskController.GetTaskByID)

		// Admin routes
		admin := protected.Group("/")
		admin.Use(Infrastructure.AdminMiddleware())
		{
			admin.POST("/tasks", taskController.CreateTask)
			admin.PUT("/tasks/:id", taskController.UpdateTask)
			admin.DELETE("/tasks/:id", taskController.DeleteTask)
			admin.POST("/promote", userController.PromoteUser)
		}
	}

	return r
}
