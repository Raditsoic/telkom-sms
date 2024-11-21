package service

import (
	"fmt"
	"strconv"

	"gtihub.com/raditsoic/telkom-storage-ms/src/database/repository"
	"gtihub.com/raditsoic/telkom-storage-ms/src/model"
	"gtihub.com/raditsoic/telkom-storage-ms/src/utils"
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

func (service *ItemService) DeleteItem(id string) (*model.DeleteItemResponse, error) {
	if _, err := service.itemRepository.GetItemByID(id); err != nil {
		return nil, utils.ErrItemNotFound
	}

	err := service.itemRepository.DeleteItem(id)
	if err != nil {
		return nil, err
	}

	return &model.DeleteItemResponse{
		Message: "Item deleted successfully",
		ID:      id,
	}, nil
}

func (service *ItemService) UpdateItemName(id, new_name string) (*struct {
	NewName string `json:"new_name"`
	OldName string `json:"old_name"`
}, error) {
	item, err := service.itemRepository.GetItemByID(id)
	if err != nil {
		fmt.Println("Item not found")
		return nil, utils.ErrItemNotFound
	}

	old_name := item.Name

	item.Name = new_name

	err = service.itemRepository.UpdateItem(*item)
	if err != nil {
		fmt.Println("Failed to update item")
		return nil, err
	}

	update := struct {
		NewName string `json:"new_name"`
		OldName string `json:"old_name"`
	} {
		NewName: new_name,
		OldName: old_name,
	}

	return &update, nil
}
