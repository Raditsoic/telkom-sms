package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gtihub.com/raditsoic/telkom-storage-ms/src/model"
	"gtihub.com/raditsoic/telkom-storage-ms/src/service"
	"gtihub.com/raditsoic/telkom-storage-ms/src/utils"
)

func ItemRoutes(r *mux.Router, itemService *service.ItemService, jwtUtils *utils.JWTUtils) {
	r.HandleFunc("/api/items", func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		limit := r.URL.Query().Get("limit")

		items, err := itemService.GetItems(page, limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(items); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}).Methods("GET")

	r.HandleFunc("/api/item/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		item, err := itemService.GetItemByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(item); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}).Methods("GET")

	r.HandleFunc("/api/item/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		response, err := itemService.DeleteItem(id)
		if err != nil {
			if err == utils.ErrItemNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}).Methods("DELETE")

	r.HandleFunc("/api/item", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var item model.Item
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		createdItem, err := itemService.CreateItem(&item)
		if err != nil {
			log.Printf("Error creating item: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(createdItem); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})).Methods("POST")

	r.HandleFunc("/api/item/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var req model.UpdateItemNameRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		update, err := itemService.UpdateItemName(id, req.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := struct {
			Message string `json:"message"`
			ID      string `json:"id"`
			OldName string `json:"old_name"`
			NewName string `json:"new_name"`
		}{
			Message: "Item updated successfully",
			ID:      id,
			OldName: update.OldName,
			NewName: update.NewName,
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})).Methods("PATCH")
}
