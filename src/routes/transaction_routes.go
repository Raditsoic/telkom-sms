package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gtihub.com/raditsoic/telkom-storage-ms/src/middleware"
	"gtihub.com/raditsoic/telkom-storage-ms/src/model"
	"gtihub.com/raditsoic/telkom-storage-ms/src/service"
	"gtihub.com/raditsoic/telkom-storage-ms/src/utils"
)

func TransactionRoutes(r *mux.Router, transactionService *service.TransactionService, jwtUtils *utils.JWTUtils) {
	r.HandleFunc("/api/transactions", func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit < 1 {
			limit = 10
		}

		transactions, err := transactionService.GetTransactions(page, limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(transactions); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}).Methods("GET")

	r.HandleFunc("/api/transaction/loan/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		transaction, err := transactionService.GetLoanTransactionByID(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(transaction); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}).Methods("GET")

	r.HandleFunc("/api/transaction/inquiry/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		transaction, err := transactionService.GetInquiryTransactionByID(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(transaction); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}).Methods("GET")

	r.HandleFunc("/api/transaction/insert/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		transaction, err := transactionService.GetInsertionTransactionByID(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(transaction); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}).Methods("GET")

	r.HandleFunc("/api/transaction/loan", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		requiredFields := map[string]string{
			"employee_name":       r.FormValue("employee_name"),
			"employee_department": r.FormValue("employee_department"),
			"employee_position":   r.FormValue("employee_position"),
			"notes":               r.FormValue("notes"),
			"item_id":             r.FormValue("item_id"),
			"quantity":            r.FormValue("quantity"),
			"return_time":         r.FormValue("return_time"),
		}

		for field, value := range requiredFields {
			if value == "" {
				http.Error(w, fmt.Sprintf("Field '%s' is required", field), http.StatusBadRequest)
				return
			}
		}

		item_id, err := strconv.ParseUint(r.FormValue("item_id"), 10, 32)
		if err != nil {
			http.Error(w, "Invalid item ID", http.StatusBadRequest)
			return
		}

		quantity, err := strconv.Atoi(r.FormValue("quantity"))
		if err != nil {
			http.Error(w, "Invalid quantity", http.StatusBadRequest)
			return
		}

		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Image file is required", http.StatusBadRequest)
		}
		defer file.Close()

		imageData, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Could not read image file", http.StatusInternalServerError)
			return
		}

		returnTime, err := time.Parse(time.RFC3339, r.FormValue("return_time"))
		if err != nil {
			http.Error(w, "Invalid return time format. Expected RFC3339 format: "+err.Error(), http.StatusBadRequest)
			return
		}

		now := time.Now()
		if returnTime.Before(now) {
			http.Error(w, "Return time must be after loan time", http.StatusBadRequest)
			return
		}

		req := model.LoanTransaction{
			EmployeeName:       r.FormValue("employee_name"),
			EmployeeDepartment: r.FormValue("employee_department"),
			EmployeePosition:   r.FormValue("employee_position"),
			Quantity:           quantity,
			Notes:              r.FormValue("notes"),
			LoanTime:           now,
			ReturnTime:         returnTime,
			Image:              imageData,
			ItemID:             uint(item_id),
		}

		transaction, err := transactionService.CreateLoanTransaction(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(transaction); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}).Methods("POST")

	r.HandleFunc("/api/transaction/inquiry", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		employee_name := r.FormValue("employee_name")
		employee_department := r.FormValue("employee_department")
		employee_position := r.FormValue("employee_position")
		notes := r.FormValue("notes")

		if employee_name == "" || employee_department == "" || employee_position == "" || notes == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		item_id, err := strconv.ParseUint(r.FormValue("item_id"), 10, 32)
		if err != nil {
			http.Error(w, "Invalid item ID", http.StatusBadRequest)
			return
		}

		quantity, err := strconv.Atoi(r.FormValue("quantity"))
		if err != nil {
			http.Error(w, "Invalid quantity", http.StatusBadRequest)
			return
		}

		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Image file is required", http.StatusBadRequest)
		}
		defer file.Close()

		imageData, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Could not read image file", http.StatusInternalServerError)
			return
		}

		req := model.InquiryTransaction{
			EmployeeName:       employee_name,
			EmployeeDepartment: employee_department,
			EmployeePosition:   employee_position,
			Quantity:           quantity,
			Notes:              notes,
			Image:              imageData,
			ItemID:             uint(item_id),
		}

		transaction, err := transactionService.CreateInquiryTransaction(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(transaction); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}).Methods("POST")

	r.HandleFunc("/api/transaction/insert", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		requiredFields := map[string]string{
			"employee_name":       r.FormValue("employee_name"),
			"employee_department": r.FormValue("employee_department"),
			"employee_position":   r.FormValue("employee_position"),
			"notes":               r.FormValue("notes"),
			"item_name":           r.FormValue("item_name"),
			"quantity":            r.FormValue("quantity"),
			"shelf":               r.FormValue("shelf"),
			"category_id":         r.FormValue("category_id"),
		}

		for field, value := range requiredFields {
			if value == "" {
				http.Error(w, fmt.Sprintf("Field '%s' is required", field), http.StatusBadRequest)
				return
			}
		}

		categoryID, err := strconv.ParseUint(r.FormValue("category_id"), 10, 32)
		if err != nil {
			http.Error(w, "Invalid category ID: "+err.Error(), http.StatusBadRequest)
			return
		}

		quantity, err := strconv.Atoi(r.FormValue("quantity"))
		if err != nil {
			http.Error(w, "Invalid quantity: "+err.Error(), http.StatusBadRequest)
			return
		}
		if quantity <= 0 {
			http.Error(w, "Quantity must be greater than 0", http.StatusBadRequest)
			return
		}

		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Image file is required", http.StatusBadRequest)
		}
		defer file.Close()

		imageData, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Could not read image file", http.StatusInternalServerError)
			return
		}

		item := model.Item{
			Name:       r.FormValue("item_name"),
			Quantity:   quantity,
			Shelf:      r.FormValue("shelf"),
			CategoryID: uint(categoryID),
		}

		req := model.InsertionTransaction{
			EmployeeName:       r.FormValue("employee_name"),
			EmployeeDepartment: r.FormValue("employee_department"),
			EmployeePosition:   r.FormValue("employee_position"),
			Notes:              r.FormValue("notes"),
			Image:              imageData,
			Item:               item,
		}

		transaction, err := transactionService.CreateInsertionTransaction(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(transaction); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}).Methods("POST")

	r.Handle("/api/transaction/loan/{id}/{status}", middleware.AuthMiddleware(jwtUtils, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		status := vars["status"]
		if status != "Approved" && status != "Rejected" {
			http.Error(w, "Invalid status", http.StatusBadRequest)
			return
		}

		transaction, err := transactionService.UpdateLoanTransaction(uint(id), status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(transaction); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}))).Methods("PUT")

	r.Handle("/api/transaction/inquiry/{id}/{status}", middleware.AuthMiddleware(jwtUtils, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		status := vars["status"]
		if status != "Approved" && status != "Rejected" {
			http.Error(w, "Invalid status", http.StatusBadRequest)
			return
		}

		transaction, err := transactionService.UpdateInquiryTransaction(uint(id), status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(transaction); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}))).Methods("PUT")

	r.Handle("/api/transaction/insert/{id}/{status}", middleware.AuthMiddleware(jwtUtils, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		status := vars["status"]
		if status != "Approved" && status != "Rejected" {
			http.Error(w, "Invalid status", http.StatusBadRequest)
			return
		}

		transaction, err := transactionService.UpdateInsertionTransaction(uint(id), status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(transaction); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}))).Methods("PUT")

}
