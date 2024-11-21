package repository

import (
	"fmt"

	"gorm.io/gorm"
	"gtihub.com/raditsoic/telkom-storage-ms/src/model"
)

type ItemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (repo *ItemRepository) CreateItem(item *model.Item) (*model.Item, error) {
	if err := repo.db.Create(item).Error; err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}
	return item, nil
}
func (repo *ItemRepository) GetItemByID(id string) (*model.Item, error) {
	var item model.Item
	if err := repo.db.Where("id = ?", id).First(&item).Error; err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	return &item, nil
}

func (repo *ItemRepository) GetItems(limit, offset int) ([]model.Item, error) {
	var items []model.Item

	if err := repo.db.Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}

	return items, nil
}

func (repo *ItemRepository) UpdateItem(item model.Item) error {
	if err := repo.db.Save(&item).Error; err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return nil
}

func (repo *ItemRepository) DeleteItem(id string) error {
	if err := repo.db.Where("id = ?", id).Delete(&model.Item{}).Error; err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}

func (repo *ItemRepository) AddItem(item *model.Item) error {
	if err := repo.db.Create(item).Error; err != nil {
		return fmt.Errorf("failed to add item: %w", err)
	}

	return nil
}

func (repo *ItemRepository) GetItemByName(name string) (*model.Item, error) {
	var item model.Item
	if err := repo.db.Where("name = ?", name).First(&item).Error; err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	return &item, nil
}


