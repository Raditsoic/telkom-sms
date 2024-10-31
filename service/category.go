package service

import (
	"encoding/json"
	"strconv"

	"gtihub.com/raditsoic/telkom-storage-ms/database/repository"
	"gtihub.com/raditsoic/telkom-storage-ms/model"
)

type CategoryService struct {
	repository repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) *CategoryService {
	return &CategoryService{repository: repo}
}

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

func (s *CategoryService) CreateCategory(categoryData []byte) (*model.Category, error) {
	var category model.Category
	if err := json.Unmarshal(categoryData, &category); err != nil {
		return nil, err
	}

	if err := s.repository.CreateCategory(category); err != nil {
		return nil, err
	}

	return &category, nil
}

func (s *CategoryService) GetCategoryByID(id string) (*model.CategoryByIDResponse, error) {
	return s.repository.GetCategoryByID(id)
}

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
	}

	return response, nil
}

func (service *CategoryService) DeleteCategory(id int) error {
	_, err := service.repository.GetCategoryByID(strconv.Itoa(id))
	if err != nil {
		return err
	}

	return service.repository.DeleteCategory(id)
}
