package service

import (
	"fmt"
	"strconv"

	"gorm.io/gorm"
	"gtihub.com/raditsoic/telkom-storage-ms/src/database/repository"
	"gtihub.com/raditsoic/telkom-storage-ms/src/model"
	"gtihub.com/raditsoic/telkom-storage-ms/src/utils"
)

type CategoryService struct {
	repository repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) *CategoryService {
	return &CategoryService{repository: repo}
}

// Get All Categories
func (s *CategoryService) GetCategories(pageParam, limitParam string) ([]model.AllCategoryResponse, error) {
	page, limit := 1, 10

	if parsedPage, err := strconv.Atoi(pageParam); err == nil && parsedPage > 0 {
		page = parsedPage
	}
	if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
		limit = parsedLimit
	}

	offset := (page - 1) * limit
	return s.repository.GetCategories(limit, offset)
}

// Create Category
func (s *CategoryService) CreateCategory(category *model.Category) (*model.CreateCategoryResponse, error) {
	createdCategory, err := s.repository.CreateCategory(category)
	if err != nil {
		return nil, err
	}

	response := &model.CreateCategoryResponse{
		Message: "Category created successfully",
		ID:      strconv.FormatUint(uint64(createdCategory.ID), 10),
		Name:    createdCategory.Name,
	}

	return response, nil
}

// Get Category By ID
func (s *CategoryService) GetCategoryByID(id string) (*model.CategoryByIDResponse, error) {
	return s.repository.GetCategoryByIDStorage(id)
}

// Get Category With Items
func (service *CategoryService) GetCategoryWithItems(categoryID uint) (*model.CategoryWithItemsResponse, error) {
	category, err := service.repository.GetCategoryWithItems(categoryID)
	if err != nil {
		return nil, err
	}

	response := &model.CategoryWithItemsResponse{
		ID:        category.ID,
		Name:      category.Name,
		StorageID: category.StorageID,
		Storage:   category.Storage,
		Items:     category.Items,
		Image:     category.Image,
	}

	return response, nil
}

// Delete Category
func (service *CategoryService) DeleteCategory(id string) (*model.DeleteCategoryResponse, error) {
	if _, err := service.repository.GetCategoryByID(id); err != nil {
		return nil, gorm.ErrRecordNotFound
	}

	if err := service.repository.DeleteCategory(id); err != nil {
		return nil, err
	}

	response := &model.DeleteCategoryResponse{
		Message: "Category deleted successfully",
		ID:      id,
	}

	return response, nil
}

// Update Category Name
func (service *CategoryService) UpdateCategoryName(id, new_name string) (*model.UpdateCategoryNameResponse, error) {
	category, err := service.repository.GetCategoryByID(id)
	if err != nil {
		fmt.Println("Item not found")
		return nil, utils.ErrItemNotFound
	}

	old_name := category.Name
	category.Name = new_name

	if err := service.repository.UpdateCategory(*category); err != nil {
		return nil, err
	}

	return &model.UpdateCategoryNameResponse{
		Message: "Category name updated successfully",
		ID:      id,
		NewName: new_name,
		OldName: old_name,
	}, nil

}
