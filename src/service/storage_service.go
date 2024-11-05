package service

import (
	"encoding/json"

	"gtihub.com/raditsoic/telkom-storage-ms/src/database/repository"
	"gtihub.com/raditsoic/telkom-storage-ms/src/model"
)

type StorageService struct {
	repository repository.StorageRepository
}

func NewStorageService(repo repository.StorageRepository) *StorageService {
	return &StorageService{repository: repo}
}

// GetStorages handles GET requests to retrieve storage items.
func (service *StorageService) GetStorages() ([]model.Storage, error) {
	return service.repository.GetStorages()
}

// CreateStorage handles POST requests to create a new storage item.
func (service *StorageService) CreateStorage(storageData []byte) (*model.Storage, error) {
	var storage model.Storage
	if err := json.Unmarshal(storageData, &storage); err != nil {
		return nil, err
	}

	if err := service.repository.CreateStorage(storage); err != nil {
		return nil, err
	}

	return &storage, nil
}

func (service *StorageService) DeleteStorage(id int) error {
	_, err := service.repository.GetStorageByID(id)
	if err != nil {
		return err
	}

	return service.repository.DeleteStorage(id)
}

func (service *StorageService) GetStorageByID(id int) (*model.StorageByIDResponse, error) {

	storage, err := service.repository.GetStorageByIDwithCategories(id)
	if err != nil {
		return nil, err
	}

	var categories []model.StorageCategoryResponse
	for _, category := range storage.Categories {
		categories = append(categories, model.StorageCategoryResponse{
			ID:        category.ID,
			Name:      category.Name,
		})
	}

	return &model.StorageByIDResponse{
		ID:         storage.ID,
		Name:       storage.Name,
		Location:   storage.Location,
		Categories: categories,
	}, nil
}
