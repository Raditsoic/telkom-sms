package model

type Item struct {
	ID         uint     `gorm:"primaryKey" json:"id"`
	Name       string   `json:"name"`
	Quantity   int      `json:"quantity"`
	Shelf      string   `json:"shelf"`
	CategoryID uint     `json:"category_id"`
	Category   Category `json:"-" gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	LoanTransactions      []LoanTransaction      `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	InquiryTransactions   []InquiryTransaction   `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	InsertionTransactions []InsertionTransaction `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
}

// Delete Item
type DeleteItemResponse struct {
	Message string `json:"message"`
	ID      string   `json:"id"`
}

// Update Item Name
type UpdateItemNameRequest struct {
	Name string `json:"name"`
}

type UpdateItemNameResponse struct {
	NewName string `json:"new_name"`
	OldName string `json:"old_name"`
	Message string `json:"message"`
	ID      string `json:"id"`
}

/*
	name: Biru
	quantity: 10
	category_id: 1
*/
