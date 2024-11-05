package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gtihub.com/raditsoic/telkom-storage-ms/src/database"
	"gtihub.com/raditsoic/telkom-storage-ms/src/database/repository"
	"gtihub.com/raditsoic/telkom-storage-ms/src/model"
	"gtihub.com/raditsoic/telkom-storage-ms/src/service"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	CategoryRepository := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(*CategoryRepository)

	StorageRepository := repository.NewStorageRepository(db)
	StorageService := service.NewStorageService(*StorageRepository)

	ItemRepository := repository.NewItemRepository(db)
	itemService := service.NewItemService(*ItemRepository)

	TransactionRepository := repository.NewTransactionRepository(db)
	TransactionService := service.NewTransactionService(*TransactionRepository, *ItemRepository)

	r := mux.NewRouter()

	// Root route
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hello, World!")); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}).Methods("GET")

	// Admin routes
	r.HandleFunc("/api/admin/register", service.AdminRegister).Methods("POST")
	r.HandleFunc("/api/admin/login", service.AdminLogin).Methods("POST")
	r.HandleFunc("/api/admins", service.GetAdmin).Methods("GET")      // Endpoint testing Jgn lupa dihapus klo udha mau prod
	r.HandleFunc("/api/admin", service.DeleteAdmin).Methods("DELETE") // Endpoint testing Jgn lupa dihapus klo udha mau prod

	// Item routes
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

		if err := itemService.DeleteItem(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode("Item deleted"); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}).Methods("DELETE")
	r.HandleFunc("/api/item", func(w http.ResponseWriter, r *http.Request) {
		var item model.Item
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		createdItem, err := itemService.CreateItem(&item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(createdItem); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}).Methods("POST")

	// Category routes
	r.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		limit := r.URL.Query().Get("limit")

		categories, err := categoryService.GetCategories(page, limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(categories); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}).Methods("GET")
	r.HandleFunc("/api/category", func(w http.ResponseWriter, r *http.Request) {
		var category model.Category
		if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		newCategory, err := json.Marshal(category)
		if err != nil {
			http.Error(w, "Failed to marshal category", http.StatusInternalServerError)
			return
		}
		createdCategory, err := categoryService.CreateCategory(newCategory)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(createdCategory); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}).Methods("POST")
	r.HandleFunc("/api/category/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		category, err := categoryService.GetCategoryByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(category); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}).Methods("GET")
	r.HandleFunc("/api/category/{id}/items", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		categoryID, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		category, err := categoryService.GetCategoryWithItems(uint(categoryID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the response header and return the category with items
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(category); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}).Methods("GET")
	r.HandleFunc("/api/category/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
		if err := categoryService.DeleteCategory(idInt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode("Category deleted"); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}).Methods("DELETE")

	// Storage routes
	r.HandleFunc("/api/storage", func(w http.ResponseWriter, r *http.Request) {
		var storage model.Storage
		if err := json.NewDecoder(r.Body).Decode(&storage); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		storageData, err := json.Marshal(storage)
		if err != nil {
			http.Error(w, "Failed to marshal category", http.StatusInternalServerError)
			return
		}
		createdStorage, err := StorageService.CreateStorage(storageData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(createdStorage); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}).Methods("POST")
	r.HandleFunc("/api/storages", func(w http.ResponseWriter, r *http.Request) {
		storages, err := StorageService.GetStorages()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(storages); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}).Methods("GET")
	r.HandleFunc("/api/storage/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid storage ID", http.StatusBadRequest)
			return
		}

		if err := StorageService.DeleteStorage(idInt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode("Item deleted"); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}).Methods("DELETE")
	r.HandleFunc("/api/storage/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid storage ID", http.StatusBadRequest)
			return
		}

		storage, err := StorageService.GetStorageByID(idInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(storage); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}).Methods("GET")

	// Transaction routes
	r.HandleFunc("/api/transactions", func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit < 1 {
			limit = 10
		}

		transactions, err := TransactionService.GetTransactions(page, limit)
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

		transaction, err := TransactionService.GetLoanTransactionByID(uint(id))
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

		transaction, err := TransactionService.GetInquiryTransactionByID(uint(id))
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

		transaction, err := TransactionService.GetInsertionTransactionByID(uint(id))
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
		var req model.LoanTransaction

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		transaction, err := TransactionService.CreateLoanTransaction(req)
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
		var inquiry model.InquiryTransaction

		if err := json.NewDecoder(r.Body).Decode(&inquiry); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		transaction, err := TransactionService.CreateInquiryTransaction(inquiry)
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
		var insert model.InsertionTransaction

		if err := json.NewDecoder(r.Body).Decode(&insert); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		transaction, err := TransactionService.CreateInsertionTransaction(&insert)
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

	// Update transaction (Approve/Reject)
	// r.HandleFunc("/api/transaction/loan/{id}/{status}", func(w http.ResponseWriter, r *http.Request) {
	// 	vars := mux.Vars(r)
	// 	id, err := strconv.ParseUint(vars["id"], 10, 32)
	// 	if err != nil {
	// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
	// 		return
	// 	}

	// 	status := vars["status"]
	// 	if status != "approve" && status != "reject" {
	// 		http.Error(w, "Invalid status", http.StatusBadRequest)
	// 		return
	// 	}

	// 	transaction, err := TransactionService.UpdateLoanTransaction(uint(id), status)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}

	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusOK)
	// 	if err := json.NewEncoder(w).Encode(transaction); err != nil {
	// 		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	// 		return
	// 	}
	// }).Methods("PATCH")
	// r.HandleFunc("/api/transaction/inquiry/{id}/{status}", func(w http.ResponseWriter, r *http.Request) {
	// 	vars := mux.Vars(r)
	// 	id, err := strconv.ParseUint(vars["id"], 10, 32)
	// 	if err != nil {
	// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
	// 		return
	// 	}

	// 	status := vars["status"]
	// 	if status != "approve" && status != "reject" {
	// 		http.Error(w, "Invalid status", http.StatusBadRequest)
	// 		return
	// 	}

	// 	transaction, err := TransactionService.UpdateInquiryTransaction(uint(id), status)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}

	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusOK)
	// 	if err := json.NewEncoder(w).Encode(transaction); err != nil {
	// 		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	// 		return
	// 	}
	// }).Methods("PATCH")
	// r.HandleFunc("/api/transaction/insert/{id}/{status}", func(w http.ResponseWriter, r *http.Request) {
	// 	vars := mux.Vars(r)
	// 	id, err := strconv.ParseUint(vars["id"], 10, 32)
	// 	if err != nil {
	// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
	// 		return
	// 	}

	// 	status := vars["status"]
	// 	if status != "approve" && status != "reject" {
	// 		http.Error(w, "Invalid status", http.StatusBadRequest)
	// 		return
	// 	}

	// 	transaction, err := TransactionService.UpdateInsertionTransaction(uint(id), status)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}

	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusOK)
	// 	if err := json.NewEncoder(w).Encode(transaction); err != nil {
	// 		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	// 		return
	// 	}
	// }).Methods("PATCH")

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})

	// Apply CORS middleware
	handler := c.Handler(r)

	// Start the server
	fmt.Println("Server is starting on :8080...")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
