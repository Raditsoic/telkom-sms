package model

type Category struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
	Item []Item `gorm:"foreignKey:CategoryID" json:"items"`
}
