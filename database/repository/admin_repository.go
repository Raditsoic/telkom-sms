package repository

import (
	"fmt"

	"gtihub.com/raditsoic/telkom-storage-ms/database"
	"gtihub.com/raditsoic/telkom-storage-ms/model"
)

func GetAdminByUsername(username string) (*model.Admin, error) {
	db, err := database.Connect()
	if err != nil {
		return &model.Admin{}, fmt.Errorf("failed to connect to database: %v", err)
	}

	var admin model.Admin
	if err := db.Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, fmt.Errorf("failed to get admin: %v", err)
	}
	return &admin, nil
}

func RegisterAdmin(admin *model.Admin) error {
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.Create(admin).Error; err != nil {
		return fmt.Errorf("failed to create admin: %v", err)
	}
	return nil
}

func GetAdmins() ([]model.Admin, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	var admins []model.Admin
	if err := db.Find(&admins).Error; err != nil {
		return nil, fmt.Errorf("failed to get admins: %v", err)
	}
	return admins, nil
}

func DeleteAdmin(adminID string) error {
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.Where("id = ?", adminID).Delete(&model.Admin{}).Error; err != nil {
		return fmt.Errorf("failed to delete admin: %v", err)
	}
	return nil
}
