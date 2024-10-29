package service

import (
	"fmt"
	"strconv"
	"time"

	"gtihub.com/raditsoic/telkom-storage-ms/database/repository"
	"gtihub.com/raditsoic/telkom-storage-ms/model"
)

type TransactionService struct {
	logRepository  repository.TransactionRepository
	itemRepository repository.ItemRepository
}

func NewTransactionService(log repository.TransactionRepository, item repository.ItemRepository) *TransactionService {
	return &TransactionService{logRepository: log, itemRepository: item}
}

func (s *TransactionService) CreateTransaction(req *model.TransactionRequest) (*model.Transaction, error) {
	// Get the item
	item, err := s.itemRepository.GetItemByID(fmt.Sprintf("%d", req.ItemID))
	if err != nil {
		return nil, fmt.Errorf("item not found: %w", err)
	}

	// Calculate new quantity based on transaction type
	var newQuantity int
	switch req.Type {
	case "ADD":
		newQuantity = item.Quantity + req.Quantity
	case "SUBTRACT":
		newQuantity = item.Quantity - req.Quantity
		if newQuantity < 0 {
			return nil, fmt.Errorf("insufficient quantity")
		}
	default:
		return nil, fmt.Errorf("invalid transaction type")
	}

	// Update item quantity
	item.Quantity = newQuantity
	if err := s.itemRepository.UpdateItem(*item); err != nil {
		return nil, fmt.Errorf("failed to update item quantity: %w", err)
	}

	// Create transaction record
	transaction := &model.Transaction{
		EmployeeName:       req.EmployeeName,
		EmployeeDepartment: req.EmployeeDepartment,
		EmployeePosition:   req.EmployeePosition,
		Quantity:           req.Quantity,
		Status:             "Pending",
		Time:               time.Now(),
		Type:               req.Type,
		ItemID:             req.ItemID,
	}

	if err := s.logRepository.CreateTransaction(transaction); err != nil {
		// If transaction log fails, try to revert the item quantity
		item.Quantity = item.Quantity - req.Quantity
		_ = s.itemRepository.UpdateItem(*item) // Best effort to revert
		return nil, fmt.Errorf("failed to create transaction log: %w", err)
	}

	return transaction, nil
}

func (s *TransactionService) GetTransactions(pageParam, limitParam string) ([]model.Transaction, error) {
	page, limit := 1, 10

	if parsedPage, err := strconv.Atoi(pageParam); err == nil && parsedPage > 0 {
		page = parsedPage
	}
	if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
		limit = parsedLimit
	}

	offset := (page - 1) * limit
	return s.logRepository.GetTransactions(limit, offset)
}

func (s *TransactionService) GetTransactionByID(id uint) (*model.Transaction, error) {
	return s.logRepository.GetTransactionByID(id)
}
