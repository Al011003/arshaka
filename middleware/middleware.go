package middleware

import (
	"net/http"

	"backend/token"

	"github.com/gin-gonic/gin"
)

// JWTAuth middleware: cek access token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		claims, err := token.ValidateAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// simpan user_id & role di context
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)


		c.Next()
	}
}

// UserOnly middleware: hanya user
func UserOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "user" {
			c.JSON(http.StatusForbidden, gin.H{"error": "User access only"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// AdminOnly middleware: hanya admin
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access only"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// SuperAdminOnly middleware: hanya superadmin
func SuperAdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "superadmin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "SuperAdmin access only"})
			c.Abort()
			return
		}
		c.Next()
	}
}
