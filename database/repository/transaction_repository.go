package repository

import (
	"fmt"

	"gtihub.com/raditsoic/telkom-storage-ms/database"
	"gtihub.com/raditsoic/telkom-storage-ms/model"
)

func CreateTransaction(transaction model.Transaction) error {
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Create(&transaction).Error; err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

func GetTransactionByID(id int) (*model.Transaction, error) {
	db, err := database.Connect()
	if err != nil {
		return &model.Transaction{}, fmt.Errorf("failed to connect to database: %w", err)
	}

	var transaction model.Transaction
	if err := db.Where("id = ?", id).First(&transaction).Error; err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &transaction, nil
}

func GetTransactions() ([]model.Transaction, error) {
	db, err := database.Connect()
	if err != nil {
		return []model.Transaction{}, fmt.Errorf("failed to connect to database: %w", err)
	}

	var transactions []model.Transaction
	if err := db.Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, nil
}

func UpdateTransaction(transaction model.Transaction) error {
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Save(&transaction).Error; err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	return nil
}

func DeleteTransaction(id int) error {
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Where("id = ?", id).Delete(&model.Transaction{}).Error; err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	return nil
}
