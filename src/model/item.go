package model

type Item struct {
	ID         uint     `gorm:"primaryKey" json:"id"`
	Name       string   `json:"name"`
	Quantity   int      `json:"quantity"`
	Shelf      string   `json:"shelf"`
	CategoryID uint     `json:"category_id"`
	Category   Category `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`

	LoanTransactions      []LoanTransaction      `gorm:"constraint:OnDelete:SET NULL;" json:"-"`
	InquiryTransactions   []InquiryTransaction   `gorm:"constraint:OnDelete:SET NULL;" json:"-"`
	InsertionTransactions []InsertionTransaction `gorm:"constraint:OnDelete:SET NULL;" json:"-"`
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
