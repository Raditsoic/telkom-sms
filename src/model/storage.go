package model

type Storage struct {
	ID         int        `json:"id" gorm:"primaryKey"`
	Name       string     `json:"name"`
	Location   string     `json:"location"`
	Categories []Category `gorm:"foreignKey:StorageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"categories"`
}

type StorageByIDResponse struct {
	ID         int                       `json:"id"`
	Name       string                    `json:"name"`
	Location   string                    `json:"location"`
	Categories []StorageCategoryResponse `json:"categories"`
}

/*
 	"name":"ATK",
	"location":"TSO Manyar"
*/
