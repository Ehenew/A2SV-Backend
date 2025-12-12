package mocks

import (
	"a2sv-backend/task_manager_v3/Domain"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) Register(user Domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserUsecase) Login(username, password string) (string, Domain.User, error) {
	args := m.Called(username, password)
	return args.String(0), args.Get(1).(Domain.User), args.Error(2)
}

func (m *MockUserUsecase) Promote(userID primitive.ObjectID) error {
	args := m.Called(userID)
	return args.Error(0)
}
