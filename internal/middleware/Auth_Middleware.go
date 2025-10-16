package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
)
func FirebaseAuth(app *firebase.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			fmt.Println("No Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		fmt.Println("ot header:", authHeader)

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			fmt.Println("Invalid header format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header"})
			c.Abort()
			return
		}

		idToken := parts[1]
		client, err := app.Auth(context.Background())
		if err != nil {
			fmt.Println("Auth client error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error initializing auth"})
			c.Abort()
			return
		}

		token, err := client.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			fmt.Println("nvalid token:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		fmt.Println("uth success UID:", token.UID)
		c.Set("user_id", token.UID)
		c.Next()
	}
}
