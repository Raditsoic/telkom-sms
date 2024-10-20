package model

type Storage struct {
	ID       int        `json:"id" gorm:"primaryKey"`
	Name     string     `json:"name"`
	Location string     `json:"location"`
	Categories []Category `gorm:"foreignKey:StorageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"categories"`
}

/*
 	"name":"ATK",
	"location":"TSO Manyar"
*/