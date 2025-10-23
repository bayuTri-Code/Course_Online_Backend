package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("roles")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Role not found"})
			c.Abort()
			return
		}


		isAllowed := false

		switch roles := roleValue.(type) {
		case string:
			for _, allowed := range allowedRoles {
				if roles == allowed {
					isAllowed = true
					break
				}
			}

		case []string: 
			for _, userRole := range roles {
				for _, allowed := range allowedRoles {
					if userRole == allowed {
						isAllowed = true
						break
					}
				}
				if isAllowed {
					break
				}
			}

		default:
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid role format"})
			c.Abort()
			return
		}

		if !isAllowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient role"})
			c.Abort()
			return
		}

		c.Next()
	}
}
