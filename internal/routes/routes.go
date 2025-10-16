package routes

import (
	"course_online_backend/internal/handler"
	"course_online_backend/internal/services"
	"log"

	"firebase.google.com/go/v4"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routes(db *gorm.DB, app *firebase.App) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173", "http://192.168.100.247:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Panicf("Failed to set trusted proxies: %v", err)
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success api test"})
	})

	authService := services.NewAuthService(db, app)
	authHandler := handlers.NewAuthHandler(authService)

	authRoutes := r.Group("/api/auth")
	{
		authRoutes.POST("/firebase-login", authHandler.FirebaseLoginHandler)
	}

	// protectedRoutes := r.Group("/api")
	// protectedRoutes.Use(middleware.FirebaseAuth(app))e
	// {

	// }

	return r
}