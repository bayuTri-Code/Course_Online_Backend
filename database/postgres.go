package database

import (
	"course_online_backend/internal/config"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func PostgresConn(){
	configDb := config.DbConfig

	SetDb := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		configDb.DBHost,
		configDb.DBPort,
		configDb.DBUser,
		configDb.DBPassword,
		configDb.DBName,
		configDb.DBSslmode,
	)
	db, err := gorm.Open(postgres.Open(SetDb), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	Db = db
	log.Println("Database Connected")
	autoMigrate()
}

func autoMigrate() {
	err := Db.AutoMigrate(
		
	)

	if err != nil {
		log.Fatalf("Auto Migration Failed: %v", err)
	}
	log.Println("Auto Migration Complete!")
}