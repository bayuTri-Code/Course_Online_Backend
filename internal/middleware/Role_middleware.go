package middleware

import (
	"course_online_backend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func CheckEnrollmentAccess(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: user not found",
			})
			c.Abort()
			return
		}

		var courseID uuid.UUID
		var err error

		if id := c.Param("id"); id != "" {
			courseID, err = uuid.Parse(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
				c.Abort()
				return
			}
		}

		if courseID == uuid.Nil {
			if moduleID := c.Param("module_id"); moduleID != "" {
				var module models.Module
				if err := db.First(&module, "id = ?", moduleID).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
					c.Abort()
					return
				}
				courseID = module.CourseID
			}
		}

		if courseID == uuid.Nil {
			if lessonID := c.Param("lesson_id"); lessonID != "" || c.Param("id") != "" {
				target := c.Param("lesson_id")
				if target == "" {
					target = c.Param("id")
				}

				var lesson models.Lesson
				if err := db.First(&lesson, "id = ?", target).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Lesson not found"})
					c.Abort()
					return
				}

				var module models.Module
				db.First(&module, "id = ?", lesson.ModuleID)
				courseID = module.CourseID
			}
		}

		if courseID == uuid.Nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot determine course for this content"})
			c.Abort()
			return
		}

		var enrollment models.Enrollment
		err = db.
			Where("user_id = ? AND course_id = ?", userID, courseID).
			First(&enrollment).Error

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied: You are not enrolled in this course",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
