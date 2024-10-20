package model

type Category struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `json:"name"`
	StorageID uint   `json:"storage_id"`
	Items     []Item `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items"`
}

/*
{
    "name":"Pulpen",
    "storage_id":1
}
*/