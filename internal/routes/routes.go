package routes

import (
	"github.com/gin-contrib/cors"
	"log"

	"github.com/gin-gonic/gin"
)

func Routes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	err := r.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		log.Panicf("Failed to set trusted proxies: %v", err)
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "success api test",
		})
	})

	return r
}