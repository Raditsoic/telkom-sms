package model

type Category struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `json:"name"`
	StorageID uint   `json:"storage_id"`
	Items     []Item `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items"`
	Storage   Storage `gorm:"foreignKey:StorageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"storage"`
}

type AllCategoryResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	StorageID uint   `json:"storage_id"`
	Storage   Storage `json:"storage"`
}

type CategoryByIDResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	StorageID uint   `json:"storage_id"`
}

type CategoryWithItemsResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	StorageID uint   `json:"storage_id"`
	Storage   Storage `json:"storage"`
	Items     []Item `json:"items"`
}
	
/*
{
    "name":"Pulpen",
    "storage_id":1
}
*/
