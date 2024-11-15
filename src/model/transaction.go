package model

import (
	"time"

	"github.com/google/uuid"
)

type LoanTransaction struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	UUID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();uniqueIndex" json:"uuid"`
	TransactionType    string    `json:"transaction_type"`
	EmployeeName       string    `json:"employee_name"`
	EmployeeDepartment string    `json:"employee_department"`
	EmployeePosition   string    `json:"employee_position"`
	Quantity           int       `json:"quantity"`
	Status             string    `json:"status"`
	Time               time.Time `json:"time"`
	Notes              string    `json:"notes"`
	ItemID             uint      `json:"item_id"`
	Item               *Item     `gorm:"foreignKey:ItemID" json:"item"`
	LoanTime           time.Time `json:"loan_time"`
	ReturnTime         time.Time `json:"return_time"`
}

type InquiryTransaction struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	UUID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();uniqueIndex" json:"uuid"`
	TransactionType    string    `json:"transaction_type"`
	EmployeeName       string    `json:"employee_name"`
	EmployeeDepartment string    `json:"employee_department"`
	EmployeePosition   string    `json:"employee_position"`
	Quantity           int       `json:"quantity"`
	Status             string    `json:"status"`
	Notes              string    `json:"notes"`
	Time               time.Time `json:"time"`
	ItemID             uint      `json:"item_id"`
	Item               *Item     `gorm:"foreignKey:ItemID" json:"item"`
}

type InsertionTransaction struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	UUID               uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();uniqueIndex" json:"uuid"`
	TransactionType    string         `json:"transaction_type"`
	EmployeeName       string         `json:"employee_name"`
	EmployeeDepartment string         `json:"employee_department"`
	EmployeePosition   string         `json:"employee_position"`
	Status             string         `json:"status"`
	Notes              string         `json:"notes"`
	Time               time.Time      `json:"time"`
	Image              []byte         `json:"image"`
	ItemID             *uint          `json:"item_id"`                                                                     // Make it nullable
	Item               *Item          `json:"item" gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Make it nullable
	ItemRequest        ItemRequestDTO `gorm:"embedded;embeddedPrefix:item_request_" json:"item_request"`                   // Embed the request data
}

type CreateInsertionTransactionDTO struct {
	EmployeeName       string         `json:"employee_name" validate:"required"`
	EmployeeDepartment string         `json:"employee_department" validate:"required"`
	EmployeePosition   string         `json:"employee_position" validate:"required"`
	Notes              string         `json:"notes"`
	Image              []byte         `json:"image" validate:"required"`
	ItemRequest        ItemRequestDTO `json:"item_request" validate:"required"`
}

type ItemRequestDTO struct {
	Name       string `json:"name" validate:"required"`
	Quantity   int    `json:"quantity" validate:"required,gt=0"`
	Shelf      string `json:"shelf" validate:"required"`
	CategoryID uint   `json:"category_id" validate:"required"`
}

type InsertionTransactionRequest struct {
	EmployeeName       string `json:"employee_name"`
	EmployeeDepartment string `json:"employee_department"`
	EmployeePosition   string `json:"employee_position"`
	Notes              string `json:"notes"`
	ItemName           string `json:"item_name"`
	Quantity           int    `json:"quantity"`
	Shelf              string `json:"shelf"`
	CategoryID         uint   `json:"category_id"`
	Image              []byte `json:"image"`
}

type GetAllTransactionsResponse struct {
	UUID               string          `json:"uuid"`
	TransactionType    string          `json:"transaction_type"`
	EmployeeName       string          `json:"employee_name"`
	EmployeeDepartment string          `json:"employee_department"`
	EmployeePosition   string          `json:"employee_position"`
	Quantity           int             `json:"quantity"`
	Status             string          `json:"status"`
	Notes              string          `json:"notes"`
	Time               time.Time       `json:"time"`
	Image              *[]byte         `json:"image"`
	LoanTime           *time.Time      `json:"loan_time,omitempty"`
	ReturnTime         *time.Time      `json:"return_time,omitempty"`
	ItemRequest        *ItemRequestDTO `json:"item_request"`
}

type UpdateTransactionResponse struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}

type DeleteTransactionResponse struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}
