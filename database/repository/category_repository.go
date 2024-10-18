package repository

import (
	"fmt"

	"gtihub.com/raditsoic/telkom-storage-ms/database"
	"gtihub.com/raditsoic/telkom-storage-ms/model"
)

func CreateCategory(category model.Category) error {
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.Create(&category).Error; err != nil {
		return fmt.Errorf("failed to create category: %v", err)
	}

	return nil
}

func GetCategoryByID(id int) (*model.Category, error) {
	db, err := database.Connect()
	if err != nil {
		return &model.Category{}, fmt.Errorf("failed to connect to database: %v", err)
	}

	var category model.Category
	if err := db.Where("id = ?", id).First(&category).Error; err != nil {
		return nil, fmt.Errorf("failed to get category: %v", err)
	}

	return &category, nil
}

func GetCategories() ([]model.Category, error) {
	db, err := database.Connect()
	if err != nil {
		return []model.Category{}, fmt.Errorf("failed to connect to database: %v", err)
	}

	var categories []model.Category
	if err := db.Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to get categories: %v", err)
	}

	return categories, nil
}

func UpdateCategory(category model.Category) error {
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.Save(&category).Error; err != nil {
		return fmt.Errorf("failed to update category: %v", err)
	}

	return nil
}

func DeleteCategory(id int) error {
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.Where("id = ?", id).Delete(&model.Category{}).Error; err != nil {
		return fmt.Errorf("failed to delete category: %v", err)
	}

	return nil
}

