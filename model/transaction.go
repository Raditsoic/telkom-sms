package model

import "time"

type Transaction struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	EmployeeName       string    `json:"employee_name"`
	EmployeeDepartment string    `json:"employee_department"`
	EmployeePosition   string    `json:"employee_position"`
	Quantity           int       `json:"quantity"`
	Status             string    `json:"status"`
	Time               time.Time `json:"time"`
	Type               string    `json:"type"`
	ItemID             uint      `json:"item_id"`
	Item               Item      `json:"item" gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type TransactionRequest struct {
	EmployeeName       string `json:"employee_name"`
	EmployeeDepartment string `json:"employee_department"`
	EmployeePosition   string `json:"employee_position"`
	Quantity           int    `json:"quantity"`
	Type               string `json:"type"`
	ItemID             uint   `json:"item_id"`
}
