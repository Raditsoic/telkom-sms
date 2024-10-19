package main

import (
	"fmt"
	"log"
	"net/http"

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

	if err := db.AutoMigrate(&model.Admin{}, &model.Storage{}, &model.Item{}); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}

	app := http.NewServeMux()
	app.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hello, World!")); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

	app.HandleFunc("POST /api/admin", service.AdminLogin)

	app.HandleFunc("GET /api/items", service.GetItems)

	app.HandleFunc("POST /api/storage", service.CreateStorage)
	app.HandleFunc("GET /api/storages", service.GetStorages)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow your Next.js app
		AllowCredentials: true,
	})

	handler := c.Handler(app)

	if err := http.ListenAndServe(":8080", handler); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
