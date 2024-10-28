package service

import (
	"encoding/json"
	"net/http"

	"gtihub.com/raditsoic/telkom-storage-ms/database/repository"
	"gtihub.com/raditsoic/telkom-storage-ms/model"
	"gtihub.com/raditsoic/telkom-storage-ms/utils"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var LoginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&LoginReq); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	admin, err := repository.GetAdminByUsername(LoginReq.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if admin.Password != LoginReq.Password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(admin.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func AdminRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var RegisReq RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&RegisReq); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	admin := &model.Admin{
		Username: RegisReq.Username,
		Password: RegisReq.Password,
	}
	err := repository.RegisterAdmin(admin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Admin created"}`))
}

func GetAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	admins, err := repository.GetAdmins()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(admins)
}

func DeleteAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	adminID := r.URL.Query().Get("id")
	if adminID == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err := repository.DeleteAdmin(adminID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Admin deleted"}`))
}
