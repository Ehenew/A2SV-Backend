package infrastructure_test

import (
	"a2sv-backend/task_manager_v3/Infrastructure"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordService_HashPassword(t *testing.T) {
	service := Infrastructure.NewBcryptPasswordService()
	password := "password123"

	hashedPassword, err := service.HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, password, hashedPassword)
}

func TestPasswordService_ComparePassword(t *testing.T) {
	service := Infrastructure.NewBcryptPasswordService()
	password := "password123"
	hashedPassword, _ := service.HashPassword(password)

	t.Run("Success", func(t *testing.T) {
		err := service.ComparePassword(hashedPassword, password)
		assert.NoError(t, err)
	})

	t.Run("Failure", func(t *testing.T) {
		err := service.ComparePassword(hashedPassword, "wrongpassword")
		assert.Error(t, err)
	})
}
