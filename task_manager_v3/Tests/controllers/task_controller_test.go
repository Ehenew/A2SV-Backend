package controllers_test

import (
	"a2sv-backend/task_manager_v3/Delivery/controllers"
	"a2sv-backend/task_manager_v3/Domain"
	"a2sv-backend/task_manager_v3/Tests/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTaskController_CreateTask(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockTaskUsecase := new(mocks.MockTaskUsecase)
	taskController := controllers.NewTaskController(mockTaskUsecase)

	t.Run("Success", func(t *testing.T) {
		task := Domain.Task{
			Title:       "Test Task",
			Description: "Test Description",
			DueDate:     time.Now(),
			Status:      "pending",
		}
		
		// Mock expects the task as it is sent in the request
		// We use MatchedBy to handle potential time precision issues or monotonic clock differences
		mockTaskUsecase.On("Create", mock.MatchedBy(func(t Domain.Task) bool {
			return t.Title == task.Title && t.Description == task.Description && t.Status == task.Status
		})).Return(task, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonValue, _ := json.Marshal(task)
		c.Request, _ = http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonValue))

		taskController.CreateTask(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockTaskUsecase.AssertExpectations(t)
	})

	t.Run("BadRequest", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest("POST", "/tasks", bytes.NewBuffer([]byte("invalid json")))

		taskController.CreateTask(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestTaskController_GetAllTasks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockTaskUsecase := new(mocks.MockTaskUsecase)
	taskController := controllers.NewTaskController(mockTaskUsecase)

	t.Run("Success", func(t *testing.T) {
		tasks := []Domain.Task{
			{Title: "Task 1"},
			{Title: "Task 2"},
		}

		mockTaskUsecase.On("GetAll").Return(tasks, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest("GET", "/tasks", nil)

		taskController.GetAllTasks(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockTaskUsecase.AssertExpectations(t)
	})
}

func TestTaskController_GetTaskByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockTaskUsecase := new(mocks.MockTaskUsecase)
	taskController := controllers.NewTaskController(mockTaskUsecase)

	t.Run("Success", func(t *testing.T) {
		taskID := primitive.NewObjectID()
		task := Domain.Task{
			ID:    taskID,
			Title: "Test Task",
		}

		mockTaskUsecase.On("GetByID", taskID).Return(task, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = gin.Params{{Key: "id", Value: taskID.Hex()}}
		c.Request, _ = http.NewRequest("GET", "/tasks/"+taskID.Hex(), nil)

		taskController.GetTaskByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockTaskUsecase.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		taskID := primitive.NewObjectID()

		mockTaskUsecase.On("GetByID", taskID).Return(Domain.Task{}, errors.New("task not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = gin.Params{{Key: "id", Value: taskID.Hex()}}
		c.Request, _ = http.NewRequest("GET", "/tasks/"+taskID.Hex(), nil)

		taskController.GetTaskByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
