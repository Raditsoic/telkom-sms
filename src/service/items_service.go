package service

import (
	"strconv"

	"gtihub.com/raditsoic/telkom-storage-ms/src/database/repository"
	"gtihub.com/raditsoic/telkom-storage-ms/src/model"
)

type ItemService struct {
	itemRepository repository.ItemRepository
}

func NewItemService(repo repository.ItemRepository) *ItemService {
	return &ItemService{itemRepository: repo}
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
	return service.itemRepository.GetItems(limit, offset)
}

func (service *ItemService) CreateItem(item *model.Item) (*model.Item, error) {
	item, err := service.itemRepository.CreateItem(item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (service *ItemService) GetItemByID(id string) (*model.Item, error) {
	return service.itemRepository.GetItemByID(id)
}

func (service *ItemService) DeleteItem(id string) error {
	_, err := service.GetItemByID(id)
	if err != nil {
		return err
	}

	err = service.itemRepository.DeleteItem(id)
	if err != nil {
		return err
	}

	return service.itemRepository.DeleteItem(id)
}
