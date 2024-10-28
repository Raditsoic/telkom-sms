package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gtihub.com/raditsoic/telkom-storage-ms/database"
	"gtihub.com/raditsoic/telkom-storage-ms/database/repository"
	"gtihub.com/raditsoic/telkom-storage-ms/model"
	"gtihub.com/raditsoic/telkom-storage-ms/service"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&model.Admin{}, &model.Storage{}, &model.Item{}, &model.Category{}); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}

	CategoryRepository := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(*CategoryRepository)
	StorageRepository := repository.NewStorageRepository(db)
	StorageService := service.NewStorageService(*StorageRepository)

	r := mux.NewRouter()

	// Root route
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hello, World!")); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}).Methods("GET")

	// Admin routes
	r.HandleFunc("/api/admin", service.AdminLogin).Methods("POST")

	// Item routes
	r.HandleFunc("/api/items", service.GetItems).Methods("GET")
	r.HandleFunc("/api/item", service.CreateItem).Methods("POST")

	// Category routes
	r.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		limit := r.URL.Query().Get("limit")

		categories, err := categoryService.GetCategories(page, limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(categories)
	}).Methods("GET")
	r.HandleFunc("/api/category", func(w http.ResponseWriter, r *http.Request) {
		var category model.Category
		if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		categoryData, err := json.Marshal(category)
		if err != nil {
			http.Error(w, "Failed to marshal category", http.StatusInternalServerError)
			return
		}
		createdCategory, err := categoryService.CreateCategory(categoryData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdCategory)
	}).Methods("POST")

	// Storage routes
	r.HandleFunc("/api/storage", func(w http.ResponseWriter, r *http.Request) {
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
		createdStorage, err := StorageService.CreateStorage(storageData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdStorage)
	}).Methods("POST")
	r.HandleFunc("/api/storages", func(w http.ResponseWriter, r *http.Request) {
		storages, err := StorageService.GetStorages()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(storages)
	}).Methods("GET")
	// r.HandleFunc("/api/storage/{id}", service.GetStorageByID).Methods("GET")

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})

	// Apply CORS middleware
	handler := c.Handler(r)

	// Start the server
	fmt.Println("Server is starting on :8080...")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
