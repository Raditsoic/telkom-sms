package repository

import (
	"fmt"

	"gtihub.com/raditsoic/telkom-storage-ms/database"
	"gtihub.com/raditsoic/telkom-storage-ms/model"
)

func CreateItem(item model.Item) error {
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.Create(&item).Error; err != nil {
		return fmt.Errorf("failed to create item: %v", err)
	}

	return nil
}

func GetItemByID(id int) (*model.Item, error) {
	db, err := database.Connect()
	if err != nil {
		return &model.Item{}, fmt.Errorf("failed to connect to database: %v", err)
	}

	var item model.Item
	if err := db.Where("id = ?", id).First(&item).Error; err != nil {
		return nil, fmt.Errorf("failed to get item: %v", err)
	}

	return &item, nil
}

func GetItems() ([]model.Item, error) {
	db, err := database.Connect()
	if err != nil {
		return []model.Item{}, fmt.Errorf("failed to connect to database: %v", err)
	}

	var items []model.Item
	if err := db.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get items: %v", err)
	}

	return items, nil
}

func UpdateItem(item model.Item) error {
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.Save(&item).Error; err != nil {
		return fmt.Errorf("failed to update item: %v", err)
	}

	return nil
}

func DeleteItem(id int) error {
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.Where("id = ?", id).Delete(&model.Item{}).Error; err != nil {
		return fmt.Errorf("failed to delete item: %v", err)
	}

	return nil
}
