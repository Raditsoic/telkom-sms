package service

import (
	"fmt"

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
	// Create the new item
	item, err := s.itemRepository.CreateItem(&insertion.Item)
	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	// Set the item ID in the insertion transaction
	insertion.ItemID = item.ID

	// Create the insertion transaction
	if err := s.logRepository.CreateInsertionTransaction(insertion); err != nil {
		return nil, fmt.Errorf("failed to create insertion transaction: %w", err)
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
			TransactionType:    "Peminjaman",
			GlobalID:           fmt.Sprintf("loan_%d", loan.ID),
			EmployeeName:       loan.EmployeeName,
			EmployeeDepartment: loan.EmployeeDepartment,
			EmployeePosition:   loan.EmployeePosition,
			Quantity:           loan.Quantity,
			Status:             loan.Status,
			Time:               loan.Time,
			ItemID:             loan.ItemID,
			Item:               loan.Item,
			Notes:              loan.Notes,
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
			GlobalID:           fmt.Sprintf("inquiry_%d", inquiry.ID),
			TransactionType:    "Permintaan",
			EmployeeName:       inquiry.EmployeeName,
			EmployeeDepartment: inquiry.EmployeeDepartment,
			EmployeePosition:   inquiry.EmployeePosition,
			Quantity:           inquiry.Quantity,
			Status:             inquiry.Status,
			Time:               inquiry.Time,
			Notes:              inquiry.Notes,
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
			GlobalID:           fmt.Sprintf("insertion_%d", insertion.ID),
			TransactionType:    "Insertion",
			EmployeeName:       insertion.EmployeeName,
			EmployeeDepartment: insertion.EmployeeDepartment,
			EmployeePosition:   insertion.EmployeePosition,
			Status:             insertion.Status,
			Time:               insertion.Time,
			Notes:              insertion.Notes,
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

	newQuantity := item.Quantity - loan.Quantity
	if newQuantity < 0 {
		return &model.LoanTransaction{}, fmt.Errorf("insufficient quantity")
	}

	item.Quantity = newQuantity
	if err := s.itemRepository.UpdateItem(*item); err != nil {
		return nil, fmt.Errorf("failed to update item quantity: %w", err)
	}

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

	newQuantity := item.Quantity - inquiry.Quantity
	if newQuantity < 0 {
		return &model.InquiryTransaction{}, fmt.Errorf("insufficient quantity")
	}

	item.Quantity = newQuantity
	if err := s.itemRepository.UpdateItem(*item); err != nil {
		return nil, fmt.Errorf("failed to update item quantity: %w", err)
	}

	if err := s.logRepository.CreateInquiryTransaction(inquiry); err != nil {
		item.Quantity += inquiry.Quantity
		_ = s.itemRepository.UpdateItem(*item)
		return nil, fmt.Errorf("failed to create inquiry transaction log: %w", err)
	}
	return &inquiry, nil
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
