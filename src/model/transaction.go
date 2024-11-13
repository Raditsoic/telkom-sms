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
	ID                 uint      `gorm:"primaryKey" json:"id"`
	UUID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();uniqueIndex" json:"uuid"`
	TransactionType    string    `json:"transaction_type"`
	EmployeeName       string    `json:"employee_name"`
	EmployeeDepartment string    `json:"employee_department"`
	EmployeePosition   string    `json:"employee_position"`
	Status             string    `json:"status"`
	Notes              string    `json:"notes"`
	Time               time.Time `json:"time"`
	Image              []byte    `json:"image"`
	ItemID             uint      `json:"item_id"`
	Item               Item      `json:"item" gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type GetAllTransactionsResponse struct {
	UUID               string     `json:"uuid"`
	TransactionType    string     `json:"transaction_type"`
	EmployeeName       string     `json:"employee_name"`
	EmployeeDepartment string     `json:"employee_department"`
	EmployeePosition   string     `json:"employee_position"`
	Quantity           int        `json:"quantity"`
	Status             string     `json:"status"`
	Notes              string     `json:"notes"`
	Time               time.Time  `json:"time"`
	Image              []byte     `json:"image"`
	ItemID             uint       `json:"item_id"`
	Item               *Item      `json:"item" gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LoanTime           *time.Time `json:"loan_time,omitempty"`
	ReturnTime         *time.Time `json:"return_time,omitempty"`
}

type UpdateTransactionResponse struct {
	Message string `json:"message"`
}
