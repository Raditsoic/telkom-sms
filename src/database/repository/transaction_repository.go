package repository

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gtihub.com/raditsoic/telkom-storage-ms/src/model"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repository *TransactionRepository) GetLoanTransactionByID(id uint) (*model.Transaction, error) {
	var transaction model.Transaction

	var loan model.LoanTransaction
	if err := repository.db.Preload("Item").Where("id = ?", id).First(&loan).Error; err == nil {
		transaction = model.Transaction{
			ID:                 loan.ID,
			TransactionType:    "loan",
			EmployeeName:       loan.EmployeeName,
			EmployeeDepartment: loan.EmployeeDepartment,
			EmployeePosition:   loan.EmployeePosition,
			Quantity:           loan.Quantity,
			Status:             loan.Status,
			Time:               loan.Time,
			ItemID:             loan.ItemID,
			Item:               loan.Item,
			LoanTime:           &loan.LoanTime,
			ReturnTime:         &loan.ReturnTime,
		}

		return &transaction, nil
	}

	return nil, gorm.ErrRecordNotFound
}

func (repository *TransactionRepository) GetInquiryTransactionByID(id uint) (*model.Transaction, error) {
	var transaction model.Transaction

	var inquiry model.InquiryTransaction
	if err := repository.db.Preload("Item").Where("id = ?", id).First(&inquiry).Error; err == nil {
		transaction = model.Transaction{
			ID:                 inquiry.ID,
			TransactionType:    "inquiry",
			EmployeeName:       inquiry.EmployeeName,
			EmployeeDepartment: inquiry.EmployeeDepartment,
			EmployeePosition:   inquiry.EmployeePosition,
			Quantity:           inquiry.Quantity,
			Status:             inquiry.Status,
			Time:               inquiry.Time,
			ItemID:             inquiry.ItemID,
			Item:               inquiry.Item,
		}

		return &transaction, nil
	}

	return &transaction, gorm.ErrRecordNotFound
}

func (repository *TransactionRepository) CreateLoanTransaction(loan model.LoanTransaction) error {
	loan.LoanTime = time.Now()
	loan.Time = time.Now()
	loan.Status = "Pending"

	if err := repository.db.Create(&loan).Error; err != nil {
		return fmt.Errorf("failed to create loan transaction: %w", err)
	}

	return nil
}

func (repository *TransactionRepository) CreateInquiryTransaction(inquiry model.InquiryTransaction) error {
	inquiry.Time = time.Now()
	inquiry.Status = "Pending"

	if err := repository.db.Create(&inquiry).Error; err != nil {
		return fmt.Errorf("failed to create loan transaction: %w", err)
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

// func (repository *TransactionRepository) GetTransactionByID(id int) (*model.Transaction, error) {
// 	var transaction model.Transaction

// 	var loan model.LoanTransaction
// 	if err := repository.db.Preload("Item").Where("id = ?", id).First(&loan).Error; err == nil {
// 		transaction = model.Transaction{
// 			ID:                 loan.ID,
// 			TransactionType:    "loan",
// 			EmployeeName:       loan.EmployeeName,
// 			EmployeeDepartment: loan.EmployeeDepartment,
// 			EmployeePosition:   loan.EmployeePosition,
// 			Quantity:           loan.Quantity,
// 			Status:             loan.Status,
// 			Time:               loan.Time,
// 			ItemID:             loan.ItemID,
// 			Item:               loan.Item,
// 			LoanTime:           &loan.LoanTime,
// 			ReturnTime:         &loan.ReturnTime,
// 		}
// 		return &transaction, nil
// 	}

// 	var inquiry model.InquiryTransaction
// 	if err := repository.db.Preload("Item").Where("id = ?", id).First(&inquiry).Error; err == nil {
// 		transaction = model.Transaction{
// 			ID:                 inquiry.ID,
// 			TransactionType:    "inquiry",
// 			EmployeeName:       inquiry.EmployeeName,
// 			EmployeeDepartment: inquiry.EmployeeDepartment,
// 			EmployeePosition:   inquiry.EmployeePosition,
// 			Quantity:           inquiry.Quantity,
// 			Status:             inquiry.Status,
// 			Time:               inquiry.Time,
// 			ItemID:             inquiry.ItemID,
// 			Item:               inquiry.Item,
// 		}
// 		return &transaction, nil
// 	}

// 	return nil, gorm.ErrRecordNotFound
// }
