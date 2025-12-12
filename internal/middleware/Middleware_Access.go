package middleware

import (
	"course_online_backend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CheckEnrollmentAccess(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		firebaseUID := c.GetString("user_id")
		if firebaseUID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: user not found",
			})
			c.Abort()
			return
		}

		var user models.User
		if err := db.Select("id").Where("firebase_uid = ?", firebaseUID).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "User not found in database",
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to fetch user",
				})
			}
			c.Abort()
			return
		}

		var courseID uuid.UUID
		var err error

		if id := c.Param("courseId"); id != "" {
			courseID, err = uuid.Parse(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
				c.Abort()
				return
			}
		}

		if courseID == uuid.Nil {
			if id := c.Param("id"); id != "" {
				courseID, err = uuid.Parse(id)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
					c.Abort()
					return
				}
			}
		}

		if courseID == uuid.Nil {
			if quizID := c.Param("quizId"); quizID != "" {
				quizUUID, err := uuid.Parse(quizID)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quiz ID"})
					c.Abort()
					return
				}

				var quiz models.Quiz
				if err := db.Select("course_id").Where("id = ?", quizUUID).First(&quiz).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
					} else {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch quiz"})
					}
					c.Abort()
					return
				}
				courseID = quiz.CourseID
			}
		}

		if courseID == uuid.Nil {
			if moduleID := c.Param("module_id"); moduleID != "" {
				moduleUUID, err := uuid.Parse(moduleID)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
					c.Abort()
					return
				}

				var module models.Module
				if err := db.Select("course_id").Where("id = ?", moduleUUID).First(&module).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
					} else {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch module"})
					}
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

				lessonUUID, err := uuid.Parse(target)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
					c.Abort()
					return
				}

				var lesson models.Lesson
				if err := db.Select("module_id").Where("id = ?", lessonUUID).First(&lesson).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						c.JSON(http.StatusNotFound, gin.H{"error": "Lesson not found"})
					} else {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch lesson"})
					}
					c.Abort()
					return
				}

				var module models.Module
				if err := db.Select("course_id").Where("id = ?", lesson.ModuleID).First(&module).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch module from lesson"})
					c.Abort()
					return
				}
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
			Where("user_id = ? AND course_id = ? AND status = ?", user.ID, courseID, "active").
			First(&enrollment).Error

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied: You are not enrolled in this course",
			})
			c.Abort()
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to check enrollment",
			})
			c.Abort()
			return
		}

		c.Set("course_id", courseID)
		c.Set("user_db_id", user.ID)
		c.Next()
	}
}
