package routes

import (
	"course_online_backend/internal/handler"
	"course_online_backend/internal/handler/dashboard"
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
		c.JSON(200, gin.H{"status": "success", "message": "Course Online Backend API Success"})
	})

	authService := services.NewAuthService(db, app)
	userService := services.NewUserService(db)
	activityService := services.NewActivityService(db)
	biodataService := services.NewBiodataService(db)
	dashboardService := services.NewDashboardService(db)

	authHandler := handler.NewAuthHandler(authService, activityService)
	userHandler := handler.NewUserHandler(userService, activityService)
	biodataHandler := handler.NewBiodataHandler(biodataService, activityService)
	activityHandler := handler.NewActivityHandler(activityService)

	superAdminHandler := dashboard.NewSuperAdminDashboardHandler(dashboardService)
	adminHandler := dashboard.NewAdminDashboardHandler(dashboardService)
	instructorHandler := dashboard.NewInstructorDashboardHandler(dashboardService)
	studentHandler := dashboard.NewStudentDashboardHandler(dashboardService)

	public := r.Group("/api")
	{
		auth := public.Group("/auth")
		auth.POST("/firebase-login", authHandler.FirebaseLoginHandler)
	}

	protected := r.Group("/api")
	protected.Use(middleware.FirebaseAuth(app, db))
	{
		biodata := protected.Group("/profile")
		biodata.Use(middleware.RoleMiddleware("student", "super_admin", "admin", "instructor"))
		{
			biodata.POST("/biodata", biodataHandler.CreateBiodata)
			biodata.GET("/mybiodata", biodataHandler.GetBiodata)
			biodata.PUT("/biodata", biodataHandler.UpdateBiodata)
			biodata.DELETE("/biodata", biodataHandler.DeleteBiodata)
		}

	
		activity := protected.Group("/history")
		activity.Use(middleware.RoleMiddleware("admin", "super_admin"))
		{
			activity.GET("/All-activity", activityHandler.GetAllActivity)
			activity.GET("/ByUser-activity/:id", activityHandler.GetActivityByUserId)
			activity.GET("/Recent-activity", activityHandler.GetRecentActivities)
			activity.GET("/Search-activity", activityHandler.SearchActivity)
			activity.GET("/Summary-activity/:id", activityHandler.GetActivitySummary)
			activity.GET("/ByRole-activity/:role", activityHandler.GetActivityByRole)
		}

		user := protected.Group("/user")
		{
			user.GET("/", userHandler.GetAllUsers)
			user.GET("/:id", userHandler.GetUserByID)
			user.PUT("/:id", userHandler.UpdateUser)
			user.DELETE("/:id", userHandler.DeleteUser)
		}

		dashboard := protected.Group("/dashboard")
		{
			dashboard.GET("/super_admin", middleware.RoleMiddleware("super_admin"), superAdminHandler.GetDashboardSuperAdmin)
			dashboard.GET("/admin", middleware.RoleMiddleware("admin"), adminHandler.GetDashboardAdmin)
			dashboard.GET("/instructor", middleware.RoleMiddleware("instructor"), instructorHandler.GetDashboardInstructor)
			dashboard.GET("/student", middleware.RoleMiddleware("student"), studentHandler.GetDashboardStudent)
		}
	}

	return r
}
