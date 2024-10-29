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
	if err := service.repository.CreateItem(item); err != nil {
		return nil, err
	}

	return item, nil
}

func (service *ItemService) GetItemByID(id string) (*model.Item, error) {
	return service.repository.GetItemByID(id)
}

func (service *ItemService) UpdateItem(item *model.Item) (*model.Item, error) {
	_, err := service.GetItemByID(strconv.FormatUint(uint64(item.ID), 10))
	if err != nil {
		return item, err
	}

	if err = service.repository.UpdateItem(*item); err != nil {
		return nil, err
	}

	return item, nil
}

func (service *ItemService) DeleteItem(id string) error {
	_, err := service.GetItemByID(id)
	if err != nil {
		return err
	}

	err = service.repository.DeleteItem(id)
	if err != nil {
		return err
	}

	return service.repository.DeleteItem(id)
}
