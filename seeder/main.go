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
	Shelf      string   `json:"shelf"`
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
	// Database connection setup
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	if dbPort == "" {
		dbPort = "5432"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&Storage{}, &Category{}, &Item{}, &LoanTransaction{}, &InquiryTransaction{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Clear existing data
	if err := db.Exec("TRUNCATE TABLE insertion_transactions, loan_transactions, inquiry_transactions, items, categories, storages RESTART IDENTITY CASCADE").Error; err != nil {
		log.Fatal("Failed to truncate tables:", err)
	}

	// 1. Create Storages first
	storages := []Storage{
		{Name: "Gudang ATK", Location: "TSO Manyar"},
		{Name: "Gudang A", Location: "TSO Manyar"},
	}

	if err := db.Create(&storages).Error; err != nil {
		log.Fatal("Failed to create storages:", err)
		return
	}

	// 2. Create Categories with valid StorageID references
	categories := []Category{
		{Name: "Amplop", StorageID: uint(storages[0].ID)},
		{Name: "Kertas", StorageID: uint(storages[0].ID)},
		{Name: "Baterai", StorageID: uint(storages[0].ID)},
		{Name: "Pulpen Faster", StorageID: uint(storages[0].ID)},
		{Name: "Double Tape", StorageID: uint(storages[0].ID)},
		{Name: "Isolasi", StorageID: uint(storages[0].ID)},
		{Name: "Kalkulator Joyko", StorageID: uint(storages[0].ID)},
		{Name: "Lakban", StorageID: uint(storages[0].ID)},
		{Name: "Lem", StorageID: uint(storages[0].ID)},
		{Name: "Map", StorageID: uint(storages[0].ID)},
		{Name: "Stabilo", StorageID: uint(storages[0].ID)},
		{Name: "Staples", StorageID: uint(storages[0].ID)},
		{Name: "Paper Clip", StorageID: uint(storages[0].ID)},
	}

	if err := db.Create(&categories).Error; err != nil {
		log.Fatal("Failed to create categories:", err)
		return
	}

	// Create a map to store category IDs by name for easy reference
	categoryMap := make(map[string]uint)
	for _, category := range categories {
		categoryMap[category.Name] = category.ID
	}

	// 3. Create Items with valid CategoryID references
	items := []Item{
		{Name: "Besar A3", Quantity: 50, Shelf: "1", CategoryID: categoryMap["Amplop"]},
		{Name: "Besar A4", Quantity: 30, Shelf: "1", CategoryID: categoryMap["Amplop"]},
		{Name: "Kecil", Quantity: 1000, Shelf: "1", CategoryID: categoryMap["Amplop"]},
		{Name: "A4", Quantity: 500, Shelf: "1", CategoryID: categoryMap["Kertas"]},
		{Name: "F4", Quantity: 100, Shelf: "1", CategoryID: categoryMap["Kertas"]},
		{Name: "AA", Quantity: 50, Shelf: "6", CategoryID: categoryMap["Baterai"]},
		{Name: "AAA", Quantity: 50, Shelf: "6", CategoryID: categoryMap["Baterai"]},
		{Name: "Biru", Quantity: 50, Shelf: "6", CategoryID: categoryMap["Pulpen Faster"]},
		{Name: "Hitam", Quantity: 50, Shelf: "6", CategoryID: categoryMap["Pulpen Faster"]},
		{Name: "12 mm", Quantity: 50, Shelf: "6", CategoryID: categoryMap["Double Tape"]},
		{Name: "24 mm", Quantity: 50, Shelf: "6", CategoryID: categoryMap["Double Tape"]},
		{Name: "Kecil", Quantity: 50, Shelf: "6", CategoryID: categoryMap["Isolasi"]},
		{Name: "CC-11A", Quantity: 50, Shelf: "6", CategoryID: categoryMap["Kalkulator Joyko"]},
		{Name: "CC-34A", Quantity: 50, Shelf: "6", CategoryID: categoryMap["Kalkulator Joyko"]},
		{Name: "Bening", Quantity: 50, Shelf: "6", CategoryID: categoryMap["Lakban"]},
		{Name: "Coklat", Quantity: 50, Shelf: "6", CategoryID: categoryMap["Lakban"]},
		{Name: "Hitam", Quantity: 50, Shelf: "6", CategoryID: categoryMap["Lakban"]},
		{Name: "UHU Stik", Quantity: 50, Shelf: "6", CategoryID: categoryMap["Lem"]},
		{Name: "Bening", Quantity: 50, Shelf: "6", CategoryID: categoryMap["Map"]},
		{Name: "Biru", Quantity: 50, Shelf: "6", CategoryID: categoryMap["Stabilo"]},
		{Name: "Hijau", Quantity: 50, Shelf: "6", CategoryID: categoryMap["Stabilo"]},
		{Name: "Orange", Quantity: 50, Shelf: "6", CategoryID: categoryMap["Stabilo"]},
		{Name: "Kuning", Quantity: 50, Shelf: "6", CategoryID: categoryMap["Stabilo"]},
		{Name: "Pink", Quantity: 50, Shelf: "6", CategoryID: categoryMap["Stabilo"]},
		{Name: "HD-10", Quantity: 50, Shelf: "6", CategoryID: categoryMap["Staples"]},
		{Name: "HD-100", Quantity: 50, Shelf: "6", CategoryID: categoryMap["Staples"]},
		{Name: "HD-50", Quantity: 50, Shelf: "6", CategoryID: categoryMap["Staples"]},
		{Name: "Isi No.10-1M", Quantity: 50, Shelf: "6", CategoryID: categoryMap["Staples"]},
		{Name: "Isi No.3-1M", Quantity: 50, Shelf: "6", CategoryID: categoryMap["Staples"]},
		{Name: "No. 1", Quantity: 50, Shelf: "Tanpa Rak", CategoryID: categoryMap["Staples"]},
	}

	if err := db.Create(&items).Error; err != nil {
		log.Fatal("Failed to create items:", err)
		return
	}

	// Create a map to store item IDs by name for easy reference
	itemMap := make(map[string]uint)
	for _, item := range items {
		itemMap[item.Name] = item.ID
	}

	// 4. Create transactions with valid ItemID references
	now := time.Now()
	loanTransactions := []LoanTransaction{
		{
			EmployeeName:       "Agus Setiawan",
			EmployeeDepartment: "IT",
			EmployeePosition:   "Developer",
			Quantity:           1,
			Status:             "Approved/Borrowed",
			Time:               now,
			ItemID:             itemMap["Besar A3"],
			LoanTime:           now,
			ReturnTime:         now.Add(7 * 24 * time.Hour),
		},
		{
			EmployeeName:       "Asep Sutisna",
			EmployeeDepartment: "HR",
			EmployeePosition:   "Manager",
			Quantity:           2,
			Status:             "Returned",
			Time:               now.Add(-24 * time.Hour),
			ItemID:             itemMap["Besar A4"],
			LoanTime:           now.Add(-24 * time.Hour),
			ReturnTime:         now,
		},
	}

	if err := db.Create(&loanTransactions).Error; err != nil {
		log.Fatal("Failed to create loan transactions:", err)
		return
	}

	inquiryTransactions := []InquiryTransaction{
		{
			EmployeeName:       "Yoga Hartono",
			EmployeeDepartment: "GA",
			EmployeePosition:   "Staff Magang",
			Quantity:           5,
			Status:             "Pending",
			Time:               now,
			ItemID:             itemMap["Kecil"],
		},
		{
			EmployeeName:       "Yoga Hartono",
			EmployeeDepartment: "GA",
			EmployeePosition:   "Staff Magang",
			Quantity:           10,
			Status:             "Approved",
			Time:               now.Add(-12 * time.Hour),
			ItemID:             itemMap["A4"],
		},
	}

	if err := db.Create(&inquiryTransactions).Error; err != nil {
		log.Fatal("Failed to create inquiry transactions:", err)
		return
	}

	fmt.Println("Seeding completed successfully!")
}
