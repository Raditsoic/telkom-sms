package model

import "time"

type LoanTransaction struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
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

type Transaction struct {
	ID                 uint       `json:"id"`
	TransactionType    string     `json:"transaction_type"`
	GlobalID           string     `json:"global_id"`
	EmployeeName       string     `json:"employee_name"`
	EmployeeDepartment string     `json:"employee_department"`
	EmployeePosition   string     `json:"employee_position"`
	Quantity           int        `json:"quantity"`
	Status             string     `json:"status"`
	Notes              string     `json:"notes"`
	Time               time.Time  `json:"time"`
	ItemID             uint       `json:"item_id"`
	Item               *Item      `json:"item" gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LoanTime           *time.Time `json:"loan_time,omitempty"`
	ReturnTime         *time.Time `json:"return_time,omitempty"`
}

type AllTransactionsRequest struct {
	ID                 uint       `json:"id"`
	TransactionType    string     `json:"transaction_type"`
	EmployeeName       string     `json:"employee_name"`
	EmployeeDepartment string     `json:"employee_department"`
	EmployeePosition   string     `json:"employee_position"`
	Quantity           int        `json:"quantity"`
	Status             string     `json:"status"`
	Notes              string     `json:"notes"`
	Time               time.Time  `json:"time"`
	ItemID             uint       `json:"item_id"`
	Items              []Item     `json:"items"`
	LoanTime           *time.Time `json:"loan_time,omitempty"`
	ReturnTime         *time.Time `json:"return_time,omitempty"`
}

type InsertionTransaction struct {
	ID                 uint      `json:"id"`
	TransactionType    string    `json:"transaction_type"`
	GlobalID           string    `json:"global_id"`
	EmployeeName       string    `json:"employee_name"`
	EmployeeDepartment string    `json:"employee_department"`
	EmployeePosition   string    `json:"employee_position"`
	Status             string    `json:"status"`
	Notes              string    `json:"notes"`
	Time               time.Time `json:"time"`
	ItemID             uint      `json:"item_id"`
	Item               Item      `json:"item" gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type UpdateTransactionRequest struct {
	TransactionID string `json:"global_id"`
	Status        string `json:"status"`
}

// type AdditionTransaction struct {
// 	ID                 uint      `json:"id"`
// 	TransactionType    string    `json:"transaction_type"`
// 	GlobalID           string    `json:"global_id"`
// 	EmployeeName       string    `json:"employee_name"`
// 	EmployeeDepartment string    `json:"employee_department"`
// 	EmployeePosition   string    `json:"employee_position"`
// 	Quantity           int       `json:"quantity"`
// }
