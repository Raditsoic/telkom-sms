package repository

import (
	"fmt"

	"gorm.io/gorm"
	"gtihub.com/raditsoic/telkom-storage-ms/model"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repository *TransactionRepository) CreateTransaction(transaction *model.Transaction) error {
	if err := repository.db.Create(&transaction).Error; err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

func (repository *TransactionRepository) GetTransactionByID(id uint) (*model.Transaction, error) {
	var transaction model.Transaction
	if err := repository.db.Where("id = ?", id).First(&transaction).Error; err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &transaction, nil
}

func (repository *TransactionRepository) GetTransactions(limit, offset int) ([]model.Transaction, error) {
	var transactions []model.Transaction

	if err := repository.db.Limit(limit).Offset(offset).Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}

	return transactions, nil
}

func (repository *TransactionRepository) UpdateTransaction(transaction model.Transaction) error {
	if err := repository.db.Save(&transaction).Error; err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	return nil
}

func (repository *TransactionRepository) DeleteTransaction(id int) error {
	if err := repository.db.Where("id = ?", id).Delete(&model.Transaction{}).Error; err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	return nil
}
