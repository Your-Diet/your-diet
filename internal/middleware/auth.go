package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ContextKey is a type for context keys
type ContextKey string

const (
	// TokenContextKey is the key used to store the token claims in the context
	TokenContextKey ContextKey = "token_claims"
	// UserIDContextKey is the key used to store the user ID in the context
	UserIDContextKey ContextKey = "user_id"
	// PermissionsContextKey is the key used to store the user permissions in the context
	PermissionsContextKey ContextKey = "permissions"
)

// Claims defines the JWT claims structure
type Claims struct {
	UserID      string   `json:"user_id"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

// AuthMiddleware creates a new authentication middleware
func AuthMiddleware(jwtSecretKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Check if the token is in the format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		tokenString := parts[1]

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecretKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Get the claims
		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Add claims, user ID and permissions to context
		c.Set(string(TokenContextKey), claims)
		c.Set(string(UserIDContextKey), claims.UserID)
		c.Set(string(PermissionsContextKey), claims.Permissions)

		// Continue to the next handler
		c.Next()
	}
}

// HasPermission checks if the user has the required permission
func HasPermission(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions, ok := c.Get(string(PermissionsContextKey))
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No authentication data found"})
			return
		}

		hasPermission := false
		for _, p := range permissions.([]string) {
			if p == requiredPermission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}

		c.Next()
	}
}
