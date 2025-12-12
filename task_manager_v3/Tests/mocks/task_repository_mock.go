package mocks

import (
	"a2sv-backend/task_manager_v3/Domain"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(task Domain.Task) (Domain.Task, error) {
	args := m.Called(task)
	return args.Get(0).(Domain.Task), args.Error(1)
}

func (m *MockTaskRepository) FindAll() ([]Domain.Task, error) {
	args := m.Called()
	return args.Get(0).([]Domain.Task), args.Error(1)
}

func (m *MockTaskRepository) FindByID(id primitive.ObjectID) (Domain.Task, error) {
	args := m.Called(id)
	return args.Get(0).(Domain.Task), args.Error(1)
}

func (m *MockTaskRepository) Update(task Domain.Task) (Domain.Task, error) {
	args := m.Called(task)
	return args.Get(0).(Domain.Task), args.Error(1)
}

func (m *MockTaskRepository) Delete(id primitive.ObjectID) error {
	args := m.Called(id)
	return args.Error(0)
}
