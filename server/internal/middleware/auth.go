package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	jwtpkg "github.com/lunancy1992/jianghu-server/internal/pkg/jwt"
	"github.com/lunancy1992/jianghu-server/internal/pkg/response"
)

// JWTAuth requires a valid JWT token.
func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := extractClaims(c, secret)
		if claims == nil {
			response.Unauthorized(c, "authentication required")
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Set("user_nickname", claims.Nickname)
		c.Next()
	}
}

// OptionalAuth extracts user info if token is present but doesn't require it.
func OptionalAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := extractClaims(c, secret)
		if claims != nil {
			c.Set("user_id", claims.UserID)
			c.Set("user_role", claims.Role)
			c.Set("user_nickname", claims.Nickname)
		}
		c.Next()
	}
}

// RequireAdmin checks that the authenticated user has admin role.
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists || role.(string) != "admin" {
			response.Forbidden(c, "admin access required")
			c.Abort()
			return
		}
		c.Next()
	}
}

func extractClaims(c *gin.Context, secret string) *jwtpkg.Claims {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		return nil
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil
	}
	claims, err := jwtpkg.Validate(secret, parts[1])
	if err != nil {
		return nil
	}
	return claims
}

// GetUserID returns the authenticated user ID from the context.
func GetUserID(c *gin.Context) (int64, bool) {
	val, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return val.(int64), true
}
