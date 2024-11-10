package routes

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gtihub.com/raditsoic/telkom-storage-ms/src/middleware"
	"gtihub.com/raditsoic/telkom-storage-ms/src/model"
	"gtihub.com/raditsoic/telkom-storage-ms/src/service"
	"gtihub.com/raditsoic/telkom-storage-ms/src/utils"
)

func AdminRoutes(r *mux.Router, authService *service.AuthService, jwtUtils *utils.JWTUtils) {
	r.HandleFunc("/api/admin/register", func(w http.ResponseWriter, r *http.Request) {
		var RegisReq model.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&RegisReq); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		registeredResponse, err := authService.AdminRegister(&RegisReq)
		if err != nil {
			if errors.Is(err, utils.ErrUsernameExists) {
				http.Error(w, err.Error(), http.StatusConflict) 
				return
			}
			log.Printf("Error Registering Admin: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(registeredResponse); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}).Methods("POST")

	r.HandleFunc("/api/admin/login", func(w http.ResponseWriter, r *http.Request) {
		var LoginReq model.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&LoginReq); err != nil {
			http.Error(w, "Invalid login request", http.StatusBadRequest)
			return
		}

		loginResponse, err := authService.AdminLogin(&LoginReq)
		if err != nil {
			if errors.Is(err, utils.ErrInvalidCredentials) {
				http.Error(w, "Invalid username or password", http.StatusUnauthorized)
				return
			}
			log.Printf("Error Logging in: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(loginResponse); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}).Methods("POST")

	r.Handle("/api/admins", middleware.AuthMiddleware(jwtUtils, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getAdmins, err := authService.GetAdmin()
		if err != nil {
			log.Printf("Error getting admins: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(getAdmins); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}))).Methods("GET")

	r.Handle("/api/admin/{id}", middleware.AuthMiddleware(jwtUtils, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		deleteResponse, err := authService.DeleteAdmin(id)
		if err != nil {
			if errors.Is(err, utils.ErrInvalidID) {
				http.Error(w, "Invalid ID", http.StatusBadRequest)
				return
			}
			log.Printf("Error deleting admin: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(deleteResponse); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}))).Methods("DELETE")
}
