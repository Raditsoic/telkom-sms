package main

import (
	"fmt"
	"net/http"
)

func main() {
	app := http.NewServeMux()
	app.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hello, World!")); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
	
	if err := http.ListenAndServe(":8080", app); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
