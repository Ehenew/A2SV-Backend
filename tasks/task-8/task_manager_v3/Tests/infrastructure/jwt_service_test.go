package infrastructure_test

import (
	"a2sv-backend/task_manager_v3/Domain"
	"a2sv-backend/task_manager_v3/Infrastructure"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestJWTService_GenerateToken(t *testing.T) {
	service := Infrastructure.NewJWTService()
	user := Domain.User{
		ID:       primitive.NewObjectID(),
		Username: "testuser",
		Role:     "user",
	}

	token, err := service.GenerateToken(user)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestJWTService_ValidateToken(t *testing.T) {
	service := Infrastructure.NewJWTService()
	user := Domain.User{
		ID:       primitive.NewObjectID(),
		Username: "testuser",
		Role:     "user",
	}
	tokenString, _ := service.GenerateToken(user)

	t.Run("Success", func(t *testing.T) {
		token, err := service.ValidateToken(tokenString)

		assert.NoError(t, err)
		assert.NotNil(t, token)
		assert.True(t, token.Valid)

		claims, ok := token.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, user.Username, claims["username"])
		assert.Equal(t, user.Role, claims["role"])
	})

	t.Run("InvalidToken", func(t *testing.T) {
		_, err := service.ValidateToken("invalid_token")
		assert.Error(t, err)
	})
}
