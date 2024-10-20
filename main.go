package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gtihub.com/raditsoic/telkom-storage-ms/database"
	"gtihub.com/raditsoic/telkom-storage-ms/model"
	"gtihub.com/raditsoic/telkom-storage-ms/service"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&model.Admin{}, &model.Storage{}, &model.Item{}, &model.Category{}, &model.CategoryWithStorage{}); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}

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
	r.HandleFunc("/api/categories", service.GetCategories).Methods("GET")
	r.HandleFunc("/api/category", service.CreateCategory).Methods("POST")

	// Storage routes
	r.HandleFunc("/api/storage", service.CreateStorage).Methods("POST")
	r.HandleFunc("/api/storages", service.GetStorages).Methods("GET")
	r.HandleFunc("/api/storage/{id}", service.GetStorageByID).Methods("GET")

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
