package routes

import (
	"course_online_backend/internal/middleware"
	courseroutes "course_online_backend/internal/routes/Course_management_routes"
	quizroutes "course_online_backend/internal/routes/Course_management_routes/Quiz"
	enrollroutes "course_online_backend/internal/routes/Enrollment_Management_routes"
	"log"
	"os"
	"strings"

	"firebase.google.com/go/v4"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Routes(
	db *gorm.DB,
	app *firebase.App,
	minioClient *minio.Client,
	minioBucket string,
	minioURL string,
	redisClient *redis.Client,
) *gin.Engine {

	corsOrigins := os.Getenv("CORS_ORIGINS")
	allowedOrigins := strings.Split(corsOrigins, ",")

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	

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

	AuthRoutesPublic(public, db, app, redisClient )
	courseroutes.CourseBrowsingRoutes(public, db, app)
	courseroutes.CourseTypePublic(public, db, app)
	courseroutes.LessonPublicRoutes(public, db, app)

	AuthRoutesPrivate(protected, db, app, redisClient )
	UserRoutes(protected, db, app)
	RoleRoutes(protected, db, app)
	BiodataRoutes(protected, db, app)
	ActivityRoutes(protected, db, app)
	DashboardRoutes(protected, db, app)

	courseroutes.CourseManagementRoutes(protected, db, app)
	courseroutes.CourseTypeRoutes(protected, db, app)
	courseroutes.ModuleManagementRoutes(protected, db, app)
	courseroutes.ZoomRoutes(protected, db, app)
	courseroutes.LessonRoutes(protected, db, app)
	quizroutes.QuizManagementRoutes(protected, db, app)


	enrollroutes.EnrollmentRoutes(protected, db, app)

	return r
}
