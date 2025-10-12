package main

import (
	"course_online_backend/database"
	"course_online_backend/internal/config"
	"course_online_backend/internal/routes"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConfigDb()
	db := database.PostgresConn()
	app := config.InitFirebase()

	r := routes.Routes(db, app)

	if config.DbConfig.ServerEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	serverAddr := fmt.Sprintf("%s:%s",
		config.DbConfig.ServerHost,
		config.DbConfig.ServerPort,
	)

	log.Printf("Server running on http://%s", serverAddr)
	log.Printf("Environment: %s", config.DbConfig.ServerEnv)

	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
