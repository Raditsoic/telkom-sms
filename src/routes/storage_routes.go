package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gtihub.com/raditsoic/telkom-storage-ms/src/model"
	"gtihub.com/raditsoic/telkom-storage-ms/src/service"
	"gtihub.com/raditsoic/telkom-storage-ms/src/utils"
)

func StorageRoutes(r *mux.Router, storageService *service.StorageService, jwtUtils *utils.JWTUtils) {
	r.HandleFunc("/api/storage", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var storage model.Storage
		if err := json.NewDecoder(r.Body).Decode(&storage); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		storageData, err := json.Marshal(storage)
		if err != nil {
			http.Error(w, "Failed to marshal category", http.StatusInternalServerError)
			return
		}
		createdStorage, err := storageService.CreateStorage(storageData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := struct {
			Message string `json:"message"`
			StorageName 	string    `json:"storage_name"`
		} {
			Message: "Storage created successfully",
			StorageName: createdStorage.Name,
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})).Methods("POST")

	r.HandleFunc("/api/storages", func(w http.ResponseWriter, r *http.Request) {
		storages, err := storageService.GetStorages()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(storages); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}).Methods("GET")

	// r.Handle("/api/storage/{id}", middleware.AuthMiddleware(jwtUtils, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	vars := mux.Vars(r)
	// 	id := vars["id"]
	// 	idInt, err := strconv.Atoi(id)
	// 	if err != nil {
	// 		http.Error(w, "Invalid storage ID", http.StatusBadRequest)
	// 		return
	// 	}

	// 	if err := storageService.DeleteStorage(idInt); err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}

	// 	w.WriteHeader(http.StatusOK)
	// 	if err := json.NewEncoder(w).Encode("Item deleted"); err != nil {
	// 		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	// 		return
	// 	}
	// }))).Methods("DELETE")

	r.HandleFunc("/api/storage/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		response, err := storageService.DeleteStorage(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})).Methods("DELETE")

	r.HandleFunc("/api/storage/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid storage ID", http.StatusBadRequest)
			return
		}

		storage, err := storageService.GetStorageByID(idInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(storage); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}).Methods("GET")

	r.HandleFunc("/api/storage/{id}/no-image", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid storage ID", http.StatusBadRequest)
			return
		}

		storage, err := storageService.GetStorageByIDNoImage(idInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(storage); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}).Methods("GET")
}
