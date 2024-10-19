package model

type Storage struct {
	ID   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	Item []Item `gorm:"foreignKey:StorageID" json:"items"`
}
