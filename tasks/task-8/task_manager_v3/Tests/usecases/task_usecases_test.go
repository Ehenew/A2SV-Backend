package usecases_test

import (
	"a2sv-backend/task_manager_v3/Domain"
	"a2sv-backend/task_manager_v3/Tests/mocks"
	"a2sv-backend/task_manager_v3/Usecases"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTaskUsecase_Create(t *testing.T) {
	mockTaskRepo := new(mocks.MockTaskRepository)
	taskUsecase := Usecases.NewTaskUsecase(mockTaskRepo)

	t.Run("Success", func(t *testing.T) {
		task := Domain.Task{
			Title:       "Test Task",
			Description: "Test Description",
			DueDate:     time.Now(),
			Status:      "pending",
		}

		mockTaskRepo.On("Create", task).Return(task, nil)

		createdTask, err := taskUsecase.Create(task)

		assert.NoError(t, err)
		assert.Equal(t, task, createdTask)
		mockTaskRepo.AssertExpectations(t)
	})
}

func TestTaskUsecase_GetAll(t *testing.T) {
	mockTaskRepo := new(mocks.MockTaskRepository)
	taskUsecase := Usecases.NewTaskUsecase(mockTaskRepo)

	t.Run("Success", func(t *testing.T) {
		tasks := []Domain.Task{
			{Title: "Task 1"},
			{Title: "Task 2"},
		}

		mockTaskRepo.On("FindAll").Return(tasks, nil)

		resultTasks, err := taskUsecase.GetAll()

		assert.NoError(t, err)
		assert.Equal(t, tasks, resultTasks)
		mockTaskRepo.AssertExpectations(t)
	})
}

func TestTaskUsecase_GetByID(t *testing.T) {
	mockTaskRepo := new(mocks.MockTaskRepository)
	taskUsecase := Usecases.NewTaskUsecase(mockTaskRepo)

	t.Run("Success", func(t *testing.T) {
		taskID := primitive.NewObjectID()
		task := Domain.Task{
			ID:    taskID,
			Title: "Test Task",
		}

		mockTaskRepo.On("FindByID", taskID).Return(task, nil)

		resultTask, err := taskUsecase.GetByID(taskID)

		assert.NoError(t, err)
		assert.Equal(t, task, resultTask)
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		taskID := primitive.NewObjectID()

		mockTaskRepo.On("FindByID", taskID).Return(Domain.Task{}, errors.New("task not found"))

		_, err := taskUsecase.GetByID(taskID)

		assert.Error(t, err)
		assert.Equal(t, "task not found", err.Error())
	})
}

func TestTaskUsecase_Update(t *testing.T) {
	mockTaskRepo := new(mocks.MockTaskRepository)
	taskUsecase := Usecases.NewTaskUsecase(mockTaskRepo)

	t.Run("Success", func(t *testing.T) {
		taskID := primitive.NewObjectID()
		task := Domain.Task{
			Title: "Updated Task",
		}
		
		expectedTask := task
		expectedTask.ID = taskID

		mockTaskRepo.On("Update", expectedTask).Return(expectedTask, nil)

		updatedTask, err := taskUsecase.Update(taskID, task)

		assert.NoError(t, err)
		assert.Equal(t, expectedTask, updatedTask)
		mockTaskRepo.AssertExpectations(t)
	})
}

func TestTaskUsecase_Delete(t *testing.T) {
	mockTaskRepo := new(mocks.MockTaskRepository)
	taskUsecase := Usecases.NewTaskUsecase(mockTaskRepo)

	t.Run("Success", func(t *testing.T) {
		taskID := primitive.NewObjectID()

		mockTaskRepo.On("Delete", taskID).Return(nil)

		err := taskUsecase.Delete(taskID)

		assert.NoError(t, err)
		mockTaskRepo.AssertExpectations(t)
	})
}
