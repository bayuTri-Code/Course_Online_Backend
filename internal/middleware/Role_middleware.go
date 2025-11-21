package middleware

import (
	"course_online_backend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func CheckCourseOwnership(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "Unauthorized: user not found in context",
			})
			c.Abort()
			return
		}

		roleValue, exists := c.Get("roles")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "Unauthorized: roles not found in context",
			})
			c.Abort()
			return
		}

		courseID := c.Param("id")

		var course models.Course
		if err := db.Select("created_by").Where("id = ?", courseID).First(&course).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{
					"status":  false,
					"message": "Course not found",
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to check course ownership",
			})
			c.Abort()
			return
		}

		userIDStr := userID.(string)
		isSuperAdmin := false
		isRegularAdmin := false

		switch roles := roleValue.(type) {
		case string:
			switch roles {
			case "super_admin":
				isSuperAdmin = true
			case "admin":
				isRegularAdmin = true
			}

		case []string:
			for _, r := range roles {
				if r == "super_admin" {
					isSuperAdmin = true
					break
				}
				if r == "admin" {
					isRegularAdmin = true
				}
			}
		}

		if isSuperAdmin {
			c.Next()
			return
		}

		if isRegularAdmin {
			if course.CreatedBy != nil && course.CreatedBy.String() == userIDStr {
				c.Next()
				return
			}

			c.JSON(http.StatusForbidden, gin.H{
				"status":  false,
				"message": "Forbidden: you can only manage courses you created",
			})
			c.Abort()
			return
		}

		c.JSON(http.StatusForbidden, gin.H{
			"status":  false,
			"message": "Forbidden: insufficient permissions",
		})
		c.Abort()
	}
}
