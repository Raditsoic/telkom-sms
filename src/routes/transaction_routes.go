package routes

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
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

	r.HandleFunc("/api/transaction/loan", func(w http.ResponseWriter, r *http.Request) {
		var req model.LoanTransaction
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
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
		var req model.InquiryTransaction
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
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

		if transactionService == nil {
			http.Error(w, "Transaction service not initialized", http.StatusInternalServerError)
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

		file, fileHeader, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Image file is required", http.StatusBadRequest)
			return
		}
		if file == nil {
			http.Error(w, "Image file is empty", http.StatusBadRequest)
			return
		}
		defer func() {
			if file != nil {
				file.Close()
			}
		}()

		if fileHeader.Size > 10<<20 {
			http.Error(w, "File size too large", http.StatusBadRequest)
			return
		}

		imageData, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Could not read image file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		req := model.CreateInsertionTransactionDTO{
			EmployeeName:       r.FormValue("employee_name"),
			EmployeeDepartment: r.FormValue("employee_department"),
			EmployeePosition:   r.FormValue("employee_position"),
			Notes:              r.FormValue("notes"),
			Image:              imageData,
			ItemRequest: model.ItemRequestDTO{
				Name:       r.FormValue("item_name"),
				Quantity:   quantity,
				Shelf:      r.FormValue("shelf"),
				CategoryID: uint(categoryID),
			},
		}

		transaction, err := transactionService.CreateInsertionTransaction(&req)
		if err != nil {
			log.Printf("Error creating insertion transaction: %v", err)
			http.Error(w, "Failed to create transaction: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(transaction); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}).Methods("POST")

	// r.Handle("/api/transaction/{uuid}/{status}", middleware.AuthMiddleware(jwtUtils, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	uuid := mux.Vars(r)["uuid"]
	// 	status := mux.Vars(r)["status"]
	// 	transaction, err := transactionService.UpdateTransactionStatus(status, uuid)
	// 	if err != nil {
	// 		if errors.Is(err, utils.ErrTransactionType) {
	// 			http.Error(w, "Invalid transaction type", http.StatusBadRequest)
	// 			return
	// 		}
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}

	// 	w.Header().Set("Content-Type", "application/json")
	// 	if err := json.NewEncoder(w).Encode(transaction); err != nil {
	// 		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	// 		return
	// 	}
	// }))).Methods("PATCH")

	r.Handle("/api/transaction/{uuid}/{status}", middleware.AuthMiddleware(jwtUtils, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uuid := mux.Vars(r)["uuid"]
		status := mux.Vars(r)["status"]
		transaction, err := transactionService.UpdateTransactionStatus(status, uuid)
		if err != nil {
			if errors.Is(err, utils.ErrTransactionType) {
				http.Error(w, "Invalid transaction type", http.StatusBadRequest)
				return
			}
			if errors.Is(err, utils.ErrTransactionNotFound) {
				http.Error(w, "Transaction not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(transaction); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}))).Methods("PATCH")

	r.Handle("/api/transaction/{uuid}", middleware.AuthMiddleware(jwtUtils, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uuid := mux.Vars(r)["uuid"]
		response, err := transactionService.DeleteTransaction(uuid)
		if err != nil {
			if errors.Is(err, utils.ErrTransactionNotFound) {
				http.Error(w, "Transaction not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}))).Methods("DELETE")

	r.HandleFunc("/api/transactions/export", func(w http.ResponseWriter, r *http.Request) {
		startTimeParam := r.URL.Query().Get("from")
		endTimeParam := r.URL.Query().Get("to")

		var startTime, endTime time.Time
    	var err error

		if startTimeParam == "" {
        // Default to a reasonable past time, e.g., 1 year ago
			startTime = time.Now().AddDate(-1, 0, 0)
		} else {
			startTime, err = time.Parse(time.RFC3339, startTimeParam)
			if err != nil {
				http.Error(w, "Invalid start time format. Use RFC3339 (e.g., 2024-12-05T00:00:00Z)", http.StatusBadRequest)
				return
			}
		}

		if endTimeParam == "" {
			endTime = time.Now()
		} else {
			endTime, err = time.Parse(time.RFC3339, endTimeParam)
			if err != nil {
				http.Error(w, "Invalid end time format. Use RFC3339 (e.g., 2024-12-05T23:59:59Z)", http.StatusBadRequest)
				return
			}
    	}

		transactions, err := transactionService.ExportTransactions(startTime, endTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Disposition", "attachment; filename=combined_transactions.csv")
		w.Header().Set("Content-Type", "text/csv")

		writer := csv.NewWriter(w)
		defer writer.Flush()

		header := []string{
			"TransactionType", "ID", "UUID", "EmployeeName", "EmployeeDepartment", "EmployeePosition",
			"CategoryName", "ItemName", "Quantity", "Status", "Notes", "Time", "ItemID",
			"LoanTime", "ReturnTime", "CompletedTime", "ReturnedTime", "Image",
		}
		if err := writer.Write(header); err != nil {
			http.Error(w, "Failed to write CSV header", http.StatusInternalServerError)
			return
		}

		for _, t := range transactions {
			row := []string{
				t.TransactionType,
				fmt.Sprintf("%d", t.ID),
				t.UUID,
				t.EmployeeName,
				t.EmployeeDepartment,
				t.EmployeePosition,
				t.CategoryName.String,
				t.ItemName.String,
				fmt.Sprintf("%d", t.Quantity.Int32),
				t.Status,
				t.Notes.String,
				t.Time.Time.Format(time.RFC3339),
				fmt.Sprintf("%d", t.ItemID.Int32),
				t.LoanTime.Time.Format(time.RFC3339),
				t.ReturnTime.Time.Format(time.RFC3339),
				t.CompletedTime.Time.Format(time.RFC3339),
				t.ReturnedTime.Time.Format(time.RFC3339),
				t.Image.String,
			}
			if err := writer.Write(row); err != nil {
				http.Error(w, "Failed to write CSV row", http.StatusInternalServerError)
				return
			}
		}
	}).Methods("GET")
}
