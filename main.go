package main

import (
	"fmt"
	"log"
	"net/http"

	"gtihub.com/raditsoic/telkom-storage-ms/database"
	"gtihub.com/raditsoic/telkom-storage-ms/model"
	"gtihub.com/raditsoic/telkom-storage-ms/service"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&model.Admin{}); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}

	app := http.NewServeMux()
	app.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hello, World!")); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/api/admin", service.AdminLogin)

	if err := http.ListenAndServe(":8080", app); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
