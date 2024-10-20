package model

type Transaction struct {
	ID                 uint   `gorm:"primaryKey" json:"id"`
	EmployeeName       string `json:"employee_name"`
	EmployeeDepartment string `json:"employee_department"`
	EmployeePosition   string `json:"employee_position"`
	Quantity           int    `json:"quantity"`
	Status             string `json:"status"`
	ItemID             uint   `json:"item_id"`
}
