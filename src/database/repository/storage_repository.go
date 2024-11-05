package repository

import (
	"fmt"

	"gorm.io/gorm"
	"gtihub.com/raditsoic/telkom-storage-ms/src/model"
)

type StorageRepository struct {
	db *gorm.DB
}

func NewStorageRepository(db *gorm.DB) *StorageRepository {
	return &StorageRepository{db: db}
}

func (repo *StorageRepository) CreateStorage(storage model.Storage) error {
	if err := repo.db.Create(&storage).Error; err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}

	return nil
}

func (repo *StorageRepository) GetStorageByID(id int) (*model.Storage, error) {
	var storage model.Storage
	if err := repo.db.First(&storage, id).Error; err != nil {
		return &model.Storage{}, fmt.Errorf("failed to get storage: %w", err)
	}
	return &storage, nil
}

func (repo *StorageRepository) GetStorageByIDwithCategories(id int) (*model.Storage, error) {
	var storage model.Storage
	if err := repo.db.Preload("Categories").First(&storage, id).Error; err != nil {
		return &storage, err
	}
	return &storage, nil
}

func (repo *StorageRepository) GetStorages() ([]model.Storage, error) {
	var storages []model.Storage
	if err := repo.db.Preload("Categories").Find(&storages).Error; err != nil {
		return nil, fmt.Errorf("failed to get storages: %w", err)
	}

	return storages, nil
}

func (repo *StorageRepository) UpdateStorage(storage model.Storage) error {
	if err := repo.db.Save(&storage).Error; err != nil {
		return fmt.Errorf("failed to update storage: %w", err)
	}

	return nil
}

func (repo *StorageRepository) DeleteStorage(id int) error {
	if err := repo.db.Where("id = ?", id).Delete(&model.Storage{}).Error; err != nil {
		return fmt.Errorf("failed to delete storage: %w", err)
	}

	return nil
}
