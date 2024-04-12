package ports

import "pdi/order/internal/core/domain"

type ItemRepository interface {
	GetItem(itemID string) (*domain.Item, error)
}
