package service

import (
	"encoding/json"
	"net/http"

	"gtihub.com/raditsoic/telkom-storage-ms/database/repository"
)

func GetItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	items, err := repository.GetItems()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
