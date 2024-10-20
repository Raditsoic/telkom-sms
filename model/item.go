package model

type Item struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Name       string `json:"name"`
	Quantity   int    `json:"quantity"`
	CategoryID uint   `json:"category_id"`
}

/*
	name: Biru
	quantity: 10
	category_id: 1
*/