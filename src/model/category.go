package model

type Category struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	Name      string  `json:"name"`
	StorageID uint    `json:"storage_id"`
	Image     []byte  `json:"image"`
	Items     []Item  `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items"`
	Storage   Storage `gorm:"foreignKey:StorageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"storage"`
}

type CreateCategoryRequest struct {
	Name      string `form:"name"`
	Image     []byte `json:"image"`
	StorageID uint   `form:"storage_id"`
}

type AllCategoryResponse struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	Image     []byte  `json:"image"`
	StorageID uint    `json:"storage_id"`
	Storage   Storage `json:"storage"`
}

type CategoryByIDResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Image     []byte `json:"image"`
	StorageID uint   `json:"storage_id"`
}

type CategoryWithItemsResponse struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	StorageID uint    `json:"storage_id"`
	Storage   Storage `json:"storage"`
	Items     []Item  `json:"items"`
	Image     []byte  `json:"image"`
}

type StorageCategoryResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
}

/*
{
    "name":"Pulpen",
    "storage_id":1
}
*/
