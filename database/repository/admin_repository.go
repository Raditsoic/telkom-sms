package repository

import (
	"fmt"

	"gtihub.com/raditsoic/telkom-storage-ms/database"
	"gtihub.com/raditsoic/telkom-storage-ms/model"
)

func GetAdminByUsername(username string) (*model.Admin, error) {
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

func RegisterAdmin(admin *model.Admin) error {
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Create(admin).Error; err != nil {
		return fmt.Errorf("failed to create admin: %w", err)
	}
	return nil
}

func GetAdmins() ([]model.Admin, error) {
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

func DeleteAdmin(adminID string) error {
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Where("id = ?", adminID).Delete(&model.Admin{}).Error; err != nil {
		return fmt.Errorf("failed to delete admin: %w", err)
	}
	return nil
}
