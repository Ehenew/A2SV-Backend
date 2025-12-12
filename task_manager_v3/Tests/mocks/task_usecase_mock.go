package mocks

import (
	"a2sv-backend/task_manager_v3/Domain"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockTaskUsecase struct {
	mock.Mock
}

func (m *MockTaskUsecase) Create(task Domain.Task) (Domain.Task, error) {
	args := m.Called(task)
	return args.Get(0).(Domain.Task), args.Error(1)
}

func (m *MockTaskUsecase) GetAll() ([]Domain.Task, error) {
	args := m.Called()
	return args.Get(0).([]Domain.Task), args.Error(1)
}

func (m *MockTaskUsecase) GetByID(id primitive.ObjectID) (Domain.Task, error) {
	args := m.Called(id)
	return args.Get(0).(Domain.Task), args.Error(1)
}

func (m *MockTaskUsecase) Update(id primitive.ObjectID, task Domain.Task) (Domain.Task, error) {
	args := m.Called(id, task)
	return args.Get(0).(Domain.Task), args.Error(1)
}

func (m *MockTaskUsecase) Delete(id primitive.ObjectID) error {
	args := m.Called(id)
	return args.Error(0)
}
