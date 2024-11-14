package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gtihub.com/raditsoic/telkom-storage-ms/src/database/repository"
	"gtihub.com/raditsoic/telkom-storage-ms/src/model"
	"gtihub.com/raditsoic/telkom-storage-ms/src/utils"
)

type TransactionService struct {
	logRepository  repository.TransactionRepository
	itemRepository repository.ItemRepository
}

func NewTransactionService(log repository.TransactionRepository, item repository.ItemRepository) *TransactionService {
	return &TransactionService{logRepository: log, itemRepository: item}
}

func (s *TransactionService) GetTransactions(page, limit int) ([]model.GetAllTransactionsResponse, error) {
	var transactions []model.GetAllTransactionsResponse
	offset := (page - 1) * limit

	loanTransactions, err := s.logRepository.GetLoanTransactions(limit, offset)
	if err != nil {
		return nil, err
	}
	for _, loan := range loanTransactions {
		customUUID := fmt.Sprintf("%s_%s", "loan", loan.UUID)
		transaction := model.GetAllTransactionsResponse{
			UUID:               customUUID,
			TransactionType:    loan.TransactionType,
			EmployeeName:       loan.EmployeeName,
			EmployeeDepartment: loan.EmployeeDepartment,
			EmployeePosition:   loan.EmployeePosition,
			Quantity:           loan.Quantity,
			Status:             loan.Status,
			Notes:              loan.Notes,
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
		customUUID := fmt.Sprintf("%s_%s", "inquiry", inquiry.UUID)
		transaction := model.GetAllTransactionsResponse{
			UUID:               customUUID,
			TransactionType:    inquiry.TransactionType,
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
		customUUID := fmt.Sprintf("%s_%s", "insert", insertion.UUID)
		transaction := model.GetAllTransactionsResponse{
			UUID:               customUUID,
			TransactionType:    insertion.TransactionType,
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

func (s *TransactionService) CreateInsertionTransaction(insertion *model.InsertionTransaction) (*model.InsertionTransaction, error) {
	insertion.UUID = uuid.New()

	insertion.TransactionType = "insert"
	insertion.Time = time.Now()
	insertion.Status = "pending"

	if err := s.logRepository.CreateInsertionTransaction(insertion); err != nil {
		return nil, fmt.Errorf("failed to create insertion transaction: %w", err)
	}

	return insertion, nil
}

func (s *TransactionService) CreateLoanTransaction(loan model.LoanTransaction) (*model.LoanTransaction, error) {
	item, err := s.itemRepository.GetItemByID(fmt.Sprintf("%d", loan.ItemID))
	if err != nil {
		return nil, fmt.Errorf("item not found: %w", err)
	}

	loan.UUID = uuid.New()

	loan.TransactionType = "loan"
	loan.LoanTime = time.Now()
	loan.Time = time.Now()
	loan.Status = "pending"

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

	inquiry.UUID = uuid.New()
	inquiry.TransactionType = "inquiry"
	inquiry.Time = time.Now()
	inquiry.Status = "pending"

	if err := s.logRepository.CreateInquiryTransaction(inquiry); err != nil {
		item.Quantity += inquiry.Quantity
		_ = s.itemRepository.UpdateItem(*item)
		return nil, fmt.Errorf("failed to create inquiry transaction log: %w", err)
	}
	return &inquiry, nil
}

func (s *TransactionService) UpdateTransactionStatus(status, uuidStr string) (*model.UpdateTransactionResponse, error) {
	parts := strings.Split(uuidStr, "_")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid UUID format, expected type_UUID but got: %s", uuidStr)
	}

	transaction_type := parts[0]
	uuid, err := uuid.Parse(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %w", err)
	}

	status = strings.ToLower(status)

	if transaction_type == "loan" {
		loan, err := s.logRepository.GetLoanTransactionByUUID(uuid)
		if err != nil {
			return nil, fmt.Errorf("failed to get loan transaction: %w", err)
		}

		if loan.Status == "returned" {
			return nil, fmt.Errorf("loan transaction already returned")
		}

		item := loan.Item

		switch status {
		case "returned":
			item.Quantity += loan.Quantity
			if err := s.itemRepository.UpdateItem(*item); err != nil {
				return nil, fmt.Errorf("failed to update item quantity: %w", err)
			}
		case "approved":
			if item.Quantity < loan.Quantity {
				return nil, fmt.Errorf("insufficient quantity")
			}

			item.Quantity -= loan.Quantity
			if err := s.itemRepository.UpdateItem(*item); err != nil {
				return nil, fmt.Errorf("failed to update item quantity: %w", err)
			}
		case "rejected":
		default:
			return nil, fmt.Errorf("invalid status")
		}

		loan.Status = status
		if err := s.logRepository.UpdateLoanTransaction(loan); err != nil {
			return nil, fmt.Errorf("failed to update loan transaction: %w", err)
		}

		return &model.UpdateTransactionResponse{
			Message: fmt.Sprintf("Loan transaction %s successfully", status),
		}, nil
	} else if transaction_type == "inquiry" {
		inquiry, err := s.logRepository.GetInquiryTransactionByUUID(uuid)
		if err != nil {
			return nil, fmt.Errorf("failed to get inquiry transaction: %w", err)
		}

		item := inquiry.Item
		switch status {
		case "approved":
			if item.Quantity < inquiry.Quantity {
				return nil, fmt.Errorf("insufficient quantity")
			}

			item.Quantity -= inquiry.Quantity
			if err := s.itemRepository.UpdateItem(*item); err != nil {
				return nil, fmt.Errorf("failed to update item quantity: %w", err)
			}
		case "rejected":
		default:
			return nil, fmt.Errorf("invalid status")
		}

		inquiry.Status = status
		if err := s.logRepository.UpdateInquiryTransaction(inquiry); err != nil {
			return nil, fmt.Errorf("failed to update inquiry transaction: %w", err)
		}

		return &model.UpdateTransactionResponse{
			Message: fmt.Sprintf("Inquiry transaction %s successfully", status),
		}, nil
	} else if transaction_type == "insertion" {
		insertion, err := s.logRepository.GetInsertionTransactionByUUID(uuid)
		if err != nil {
			return nil, fmt.Errorf("failed to get insertion transaction: %w", err)
		}

		item := insertion.Item

		switch status {
		case "approved":
			_, err := s.itemRepository.CreateItem(&item)
			if err != nil {
				return nil, fmt.Errorf("failed to create item: %w", err)
			}
		case "rejected":
		default:
			return nil, fmt.Errorf("invalid status")
		}

		insertion.Status = status
		if err := s.logRepository.UpdateInsertionTransaction(insertion); err != nil {
			return nil, fmt.Errorf("failed to update insertion transaction: %w", err)
		}

		return &model.UpdateTransactionResponse{
			Message: fmt.Sprintf("Insertion transaction %s successfully", status),
		}, nil
	} else {
		return nil, utils.ErrTransactionType
	}
}
