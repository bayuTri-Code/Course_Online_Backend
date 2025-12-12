package middleware

import (
	
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

func CheckCourseOwnershipDynamic(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		rolesValue, _ := c.Get("roles")
		userIDStr := userID.(string)

		isSuper := false
		isAdmin := false

		switch r := rolesValue.(type) {
		case string:
			if r == "super_admin" {
				isSuper = true
			}
			if r == "admin" {
				isAdmin = true
			}
		case []string:
			for _, role := range r {
				if role == "super_admin" {
					isSuper = true
				}
				if role == "admin" {
					isAdmin = true
				}
			}
		}

		if isSuper {
			c.Next()
			return
		}

		var courseID string

		if cid := c.Param("id"); cid != "" {
			courseID = cid
		}

		//quiz
		if qid := c.Param("quizId"); qid != "" {
			_ = db.Raw(`
        SELECT course_id FROM quizzes WHERE id = ?
    `, qid).Scan(&courseID)
		}

		if mid := c.Param("module_id"); mid != "" {
			_ = db.Raw(`
                SELECT course_id FROM modules WHERE id = ?
            `, mid).Scan(&courseID)
		}

		if lid := c.Param("id"); lid != "" && c.FullPath() != "/courses/:id" {
			_ = db.Raw(`
                SELECT m.course_id 
                FROM lessons l 
                JOIN modules m ON l.module_id = m.id
                WHERE l.id = ?
            `, lid).Scan(&courseID)
		}

		if courseID == "" {
			c.JSON(400, gin.H{"message": "Invalid or missing resource identifier"})
			c.Abort()
			return
		}

		var createdBy string
		err := db.Raw(`
            SELECT created_by FROM courses WHERE id = ?
        `, courseID).Scan(&createdBy).Error

		if err != nil || createdBy == "" {
			c.JSON(404, gin.H{"message": "Course not found"})
			c.Abort()
			return
		}

		if isAdmin {
			if createdBy == userIDStr {
				c.Next()
				return
			}
			c.JSON(403, gin.H{"message": "Forbidden: you can only manage your own course content"})
			c.Abort()
			return
		}

		c.JSON(403, gin.H{"message": "Forbidden: insufficient permissions"})
		c.Abort()
	}
}
