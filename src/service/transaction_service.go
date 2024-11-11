package service

import (
	"fmt"
	"time"

	"gtihub.com/raditsoic/telkom-storage-ms/src/database/repository"
	"gtihub.com/raditsoic/telkom-storage-ms/src/model"
)

type TransactionService struct {
	logRepository  repository.TransactionRepository
	itemRepository repository.ItemRepository
}

func NewTransactionService(log repository.TransactionRepository, item repository.ItemRepository) *TransactionService {
	return &TransactionService{logRepository: log, itemRepository: item}
}

func (s *TransactionService) CreateInsertionTransaction(insertion *model.InsertionTransaction) (*model.InsertionTransaction, error) {
	_, err := s.itemRepository.CreateItem(&insertion.Item)
	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	insertion.TransactionType = "Insertion"
	insertion.Time = time.Now()
	insertion.Status = "Pending"

	if err := s.logRepository.CreateInsertionTransaction(insertion); err != nil {
		return nil, fmt.Errorf("failed to create insertion transaction: %w", err)
	}

	return insertion, nil
}

func (s *TransactionService) GetInsertionTransactionByID(id uint) (*model.InsertionTransaction, error) {
	insertion, err := s.logRepository.GetInsertionTransactionByID(id)
	if err != nil {
		return nil, err
	}

	return insertion, nil
}

func (s *TransactionService) GetTransactions(page, limit int) ([]model.Transaction, error) {
	var transactions []model.Transaction
	offset := (page - 1) * limit

	loanTransactions, err := s.logRepository.GetLoanTransactions(limit, offset)
	if err != nil {
		return nil, err
	}
	for _, loan := range loanTransactions {
		transaction := model.Transaction{
			ID:                 loan.ID,
			TransactionType:    "Loan",
			EmployeeName:       loan.EmployeeName,
			EmployeeDepartment: loan.EmployeeDepartment,
			EmployeePosition:   loan.EmployeePosition,
			Quantity:           loan.Quantity,
			Status:             loan.Status,
			Notes:              loan.Notes,
			Image:              loan.Image,
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
		transaction := model.Transaction{
			ID:                 inquiry.ID,
			TransactionType:    "Inquiry",
			EmployeeName:       inquiry.EmployeeName,
			EmployeeDepartment: inquiry.EmployeeDepartment,
			EmployeePosition:   inquiry.EmployeePosition,
			Quantity:           inquiry.Quantity,
			Status:             inquiry.Status,
			Time:               inquiry.Time,
			Notes:              inquiry.Notes,
			Image:              inquiry.Image,
			ItemID:             inquiry.ItemID,
			Item:               inquiry.Item,
		}

		if inquiry.Item != nil {
			transaction.Item = inquiry.Item
		}

		transactions = append(transactions, transaction)
	}

	insertionTransactions, err := s.logRepository.GetInsertionTransactions(limit, offset)
	if err != nil {
		return nil, err
	}
	for _, insertion := range insertionTransactions {
		transaction := model.Transaction{
			ID:                 insertion.ID,
			TransactionType:    "Insertion",
			EmployeeName:       insertion.EmployeeName,
			EmployeeDepartment: insertion.EmployeeDepartment,
			EmployeePosition:   insertion.EmployeePosition,
			Status:             insertion.Status,
			Time:               insertion.Time,
			Notes:              insertion.Notes,
			Image:              insertion.Image,
			ItemID:             insertion.ItemID,
			Item:               &insertion.Item,
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (s *TransactionService) GetLoanTransactionByID(id uint) (*model.Transaction, error) {
	loan, err := s.logRepository.GetLoanTransactionByID(id)
	if err != nil {
		return nil, err
	}

	return loan, nil
}

func (s *TransactionService) GetInquiryTransactionByID(id uint) (*model.Transaction, error) {
	inquiry, err := s.logRepository.GetInquiryTransactionByID(id)
	if err != nil {
		return nil, err
	}

	return inquiry, nil
}

func (s *TransactionService) CreateLoanTransaction(loan model.LoanTransaction) (*model.LoanTransaction, error) {
	item, err := s.itemRepository.GetItemByID(fmt.Sprintf("%d", loan.ItemID))
	if err != nil {
		return nil, fmt.Errorf("item not found: %w", err)
	}

	loan.LoanTime = time.Now()
	loan.Time = time.Now()
	loan.Status = "Pending"

	if err := s.logRepository.CreateLoanTransaction(loan); err != nil {
		item.Quantity += loan.Quantity
		_ = s.itemRepository.UpdateItem(*item)
		return nil, fmt.Errorf("failed to create loan transaction log: %w", err)
	}
	return &loan, nil
}

func (s *TransactionService) CreateInquiryTransaction(inquiry model.InquiryTransaction) (*model.InquiryTransaction, error) {
	item, err := s.itemRepository.GetItemByID(fmt.Sprintf("%d", inquiry.ItemID))
	if err != nil {
		return nil, fmt.Errorf("item not found: %w", err)
	}

	inquiry.Time = time.Now()
	inquiry.Status = "Pending"

	if err := s.logRepository.CreateInquiryTransaction(inquiry); err != nil {
		item.Quantity += inquiry.Quantity
		_ = s.itemRepository.UpdateItem(*item)
		return nil, fmt.Errorf("failed to create inquiry transaction log: %w", err)
	}
	return &inquiry, nil
}

func (s *TransactionService) UpdateLoanTransaction(id uint, status string) (*model.Transaction, error) {
	loan, err := s.logRepository.GetLoanTransactionByID(id)
	if err != nil {
		return nil, fmt.Errorf("loan transaction not found: %w", err)
	}

	if loan.Status == "Returned" {
		return nil, fmt.Errorf("loan transaction already returned")
	}

	item := loan.Item

	switch status {
	case "Returned":
		item.Quantity += loan.Quantity
		if err := s.itemRepository.UpdateItem(*item); err != nil {
			return nil, fmt.Errorf("failed to update item quantity: %w", err)
		}
	case "Approved":
		if item.Quantity < loan.Quantity {
			return nil, fmt.Errorf("insufficient quantity")
		}

		item.Quantity -= loan.Quantity
		if err := s.itemRepository.UpdateItem(*item); err != nil {
			return nil, fmt.Errorf("failed to update item quantity: %w", err)
		}
	}

	loan.Status = status
	loanTransaction := &model.LoanTransaction{
		ID:                 loan.ID,
		EmployeeName:       loan.EmployeeName,
		EmployeeDepartment: loan.EmployeeDepartment,
		EmployeePosition:   loan.EmployeePosition,
		Quantity:           loan.Quantity,
		Status:             loan.Status,
		Time:               loan.Time,
		Notes:              loan.Notes,
		Image:              loan.Image,
		ItemID:             loan.ItemID,
		Item:               loan.Item,
		LoanTime:           *loan.LoanTime,
		ReturnTime:         *loan.ReturnTime,
	}
	if err := s.logRepository.UpdateLoanTransaction(loanTransaction); err != nil {
		return nil, fmt.Errorf("failed to update loan transaction: %w", err)
	}

	return loan, nil
}

func (s *TransactionService) UpdateInquiryTransaction(id uint, status string) (*model.Transaction, error) {
	inquiry, err := s.logRepository.GetInquiryTransactionByID(id)
	if err != nil {
		return nil, fmt.Errorf("inquiry transaction not found: %w", err)
	}

	item := inquiry.Item

	switch status {
	case "Approved":
		if item.Quantity < inquiry.Quantity {
			return nil, fmt.Errorf("insufficient quantity")
		}

		item.Quantity -= inquiry.Quantity
		if err := s.itemRepository.UpdateItem(*item); err != nil {
			return nil, fmt.Errorf("failed to update item quantity: %w", err)
		}
	}

	inquiry.Status = status
	inquiryTransaction := &model.InquiryTransaction{
		ID:                 inquiry.ID,
		EmployeeName:       inquiry.EmployeeName,
		EmployeeDepartment: inquiry.EmployeeDepartment,
		EmployeePosition:   inquiry.EmployeePosition,
		Quantity:           inquiry.Quantity,
		Status:             inquiry.Status,
		Time:               inquiry.Time,
		Image:              inquiry.Image,
		ItemID:             inquiry.ItemID,
		Item:               inquiry.Item,
		Notes:              inquiry.Notes,
	}
	if err := s.logRepository.UpdateInquiryTransaction(inquiryTransaction); err != nil {
		return nil, fmt.Errorf("failed to update inquiry transaction: %w", err)
	}

	return inquiry, nil
}

func (s *TransactionService) UpdateInsertionTransaction(id uint, status string) (*model.InsertionTransaction, error) {
	insertion, err := s.logRepository.GetInsertionTransactionByID(id)
	if err != nil {
		return nil, fmt.Errorf("insertion transaction not found: %w", err)
	}

	item := insertion.Item

	switch status {
	case "Approved":
		item.Quantity += insertion.Item.Quantity
		if err := s.itemRepository.UpdateItem(item); err != nil {
			return nil, fmt.Errorf("failed to update item quantity: %w", err)
		}
	}

	insertion.Status = status
	insertionTransaction := &model.InsertionTransaction{
		ID:                 insertion.ID,
		EmployeeName:       insertion.EmployeeName,
		EmployeeDepartment: insertion.EmployeeDepartment,
		EmployeePosition:   insertion.EmployeePosition,
		Status:             insertion.Status,
		Time:               insertion.Time,
		Notes:              insertion.Notes,
		Image:              insertion.Image,
		ItemID:             insertion.ItemID,
		Item:               insertion.Item,
	}
	if err := s.logRepository.UpdateInsertionTransaction(insertionTransaction); err != nil {
		return nil, fmt.Errorf("failed to update insertion transaction: %w", err)
	}

	return insertion, nil
}

// func (s *TransactionService) UpdateTransaction(payload model.UpdateTransactionRequest) (*model.LoanTransaction, error) {
// 	parts := strings.Split(payload.TransactionID, "_")
// 	if len(parts) != 2 {
// 		return nil, fmt.Errorf("invalid input format")
// 	}

// 	transType := parts[0]
// 	id, err := strconv.Atoi(parts[1])
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid ID format")
// 	}

// 	if transType == "loan" {
// 		return s.logRepository.UpdateLoanTransaction(id, payload.Status)
// 	} else if transType == "inquiry" {
// 		return s.logRepository.UpdateInquiryTransaction(id, payload.Status)
// 	}

// 	return nil, fmt.Errorf("invalid transaction")
// }
