package model

type Item struct {
	ID         uint     `gorm:"primaryKey" json:"id"`
	Name       string   `json:"name"`
	Quantity   int      `json:"quantity"`
	Shelf      string   `json:"shelf"`
	CategoryID uint     `json:"category_id"`
	Category   Category `json:"-" gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type DeleteItemResponse struct {
	Message string `json:"message"`
	ID      string   `json:"id"`
}

type UpdateItemNameRequestDTO struct {
	Name string `json:"name"`
}

/*
	name: Biru
	quantity: 10
	category_id: 1
*/
