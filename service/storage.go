package service

import (
	"encoding/json"

	"gtihub.com/raditsoic/telkom-storage-ms/database/repository"
	"gtihub.com/raditsoic/telkom-storage-ms/model"
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

// func (service *StorageService) GetStorageByID(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	idStr := vars["id"]

// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid ID format", http.StatusBadRequest)
// 		return
// 	}

// 	storage, err := repository.GetStorageByIDwithCategories(id)
// 	if err != nil {
// 		http.Error(w, "Failed to retrieve storage: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(storage); err != nil {
// 		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }
