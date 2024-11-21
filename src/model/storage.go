package model

type Storage struct {
	ID         int        `json:"id" gorm:"primaryKey"`
	Name       string     `json:"name"`
	Location   string     `json:"location"`
	Categories []Category `gorm:"foreignKey:StorageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"categories"`
}


// Create Storage
type CreateStorageResponse struct {
	Message string `json:"message"`
	ID      string `json:"id"`
	Name    string `json:"name"`
}

// Get Storage By ID
type StorageByIDResponse struct {
	ID         int                       `json:"id"`
	Name       string                    `json:"name"`
	Location   string                    `json:"location"`
	Categories []StorageCategoryResponse `json:"categories"`
}

type StorageByIDResponseNoImage struct {
	ID         int                            `json:"id"`
	Name       string                         `json:"name"`
	Location   string                         `json:"location"`
	Categories []StorageCategoryNoImageResponse `json:"categories"`
}

// Delete Storage
type DeleteStorageResponse struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}

/*
 	"name":"ATK",
	"location":"TSO Manyar"
*/
