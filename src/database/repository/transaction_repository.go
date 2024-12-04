package repository

import (
	"fmt"
	"time"

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

func (repository *TransactionRepository) CreateLoanTransaction(loan model.LoanTransaction) (*model.LoanTransaction, error) {
	if err := repository.db.Create(&loan).Error; err != nil {
		return nil, fmt.Errorf("failed to create loan transaction: %w", err)
	}

	return &loan, nil
}

func (repository *TransactionRepository) CreateInquiryTransaction(inquiry model.InquiryTransaction) (*model.InquiryTransaction, error) {
	if err := repository.db.Create(&inquiry).Error; err != nil {
		return nil, fmt.Errorf("failed to create loan transaction: %w", err)
	}

	return &inquiry, nil
}

func (repository *TransactionRepository) CreateInsertionTransaction(insert *model.InsertionTransaction) (*model.InsertionTransaction, error) {
	if err := repository.db.Create(insert).Error; err != nil {
		return nil, fmt.Errorf("failed to create insert transaction: %w", err)
	}

	return insert, nil
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

func (repository *TransactionRepository) DeleteLoanTransactionByUUID(uuid uuid.UUID) error {
	if err := repository.db.Where("uuid = ?", uuid).Delete(&model.LoanTransaction{}).Error; err != nil {
		return fmt.Errorf("failed to delete loan transaction: %w", err)
	}

	return nil
}

func (repository *TransactionRepository) DeleteInquiryTransactionByUUID(uuid uuid.UUID) error {
	if err := repository.db.Where("uuid = ?", uuid).Delete(&model.InquiryTransaction{}).Error; err != nil {
		return fmt.Errorf("failed to delete inquiry transaction: %w", err)
	}

	return nil
}

func (repository *TransactionRepository) DeleteInsertionTransactionByUUID(uuid uuid.UUID) error {
	if err := repository.db.Where("uuid = ?", uuid).Delete(&model.InsertionTransaction{}).Error; err != nil {
		return fmt.Errorf("failed to delete insertion transaction: %w", err)
	}

	return nil
}

func (repository *TransactionRepository) ExportTransactions(from, to time.Time) ([]model.ExportTransaction, error) {
	query := `
		SELECT 
			'LoanTransaction' AS transaction_type,
			lt.id,
			lt.uuid,
			lt.employee_name,
			lt.employee_department,
			lt.employee_position,
			c.name AS category_name,
			i.name AS item_name,
			lt.quantity,
			lt.status,
			lt.notes,
			lt.time,
			lt.item_id,
			lt.loan_time,
			lt.return_time,
			lt.completed_time,
			lt.returned_time,
			NULL::TEXT AS image
		FROM loan_transactions lt
		LEFT JOIN items i ON lt.item_id = i.id
		LEFT JOIN categories c ON i.category_id = c.id
		WHERE lt.time BETWEEN ? AND ?

		UNION ALL

		SELECT 
			'InquiryTransaction' AS transaction_type,
			it.id,
			it.uuid,
			it.employee_name,
			it.employee_department,
			it.employee_position,
			c.name AS category_name,
			i.name AS item_name,
			it.quantity,
			it.status,
			it.notes,
			it.time,
			it.item_id,
			NULL,
			NULL,
			it.completed_time,
			NULL,
			NULL::TEXT AS image
		FROM inquiry_transactions it
		LEFT JOIN items i ON it.item_id = i.id
		LEFT JOIN categories c ON i.category_id = c.id
		WHERE it.time BETWEEN ? AND ?

		UNION ALL

		SELECT 
			'InsertionTransaction' AS transaction_type,
			int.id,
			int.uuid,
			int.employee_name,
			int.employee_department,
			int.employee_position,
			c.name AS category_name,
			i.name AS item_name,
			int.item_request_quantity AS quantity,
			int.status,
			int.notes,
			int.time,
			int.item_id,
			NULL,
			NULL,
			int.completed_time,
			NULL,
			ENCODE(int.image, 'base64') AS image
		FROM insertion_transactions int
		LEFT JOIN items i ON int.item_id = i.id
		LEFT JOIN categories c ON i.category_id = c.id
		WHERE int.time BETWEEN ? AND ?
	`

	var results []model.ExportTransaction
	if err := repository.db.Raw(query, from, to, from, to, from, to).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to execute combined transactions query: %w", err)
	}

	return results, nil
}
