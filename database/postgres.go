package database

import (
	"course_online_backend/internal/config"
	"course_online_backend/internal/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func PostgresConn() *gorm.DB {
	configDb := config.DbConfig

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		configDb.DBHost,
		configDb.DBPort,
		configDb.DBUser,
		configDb.DBPassword,
		configDb.DBName,
		configDb.DBSslmode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = db.SetupJoinTable(&models.User{}, "Roles", &models.UserRole{})
	if err != nil {
		log.Fatalf("failed to setup join table User-Roles: %v", err)
	}

	err = db.SetupJoinTable(&models.Role{}, "Users", &models.UserRole{})
	if err != nil {
		log.Fatalf("failed to setup join table Role-Users: %v", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.UserRole{},
		&models.Biodata{},
	)
	if err != nil {
		log.Fatalf("auto migration failed: %v", err)
	}

	if err := SeedRoles(db); err != nil {
		log.Fatalf("failed to seed roles: %v", err)
	}

	log.Println("Database connected and migrated successfully.")
	return db
}