package database

import (
	"course_online_backend/internal/models"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) error {
	roles := []models.Role{
		{ID: uuid.New(), Name: "student"},
		{ID: uuid.New(), Name: "instructor"},
		{ID: uuid.New(), Name: "admin"},
		{ID: uuid.New(), Name: "super_admin"},
	}

	for _, role := range roles {
		var existingRole models.Role
		result := db.Where("name = ?", role.Name).First(&existingRole)

		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&role).Error; err != nil {
				log.Printf("Failed to seed role %s: %v", role.Name, err)
				return err
			}
			log.Printf("Role '%s' seeded successfully", role.Name)
		} else if result.Error != nil {
			return result.Error
		} else {
			log.Printf("ℹRole '%s' already exists, skipping", role.Name)
		}
	}

	return nil
}

func SeedSuperAdmin(db *gorm.DB) error {
	firebaseUID := os.Getenv("SUPER_ADMIN_FIREBASE_UID")
	email := os.Getenv("SUPER_ADMIN_EMAIL")
	username := os.Getenv("SUPER_ADMIN_USERNAME")

	if firebaseUID == "" || email == "" {
		log.Println("Super admin config not found in .env, skipping...")
		log.Println("Add these to .env:")
		log.Println("SUPER_ADMIN_FIREBASE_UID=your_firebase_uid")
		log.Println("SUPER_ADMIN_EMAIL=admin@example.com")
		log.Println("SUPER_ADMIN_USERNAME=Super Admin")
		return nil
	}

	log.Printf(" Looking for user: %s", email)

	var existingUser models.User
	result := db.Where("firebase_uid = ? OR email_address = ?", firebaseUID, email).First(&existingUser)

	if result.Error == gorm.ErrRecordNotFound {
		log.Printf("User not found, creating new super admin...")

		newUser := models.User{
			ID:           uuid.New(),
			FirebaseUID:  firebaseUID,
			Username:     username,
			EmailAddress: email,
			IsActive:     true,
			LastLogin:    time.Now(),
		}

		if err := db.Create(&newUser).Error; err != nil {
			log.Printf("Failed to create user: %v", err)
			return err
		}

		log.Printf("User created: %s", email)

		var superAdminRole models.Role
		if err := db.Where("name = ?", "super_admin").First(&superAdminRole).Error; err != nil {
			log.Printf(" Super admin role not found: %v", err)
			log.Println(" Make sure you run SeedRoles first!")
			return err
		}

		if err := db.Model(&newUser).Association("Roles").Append(&superAdminRole); err != nil {
			log.Printf(" Failed to assign super admin role: %v", err)
			return err
		}

		log.Printf(" Super admin role assigned to: %s", email)
		log.Printf(" Super admin setup complete!")
		return nil
	}

	if result.Error != nil {
		log.Printf("Database error: %v", result.Error)
		return result.Error
	}

	log.Printf("User found: %s (ID: %s)", existingUser.EmailAddress, existingUser.ID)

	if err := db.Preload("Roles").First(&existingUser, existingUser.ID).Error; err != nil {
		log.Printf("Failed to load user roles: %v", err)
		return err
	}

	log.Printf("Current roles: %d", len(existingUser.Roles))
	for _, role := range existingUser.Roles {
		log.Printf("   - %s", role.Name)
	}

	alreadySuperAdmin := false
	for _, role := range existingUser.Roles {
		if role.Name == "super_admin" {
			alreadySuperAdmin = true
			break
		}
	}

	if alreadySuperAdmin {
		log.Printf("User is already super admin!")
		return nil
	}

	log.Printf("Upgrading user to super admin...")

	var superAdminRole models.Role
	if err := db.Where("name = ?", "super_admin").First(&superAdminRole).Error; err != nil {
		log.Printf("Super admin role not found: %v", err)
		return err
	}

	if err := db.Model(&existingUser).Association("Roles").Clear(); err != nil {
		log.Printf("Failed to clear old roles: %v", err)
		return err
	}

	if err := db.Model(&existingUser).Association("Roles").Append(&superAdminRole); err != nil {
		log.Printf("Failed to assign super admin role: %v", err)
		return err
	}

	log.Printf("User upgraded to super admin!")
	log.Printf("Super admin setup complete!")
	return nil
}

func RunAllSeeders(db *gorm.DB) error {
	log.Println("Starting database seeding...")

	if err := SeedRoles(db); err != nil {
		return err
	}


	if err := SeedSuperAdmin(db); err != nil {
		return err
	}

	log.Println("All seeders completed!")
	return nil
}
