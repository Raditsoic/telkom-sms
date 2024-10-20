package repository

import (
	"fmt"

	"gtihub.com/raditsoic/telkom-storage-ms/database"
	"gtihub.com/raditsoic/telkom-storage-ms/model"
)

func CreateStorage(storage model.Storage) error {
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.Create(&storage).Error; err != nil {
		return fmt.Errorf("failed to create storage: %v", err)
	}

	return nil
}

func GetStorageByID(id int) (*model.Storage, error) {
	db, err := database.Connect()
	if err != nil {
		return &model.Storage{}, fmt.Errorf("failed to connect to database: %v", err)
	}

	var storage model.Storage
	if err := db.First(&storage, id).Error; err != nil {
		return &model.Storage{}, fmt.Errorf("failed to get storage: %v", err)
	}
	return &storage, nil
}

func GetStorageByIDwithCategories(id int) (*model.Storage, error) {
	db, err := database.Connect()
	if err != nil {
		return &model.Storage{}, fmt.Errorf("failed to connect to database: %v", err)
	}

	var storage model.Storage
	if err := db.Preload("Categories.Items").First(&storage, id).Error; err != nil {
		return &storage, err
	}
	return &storage, nil
}

func GetStorages() ([]model.Storage, error) {
	db, err := database.Connect()
	if err != nil {
		return []model.Storage{}, fmt.Errorf("failed to connect to database: %v", err)
	}

	var storages []model.Storage
	if err := db.Find(&storages).Error; err != nil {
		return nil, fmt.Errorf("failed to get storages: %v", err)
	}

	return storages, nil
}

func UpdateStorage(storage model.Storage) error {
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.Save(&storage).Error; err != nil {
		return fmt.Errorf("failed to update storage: %v", err)
	}

	return nil
}

func DeleteStorage(id int) error {
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.Where("id = ?", id).Delete(&model.Storage{}).Error; err != nil {
		return fmt.Errorf("failed to delete storage: %v", err)
	}

	return nil
}

