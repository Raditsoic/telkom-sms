package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gtihub.com/raditsoic/telkom-storage-ms/src/database"
	"gtihub.com/raditsoic/telkom-storage-ms/src/database/repository"
	"gtihub.com/raditsoic/telkom-storage-ms/src/routes"
	"gtihub.com/raditsoic/telkom-storage-ms/src/service"
	"gtihub.com/raditsoic/telkom-storage-ms/src/utils"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	jwtUtils := utils.NewJWTUtils()

	AuthService := service.NewAuthService(*repository.NewAdminRepository(db), jwtUtils)

	CategoryRepository := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(*CategoryRepository)

	StorageRepository := repository.NewStorageRepository(db)
	StorageService := service.NewStorageService(*StorageRepository)

	ItemRepository := repository.NewItemRepository(db)
	itemService := service.NewItemService(*ItemRepository)

	TransactionRepository := repository.NewTransactionRepository(db)
	TransactionService := service.NewTransactionService(*TransactionRepository, *ItemRepository)

	r := mux.NewRouter()

	// Root Routes
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hello, World!")); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}).Methods("GET")

	// App Routes
	routes.AdminRoutes(r, AuthService, jwtUtils)
	routes.CategoryRoutes(r, categoryService, jwtUtils)
	routes.StorageRoutes(r, StorageService, jwtUtils)
	routes.ItemRoutes(r, itemService, jwtUtils)
	routes.TransactionRoutes(r, TransactionService, jwtUtils)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "PUT"},
	})

	// Apply CORS middleware
	handler := c.Handler(r)

	// Start the server
	fmt.Println("Server is starting on :8080...")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
