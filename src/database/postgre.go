package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gtihub.com/raditsoic/telkom-storage-ms/src/model"
)

func Connect() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable&TimeZone=Asia/Jakarta",
		dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

	db.Exec(`
    UPDATE loan_transactions SET uuid = uuid_generate_v4() WHERE uuid IS NULL;
    UPDATE inquiry_transactions SET uuid = uuid_generate_v4() WHERE uuid IS NULL;
    UPDATE insertion_transactions SET uuid = uuid_generate_v4() WHERE uuid IS NULL;
	`)

	if err := db.AutoMigrate(
		&model.Admin{},
		&model.Storage{},
		&model.Item{},
		&model.Category{},
		&model.LoanTransaction{},
		&model.InquiryTransaction{},
		&model.InsertionTransaction{},
	); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}

	return db, nil
}
