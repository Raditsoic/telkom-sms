package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Admin struct {
	ID       int
	Username string
	Password string
}

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
	Image     []byte  `json:"image"`
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

type CategoryImages struct {
	Amplop          []byte
	Baterai         []byte
	DoubleTape      []byte
	PulpenFaster    []byte
	Isolasi         []byte
	KalkulatorJoyko []byte
	Kertas          []byte
	Lakban          []byte
	Lem             []byte
	Map             []byte
	PaperClip       []byte
	Stabilo         []byte
	Staples         []byte
	BoardMarker     []byte
	Spidol          []byte
	Pensil          []byte
	Penghapus       []byte
	Plastik         []byte
	Materai         []byte
	BinderClip      []byte
	CorrectionTape  []byte
	BoardEraser     []byte
	Buku            []byte
	Cutter          []byte
	Gunting         []byte
	Penggaris       []byte
	PlongKertas     []byte
}

func LoadImages() (*CategoryImages, error) {
	// Create a helper function to reduce repetition
	loadImage := func(path string) ([]byte, error) {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read image %s: %v", path, err)
		}
		return data, nil
	}

	images := &CategoryImages{}
	var err error

	// Load all images
	if images.Amplop, err = loadImage("./seeder/img/amplop.jpg"); err != nil {
		fmt.Println("Failed to load image amplop.jpg")
		return nil, err
	}
	if images.Baterai, err = loadImage("./seeder/img/baterai.jpg"); err != nil {
		fmt.Println("Failed to load image baterai.jpg")
		return nil, err
	}
	if images.DoubleTape, err = loadImage("./seeder/img/double_tape.jpg"); err != nil {
		fmt.Println("Failed to load image double_tape.jpg")
		return nil, err
	}
	if images.PulpenFaster, err = loadImage("./seeder/img/faster.jpg"); err != nil {
		fmt.Println("Failed to load image faster.jpg")
		return nil, err
	}
	if images.Isolasi, err = loadImage("./seeder/img/isolasi.jpg"); err != nil {
		fmt.Println("Failed to load image isolasi.jpg")
		return nil, err
	}
	if images.KalkulatorJoyko, err = loadImage("./seeder/img/kalkulator_joyko.jpg"); err != nil {
		fmt.Println("Failed to load image kalkulator_joyko.jpg")
		return nil, err
	}
	if images.Kertas, err = loadImage("./seeder/img/kertas.jpg"); err != nil {
		fmt.Println("Failed to load image kertas.jpg")
		return nil, err
	}
	if images.Lakban, err = loadImage("./seeder/img/lakban.jpg"); err != nil {
		fmt.Println("Failed to load image lakban.jpg")
		return nil, err
	}
	if images.Lem, err = loadImage("./seeder/img/lem.jpg"); err != nil {
		fmt.Println("Failed to load image lem.jpg")
		return nil, err
	}
	if images.Map, err = loadImage("./seeder/img/map.jpg"); err != nil {
		fmt.Println("Failed to load image map.jpg")
		return nil, err
	}
	if images.PaperClip, err = loadImage("./seeder/img/paper_clip.jpg"); err != nil {
		fmt.Println("Failed to load image paper_clip.jpg")
		return nil, err
	}
	if images.Stabilo, err = loadImage("./seeder/img/stabilo.jpg"); err != nil {
		fmt.Println("Failed to load image stabilo.jpg")
		return nil, err
	}
	if images.Staples, err = loadImage("./seeder/img/staples.jpg"); err != nil {
		fmt.Println("Failed to load image staples.jpg")
		return nil, err
	}
	if images.BoardMarker, err = loadImage("./seeder/img/board_marker.jpg"); err != nil {
		fmt.Println("Failed to load image board_marker.jpg")
		return nil, err
	}
	if images.Spidol, err = loadImage("./seeder/img/spidol.jpeg"); err != nil {
		fmt.Println("Failed to load image spidol.jpeg")
		return nil, err
	}
	if images.Pensil, err = loadImage("./seeder/img/pensil.jpg"); err != nil {
		fmt.Println("Failed to load image pensil.jpg")
		return nil, err
	}
	if images.Penghapus, err = loadImage("./seeder/img/penghapus.jpg"); err != nil {
		fmt.Println("Failed to load image penghapus.jpg")
		return nil, err
	}
	if images.Plastik, err = loadImage("./seeder/img/plastik.jpg"); err != nil {
		fmt.Println("Failed to load image plastik.jpg")
		return nil, err
	}
	if images.Materai, err = loadImage("./seeder/img/materai.png"); err != nil {
		fmt.Println("Failed to load image materai.png")
		return nil, err
	}
	if images.BinderClip, err = loadImage("./seeder/img/binder_clip.jpeg"); err != nil {
		fmt.Println("Failed to load image binder_clip.jpeg")
		return nil, err
	}
	if images.CorrectionTape, err = loadImage("./seeder/img/correction_tape.jpg"); err != nil {
		fmt.Println("Failed to load image correction_tape.jpg")
		return nil, err
	}
	if images.BoardEraser, err = loadImage("./seeder/img/board_eraser.jpeg"); err != nil {
		fmt.Println("Failed to load image board_eraser.jpeg")
		return nil, err
	}
	if images.Buku, err = loadImage("./seeder/img/buku.jpg"); err != nil {
		fmt.Println("Failed to load image buku.jpg")
		return nil, err
	}
	if images.Cutter, err = loadImage("./seeder/img/cutter.jpeg"); err != nil {
		fmt.Println("Failed to load image cutter.jpeg")
		return nil, err
	}
	if images.Gunting, err = loadImage("./seeder/img/gunting.jpeg"); err != nil {
		fmt.Println("Failed to load image gunting.jpeg")
		return nil, err
	}
	if images.Penggaris, err = loadImage("./seeder/img/penggaris.jpg"); err != nil {
		fmt.Println("Failed to load image penggaris.jpeg")
		return nil, err
	}
	if images.PlongKertas, err = loadImage("./seeder/img/plong_kertas.jpg"); err != nil {
		fmt.Println("Failed to load image plong_kertas.jpg")
		return nil, err
	}

	return images, nil
}

func AdminSeeder(db *gorm.DB) error {
	// Define a list of admin users to seed
	admins := []Admin{
		{Username: "soic", Password: "123"},    // Replace with a secure initial password
		{Username: "admin", Password: "admin"}, // Replace with a secure initial password
	}

	// Loop through each admin and add them to the database
	for _, admin := range admins {
		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		admin.Password = string(hashedPassword)

		// Check if the admin user already exists
		var existingAdmin Admin
		if err := db.Where("username = ?", admin.Username).First(&existingAdmin).Error; err == nil {
			log.Printf("Admin user %s already exists. Skipping seeding.\n", admin.Username)
			continue
		} else if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if err := db.Create(&admin).Error; err != nil {
			return err
		}
		log.Printf("Admin user %s seeded successfully.\n", admin.Username)
	}

	log.Println("All admin users seeded successfully.")
	return nil
}

func main() {
	fmt.Println("Seeding database...")
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
	err = db.AutoMigrate(&Storage{}, &Category{}, &Item{}, &LoanTransaction{}, &InquiryTransaction{}, &Admin{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Clear existing data
	if err := db.Exec("TRUNCATE TABLE insertion_transactions, loan_transactions, inquiry_transactions, items, categories, storages, admins RESTART IDENTITY CASCADE").Error; err != nil {
		log.Fatal("Failed to truncate tables:", err)
	}

	if err := AdminSeeder(db); err != nil {
		log.Fatalf("Error seeding admins: %v", err)
	}

	fmt.Println("Creating admins...")

	// 1. Create Storages first
	storages := []Storage{
		{Name: "ATK", Location: "TSO Manyar"},
		{Name: "Lantai 8", Location: "TSO"},
		{Name: "Storage 1(8A)", Location: "TSO"},
		{Name: "Storage 2(8B)", Location: "TSO"},
		{Name: "Storage 3(8C)", Location: "TSO"},
		{Name: "Storage 4(8D)", Location: "TSO"},
		{Name: "Storage 5(8E)", Location: "TSO"},
		{Name: "Storage 6(8F)", Location: "TSO"},
		{Name: "Storage 7(8G)", Location: "TSO"},
		{Name: "Storage 8(8H)", Location: "TSO"},
	}

	fmt.Println("Creating storages...")

	if err := db.Create(&storages).Error; err != nil {
		log.Fatal("Failed to create storages:", err)
		return
	}

	fmt.Println("Loading category images...")

	images, err := LoadImages()
	if err != nil {
		return
	}

	fmt.Println("Creating categories...")

	// 2. Create Categories with valid StorageID references
	categories := []Category{
		{Name: "Amplop", StorageID: uint(storages[0].ID), Image: images.Amplop},
		{Name: "Kertas", StorageID: uint(storages[0].ID), Image: images.Kertas},
		{Name: "Baterai", StorageID: uint(storages[0].ID), Image: images.Baterai},
		{Name: "Board Marker", StorageID: uint(storages[0].ID), Image: images.BoardMarker},
		{Name: "Pulpen", StorageID: uint(storages[0].ID), Image: images.PulpenFaster},
		{Name: "Double Tape", StorageID: uint(storages[0].ID), Image: images.DoubleTape},
		{Name: "Isolasi", StorageID: uint(storages[0].ID), Image: images.Isolasi},
		{Name: "Kalkulator Joyko", StorageID: uint(storages[0].ID), Image: images.KalkulatorJoyko},
		{Name: "Lakban", StorageID: uint(storages[0].ID), Image: images.Lakban},
		{Name: "Spidol", StorageID: uint(storages[0].ID), Image: images.Spidol},
		{Name: "Lem", StorageID: uint(storages[0].ID), Image: images.Lem},
		{Name: "Map", StorageID: uint(storages[0].ID), Image: images.Map},
		{Name: "Stabilo", StorageID: uint(storages[0].ID), Image: images.Stabilo},
		{Name: "Staples", StorageID: uint(storages[0].ID), Image: images.Staples},
		{Name: "Paper Clip", StorageID: uint(storages[0].ID), Image: images.PaperClip},
		{Name: "Pensil", StorageID: uint(storages[0].ID), Image: images.Pensil},
		{Name: "Penghapus", StorageID: uint(storages[0].ID), Image: images.Penghapus},
		{Name: "Plastik", StorageID: uint(storages[0].ID), Image: images.Plastik},
		{Name: "Materai", StorageID: uint(storages[0].ID), Image: images.Materai},
		{Name: "Binder Clip", StorageID: uint(storages[0].ID), Image: images.BinderClip},
		{Name: "Correction Tape", StorageID: uint(storages[0].ID), Image: images.CorrectionTape},
		{Name: "Board Eraser", StorageID: uint(storages[0].ID), Image: images.BoardEraser},
		{Name: "Buku", StorageID: uint(storages[0].ID), Image: images.Buku},
		{Name: "Cutter", StorageID: uint(storages[0].ID), Image: images.Cutter},
		{Name: "Gunting", StorageID: uint(storages[0].ID), Image: images.Gunting},
		{Name: "Penggaris", StorageID: uint(storages[0].ID), Image: images.Penggaris},
		{Name: "Plong Kertas", StorageID: uint(storages[0].ID), Image: images.PlongKertas},
	}

	fmt.Println("Creating categories...")

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
		{Name: "Besar A3", Quantity: 100, Shelf: "1", CategoryID: categoryMap["Amplop"]},
		{Name: "Besar A4", Quantity: 100, Shelf: "1", CategoryID: categoryMap["Amplop"]},
		{Name: "Kecil Logo New", Quantity: 100, Shelf: "1", CategoryID: categoryMap["Amplop"]},
		{Name: "Kop Logo New", Quantity: 2500, Shelf: "1", CategoryID: categoryMap["Kertas"]},
		{Name: "A4", Quantity: 54, Shelf: "1", CategoryID: categoryMap["Kertas"]},
		{Name: "F4", Quantity: 25, Shelf: "1", CategoryID: categoryMap["Kertas"]},
		{Name: "AA", Quantity: 8, Shelf: "6", CategoryID: categoryMap["Baterai"]},
		{Name: "AAA", Quantity: 80, Shelf: "6", CategoryID: categoryMap["Baterai"]},
		{Name: "Kotak 9V", Quantity: 14, Shelf: "6", CategoryID: categoryMap["Baterai"]},
		{Name: "Biru", Quantity: 28, Shelf: "6", CategoryID: categoryMap["Board Marker"]},
		{Name: "Merah", Quantity: 30, Shelf: "6", CategoryID: categoryMap["Board Marker"]},
		{Name: "Hitam", Quantity: 10, Shelf: "6", CategoryID: categoryMap["Board Marker"]},
		{Name: "Biru", Quantity: 31, Shelf: "6", CategoryID: categoryMap["Pulpen"]},
		{Name: "Hitam", Quantity: 0, Shelf: "6", CategoryID: categoryMap["Pulpen"]},
		{Name: "12 mm", Quantity: 3, Shelf: "6", CategoryID: categoryMap["Double Tape"]},
		{Name: "24 mm", Quantity: 19, Shelf: "6", CategoryID: categoryMap["Double Tape"]},
		{Name: "Kecil", Quantity: 5, Shelf: "6", CategoryID: categoryMap["Isolasi"]},
		{Name: "CC-11A", Quantity: 2, Shelf: "6", CategoryID: categoryMap["Kalkulator Joyko"]},
		{Name: "CC-34A", Quantity: 0, Shelf: "6", CategoryID: categoryMap["Kalkulator Joyko"]},
		{Name: "Bening", Quantity: 16, Shelf: "6", CategoryID: categoryMap["Lakban"]},
		{Name: "Coklat", Quantity: 23, Shelf: "6", CategoryID: categoryMap["Lakban"]},
		{Name: "Hitam", Quantity: 23, Shelf: "6", CategoryID: categoryMap["Lakban"]},
		{Name: "Kecil Hitam", Quantity: 15, Shelf: "7", CategoryID: categoryMap["Spidol"]},
		{Name: "Kecil Merah", Quantity: 30, Shelf: "7", CategoryID: categoryMap["Spidol"]},
		{Name: "Kecil Ungu", Quantity: 8, Shelf: "7", CategoryID: categoryMap["Spidol"]},
		{Name: "Kecil Biru", Quantity: 36, Shelf: "Tanpa Rak", CategoryID: categoryMap["Spidol"]},
		{Name: "Kecil Hijau", Quantity: 2, Shelf: "7", CategoryID: categoryMap["Spidol"]},
		{Name: "UHU Stik", Quantity: 0, Shelf: "6", CategoryID: categoryMap["Lem"]},
		{Name: "Bening", Quantity: 35, Shelf: "6", CategoryID: categoryMap["Map"]},
		{Name: "No. 5", Quantity: 1, Shelf: "6", CategoryID: categoryMap["Paper Clip"]},
		{Name: "Biru", Quantity: 22, Shelf: "6", CategoryID: categoryMap["Stabilo"]},
		{Name: "Hijau", Quantity: 0, Shelf: "6", CategoryID: categoryMap["Stabilo"]},
		{Name: "Orange", Quantity: 20, Shelf: "6", CategoryID: categoryMap["Stabilo"]},
		{Name: "Kuning", Quantity: 20, Shelf: "6", CategoryID: categoryMap["Stabilo"]},
		{Name: "Pink", Quantity: 0, Shelf: "6", CategoryID: categoryMap["Stabilo"]},
		{Name: "Ungu", Quantity: 11, Shelf: "6", CategoryID: categoryMap["Stabilo"]},
		{Name: "HD-10", Quantity: 3, Shelf: "6", CategoryID: categoryMap["Staples"]},
		{Name: "HD-100", Quantity: 0, Shelf: "6", CategoryID: categoryMap["Staples"]},
		{Name: "HD-50", Quantity: 7, Shelf: "6", CategoryID: categoryMap["Staples"]},
		{Name: "Isi No.10-1M", Quantity: 50, Shelf: "6", CategoryID: categoryMap["Staples"]},
		{Name: "Isi No.3-1M", Quantity: 50, Shelf: "6", CategoryID: categoryMap["Staples"]},
		{Name: "No. 1", Quantity: 1, Shelf: "Tanpa Rak", CategoryID: categoryMap["Paper Clip"]},
		{Name: "Kenko Berwarna", Quantity: 0, Shelf: "Tanpa Rak", CategoryID: categoryMap["Paper Clip"]},
		{Name: "No. 3", Quantity: 0, Shelf: "Tanpa Rak", CategoryID: categoryMap["Paper Clip"]},
		{Name: "Gelpen Zebra Biru", Quantity: 20, Shelf: "Tanpa Rak", CategoryID: categoryMap["Pulpen"]},
		{Name: "Gelpen Zebra Hitam", Quantity: 0, Shelf: "Tanpa Rak", CategoryID: categoryMap["Pulpen"]},
		{Name: "10.000", Quantity: 0, Shelf: "Tanpa Rak", CategoryID: categoryMap["Materai"]},
		{Name: "Joyko CT-520", Quantity: 2, Shelf: "6", CategoryID: categoryMap["Correction Tape"]},
		{Name: "Binder Clip 260", Quantity: 5, Shelf: "7", CategoryID: categoryMap["Paper Clip"]},
		{Name: "Binder Clip No. 105", Quantity: 102, Shelf: "7", CategoryID: categoryMap["Paper Clip"]},
		{Name: "Binder Clip No. 107", Quantity: 24, Shelf: "7", CategoryID: categoryMap["Paper Clip"]},
		{Name: "Binder Clip No. 111", Quantity: 41, Shelf: "7", CategoryID: categoryMap["Paper Clip"]},
		{Name: "Binder Clip No. 200", Quantity: 14, Shelf: "7", CategoryID: categoryMap["Paper Clip"]},
		{Name: "Board Eraser", Quantity: 19, Shelf: "7", CategoryID: categoryMap["Board Eraser"]},
		{Name: "Jurnal Besar", Quantity: 2, Shelf: "7", CategoryID: categoryMap["Buku"]},
		{Name: "Nota", Quantity: 13, Shelf: "7", CategoryID: categoryMap["Buku"]},
		{Name: "Polio Ekspedisi", Quantity: 17, Shelf: "7", CategoryID: categoryMap["Buku"]},
		{Name: "Kwitansi", Quantity: 0, Shelf: "7", CategoryID: categoryMap["Buku"]},
		{Name: "Joyko A-300A Kecil", Quantity: 15, Shelf: "7", CategoryID: categoryMap["Cutter"]},
		{Name: "Joyko L-500 Besar", Quantity: 6, Shelf: "7", CategoryID: categoryMap["Cutter"]},
		{Name: "Isi Cutter A-100", Quantity: 15, Shelf: "7", CategoryID: categoryMap["Cutter"]},
		{Name: "Isi Cutter L-150", Quantity: 5, Shelf: "7", CategoryID: categoryMap["Cutter"]},
		{Name: "Joyko SC-828", Quantity: 8, Shelf: "7", CategoryID: categoryMap["Gunting"]},
		{Name: "Joyko SC-848", Quantity: 2, Shelf: "7", CategoryID: categoryMap["Gunting"]},
		{Name: "Plastik 30cm", Quantity: 24, Shelf: "7", CategoryID: categoryMap["Penggaris"]},
		{Name: "No. 30", Quantity: 24, Shelf: "7", CategoryID: categoryMap["Plong Kertas"]},
		{Name: "No. 40XL", Quantity: 8, Shelf: "7", CategoryID: categoryMap["Plong Kertas"]},
		{Name: "No. 85B", Quantity: 4, Shelf: "7", CategoryID: categoryMap["Plong Kertas"]},
	}

	fmt.Println("Creating items...")

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

	fmt.Println("Creating loan transactions...")

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

	fmt.Println("Creating inquiry transactions...")

	if err := db.Create(&inquiryTransactions).Error; err != nil {
		log.Fatal("Failed to create inquiry transactions:", err)
		return
	}

	fmt.Println("Seeding completed successfully!")
}
