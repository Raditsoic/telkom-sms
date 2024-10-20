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
