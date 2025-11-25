package middleware

import (
	"context"
	"course_online_backend/database"
	"course_online_backend/internal/models"
	"net/http"
	"strings"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FirebaseAuth(app *firebase.App, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header"})
			c.Abort()
			return
		}

		idToken := parts[1]

		isBlacklisted, err := database.RedisConn().Exists(context.Background(), "blacklist:"+idToken).Result()
		if err == nil && isBlacklisted > 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
			c.Abort()
			return
		}

		client, err := app.Auth(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error initializing Firebase auth"})
			c.Abort()
			return
		}

		uidRedis, err := database.RedisConn().Get(context.Background(), "token:"+idToken).Result()
		if err == nil && uidRedis != "" {
			var user models.User
			if err := db.Preload("Roles").Where("firebase_uid = ?", uidRedis).First(&user).Error; err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
				c.Abort()
				return
			}

			roleNames := make([]string, len(user.Roles))
			for i, role := range user.Roles {
				roleNames[i] = role.Name
			}

			c.Set("user_id", user.FirebaseUID)
			c.Set("roles", roleNames)
			c.Set("user", user)
			c.Set("token", idToken) 
			c.Next()
			return
		}

		token, err := client.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		var user models.User
		if err := db.Preload("Roles").Where("firebase_uid = ?", token.UID).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		
		_ = database.RedisConn().Set(context.Background(), "token:"+idToken, token.UID, time.Hour).Err()

		roleNames := make([]string, len(user.Roles))
		for i, role := range user.Roles {
			roleNames[i] = role.Name
		}

		c.Set("user_id", user.FirebaseUID)
		c.Set("roles", roleNames)
		c.Set("user", user)
		c.Set("token", idToken) 

		c.Next()
	}
}