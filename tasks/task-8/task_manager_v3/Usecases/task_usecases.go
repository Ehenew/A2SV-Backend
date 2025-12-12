package Usecases

import (
	"a2sv-backend/task_manager_v3/Domain"
	"a2sv-backend/task_manager_v3/Repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskUsecase interface {
	Create(task Domain.Task) (Domain.Task, error)
	GetAll() ([]Domain.Task, error)
	GetByID(id primitive.ObjectID) (Domain.Task, error)
	Update(id primitive.ObjectID, task Domain.Task) (Domain.Task, error)
	Delete(id primitive.ObjectID) error
}

type taskUsecase struct {
	taskRepo Repositories.TaskRepository
}

func NewTaskUsecase(taskRepo Repositories.TaskRepository) TaskUsecase {
	return &taskUsecase{
		taskRepo: taskRepo,
	}
}

func (u *taskUsecase) Create(task Domain.Task) (Domain.Task, error) {
	return u.taskRepo.Create(task)
}

func (u *taskUsecase) GetAll() ([]Domain.Task, error) {
	return u.taskRepo.FindAll()
}

func (u *taskUsecase) GetByID(id primitive.ObjectID) (Domain.Task, error) {
	return u.taskRepo.FindByID(id)
}

func (u *taskUsecase) Update(id primitive.ObjectID, task Domain.Task) (Domain.Task, error) {
	task.ID = id
	return u.taskRepo.Update(task)
}

func (u *taskUsecase) Delete(id primitive.ObjectID) error {
	return u.taskRepo.Delete(id)
}
