package middleware_test

import (
	"a2sv-backend/task_manager_v3/Infrastructure"
	"a2sv-backend/task_manager_v3/Tests/mocks"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("NoHeader", func(t *testing.T) {
		mockJWTService := new(mocks.MockJWTService)
		middleware := Infrastructure.AuthMiddleware(mockJWTService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/", nil)

		middleware(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authorization header is required")
	})

	t.Run("InvalidFormat", func(t *testing.T) {
		mockJWTService := new(mocks.MockJWTService)
		middleware := Infrastructure.AuthMiddleware(mockJWTService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", "InvalidFormat")

		middleware(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid authorization header format")
	})

	t.Run("InvalidToken", func(t *testing.T) {
		mockJWTService := new(mocks.MockJWTService)
		middleware := Infrastructure.AuthMiddleware(mockJWTService)

		mockJWTService.On("ValidateToken", "invalid_token").Return(&jwt.Token{}, errors.New("invalid token"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", "Bearer invalid_token")

		middleware(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid JWT token")
	})

	t.Run("Success", func(t *testing.T) {
		mockJWTService := new(mocks.MockJWTService)
		middleware := Infrastructure.AuthMiddleware(mockJWTService)

		token := &jwt.Token{
			Valid: true,
			Claims: jwt.MapClaims{
				"user_id":  "123",
				"username": "testuser",
				"role":     "user",
			},
		}

		mockJWTService.On("ValidateToken", "valid_token").Return(token, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", "Bearer valid_token")

		middleware(c)

		assert.Equal(t, http.StatusOK, w.Code)
		// Check if context keys are set
		userID, _ := c.Get("user_id")
		assert.Equal(t, "123", userID)
	})
}
