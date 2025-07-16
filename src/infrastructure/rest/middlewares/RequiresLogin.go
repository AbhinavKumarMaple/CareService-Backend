package middlewares

import (
	"github.com/gin-gonic/gin"
)

func AuthJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Disable all authentication and authorization checks for experimental phase.
		// All requests will be allowed to pass through without token validation.
		c.Next()
	}
}
