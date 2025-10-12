package database

import (
	"course_online_backend/internal/models"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) error {
    roles := []models.Role{
        {ID: uuid.New(), Name: "user"},
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
			log.Printf("✓ Role '%s' seeded successfully", role.Name)
		} else if result.Error != nil {
			return result.Error
		} else {
			log.Printf("○ Role '%s' already exists, skipping", role.Name)
		}
	}

	return nil
}
