package services

import (
	"errors"
	"course_online_backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRoleService struct {
	DB *gorm.DB
}

func NewUserRoleService(db *gorm.DB) *UserRoleService {
	return &UserRoleService{DB: db}
}

func (s *UserRoleService) AssignRole(targetUserID, roleID, currentUserID uuid.UUID) error {
	var currentUser models.User
	if err := s.DB.Preload("Roles").First(&currentUser, currentUserID).Error; err != nil {
		return errors.New("failed to load current user")
	}

	var newRole models.Role
	if err := s.DB.First(&newRole, roleID).Error; err != nil {
		return errors.New("role not found")
	}

	var targetUser models.User
	if err := s.DB.Preload("Roles").First(&targetUser, targetUserID).Error; err != nil {
		return errors.New("target user not found")
	}

	if targetUser.ID == currentUser.ID {
		return errors.New("cannot change your own role")
	}

	currentUserRole := currentUser.GetRoleName()

	switch currentUserRole {
	case "super_admin":
	case "admin":
		if newRole.Name == "super_admin" || newRole.Name == "admin" {
			return errors.New("admin cannot assign super_admin or admin roles")
		}
		if targetUser.HasRole("super_admin") {
			return errors.New("admin cannot modify super admin user")
		}
	default:
		return errors.New("insufficient permissions to assign roles")
	}


	tx := s.DB.Begin()

	if err := tx.Model(&targetUser).Association("Roles").Clear(); err != nil {
		tx.Rollback()
		return errors.New("failed to clear existing roles")
	}

	if err := tx.Model(&targetUser).Association("Roles").Append(&newRole); err != nil {
		tx.Rollback()
		return errors.New("failed to assign new role")
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New("failed to commit transaction")
	}

	return nil
}




func (s *UserRoleService) GetAllRoles() ([]models.Role, error) {
	var roles []models.Role
	if err := s.DB.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (s *UserRoleService) GetAssignableRoles(currentUser models.User) ([]models.Role, error) {
	var roles []models.Role

	roleName := currentUser.GetRoleName()

	switch roleName {
	case "super_admin":
		if err := s.DB.Find(&roles).Error; err != nil {
			return nil, err
		}
	case "admin":
		if err := s.DB.Where("name IN ?", []string{"instructor", "student"}).Find(&roles).Error; err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("you do not have permission to assign roles")
	}

	return roles, nil
}

func (s *UserRoleService) GetRoleByID(roleID uuid.UUID) (*models.Role, error) {
	var role models.Role
	if err := s.DB.First(&role, "id = ?", roleID).Error; err != nil {
		return nil, err
	}
	return &role, nil
}
