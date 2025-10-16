package routes

import (
	"course_online_backend/internal/handler"
	"course_online_backend/internal/middleware"
	"course_online_backend/internal/services"
	"log"

	"firebase.google.com/go/v4"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

func Routes(db *gorm.DB, app *firebase.App, minioClient *minio.Client, minioBucket, minioURL string) *gin.Engine {
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
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Course Online Backend API Succes",
		})
	})

	

	authService := services.NewAuthService(db, app)
	biodataService := services.NewBiodataService(db)

	authHandler := handler.NewAuthHandler(authService)
	biodataHandler := handler.NewBiodataHandler(biodataService)

	publicRoutes := r.Group("/api")
	{
		auth := publicRoutes.Group("/auth")
		{
			auth.POST("/firebase-login", authHandler.FirebaseLoginHandler)
		}
	}

	protectedRoutes := r.Group("/api")
	protectedRoutes.Use(middleware.FirebaseAuth(app))
	{
		biodata := protectedRoutes.Group("/profile")
		{
			biodata.POST("/biodata", biodataHandler.CreateBiodata)    
			biodata.GET("/biodata", biodataHandler.GetBiodata)        
			biodata.PUT("/biodata", biodataHandler.UpdateBiodata)      
			biodata.DELETE("/biodata", biodataHandler.DeleteBiodata)  
		}

	
	}

	return r
}