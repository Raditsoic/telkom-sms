package service

import (
	"fmt"
	"strconv"
	"strings"

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

func (s *TransactionService) CreateTransaction(req interface{}, transactionType string) (interface{}, error) {
	var item *model.Item
	var err error

	switch transactionType {
	case "inquiry":
		inquiryReq, ok := req.(*model.InquiryTransaction)
		if !ok {
			return nil, fmt.Errorf("invalid inquiry transaction request")
		}
		item, err = s.itemRepository.GetItemByID(fmt.Sprintf("%d", inquiryReq.ItemID))
		if err != nil {
			return nil, fmt.Errorf("item not found: %w", err)
		}

		newQuantity := item.Quantity - inquiryReq.Quantity
		if newQuantity < 0 {
			return &model.InquiryTransaction{}, fmt.Errorf("insufficient quantity")
		}

		item.Quantity = newQuantity
		if err := s.itemRepository.UpdateItem(*item); err != nil {
			return nil, fmt.Errorf("failed to update item quantity: %w", err)
		}

		if err := s.logRepository.CreateInquiryTransaction(*inquiryReq); err != nil {
			item.Quantity += inquiryReq.Quantity
			_ = s.itemRepository.UpdateItem(*item)
			return nil, fmt.Errorf("failed to create transaction log: %w", err)
		}
		return inquiryReq, nil

	case "loan":
		loanReq, ok := req.(*model.LoanTransaction)
		if !ok {
			return nil, fmt.Errorf("invalid loan transaction request")
		}
		item, err = s.itemRepository.GetItemByID(fmt.Sprintf("%d", loanReq.ItemID))
		if err != nil {
			return nil, fmt.Errorf("item not found: %w", err)
		}
		newQuantity := item.Quantity - loanReq.Quantity
		if newQuantity < 0 {
			return &model.LoanTransaction{}, fmt.Errorf("insufficient quantity")
		}
		item.Quantity = newQuantity
		if err := s.itemRepository.UpdateItem(*item); err != nil {
			return nil, fmt.Errorf("failed to update item quantity: %w", err)
		}
		if err := s.logRepository.CreateLoanTransaction(*loanReq); err != nil {
			item.Quantity += loanReq.Quantity
			_ = s.itemRepository.UpdateItem(*item)
			return nil, fmt.Errorf("failed to create loan transaction log: %w", err)
		}
		return loanReq, nil

	default:
		return nil, fmt.Errorf("unsupported transaction type")
	}
}

func (s *TransactionService) GetTransactions(page, limit int) ([]model.UnifiedTransaction, error) {
	var transactions []model.UnifiedTransaction
	offset := (page - 1) * limit

	loanTransactions, err := s.logRepository.GetLoanTransactions(limit, offset)
	if err != nil {
		return nil, err
	}
	for _, loan := range loanTransactions {
		transaction := model.UnifiedTransaction{
			ID:                 loan.ID,
			TransactionType:    "loan",
			GlobalID:           fmt.Sprintf("loan_%d", loan.ID),
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

		if loan.Item != nil {
			transaction.Item = loan.Item
		}

		transactions = append(transactions, transaction)
	}

	inquiryTransactions, err := s.logRepository.GetInquiryTransactions(limit, offset)
	if err != nil {
		return nil, err
	}
	for _, inquiry := range inquiryTransactions {
		transaction := model.UnifiedTransaction{
			ID:                 inquiry.ID,
			GlobalID:           fmt.Sprintf("inquiry_%d", inquiry.ID),
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

		if inquiry.Item != nil {
			transaction.Item = inquiry.Item
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (s *TransactionService) GetTransactionByID(globalID string) (*model.UnifiedTransaction, error) {
	parts := strings.Split(globalID, "_")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid global ID format")
	}

	transactionType := parts[0]
	id, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid ID in global ID")
	}

	switch transactionType {
	case "loan":
		loan, err := s.logRepository.GetLoanTransactionByID(int(id))
		if err != nil {
			return nil, err
		}
		return loan, nil

	case "inquiry":
		inquiry, err := s.logRepository.GetInquiryTransactionByID(int(id))
		if err != nil {
			return nil, err
		}
		return inquiry, nil
	default:
		return nil, fmt.Errorf("unknown transaction type")
	}
}
