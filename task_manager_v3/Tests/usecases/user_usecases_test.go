package usecases_test

import (
	"a2sv-backend/task_manager_v3/Domain"
	"a2sv-backend/task_manager_v3/Tests/mocks"
	"a2sv-backend/task_manager_v3/Usecases"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUserUsecase_Register(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockPasswordService := new(mocks.MockPasswordService)
	mockJWTService := new(mocks.MockJWTService)

	userUsecase := Usecases.NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)

	t.Run("Success", func(t *testing.T) {
		user := Domain.User{
			Username: "testuser",
			Password: "password",
		}

		mockUserRepo.On("FindByUsername", user.Username).Return(Domain.User{}, errors.New("user not found"))
		mockPasswordService.On("HashPassword", user.Password).Return("hashed_password", nil)
		mockUserRepo.On("Count").Return(int64(1), nil)
		mockUserRepo.On("Create", mock.AnythingOfType("Domain.User")).Return(user, nil)

		err := userUsecase.Register(user)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
	})

	t.Run("UsernameExists", func(t *testing.T) {
		user := Domain.User{
			Username: "existinguser",
			Password: "password",
		}

		existingUser := Domain.User{
			Username: "existinguser",
		}

		mockUserRepo.On("FindByUsername", user.Username).Return(existingUser, nil)

		err := userUsecase.Register(user)

		assert.Error(t, err)
		assert.Equal(t, "username already exists", err.Error())
	})
}

func TestUserUsecase_Login(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockPasswordService := new(mocks.MockPasswordService)
	mockJWTService := new(mocks.MockJWTService)

	userUsecase := Usecases.NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)

	t.Run("Success", func(t *testing.T) {
		username := "testuser"
		password := "password"
		hashedPassword := "hashed_password"
		user := Domain.User{
			Username: username,
			Password: hashedPassword,
		}
		token := "test_token"

		mockUserRepo.On("FindByUsername", username).Return(user, nil)
		mockPasswordService.On("ComparePassword", hashedPassword, password).Return(nil)
		mockJWTService.On("GenerateToken", user).Return(token, nil)

		resultToken, resultUser, err := userUsecase.Login(username, password)

		assert.NoError(t, err)
		assert.Equal(t, token, resultToken)
		assert.Equal(t, user, resultUser)
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		username := "testuser"
		password := "wrongpassword"
		hashedPassword := "hashed_password"
		user := Domain.User{
			Username: username,
			Password: hashedPassword,
		}

		mockUserRepo.On("FindByUsername", username).Return(user, nil)
		mockPasswordService.On("ComparePassword", hashedPassword, password).Return(errors.New("invalid password"))

		_, _, err := userUsecase.Login(username, password)

		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
	})
}

func TestUserUsecase_Promote(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockPasswordService := new(mocks.MockPasswordService)
	mockJWTService := new(mocks.MockJWTService)

	userUsecase := Usecases.NewUserUsecase(mockUserRepo, mockPasswordService, mockJWTService)

	t.Run("Success", func(t *testing.T) {
		userID := primitive.NewObjectID()
		user := Domain.User{
			ID:   userID,
			Role: "user",
		}

		mockUserRepo.On("FindByID", userID).Return(user, nil)
		mockUserRepo.On("Update", mock.MatchedBy(func(u Domain.User) bool {
			return u.ID == userID && u.Role == "admin"
		})).Return(user, nil)

		err := userUsecase.Promote(userID)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})
}
