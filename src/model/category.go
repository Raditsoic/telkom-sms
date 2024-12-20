package model

type Category struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	Name      string  `json:"name"`
	StorageID uint    `json:"storage_id"`
	Image     []byte  `json:"image"`
	Items     []Item  `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"items"`
	Storage   Storage `gorm:"foreignKey:StorageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"storage"`
}

// Create Category
type CreateCategoryRequest struct {
	Name      string `form:"name"`
	Image     []byte `json:"image"`
	StorageID uint   `form:"storage_id"`
}

type CreateCategoryResponse struct {
	Message string `json:"message"`
	ID      string `json:"id"`
	Name    string `json:"name"`
}

// Get All Category
type AllCategoryResponse struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	Image     []byte  `json:"image"`
	StorageID uint    `json:"storage_id"`
	Storage   Storage `json:"storage"`
}

// Get Category By ID
type CategoryByIDResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Image     []byte `json:"image"`
	StorageID uint   `json:"storage_id"`
}

// Get Category With Items
type CategoryWithItemsResponse struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	StorageID uint    `json:"storage_id"`
	Storage   Storage `json:"storage"`
	Items     []Item  `json:"items"`
	Image     []byte  `json:"image"`
}

type StorageCategoryResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Image []byte `json:"image"`
}

type StorageCategoryNoImageResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}


// Update Category Name
type UpdateCategoryNameRequest struct {
	Name string `json:"name"`
}

type UpdateCategoryNameResponse struct {
	Message string `json:"message"`
	ID      string `json:"id"`
	NewName string `json:"new_name"`
	OldName string `json:"old_name"`
}

// Delete Category
type DeleteCategoryResponse struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}

/*
{
    "name":"Pulpen",
    "storage_id":1
}
*/
