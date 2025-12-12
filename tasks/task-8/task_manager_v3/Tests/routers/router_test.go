package routers_test

import (
	"a2sv-backend/task_manager_v3/Delivery/controllers"
	"a2sv-backend/task_manager_v3/Delivery/routers"
	"a2sv-backend/task_manager_v3/Tests/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockTaskUsecase := new(mocks.MockTaskUsecase)
	mockUserUsecase := new(mocks.MockUserUsecase)
	mockJWTService := new(mocks.MockJWTService)

	taskController := controllers.NewTaskController(mockTaskUsecase)
	userController := controllers.NewUserController(mockUserUsecase)

	router := routers.SetupRouter(taskController, userController, mockJWTService)

	t.Run("RegisterRoute", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/register", nil)
		router.ServeHTTP(w, req)

		// We expect 400 because body is empty, but it proves the route exists
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("LoginRoute", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("TasksRoute_Protected", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tasks", nil)
		router.ServeHTTP(w, req)

		// Should be 401 because of missing auth header (middleware check)
		// If middleware wasn't there or wasn't working, it might be something else.
		// But since we are mocking JWTService, we need to see how AuthMiddleware uses it.
		// If AuthMiddleware calls ValidateToken, we need to mock it or expect 401 if token is missing.
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
