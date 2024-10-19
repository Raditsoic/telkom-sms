package service

import (
	"encoding/json"
	"net/http"

	"gtihub.com/raditsoic/telkom-storage-ms/database/repository"
	"gtihub.com/raditsoic/telkom-storage-ms/model"
)

// GetStorages handles GET requests to retrieve storage items.
func GetStorages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	storages, err := repository.GetStorages()
	if err != nil {
		http.Error(w, "Failed to retrieve storages: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the content type header and encode storages to JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(storages); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateStorage handles POST requests to create a new storage item.
func CreateStorage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var storage model.Storage
	if err := json.NewDecoder(r.Body).Decode(&storage); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Call the repository to create the storage item
	if err := repository.CreateStorage(storage); err != nil {
		http.Error(w, "Failed to create storage: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with 201 Created status
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "Storage created successfully"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}
