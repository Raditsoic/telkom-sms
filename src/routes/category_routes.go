package routes

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gtihub.com/raditsoic/telkom-storage-ms/src/middleware"
	"gtihub.com/raditsoic/telkom-storage-ms/src/model"
	"gtihub.com/raditsoic/telkom-storage-ms/src/service"
	"gtihub.com/raditsoic/telkom-storage-ms/src/utils"
)

func CategoryRoutes(r *mux.Router, categoryService *service.CategoryService, jwtUtils *utils.JWTUtils) {
	r.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		limit := r.URL.Query().Get("limit")

		response, err := categoryService.GetCategories(page, limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}).Methods("GET")

	r.Handle("/api/category", middleware.AuthMiddleware(jwtUtils, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		storage_id := r.FormValue("storage_id")

		if name == "" || storage_id == "" {
			http.Error(w, "Name and storage_id are required", http.StatusBadRequest)
			return
		}

		storageID, err := strconv.ParseUint(storage_id, 10, 32)
		if err != nil {
			http.Error(w, "Invalid storage_id format", http.StatusBadRequest)
			return
		}

		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Image file is required", http.StatusBadRequest)
			return
		}
		defer file.Close()

		imageData, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Could not read image file", http.StatusInternalServerError)
			return
		}

		category := model.Category{
			Name:      name,
			StorageID: uint(storageID),
			Image:     imageData,
		}

		response, err := categoryService.CreateCategory(&category)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}))).Methods("POST")

	r.HandleFunc("/api/category/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		response, err := categoryService.GetCategoryByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}).Methods("GET")

	r.HandleFunc("/api/category/{id}/items", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		categoryID, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		response, err := categoryService.GetCategoryWithItems(uint(categoryID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}).Methods("GET")

	r.Handle("/api/category/{id}", middleware.AuthMiddleware(jwtUtils, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		response, err := categoryService.DeleteCategory(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}))).Methods("DELETE")

	r.Handle("/api/category/{id}", middleware.AuthMiddleware(jwtUtils, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var req model.UpdateCategoryNameRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		response, err := categoryService.UpdateCategoryName(id, req.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response = &model.UpdateCategoryNameResponse{
			Message: "Category updated successfully",
			ID:      id,
			NewName: response.NewName,
			OldName: response.OldName,
		}


		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}))).Methods("PATCH")
}
