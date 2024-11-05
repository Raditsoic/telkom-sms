package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gtihub.com/raditsoic/telkom-storage-ms/src/model"
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
	// Check if storage exists
	var storage model.Storage
	if err := repo.db.First(&storage, category.StorageID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("storage with ID %d not found", category.StorageID)
		}
		return fmt.Errorf("failed to check storage: %w", err)
	}

	// Create category
	if err := repo.db.Create(&category).Error; err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}

	return nil
}

func (repo *CategoryRepository) GetCategoryByID(id string) (*model.CategoryByIDResponse, error) {
	var category model.CategoryByIDResponse
	if err := repo.db.Model(&model.Category{}).
		First(&category, id).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}

	return &category, nil
}

func (repo *CategoryRepository) GetCategories(limit, offset int) ([]model.AllCategoryResponse, error) {
	var categories []model.AllCategoryResponse

	if err := repo.db.Model(&model.Category{}).
		Preload("Storage").
		Limit(limit).
		Offset(offset).
		Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}

	return categories, nil
}

func (repo *CategoryRepository) GetCategoryWithItems(categoryID uint) (*model.Category, error) {
	var category model.Category

	if err := repo.db.Preload("Items").Preload("Storage").
		First(&category, categoryID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch category: %w", err)
	}

	return &category, nil
}

// func (repo *CategoryRepository) GetCategoriesWithStorage() ([]model.CategoryWithStorage, error) {
// 	db, err := database.Connect()
// 	if err != nil {
// 		return []model.CategoryWithStorage{}, fmt.Errorf("failed to connect to database: %w", err)
// 	}

// 	var categories []model.CategoryWithStorage
// 	if err := db.Preload("Storage").Find(&categories).Error; err != nil {
// 		return nil, fmt.Errorf("failed to get categories: %w", err)
// 	}

// 	return categories, nil
// }

func (repo *CategoryRepository) UpdateCategory(category model.Category) error {
	if err := repo.db.Save(&category).Error; err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}

	return nil
}

func (repo *CategoryRepository) DeleteCategory(id int) error {
	if err := repo.db.Where("id = ?", id).Delete(&model.Category{}).Error; err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}
