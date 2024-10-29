package service

import (
	"strconv"

	"gtihub.com/raditsoic/telkom-storage-ms/database/repository"
	"gtihub.com/raditsoic/telkom-storage-ms/model"
)

type ItemService struct {
	repository repository.ItemRepository
}

func NewItemService(repo repository.ItemRepository) *ItemService {
	return &ItemService{repository: repo}
}

func (service *ItemService) GetItems(pageParam, limitParam string) ([]model.Item, error) {
	page, limit := 1, 10

	if parsedPage, err := strconv.Atoi(pageParam); err == nil && parsedPage > 0 {
		page = parsedPage
	}
	if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
		limit = parsedLimit
	}

	offset := (page - 1) * limit
	return service.repository.GetItems(limit, offset)
}

func (service *ItemService) CreateItem(item *model.Item) (*model.Item, error) {
	// Use GORM to create the item in the database
	if err := service.repository.CreateItem(item); err != nil {
		return nil, err
	}

	return item, nil
}

func (service *ItemService) GetItemByID(id string) (*model.Item, error) {
	return service.repository.GetItemByID(id)
}
