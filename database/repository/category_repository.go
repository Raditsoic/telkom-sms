package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gtihub.com/raditsoic/telkom-storage-ms/database"
	"gtihub.com/raditsoic/telkom-storage-ms/model"
)

type CategoryWithStorage struct {
	CategoryID   uint   `json:"category_id"`
	CategoryName string `json:"category_name"`
	StorageName  string `json:"storage_name"`
}

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (repo *CategoryRepository) CreateCategory(category model.Category) error {
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Check if storage exists
	var storage model.Storage
	if err := db.First(&storage, category.StorageID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("storage with ID %d not found", category.StorageID)
		}
		return fmt.Errorf("failed to check storage: %v", err)
	}

	// Create category
	if err := db.Create(&category).Error; err != nil {
		return fmt.Errorf("failed to create category: %v", err)
	}

	return nil
}

func (repo *CategoryRepository) GetCategoryByID(id int) (*model.Category, error) {
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

func (repo *CategoryRepository) GetCategories(limit, offset int) ([]model.Category, error) {
	var categories []model.Category

	if err := repo.db.Preload("Storage").Limit(limit).Offset(offset).Find(&categories).Error; err != nil {
        return nil, fmt.Errorf("failed to get categories: %v", err)
    }

	return categories, nil
}

// func (repo *CategoryRepository) GetCategoriesWithStorage() ([]model.CategoryWithStorage, error) {
// 	db, err := database.Connect()
// 	if err != nil {
// 		return []model.CategoryWithStorage{}, fmt.Errorf("failed to connect to database: %v", err)
// 	}

// 	var categories []model.CategoryWithStorage
// 	if err := db.Preload("Storage").Find(&categories).Error; err != nil {
// 		return nil, fmt.Errorf("failed to get categories: %v", err)
// 	}

// 	return categories, nil
// }

func (repo *CategoryRepository) UpdateCategory(category model.Category) error {
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.Save(&category).Error; err != nil {
		return fmt.Errorf("failed to update category: %v", err)
	}

	return nil
}

func (repo *CategoryRepository) DeleteCategory(id int) error {
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.Where("id = ?", id).Delete(&model.Category{}).Error; err != nil {
		return fmt.Errorf("failed to delete category: %v", err)
	}

	return nil
}
