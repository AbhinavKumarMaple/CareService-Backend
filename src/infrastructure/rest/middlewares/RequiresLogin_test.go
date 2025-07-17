package middlewares

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func setupGinContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	router := gin.New()
	// Add error handling middleware
	router.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			c.JSON(500, gin.H{"error": err.Error()})
		}
	})
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	return c, w
}

func TestAuthJWTMiddleware_NoToken(t *testing.T) {
	c, w := setupGinContext()
	c.Request = httptest.NewRequest("GET", "/protected", nil)

	middleware := AuthJWTMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Token not provided")
}

func TestAuthJWTMiddleware_NoJWTSecret(t *testing.T) {
	originalSecret := os.Getenv("JWT_ACCESS_SECRET_KEY")
	os.Unsetenv("JWT_ACCESS_SECRET_KEY")
	defer os.Setenv("JWT_ACCESS_SECRET_KEY", originalSecret)

	c, w := setupGinContext()
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", "Bearer valid-token")

	middleware := AuthJWTMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "JWT_ACCESS_SECRET_KEY not configured")
}

func TestAuthJWTMiddleware_InvalidToken(t *testing.T) {
	originalSecret := os.Getenv("JWT_ACCESS_SECRET_KEY")
	os.Setenv("JWT_ACCESS_SECRET_KEY", "test-secret")
	defer os.Setenv("JWT_ACCESS_SECRET_KEY", originalSecret)

	c, w := setupGinContext()
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid-token")

	middleware := AuthJWTMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token")
}

func TestAuthJWTMiddleware_ExpiredToken(t *testing.T) {
	originalSecret := os.Getenv("JWT_ACCESS_SECRET_KEY")
	os.Setenv("JWT_ACCESS_SECRET_KEY", "test-secret")
	defer os.Setenv("JWT_ACCESS_SECRET_KEY", originalSecret)

	// Create expired token
	claims := jwt.MapClaims{
		"exp":  time.Now().Add(-1 * time.Hour).Unix(), 
		"type": "access",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test-secret"))

	c, w := setupGinContext()
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", "Bearer "+tokenString)

	middleware := AuthJWTMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Token expired")
}

func TestAuthJWTMiddleware_InvalidTokenClaims(t *testing.T) {
	originalSecret := os.Getenv("JWT_ACCESS_SECRET_KEY")
	os.Setenv("JWT_ACCESS_SECRET_KEY", "test-secret")
	defer os.Setenv("JWT_ACCESS_SECRET_KEY", originalSecret)

	// Create token without exp claim
	claims := jwt.MapClaims{
		"type": "access",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test-secret"))

	c, w := setupGinContext()
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", "Bearer "+tokenString)

	middleware := AuthJWTMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token claims")
}

func TestAuthJWTMiddleware_WrongTokenType(t *testing.T) {
	originalSecret := os.Getenv("JWT_ACCESS_SECRET_KEY")
	os.Setenv("JWT_ACCESS_SECRET_KEY", "test-secret")
	defer os.Setenv("JWT_ACCESS_SECRET_KEY", originalSecret)

	claims := jwt.MapClaims{
		"exp":  time.Now().Add(1 * time.Hour).Unix(),
		"type": "refresh", 
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test-secret"))

	c, w := setupGinContext()
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", "Bearer "+tokenString)

	middleware := AuthJWTMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Token type mismatch")
}

func TestAuthJWTMiddleware_MissingTokenType(t *testing.T) {
	originalSecret := os.Getenv("JWT_ACCESS_SECRET_KEY")
	os.Setenv("JWT_ACCESS_SECRET_KEY", "test-secret")
	defer os.Setenv("JWT_ACCESS_SECRET_KEY", originalSecret)

	claims := jwt.MapClaims{
		"exp": time.Now().Add(1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test-secret"))

	c, w := setupGinContext()
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", "Bearer "+tokenString)

	middleware := AuthJWTMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Missing token type")
}

func TestAuthJWTMiddleware_ValidToken(t *testing.T) {
	originalSecret := os.Getenv("JWT_ACCESS_SECRET_KEY")
	os.Setenv("JWT_ACCESS_SECRET_KEY", "test-secret")
	defer os.Setenv("JWT_ACCESS_SECRET_KEY", originalSecret)

	// Create valid token
	claims := jwt.MapClaims{
		"exp":  time.Now().Add(1 * time.Hour).Unix(),
		"type": "access",
		"id":   123,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test-secret"))

	c, w := setupGinContext()
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", "Bearer "+tokenString)

	middleware := AuthJWTMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestAuthJWTMiddleware_TokenWithoutBearer(t *testing.T) {
	originalSecret := os.Getenv("JWT_ACCESS_SECRET_KEY")
	os.Setenv("JWT_ACCESS_SECRET_KEY", "test-secret")
	defer os.Setenv("JWT_ACCESS_SECRET_KEY", originalSecret)

	// Create valid token
	claims := jwt.MapClaims{
		"exp":  time.Now().Add(1 * time.Hour).Unix(),
		"type": "access",
		"id":   123,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test-secret"))

	c, w := setupGinContext()
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", tokenString) 

	middleware := AuthJWTMiddleware()
	middleware(c)


	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token format")
}
