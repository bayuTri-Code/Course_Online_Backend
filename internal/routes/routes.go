package routes

import (
	"course_online_backend/internal/middleware"
	"log"
	"os"
	"strings"

	"firebase.google.com/go/v4"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

func Routes(db *gorm.DB, app *firebase.App, minioClient *minio.Client, minioBucket, minioURL string) *gin.Engine {

	corsOrigins := os.Getenv("CORS_ORIGINS")
	allowedOrigins := strings.Split(corsOrigins, ",")

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Panicf("Failed to set trusted proxies: %v", err)
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "success", "message": "Course Online Backend API Success"})
	})

	public := r.Group("/api")
	protected := r.Group("/api")
	protected.Use(middleware.FirebaseAuth(app, db))

	// public routes
	AuthRoutes(public, db, app)
	CoursePublicRoutes(public, db, app)

	//private routes
	UserRoutes(protected, db, app)
	RoleRoutes(protected, db, app)
	BiodataRoutes(protected, db, app)
	ActivityRoutes(protected, db, app)
	DashboardRoutes(protected, db, app)
	CourseProtectedRoutes(protected, db, app)
	EnrollmentProtectedRoutes(protected, db, app)

	return r
}
