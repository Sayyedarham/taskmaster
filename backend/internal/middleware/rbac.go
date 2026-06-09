package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var rolePermissions = map[string][]string{
	"admin":   {"task:create", "task:read", "task:update", "task:delete", "team:manage"},
	"manager": {"task:create", "task:read", "task:update", "task:delete"},
	"member":  {"task:read", "task:update"},
}

func RBAC() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "role not found"})
			return
		}

		// Simple check: all protected routes need at least task:read
		perms, ok := rolePermissions[role]
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "unknown role"})
			return
		}

		required := "task:read"
		if c.Request.Method == "POST" {
			required = "task:create"
		} else if c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			required = "task:update"
		} else if c.Request.Method == "DELETE" {
			required = "task:delete"
		}

		hasPerm := false
		for _, p := range perms {
			if p == required {
				hasPerm = true
				break
			}
		}

		if !hasPerm {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}
		c.Next()
	}
}
