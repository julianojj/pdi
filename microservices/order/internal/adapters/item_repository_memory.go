package adapters

import (
	"pdi/order/internal/core/domain"
	"pdi/order/internal/ports"
)

type ItemRepositoryMemory struct {
	items []*domain.Item
}

func NewItemRepositoryMemory() ports.ItemRepository {
	items := []*domain.Item{
		{
			ID:     "1",
			Name:   "Iphone",
			Amount: 500,
		},
		{
			ID:     "2",
			Name:   "Notebook Dell",
			Amount: 500,
		},
	}
	return &ItemRepositoryMemory{
		items,
	}
}

func (i *ItemRepositoryMemory) GetItem(itemID string) (*domain.Item, error) {
	for _, item := range i.items {
		if item.ID == itemID {
			return item, nil
		}
	}
	return nil, nil
}
