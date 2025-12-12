package mocks

import (
	"a2sv-backend/task_manager_v3/Domain"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user Domain.User) (Domain.User, error) {
	args := m.Called(user)
	return args.Get(0).(Domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(username string) (Domain.User, error) {
	args := m.Called(username)
	return args.Get(0).(Domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id primitive.ObjectID) (Domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(Domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(user Domain.User) (Domain.User, error) {
	args := m.Called(user)
	return args.Get(0).(Domain.User), args.Error(1)
}

func (m *MockUserRepository) Count() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
