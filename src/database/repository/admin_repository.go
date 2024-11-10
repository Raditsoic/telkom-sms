package repository

import (
	"fmt"

	"gorm.io/gorm"
	"gtihub.com/raditsoic/telkom-storage-ms/src/database"
	"gtihub.com/raditsoic/telkom-storage-ms/src/model"
)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (repository *AdminRepository) GetAdminByUsername(username string) (*model.Admin, error) {
	db, err := database.Connect()
	if err != nil {
		return &model.Admin{}, fmt.Errorf("failed to connect to database: %w", err)
	}

	var admin model.Admin
	if err := db.Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, fmt.Errorf("failed to get admin: %w", err)
	}
	return &admin, nil
}

func (repository *AdminRepository) RegisterAdmin(admin *model.Admin) error {
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Create(admin).Error; err != nil {
		return fmt.Errorf("failed to create admin: %w", err)
	}
	return nil
}

func (repository *AdminRepository) GetAdmins() ([]model.Admin, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	var admins []model.Admin
	if err := db.Find(&admins).Error; err != nil {
		return nil, fmt.Errorf("failed to get admins: %w", err)
	}
	return admins, nil
}

func(repository *AdminRepository) DeleteAdmin(adminID string) error {
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Where("id = ?", adminID).Delete(&model.Admin{}).Error; err != nil {
		return fmt.Errorf("failed to delete admin: %w", err)
	}
	return nil
}
