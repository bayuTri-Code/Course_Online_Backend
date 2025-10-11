package main

import (
	"course_online_backend/database"
	"course_online_backend/internal/config"
	"course_online_backend/internal/routes"
	"fmt"
)

func main() {
	config.ConfigDb()
	database.PostgresConn()

	r := routes.Routes()
	host := "0.0.0.0"
	port := "7070"

	fmt.Printf("server is running in http://%s:%s\n", host, port)
	r.Run(host + ":" + port)
}