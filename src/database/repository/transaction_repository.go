package repository

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gtihub.com/raditsoic/telkom-storage-ms/src/model"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repository *TransactionRepository) CreateLoanTransaction(loan model.LoanTransaction) error {
	if err := repository.db.Create(&loan).Error; err != nil {
		return fmt.Errorf("failed to create loan transaction: %w", err)
	}

	return nil
}

func (repository *TransactionRepository) CreateInquiryTransaction(inquiry model.InquiryTransaction) error {
	if err := repository.db.Create(&inquiry).Error; err != nil {
		return fmt.Errorf("failed to create loan transaction: %w", err)
	}

	return nil
}

func (repository *TransactionRepository) CreateInsertionTransaction(insert *model.InsertionTransaction) error {
	if err := repository.db.Create(insert).Error; err != nil {
		return fmt.Errorf("failed to create insert transaction: %w", err)
	}

	return nil
}

func (repository *TransactionRepository) GetLoanTransactions(limit, offset int) ([]model.LoanTransaction, error) {
	var loanTransactions []model.LoanTransaction
	if err := repository.db.Preload("Item").Limit(limit).Offset(offset).Find(&loanTransactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get loan transactions: %w", err)
	}

	return loanTransactions, nil
}

func (repository *TransactionRepository) GetInquiryTransactions(limit, offset int) ([]model.InquiryTransaction, error) {
	var inquiryTransactions []model.InquiryTransaction
	if err := repository.db.Preload("Item").Limit(limit).Offset(offset).Find(&inquiryTransactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get inquiry transactions: %w", err)
	}

	return inquiryTransactions, nil
}

func (repository *TransactionRepository) GetInsertionTransactions(limit, offset int) ([]model.InsertionTransaction, error) {
	var insertTransactions []model.InsertionTransaction
	if err := repository.db.Preload("Item").Limit(limit).Offset(offset).Find(&insertTransactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get insert transactions: %w", err)
	}

	return insertTransactions, nil
}

func (repository *TransactionRepository) UpdateLoanTransaction(loan *model.LoanTransaction) error {
	if err := repository.db.Save(loan).Error; err != nil {
		return fmt.Errorf("failed to update loan transaction: %w", err)
	}
	return nil
}

func (repository *TransactionRepository) UpdateInquiryTransaction(inquiry *model.InquiryTransaction) error {
	if err := repository.db.Save(inquiry).Error; err != nil {
		return fmt.Errorf("failed to update inquiry transaction: %w", err)
	}
	return nil
}

func (repository *TransactionRepository) UpdateInsertionTransaction(insert *model.InsertionTransaction) error {
	if err := repository.db.Save(insert).Error; err != nil {
		return fmt.Errorf("failed to update insertion transaction: %w", err)
	}
	return nil
}

func (repository *TransactionRepository) GetInsertionTransactionByUUID(uuid uuid.UUID) (*model.InsertionTransaction, error) {
	var insert model.InsertionTransaction
	if err := repository.db.Preload("Item").Where("uuid = ?", uuid).First(&insert).Error; err != nil {
		return nil, fmt.Errorf("failed to get insertion transaction: %w", err)
	}

	return &insert, nil
}

func (repository *TransactionRepository) GetLoanTransactionByUUID(uuid uuid.UUID) (*model.LoanTransaction, error) {
	var loan model.LoanTransaction
	if err := repository.db.Preload("Item").Where("uuid = ?", uuid).First(&loan).Error; err != nil {
		return nil, fmt.Errorf("failed to get loan transaction: %w", err)
	}

	return &loan, nil
}

func (repository *TransactionRepository) GetInquiryTransactionByUUID(uuid uuid.UUID) (*model.InquiryTransaction, error) {
	var inquiry model.InquiryTransaction
	if err := repository.db.Preload("Item").Where("uuid = ?", uuid).First(&inquiry).Error; err != nil {
		return nil, fmt.Errorf("failed to get inquiry transaction: %w", err)
	}

	return &inquiry, nil
}
