package model

type Item struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Name       string `json:"name"`
	Quantity   int    `json:"quantity"`
	CategoryID uint   `json:"category_id"`
	StorageID  uint   `json:"storage_id"`
}