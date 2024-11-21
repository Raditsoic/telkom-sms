package service

import (
	"encoding/json"

	"strconv"

	"gtihub.com/raditsoic/telkom-storage-ms/src/database/repository"
	"gtihub.com/raditsoic/telkom-storage-ms/src/model"
	"gtihub.com/raditsoic/telkom-storage-ms/src/utils"
)

type StorageService struct {
	repository repository.StorageRepository
}

func NewStorageService(repo repository.StorageRepository) *StorageService {
	return &StorageService{repository: repo}
}

// Get All Storages
func (service *StorageService) GetStorages() ([]model.Storage, error) {
	return service.repository.GetStorages()
}

// Create Storage
func (service *StorageService) CreateStorage(storageData []byte) (*model.CreateStorageResponse, error) {
	var storage model.Storage
	if err := json.Unmarshal(storageData, &storage); err != nil {
		return nil, err
	}

	createdStorage, err := service.repository.CreateStorage(storage)
	if err != nil {
		return nil, err
	}

	response := &model.CreateStorageResponse{
		Message: "Storage created successfully",
		ID:      strconv.Itoa(createdStorage.ID),
		Name:    createdStorage.Name,
	}

	return response, nil
}

// Get Storage By ID
func (service *StorageService) GetStorageByID(id int) (*model.StorageByIDResponse, error) {

	storage, err := service.repository.GetStorageByIDwithCategories(id)
	if err != nil {
		return nil, err
	}

	var categories []model.StorageCategoryResponse
	for _, category := range storage.Categories {
		categories = append(categories, model.StorageCategoryResponse{
			ID:    category.ID,
			Name:  category.Name,
			Image: category.Image,
		})
	}

	return &model.StorageByIDResponse{
		ID:         storage.ID,
		Name:       storage.Name,
		Location:   storage.Location,
		Categories: categories,
	}, nil
}

// Get Storage By ID No Image
func (service *StorageService) GetStorageByIDNoImage(id int) (*model.StorageByIDResponseNoImage, error) {

	storage, err := service.repository.GetStorageByIDwithCategories(id)
	if err != nil {
		return nil, err
	}

	var categories []model.StorageCategoryNoImageResponse
	for _, category := range storage.Categories {
		categories = append(categories, model.StorageCategoryNoImageResponse{
			ID:   category.ID,
			Name: category.Name,
		})
	}

	return &model.StorageByIDResponseNoImage{
		ID:         storage.ID,
		Name:       storage.Name,
		Location:   storage.Location,
		Categories: categories,
	}, nil
}

// Delete Storage
func (service *StorageService) DeleteStorage(id string) (*model.DeleteStorageResponse, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	_, err = service.repository.GetStorageByID(idInt)
	if err != nil {
		return nil, utils.ErrStorageNotFound
	}

	if err := service.repository.DeleteStorage(id); err != nil {
		return nil, err
	}

	response := &model.DeleteStorageResponse{
		Message: "Storage deleted successfully",
		ID:      id,
	}

	return response, nil
}
