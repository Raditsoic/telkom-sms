// seeder/main.go
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	ID         int        `json:"id" gorm:"primaryKey"`
	Name       string     `json:"name"`
	Location   string     `json:"location"`
	Categories []Category `gorm:"foreignKey:StorageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

type Category struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	Name      string  `json:"name"`
	StorageID uint    `json:"storage_id"`
	Items     []Item  `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items"`
	Storage   Storage `gorm:"foreignKey:StorageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"storage"`
}

type Item struct {
	ID         uint     `gorm:"primaryKey" json:"id"`
	Name       string   `json:"name"`
	Quantity   int      `json:"quantity"`
	CategoryID uint     `json:"category_id"`
	Category   Category `json:"-" gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type LoanTransaction struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	EmployeeName       string    `json:"employee_name"`
	EmployeeDepartment string    `json:"employee_department"`
	EmployeePosition   string    `json:"employee_position"`
	Quantity           int       `json:"quantity"`
	Status             string    `json:"status"`
	Time               time.Time `json:"time"`
	ItemID             uint      `json:"item_id"`
	Item               *Item     `gorm:"foreignKey:ItemID" json:"item"`
	LoanTime           time.Time `json:"loan_time"`
	ReturnTime         time.Time `json:"return_time"`
}

type InquiryTransaction struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	EmployeeName       string    `json:"employee_name"`
	EmployeeDepartment string    `json:"employee_department"`
	EmployeePosition   string    `json:"employee_position"`
	Quantity           int       `json:"quantity"`
	Status             string    `json:"status"`
	Time               time.Time `json:"time"`
	ItemID             uint      `json:"item_id"`
	Item               *Item     `gorm:"foreignKey:ItemID" json:"item"`
}

func main() {
	// Get database configuration from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	if dbPort == "" {
		dbPort = "5432" // default postgres port
	}

	// Database connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto Migrate
	err = db.AutoMigrate(&Storage{}, &Category{}, &Item{}, &LoanTransaction{}, &InquiryTransaction{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Clean existing data
	db.Exec("TRUNCATE TABLE loan_transactions, inquiry_transactions, items, categories, storages CASCADE")

	// Seed Data
	// Create Storages
	storages := []Storage{
		{Name: "Main Warehouse", Location: "Building A"},
		{Name: "Secondary Storage", Location: "Building B"},
	}

	for _, storage := range storages {
		result := db.Create(&storage)
		if result.Error != nil {
			log.Printf("Error creating storage %s: %v", storage.Name, result.Error)
		}
	}

	// Create Categories
	categories := []Category{
		{Name: "Electronics", StorageID: 1},
		{Name: "Office Supplies", StorageID: 1},
		{Name: "Furniture", StorageID: 2},
	}

	for _, category := range categories {
		result := db.Create(&category)
		if result.Error != nil {
			log.Printf("Error creating category %s: %v", category.Name, result.Error)
		}
	}

	// Create Items
	items := []Item{
		{Name: "Laptop", Quantity: 50, CategoryID: 1},
		{Name: "Monitor", Quantity: 30, CategoryID: 1},
		{Name: "Pen", Quantity: 1000, CategoryID: 2},
		{Name: "Paper", Quantity: 500, CategoryID: 2},
		{Name: "Chair", Quantity: 100, CategoryID: 3},
		{Name: "Desk", Quantity: 50, CategoryID: 3},
	}

	for _, item := range items {
		result := db.Create(&item)
		if result.Error != nil {
			log.Printf("Error creating item %s: %v", item.Name, result.Error)
		}
	}

	// Create Sample Transactions
	now := time.Now()
	loanTransactions := []LoanTransaction{
		{
			EmployeeName:       "John Doe",
			EmployeeDepartment: "IT",
			EmployeePosition:   "Developer",
			Quantity:          1,
			Status:            "Borrowed",
			Time:             now,
			ItemID:           1,
			LoanTime:         now,
			ReturnTime:       now.Add(7 * 24 * time.Hour),
		},
		{
			EmployeeName:       "Jane Smith",
			EmployeeDepartment: "HR",
			EmployeePosition:   "Manager",
			Quantity:          2,
			Status:            "Returned",
			Time:             now.Add(-24 * time.Hour),
			ItemID:           2,
			LoanTime:         now.Add(-24 * time.Hour),
			ReturnTime:       now,
		},
	}

	for _, transaction := range loanTransactions {
		result := db.Create(&transaction)
		if result.Error != nil {
			log.Printf("Error creating loan transaction for %s: %v", transaction.EmployeeName, result.Error)
		}
	}

	inquiryTransactions := []InquiryTransaction{
		{
			EmployeeName:       "Alice Johnson",
			EmployeeDepartment: "Marketing",
			EmployeePosition:   "Coordinator",
			Quantity:          5,
			Status:            "Pending",
			Time:             now,
			ItemID:           3,
		},
		{
			EmployeeName:       "Bob Wilson",
			EmployeeDepartment: "Sales",
			EmployeePosition:   "Representative",
			Quantity:          10,
			Status:            "Approved",
			Time:             now.Add(-12 * time.Hour),
			ItemID:           4,
		},
	}

	for _, transaction := range inquiryTransactions {
		result := db.Create(&transaction)
		if result.Error != nil {
			log.Printf("Error creating inquiry transaction for %s: %v", transaction.EmployeeName, result.Error)
		}
	}

	fmt.Println("Seeding completed successfully!")
}